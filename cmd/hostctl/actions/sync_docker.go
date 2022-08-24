package actions

import (
	"github.com/spf13/cobra"

	"github.com/guumaster/hostctl/pkg/docker"
	"github.com/guumaster/hostctl/pkg/file"
	"github.com/guumaster/hostctl/pkg/parser"
	"github.com/guumaster/hostctl/pkg/types"
)

func newSyncDockerCmd(removeCmd *cobra.Command, optsFn getOptionsFn) *cobra.Command {
	if optsFn == nil {
		optsFn = getDockerOptions
	}

	return &cobra.Command{
		Use:   "docker [profile] [flags]",
		Short: "Sync your Docker containers IPs with a profile.",
		Long: `
Reads from Docker the list of containers and add names and IPs to a profile in your hosts file.
`,
		Args: commonCheckArgs,
		RunE: func(cmd *cobra.Command, profiles []string) error {
			src, _ := cmd.Flags().GetString("host-file")

			opts, err := optsFn(cmd, nil)
			if err != nil {
				return err
			}

			p, err := parser.NewProfileFromDocker(opts)
			if err != nil {
				return err
			}

			h, err := file.NewFile(src)
			if err != nil {
				return err
			}

			p.Name = profiles[0]
			p.Status = types.Enabled

			err = h.ReplaceProfile(p)
			if err != nil {
				return err
			}

			return h.Flush()
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return postActionCmd(cmd, args, removeCmd, true)
		},
	}
}

func getDockerOptions(cmd *cobra.Command, _ []string) (*docker.Options, error) {
	domain, _ := cmd.Flags().GetString("domain")
	network, _ := cmd.Flags().GetString("network")

	if domain == "" {
		domain = "loc"
	}

	return &docker.Options{
		Domain:  domain,
		Network: network,
		Cli:     nil,
	}, nil
}
