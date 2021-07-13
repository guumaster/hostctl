package actions

import (
	"strings"

	"github.com/guumaster/cligger"
	"github.com/guumaster/hostctl/pkg/file"
	"github.com/spf13/cobra"
)

func newRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "remove [profiles] [flags]",
		Aliases: []string{"rm [profiles] [flags]"},
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

			h, err := file.NewFile(src)
			if err != nil {
				return err
			}

			if all {
				profiles = h.GetProfileNames()
			}
			if len(profiles) == 0 {
				return ErrEmptyProfiles
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
				cligger.Success("Profile(s) '%s' removed.\n\n", strings.Join(profiles, ", "))
			}

			return nil
		},
	}
}
