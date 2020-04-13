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
	Args: commonCheckArgsWithAll,
	RunE: func(cmd *cobra.Command, profiles []string) error {
		src, _ := cmd.Flags().GetString("host-file")
		enableOnly, _ := cmd.Flags().GetBool("only")
		all, _ := cmd.Flags().GetBool("all")

		h, err := host.NewFile(src)
		if err != nil {
			return err
		}

		if enableOnly {
			err = h.EnableOnly(profiles)
		} else if all {
			err = h.EnableAll()
		} else {
			err = h.Enable(profiles)
		}
		if err != nil {
			return err
		}

		return h.Flush()
	},
}

func init() {
	rootCmd.AddCommand(enableCmd)

	// NOTE: Added here to avoid circular references
	enableCmd.PostRunE = func(cmd *cobra.Command, args []string) error {
		return postActionCmd(cmd, args, disableCmd, true)
	}

	enableCmd.Flags().BoolP("all", "", false, "Enable all profiles")
	enableCmd.Flags().Bool("only", false, "Disable all other profiles")
	enableCmd.Flags().DurationP("wait", "w", -1, "Enables a profile for a specific amount of time")
}
