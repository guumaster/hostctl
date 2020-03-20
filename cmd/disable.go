package cmd

import (
	"github.com/spf13/cobra"

	"github.com/guumaster/hostctl/pkg/host"
)

// disableCmd represents the disable command
var disableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable a profile from your hosts file.",
	Long: `
Disable a profile from your hosts file without removing it.
It will be  listed as "off" while it is disabled.
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		profile, _ := cmd.Flags().GetString("profile")
		all, _ := cmd.Flags().GetBool("all")

		if !all && profile == "" {
			return host.MissingProfileError
		}

		if profile == "default" {
			return host.DefaultProfileError
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		src, _ := cmd.Flags().GetString("host-file")
		profile, _ := cmd.Flags().GetString("profile")

		all, _ := cmd.Flags().GetBool("all")

		var err error
		if all {
			profile = ""
		}
		err = host.Disable(src, profile)
		if err != nil {
			return err
		}

		return host.ListProfiles(src, &host.ListOptions{
			Profile: profile,
		})
	},
}

func init() {
	rootCmd.AddCommand(disableCmd)

	disableCmd.Flags().BoolP("all", "", false, "Disable all profiles")
}
