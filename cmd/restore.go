package cmd

import (
	"github.com/spf13/cobra"

	"github.com/guumaster/hostctl/pkg/host"
)

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore hosts file content from a backup file.",
	Long: `
Reads from a file and replace the content of your hosts file.

WARNING: the complete hosts file will be overwritten with the backup data.
`,
	RunE: func(cmd *cobra.Command, args []string) error {

		dst, _ := cmd.Flags().GetString("host-file")

		from, _ := cmd.Flags().GetString("from")

		return host.RestoreFile(from, dst)
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)

	restoreCmd.Flags().String("from", "", "The file to restore from")
}
