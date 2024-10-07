/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os/exec"
	"time"

	"github.com/spf13/cobra"
	"github.com/thin-edge/c8y-tedge/pkg/discovery"
)

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan for thin-edge.io devices",
	Long:  `Use mdns-sd / zeroconf to discover thin-edge.io instances`,

	RunE: func(cmd *cobra.Command, args []string) error {
		filter := &discovery.FilterOptions{}
		err := WithOptions(
			filter,
			func(options *discovery.FilterOptions) error {
				v, err := cmd.Flags().GetDuration("timeout")
				options.Timeout = v
				return err
			},
			func(options *discovery.FilterOptions) error {
				v, err := cmd.Flags().GetDuration("after")
				options.After = v
				return err
			},
			func(options *discovery.FilterOptions) error {
				v, err := cmd.Flags().GetString("pattern")
				options.Pattern = v
				return err
			},
			func(options *discovery.FilterOptions) error {
				v, err := cmd.Flags().GetBool("native")
				options.UseNative = v
				return err
			},
			func(options *discovery.FilterOptions) error {
				v, err := cmd.Flags().GetString("domain")
				options.Domain = v
				return err
			},
			func(options *discovery.FilterOptions) error {
				v, err := cmd.Flags().GetString("type")
				options.ServiceType = v
				return err
			},
		)
		if err != nil {
			return err
		}
		cmd.Printf("scan called. %v\n", filter.Timeout)

		filter.UseNative = commandExists("dns-sd")
		if filter.UseNative {
			cmd.Printf("Using dns-sd\n")
			if err := discovery.NativeDiscovery(*filter); err != nil {
				return err
			}
		} else {
			cmd.Printf("Using golang zeroconf\n")
			// if err := discovery.Discover(filter.ServiceType, filter.Domain, filter.Timeout); err != nil {
			if err := discovery.DiscoverHashicorp(filter.ServiceType, filter.Domain, filter.Timeout); err != nil {
				return err
			}
		}

		return nil
	},
}

type ScanOption func(options *discovery.FilterOptions) error

func WithOptions(filter *discovery.FilterOptions, options ...ScanOption) error {
	for _, opt := range options {
		if err := opt(filter); err != nil {
			return nil
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(scanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	scanCmd.Flags().String("pattern", "", "Filter by pattern. Only include the devices matching the given pattern")
	scanCmd.Flags().String("type", discovery.ThinEdgeServiceType, "Service type")
	scanCmd.Flags().String("domain", discovery.DefaultDomain, "Domain")
	scanCmd.Flags().DurationP("timeout", "t", 5*time.Second, "Timeout (duration)")
	scanCmd.Flags().Duration("after", 0, "Only process messages after the given duration")
	scanCmd.Flags().Bool("native", true, "Use native tools")
}
