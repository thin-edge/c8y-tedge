package discovery

import (
	"bufio"
	"cmp"
	"context"
	"fmt"
	"log"
	"os/exec"
	"slices"
	"strings"
	"time"

	"github.com/grandcat/zeroconf"
)

var ThinEdgeServiceType = "_thin-edge_mqtt._tcp"
var DefaultDomain = ".local"

func Discover(timeout time.Duration) error {
	// Discover all services on the network (e.g. _workstation._tcp)
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatalln("Failed to initialize resolver:", err.Error())
	}

	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			if device_id := GetInstanceName(entry.Instance); device_id != "" {
				fmt.Printf("%s%s\n", device_id, DefaultDomain)
			}
		}
		log.Println("No more entries.")
	}(entries)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err = resolver.Browse(ctx, ThinEdgeServiceType, "local.", entries)
	if err != nil {
		log.Fatalln("Failed to browse:", err.Error())
	}

	<-ctx.Done()
	return nil
}

func GetInstanceName(line string) (device_id string) {
	i := strings.Index(line, "(")
	if i == -1 {
		return
	}
	j := strings.Index(line, ")")
	if j > i {
		device_id = line[i+1 : j]
	}
	return
}

func Insert[T cmp.Ordered](ts []T, t T) []T {
	// Insert into sorted list
	i, _ := slices.BinarySearch(ts, t)
	return slices.Insert(ts, i, t)
}

type FilterOptions struct {
	After   time.Duration
	Pattern string
	Timeout time.Duration
}

func NativeDiscovery(options FilterOptions) error {
	// Calling dns-sd has an advantage as it does not require a special firewall
	// rule, so the user won't be prompted with a notification
	timeoutSec := fmt.Sprintf("%d", options.Timeout/time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cmd := exec.CommandContext(ctx, "dns-sd", "-t", timeoutSec, "-B", ThinEdgeServiceType, strings.Trim(DefaultDomain, "."))
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	cmd.Start()
	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)

	names := []string{}
	startWatching := time.Now().Add(options.After)
	for scanner.Scan() {
		m := scanner.Text()
		if device_id := GetInstanceName(m); device_id != "" {
			names = Insert[string](names, device_id)

			if options.After > 0 && time.Now().Before(startWatching) {
				// ignore
				continue
			}

			if options.Pattern != "" && !strings.Contains(strings.ToLower(device_id), strings.ToLower(options.Pattern)) {
				// ignore
				continue
			}

			fmt.Printf("%s%s\n", device_id, DefaultDomain)
		}
	}
	err = cmd.Wait()

	// fmt.Printf("Sorted list\n%s", strings.Join(names, "\n"))
	return err
}
