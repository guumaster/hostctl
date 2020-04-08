package cmd

import (
	"github.com/spf13/cobra"

	"github.com/guumaster/hostctl/pkg/host"
)

// setDomainsCmd represents the fromFile command
var setDomainsCmd = &cobra.Command{
	Use:   "domains",
	Short: "Add content in your hosts file.",
	Long: `
Set content in your hosts file.
If the profile already exists it will be added to it.`,
	PreRunE: func(cmd *cobra.Command, domains []string) error {
		profile, _ := cmd.Flags().GetString("profile")

		if profile == "" {
			return host.MissingProfileError
		}

		if profile == "default" {
			return host.DefaultProfileError
		}

		if len(domains) == 0 {
			return host.MissingDomainsError
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		src, _ := cmd.Flags().GetString("host-file")
		ip, _ := cmd.Flags().GetString("ip")
		profile, _ := cmd.Flags().GetString("profile")
		h, _ := cmd.Flags().GetString("host-file")
		quiet, _ := cmd.Flags().GetBool("quiet")

		err := host.AddFromArgs(&host.AddFromArgsOptions{
			Domains: args,
			IP:      ip,
			Dst:     h,
			Profile: profile,
			Reset:   true,
		})
		if err != nil {
			return err
		}

		if quiet {
			return nil
		}
		return host.ListProfiles(src, &host.ListOptions{
			Profile: profile,
		})
	},
}
