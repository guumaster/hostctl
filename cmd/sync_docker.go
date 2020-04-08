package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/guumaster/hostctl/pkg/host"
)

// syncDockerCmd represents the sync docker command
var syncDockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Sync your Docker containers IPs with a profile.",
	Long: `
Reads from Docker the list of containers and add names and IPs to a profile in your hosts file.
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		profile, _ := cmd.Flags().GetString("profile")

		if profile == "" {
			return host.MissingProfileError
		}

		if profile == "default" {
			return host.DefaultProfileError
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		hostFile, _ := cmd.Flags().GetString("host-file")
		profile, _ := cmd.Flags().GetString("profile")
		domain, _ := cmd.Flags().GetString("domain")
		network, _ := cmd.Flags().GetString("network")

		ctx := context.Background()
		return host.AddFromDocker(ctx, &host.AddFromDockerOptions{
			Dst:     hostFile,
			Domain:  domain,
			Profile: profile,
			Watch:   false,
			Docker: &host.DockerOptions{
				Network: network,
			},
		})
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		return postActionCmd(cmd, args, removeCmd)
	},
}
