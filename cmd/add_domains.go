package cmd

import (
	"github.com/spf13/cobra"

	"github.com/guumaster/hostctl/pkg/host"
)

// addDomainsCmd represents the fromFile command
var addDomainsCmd = &cobra.Command{
	Use:   "domains",
	Short: "Add content in your hosts file.",
	Long: `
Set content in your hosts file.
If the profile already exists it will be added to it.`,
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
		src, _ := cmd.Flags().GetString("host-file")
		ip, _ := cmd.Flags().GetString("ip")
		profile, _ := cmd.Flags().GetString("profile")

		err := host.AddFromArgs(&host.AddFromArgsOptions{
			Domains: args,
			IP:      ip,
			Dst:     src,
			Profile: profile,
			Reset:   false,
		})
		if err != nil {
			return err
		}

		return host.Enable(src, profile)
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		return postActionCmd(cmd, args, removeDomainsCmd)
	},
}
