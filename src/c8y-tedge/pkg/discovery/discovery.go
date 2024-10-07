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

	"github.com/hashicorp/mdns"

	"github.com/grandcat/zeroconf"
)

// var ThinEdgeServiceType = "_thin-edge_mqtt._tcp"
var ThinEdgeServiceType = "_tedge._tcp"
var DefaultDomain = ".local"

func DiscoverHashicorp(serviceType string, domain string, timeout time.Duration) error {
	// Discover all services on the network (e.g. _workstation._tcp)
	var err error

	entriesCh := make(chan *mdns.ServiceEntry, 4)
	go func() {
		for entry := range entriesCh {
			fmt.Printf("Got new entry: %v\n", entry)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	log.Printf("Starting to scan. type=%s, domain=%s\n", serviceType, domain)
	defer cancel()

	// Start the lookup
	// mdns.Query(&mdns.QueryParam{
	// 	DisableIPv6: true,
	// })
	err = mdns.Lookup(serviceType, entriesCh)
	close(entriesCh)
	if err != nil {
		log.Fatalln("Failed to browse:", err.Error())
	}

	<-ctx.Done()
	// Wait some additional time to see debug messages on go routine shutdown.
	time.Sleep(1 * time.Second)
	return nil
}

func Discover(serviceType string, domain string, timeout time.Duration) error {
	// Discover all services on the network (e.g. _workstation._tcp)
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatalln("Failed to initialize resolver:", err.Error())
	}

	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			log.Println("Instance: ", entry.Instance)
			if device_id := ParseInstance(entry.Instance); device_id != "" {
				fmt.Printf("%s%s\n", device_id, DefaultDomain)
			}
		}
		log.Println("No more entries.")
	}(entries)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	log.Printf("Starting to scan. type=%s, domain=%s\n", serviceType, domain)
	defer cancel()
	err = resolver.Browse(ctx, serviceType, domain, entries)
	if err != nil {
		log.Fatalln("Failed to browse:", err.Error())
	}

	<-ctx.Done()
	// Wait some additional time to see debug messages on go routine shutdown.
	time.Sleep(1 * time.Second)
	return nil
}

func ParseInstance(instance string) (device_id string) {
	// The Instance Name is the device_id
	if !strings.Contains(instance, "(") {
		return instance
	}

	// Parse the hostname from the instance
	i := strings.Index(instance, "(")
	j := strings.Index(instance, ")")
	if i != -1 && j != -1 && j > i {
		return instance[i+1 : j]
	}
	return
}

func GetInstanceName(line string) (device_id string) {
	if !strings.Contains(line, "local.") {
		return
	}
	fields := strings.Fields(line)
	if len(fields) < 7 {
		return
	}
	instanceName := strings.Join(fields[6:], " ")

	// The Instance Name is the device_id
	if !strings.Contains(instanceName, "(") {
		return instanceName
	}

	// Parse the hostname from the instance
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
	After       time.Duration
	Pattern     string
	Timeout     time.Duration
	ServiceType string
	Domain      string
	UseNative   bool
}

func NativeDiscovery(options FilterOptions) error {
	// Calling dns-sd has an advantage as it does not require a special firewall
	// rule, so the user won't be prompted with a notification
	timeoutSec := fmt.Sprintf("%d", options.Timeout/time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cmd := exec.CommandContext(ctx, "dns-sd", "-t", timeoutSec, "-B", ThinEdgeServiceType, strings.Trim(DefaultDomain, "."))
	// cmd := exec.CommandContext(ctx, "dns-sd", "-t", timeoutSec, "-B", ThinEdgeServiceType, DefaultDomain)
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
