package actions

import (
	"strings"

	"github.com/guumaster/cligger"
	"github.com/guumaster/hostctl/pkg/file"
	"github.com/spf13/cobra"
)

func newRemoveDomainsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "domains [profile] [domains] [flags]",
		Aliases: []string{"domain [profile] [domains] [flags]"},
		Short:   "Remove domains from your hosts file.",
		Long: `
Completely remove domains from your hosts file.
It cannot be undone unless you have a backup and restore it.
`,
		Args: commonCheckProfileOnly,
		RunE: func(cmd *cobra.Command, args []string) error {
			src, _ := cmd.Flags().GetString("host-file")
			quiet, _ := cmd.Flags().GetBool("quiet")

			name := args[0]
			domains := args[1:]

			h, err := file.NewFile(src)
			if err != nil {
				return err
			}

			removed, err := h.RemoveHostnames(name, domains)
			if err != nil {
				return err
			}

			err = h.Flush()
			if err != nil {
				return err
			}
			if !quiet {
				if removed {
					cligger.Success("Profile '%s' removed.\n", name)
				} else {
					cligger.Success("Domains '%s' removed.\n", strings.Join(args[1:], ", "))
				}
			}
			return nil
		},
		PostRunE: postRunListOnly,
	}
}
