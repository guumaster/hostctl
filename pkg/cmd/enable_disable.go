package cmd

import (
	"github.com/spf13/cobra"

	"github.com/guumaster/hostctl/pkg/host"
)

// enableCmd represents the enable command
var enableCmd = &cobra.Command{
	Use:   "enable [profiles] [flags]",
	Short: "Enable a profile on your hosts file.",
	Long: `
Enables an existing profile.
It will be listed as "on" while it is enabled.
`,
	Args: commonCheckArgsWithAll,
	RunE: makeEnableDisable("enable"),
}

func makeEnableDisable(action string) func(cmd *cobra.Command, profiles []string) error {
	actionFn := func(h *host.File, profiles []string, only, all bool) error {
		switch {
		case only:
			return h.EnableOnly(profiles)
		case all:
			return h.EnableAll()
		default:
			return h.Enable(profiles)
		}
	}
	if action == "disable" {
		actionFn = func(h *host.File, profiles []string, only, all bool) error {
			switch {
			case only:
				return h.DisableOnly(profiles)
			case all:
				return h.DisableAll()
			default:
				return h.Disable(profiles)
			}
		}
	}

	return func(cmd *cobra.Command, profiles []string) error {
		src, _ := cmd.Flags().GetString("host-file")
		only, _ := cmd.Flags().GetBool("only")
		all, _ := cmd.Flags().GetBool("all")

		h, err := host.NewFile(src)
		if err != nil {
			return err
		}

		err = actionFn(h, profiles, only, all)
		if err != nil {
			return err
		}

		return h.Flush()
	}
}

// disableCmd represents the disable command
var disableCmd = &cobra.Command{
	Use:   "disable [profiles] [flags]",
	Short: "Disable a profile from your hosts file.",
	Long: `
Disable a profile from your hosts file without removing it.
It will be listed as "off" while it is disabled.
`,
	Args: commonCheckArgsWithAll,
	RunE: makeEnableDisable("disable"),
}

func init() {
	rootCmd.AddCommand(disableCmd)
	rootCmd.AddCommand(enableCmd)

	// NOTE: Added here to avoid circular references
	enableCmd.PostRunE = func(cmd *cobra.Command, args []string) error {
		return postActionCmd(cmd, args, disableCmd, true)
	}
	disableCmd.PostRunE = func(cmd *cobra.Command, args []string) error {
		return postActionCmd(cmd, args, enableCmd, true)
	}

	enableCmd.Flags().BoolP("all", "", false, "Enable all profiles")
	enableCmd.Flags().Bool("only", false, "Disable all other profiles")
	enableCmd.Flags().DurationP("wait", "w", -1, "Enables a profile for a specific amount of time")

	disableCmd.Flags().BoolP("all", "", false, "Disable all profiles")
	disableCmd.Flags().Bool("only", false, "Enable all other profiles")
	disableCmd.Flags().DurationP("wait", "w", -1, "Enables a profile for a specific amount of time")
}
