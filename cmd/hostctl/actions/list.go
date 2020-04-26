package actions

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/guumaster/hostctl/pkg/file"
	"github.com/guumaster/hostctl/pkg/types"
)

func newListCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list [profiles] [flags]",
		Short: "Shows a detailed list of profiles on your hosts file.",
		Long: `
Shows a detailed list of profiles on your hosts file with name, ip and host name.
You can filter by profile name.

The "default" profile is all the content that is not handled by hostctl tool.
`,
		RunE: func(cmd *cobra.Command, profiles []string) error {
			src, _ := cmd.Flags().GetString("host-file")

			h, err := file.NewFile(src)
			if err != nil {
				return err
			}

			r := getRenderer(cmd, nil)

			h.List(r, &file.ListOptions{
				Profiles: profiles,
			})
			return nil
		},
	}

	listCmd.AddCommand(makeListStatusCmd(types.Enabled))
	listCmd.AddCommand(makeListStatusCmd(types.Disabled))

	return listCmd
}

func makeListStatusCmd(status types.Status) *cobra.Command {
	cmd := ""
	alias := ""

	switch status {
	case types.Enabled:
		cmd = "enabled"
		alias = "on"
	case types.Disabled:
		cmd = "disabled"
		alias = "off"
	}

	return &cobra.Command{
		Use:     cmd,
		Aliases: []string{alias},
		Short:   fmt.Sprintf("Shows list of %s profiles on your hosts file.", cmd),
		Long: fmt.Sprintf(`
Shows a detailed list of %s profiles on your hosts file with name, ip and host name.
`, cmd),
		RunE: func(cmd *cobra.Command, profiles []string) error {
			src, _ := cmd.Flags().GetString("host-file")

			h, err := file.NewFile(src)
			if err != nil {
				return err
			}

			r := getRenderer(cmd, nil)

			h.List(r, &file.ListOptions{
				Profiles:     profiles,
				StatusFilter: status,
			})

			return nil
		},
	}
}
