package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/guumaster/hostctl/pkg/host"
)

// removeDomainsCmd represents the remove command
var removeDomainsCmd = &cobra.Command{
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

		h, err := host.NewFile(src)
		if err != nil {
			return err
		}

		removed, err := h.RemoveRoutes(name, domains)
		if err != nil {
			return err
		}

		err = h.Flush()
		if err != nil {
			return err
		}
		if !quiet {
			if removed {
				fmt.Fprintf(cmd.OutOrStdout(), "Profile '%s' removed.\n\n", name)
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "Domains '%s' removed.\n\n", strings.Join(args[1:], ", "))
			}
		}
		return nil
	},
	PostRunE: postRunListOnly,
}
