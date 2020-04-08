package cmd

import (
	"github.com/spf13/cobra"

	"github.com/guumaster/hostctl/pkg/host"
)

var disableCmd *cobra.Command

func init() {

	// disableCmd represents the disable command
	disableCmd := &cobra.Command{
		Use:   "disable",
		Short: "Disable a profile from your hosts file.",
		Long: `
Disable a profile from your hosts file without removing it.
It will be listed as "off" while it is disabled.
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

			if all {
				profile = ""
			}
			return host.Disable(src, profile)

		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return postActionCmd(cmd, args, enableCmd)
		},
	}

	rootCmd.AddCommand(disableCmd)

	disableCmd.Flags().BoolP("all", "", false, "Disable all profiles")
	disableCmd.Flags().DurationP("wait", "w", -1, "Enables a profile for a specific amount of time")

}
