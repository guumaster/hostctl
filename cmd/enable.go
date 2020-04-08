package cmd

import (
	"github.com/spf13/cobra"

	"github.com/guumaster/hostctl/pkg/host"
)

// enableCmd represents the enable command
var enableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable a profile on your hosts file.",
	Long: `
Enables an existing profile.
It will be listed as "on" while it is enabled.
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
		profile, _ := cmd.Flags().GetString("profile")

		all, _ := cmd.Flags().GetBool("all")
		if all {
			profile = ""
		}

		src, _ := cmd.Flags().GetString("host-file")
		enableOnly, _ := cmd.Flags().GetBool("only")
		quiet, _ := cmd.Flags().GetBool("quiet")

		var err error
		if enableOnly {
			err = host.EnableOnly(src, profile)
		} else {
			err = host.Enable(src, profile)
		}
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

func init() {
	rootCmd.AddCommand(enableCmd)

	enableCmd.Flags().BoolP("all", "", false, "Enable all profiles")
	enableCmd.Flags().Bool("only", false, "Disable all other profiles")
}
