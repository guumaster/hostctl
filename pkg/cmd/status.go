package cmd

import (
	"github.com/spf13/cobra"

	"github.com/guumaster/hostctl/pkg/host"
)

// statusCmd represents the list command
var statusCmd = &cobra.Command{
	Use:   "status [profiles] [flags]",
	Short: "Shows a list of profile names and statuses on your hosts file.",
	Long: `
Shows a list of unique profile names on your hosts file with its status.

The "default" profile is always on and will be skipped.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		src, _ := cmd.Flags().GetString("host-file")
		raw, _ := cmd.Flags().GetBool("raw")

		h, err := host.NewFile(src)
		if err != nil {
			return err
		}

		h.ProfileStatus(&host.ListOptions{
			Writer:   cmd.OutOrStdout(),
			Profiles: args,
			RawTable: raw,
		})

		return err
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	statusCmd.Flags().Bool("raw", false, "Output without table borders")
}
