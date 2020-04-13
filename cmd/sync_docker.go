package cmd

import (
	"context"
	"fmt"

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
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return host.MissingProfileError
		} else if len(args) > 1 {
			return fmt.Errorf("specify only one profile")
		}
		if err := containsDefault(args); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, profiles []string) error {
		src, _ := cmd.Flags().GetString("host-file")
		domain, _ := cmd.Flags().GetString("domain")
		network, _ := cmd.Flags().GetString("network")

		ctx := context.Background()

		p, err := host.NewProfileFromDocker(ctx, &host.DockerOptions{
			Domain:  domain,
			Network: network,
			Cli:     nil,
		})
		if err != nil {
			return err
		}

		h, err := host.NewFile(src)
		if err != nil {
			return err
		}

		p.Name = profiles[0]
		p.Status = host.Enabled

		err = h.AddProfile(*p)
		if err != nil {
			return err
		}

		return h.Flush()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		return postActionCmd(cmd, args, removeCmd, false)
	},
}
