/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/docker/docker/api/types/filters"
	"github.com/spf13/cobra"
	"github.com/thin-edge/c8y-tedge/pkg/certificates"
	"github.com/thin-edge/c8y-tedge/pkg/docker"
)

func IsValidHostname(s string) bool {
	for _, r := range s {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') && (r != '-') && (r != '.') {
			return false
		}
	}
	return true
}

func ExpandHomeDir(p string) string {
	homedir, err := os.UserHomeDir()

	if err != nil {
		return p
	}
	return strings.ReplaceAll(p, "~", homedir)
}

// bootstrapCmd represents the bootstrap command
var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap <device|container|compose-service>",
	Args:  cobra.MinimumNArgs(1),
	Short: "Bootstrap a thin-edge.io device",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		switch len(args) {
		case 0:
			return []string{"container", "service", "device"}, cobra.ShellCompDirectiveNoFileComp
		default:
			switch args[0] {
			case "service":
				dockerClient, err := docker.NewDockerRunner(args[0])
				if err != nil {
					return []string{err.Error()}, cobra.ShellCompDirectiveError
				}
				containers, err := dockerClient.ListContainers(context.Background(), filters.KeyValuePair{
					Key:   "label",
					Value: "com.docker.compose.service",
				})
				if err != nil {
					return []string{err.Error()}, cobra.ShellCompDirectiveError
				}
				return containers, cobra.ShellCompDirectiveNoFileComp
			case "container":
				dockerClient, err := docker.NewDockerRunner(args[0])
				if err != nil {
					return []string{err.Error()}, cobra.ShellCompDirectiveError
				}
				containers, err := dockerClient.ListContainers(context.Background())
				if err != nil {
					return []string{err.Error()}, cobra.ShellCompDirectiveError
				}
				return containers, cobra.ShellCompDirectiveNoFileComp
			case "device":
				homedir, err := os.UserHomeDir()
				if err != nil {
					return nil, cobra.ShellCompDirectiveNoFileComp
				}
				knownHostsFile := filepath.Join(homedir, ".ssh", "known_hosts")
				file, err := os.Open(knownHostsFile)
				if err != nil {
					return []string{err.Error()}, cobra.ShellCompDirectiveNoFileComp
				}
				defer file.Close()
				fileScanner := bufio.NewScanner(file)
				hostnames := []string{}
				for fileScanner.Scan() {
					if hostname, _, ok := strings.Cut(fileScanner.Text(), " "); ok {
						if IsValidHostname(hostname) {
							hostnames = append(hostnames, hostname)
						}
					}
				}
				sort.Strings(hostnames)
				return hostnames, cobra.ShellCompDirectiveNoFileComp
			}
			return []string{}, cobra.ShellCompDirectiveNoFileComp
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("bootstrap called")

		if len(args) == 0 {
			return fmt.Errorf("not enough arguments")
		}

		caPublicFile := ExpandHomeDir(filepath.Join("~", "tedge-ca.crt"))
		caPrivateFile := ExpandHomeDir(filepath.Join("~", "tedge-ca.key"))

		deviceCertFile := "device.crt"
		deviceCert, err := os.Create(deviceCertFile)
		if err != nil {
			return err
		}
		defer deviceCert.Close()
		err = certificates.SignCertificateRequestFile("device.csr", caPublicFile, caPrivateFile, deviceCert)
		if err != nil {
			return err
		}

		cmd.Printf("Created cert: %s\n", deviceCertFile)
		return nil

		bootstrapType := args[0]

		cmd.Printf("Bootstrapping %s\n", bootstrapType)
		switch bootstrapType {
		case "container":
		case "device":
		case "compose-service":
		}

		runner, err := docker.NewDockerRunner(args[0])
		if err != nil {
			return err
		}
		result, err := runner.Execute("sh", "-c", "printf 'hello world'")
		if err != nil {
			return err
		}

		cmd.Printf("%s\n", result.Stdout())
		return err
	},
}

func init() {
	// bootstrapCmd.RegisterFlagCompletionFunc("")

	rootCmd.AddCommand(bootstrapCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// bootstrapCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// bootstrapCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
