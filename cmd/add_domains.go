package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/guumaster/hostctl/pkg/host"
)

// addDomainsCmd represents the fromFile command
var addDomainsCmd = &cobra.Command{
	Use:     "domains [domains] [flag]",
	Aliases: []string{"domain"},
	Short:   "Add content in your hosts file.",
	Long: `
Set content in your hosts file.
If the profile already exists it will be added to it.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return host.MissingProfileError
		}
		if err := containsDefault(args); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		src, _ := cmd.Flags().GetString("host-file")
		ip, _ := cmd.Flags().GetString("ip")
		quiet, _ := cmd.Flags().GetBool("quiet")
		name := args[0]
		routes := args[1:]

		h, err := host.NewFile(src)
		if err != nil {
			return err
		}

		err = h.AddRoutes(name, ip, routes)
		if err != nil {
			return err
		}

		err = h.WriteTo(src)
		if err != nil {
			return err
		}
		if !quiet {
			fmt.Printf("Domains '%s' added.\n\n", strings.Join(args[1:], ", "))
		}
		return nil
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		return postActionCmd(cmd, args, removeDomainsCmd, true)
	},
}
