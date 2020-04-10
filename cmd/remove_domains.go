package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/guumaster/hostctl/pkg/host"
)

// removeDomainsCmd represents the remove command
var removeDomainsCmd = &cobra.Command{
	Use:   "domains",
	Short: "Remove domains from your hosts file.",
	Long: `
Completely remove domains from your hosts file.
It cannot be undone unless you have a backup and restore it.
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		profile, _ := cmd.Flags().GetString("profile")

		if profile == "" {
			return host.MissingProfileError
		}

		if profile == "default" {
			return host.DefaultProfileError
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		profile, _ := cmd.Flags().GetString("profile")
		dst, _ := cmd.Flags().GetString("host-file")
		quiet, _ := cmd.Flags().GetBool("quiet")

		err := host.RemoveDomains(dst, profile, args)
		if err != nil {
			return err
		}
		if !quiet {
			fmt.Printf("Domains '%s' removed.\n\n", strings.Join(args, ", "))
		}
		return nil
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		return postActionCmd(cmd, args, nil)
	},
}
