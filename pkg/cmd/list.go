package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/guumaster/hostctl/pkg/host"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [profiles] [flags]",
	Short: "Shows a detailed list of profiles on your hosts file.",
	Long: `
Shows a detailed list of profiles on your hosts file with name, ip and host name.
You can filter by profile name.

The "default" profile is all the content that is not handled by hostctl tool.
`,
	RunE: func(cmd *cobra.Command, profiles []string) error {
		src, _ := cmd.Flags().GetString("host-file")
		raw, _ := cmd.Flags().GetBool("raw")
		cols, _ := cmd.Flags().GetStringSlice("column")

		h, err := host.NewFile(src)
		if err != nil {
			return err
		}

		h.List(&host.ListOptions{
			Writer:   cmd.OutOrStdout(),
			Profiles: profiles,
			RawTable: raw,
			Columns:  cols,
		})
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.AddCommand(makeListStatusCmd(host.Enabled))
	listCmd.AddCommand(makeListStatusCmd(host.Disabled))
}

// makeListStatusCmd represents the list enabled command
var makeListStatusCmd = func(status host.ProfileStatus) *cobra.Command {
	cmd := ""
	alias := ""
	switch status {
	case host.Enabled:
		cmd = "enabled"
		alias = "on"
	case host.Disabled:
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
		RunE: func(cmd *cobra.Command, args []string) error {
			src, _ := cmd.Flags().GetString("host-file")
			raw, _ := cmd.Flags().GetBool("raw")
			cols, _ := cmd.Flags().GetStringSlice("column")

			h, err := host.NewFile(src)
			if err != nil {
				return err
			}

			h.List(&host.ListOptions{
				Writer:       cmd.OutOrStdout(),
				RawTable:     raw,
				Columns:      cols,
				StatusFilter: status,
			})

			return nil
		},
	}
}
