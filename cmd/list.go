package cmd

import (
	"github.com/spf13/cobra"

	"github.com/guumaster/hostctl/pkg/host"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Shows a detailed list of profiles on your hosts file.",
	Long: `
Shows a detailed list of profiles on your hosts file with name, ip and host name.
You can filter by profile name.

The "default"" profile is all the content that is not handled by hostctl tool.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		profile, _ := cmd.Flags().GetString("profile")

		src, _ := cmd.Flags().GetString("host-file")

		err := host.ListProfiles(src, &host.ListOptions{
			Profile: profile,
		})

		return err
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
