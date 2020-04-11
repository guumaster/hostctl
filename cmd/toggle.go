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
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return host.MissingProfileError
		}
		if err := containsDefault(args); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, profiles []string) error {
		src, _ := cmd.Flags().GetString("host-file")

		h, err := host.NewFile(src)
		if err != nil {
			return err
		}

		err = h.Toggle(profiles)
		if err != nil {
			return err
		}

		return h.WriteTo(src)
	},
}

func init() {
	rootCmd.AddCommand(toggleCmd)

	// NOTE: Added here to avoid circular references
	toggleCmd.PostRunE = func(cmd *cobra.Command, args []string) error {
		return postActionCmd(cmd, args, toggleCmd, true)
	}

	toggleCmd.Flags().DurationP("wait", "w", -1, "Toggles a profile for a specific amount of time")
}
