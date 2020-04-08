package cmd

import (
	"github.com/spf13/cobra"

	"github.com/guumaster/hostctl/pkg/host"
)

// toggleCmd represents the enable command
var toggleCmd = &cobra.Command{
	Use:   "toggle",
	Short: "Change status of a profile on your hosts file.",
	Long: `
Alternates between on/off status of an existing profile.
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
		profile, _ := cmd.Flags().GetString("profile")

		src, _ := cmd.Flags().GetString("host-file")
		quiet, _ := cmd.Flags().GetBool("quiet")

		err := host.Toggle(src, profile)
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
	rootCmd.AddCommand(toggleCmd)
}
