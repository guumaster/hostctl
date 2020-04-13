package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/guumaster/hostctl/pkg/host"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm"},
	Short:   "Remove a profile from your hosts file.",
	Long: `
Completely remove a profile content from your hosts file.
It cannot be undone unless you have a backup and restore it.

If you want to remove a profile but would like to use it later,
use 'hosts disable' instead.
`,
	Args: commonCheckArgsWithAll,
	RunE: func(cmd *cobra.Command, profiles []string) error {
		src, _ := cmd.Flags().GetString("host-file")
		quiet, _ := cmd.Flags().GetBool("quiet")
		all, _ := cmd.Flags().GetBool("all")

		h, err := host.NewFile(src)
		if err != nil {
			return err
		}

		if all {
			profiles = h.GetProfileNames()
		}
		if len(profiles) == 0 {
			return fmt.Errorf("no profiles to remove")
		}

		err = h.RemoveProfiles(profiles)
		if err != nil {
			return err
		}

		err = h.Flush()
		if err != nil {
			return err
		}

		if !quiet {
			fmt.Fprintf(cmd.OutOrStdout(), "Profile(s) '%s' removed.\n\n", strings.Join(profiles, ", "))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
	removeCmd.Flags().Bool("all", false, "Remove all profiles")

	removeCmd.AddCommand(removeDomainsCmd)
}
