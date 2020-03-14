package cmd

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/guumaster/hostctl/pkg/host"
)

// enableCmd represents the enable command
var enableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable a profile on your hosts file.",
	Long: `
Disable an existing profile from your hosts file without removing it.
It will be  listed as "on" while it is enabled.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		profile, _ := cmd.Flags().GetString("profile")

		all, _ := cmd.Flags().GetBool("all")
		if all {
			profile = ""
		}

		src, _ := cmd.Flags().GetString("host-file")

		err := host.Enable(src, profile)
		if err != nil {
			return err
		}

		return host.ListProfiles(src, &host.ListOptions{
			Profile: profile,
		})
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		profile, _ := cmd.Flags().GetString("profile")
		all, _ := cmd.Flags().GetBool("all")

		if !all && profile == "" {
			return errors.New("missing profile name")
		}

		return host.ValidProfile(profile)
	},
}

func init() {
	rootCmd.AddCommand(enableCmd)

	enableCmd.Flags().BoolP("all", "", false, "Enable all profiles")
}
