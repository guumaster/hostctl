package cmd

import (
	"github.com/spf13/cobra"

	"github.com/guumaster/hostctl/pkg/host"
)

// disableCmd represents the disable command
var disableCmd = &cobra.Command{
	Use:   "disable [profiles] [flags]",
	Short: "Disable a profile from your hosts file.",
	Long: `
Disable a profile from your hosts file without removing it.
It will be listed as "off" while it is disabled.
`,
	Args: commonCheckArgsWithAll,
	RunE: func(cmd *cobra.Command, profiles []string) error {
		src, _ := cmd.Flags().GetString("host-file")
		all, _ := cmd.Flags().GetBool("all")

		h, err := host.NewFile(src)
		if err != nil {
			return err
		}

		if all {
			err = h.DisableAll()
		} else {
			err = h.Disable(profiles)
		}
		if err != nil {
			return err
		}

		return h.Flush()
	},
}

func init() {
	rootCmd.AddCommand(disableCmd)

	// NOTE: Added here to avoid circular references
	disableCmd.PostRunE = func(cmd *cobra.Command, args []string) error {
		return postActionCmd(cmd, args, enableCmd, true)
	}

	disableCmd.Flags().BoolP("all", "", false, "Disable all profiles")
	disableCmd.Flags().DurationP("wait", "w", -1, "Enables a profile for a specific amount of time")
}
