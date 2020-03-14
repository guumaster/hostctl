package cmd

import (
	"github.com/spf13/cobra"

	"github.com/guumaster/hostctl/pkg/host"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a profile from your hosts file.",
	Long: `
Completely remove a profile content from your hosts file.
It cannot be undone unless you have a backup and restore it.

If you want to remove a profile but would like to use it later, 
use 'hosts disable' instead.
`,
	RunE: func(cmd *cobra.Command, args []string) error {

		profile, _ := cmd.Flags().GetString("profile")

		dst, _ := cmd.Flags().GetString("host-file")

		return host.Remove(dst, profile)
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		profile, _ := cmd.Flags().GetString("profile")
		return host.CheckProfile(profile)
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	removeCmd.Flags().BoolP("all", "", false, "Remove all profiles")
}
