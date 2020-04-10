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
		return host.Toggle(src, profile)
	},
}

func init() {
	rootCmd.AddCommand(toggleCmd)

	// NOTE: Added here to avoid circular references
	toggleCmd.PostRunE = func(cmd *cobra.Command, args []string) error {
		return postActionCmd(cmd, args, toggleCmd)
	}

	toggleCmd.Flags().DurationP("wait", "w", -1, "Toggles a profile for a specific amount of time")
}
