package actions

import (
	"strings"

	"github.com/guumaster/cligger"
	"github.com/guumaster/hostctl/pkg/file"
	"github.com/guumaster/hostctl/pkg/types"
	"github.com/spf13/cobra"
)

func newAddRemoveDomainsCmd() (*cobra.Command, *cobra.Command) {
	addDomainsCmd := &cobra.Command{
		Use:     "domains [profile] [domains] [flags]",
		Aliases: []string{"domain [profile] [domains] [flags]"},
		Short:   "Add content in your hosts file.",
		Long: `
Set content in your hosts file.
If the profile already exists it will be added to it.`,
		Args: commonCheckProfileOnly,
		RunE: func(cmd *cobra.Command, args []string) error {
			src, _ := cmd.Flags().GetString("host-file")
			ip, _ := cmd.Flags().GetString("ip")
			quiet, _ := cmd.Flags().GetBool("quiet")
			name := args[0]
			hostnames := args[1:]

			h, err := file.NewFile(src)
			if err != nil {
				return err
			}

			err = h.AddRoute(name, types.NewRoute(ip, hostnames...))
			if err != nil {
				return err
			}

			err = h.Flush()
			if err != nil {
				return err
			}
			if !quiet {
				cligger.Success("Domains '%s' added.\n", strings.Join(args[1:], ", "))
			}
			return nil
		},
	}

	removeDomainsCmd := newRemoveDomainsCmd()

	addDomainsCmd.PostRunE = func(cmd *cobra.Command, args []string) error {
		return postActionCmd(cmd, args, removeDomainsCmd, true)
	}

	return addDomainsCmd, removeDomainsCmd
}
