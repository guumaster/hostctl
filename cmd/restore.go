package cmd

import (
	"errors"
	"fmt"

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
	PreRunE: func(cmd *cobra.Command, args []string) error {
		profile, _ := cmd.Flags().GetString("profile")

		if profile != "" {
			return errors.New("restore can only be done to whole file. remove profile flag")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		dst, _ := cmd.Flags().GetString("host-file")
		from, _ := cmd.Flags().GetString("from")

		err := host.RestoreFile(from, dst)
		if err != nil {
			return err
		}
		_ = host.ListProfiles(dst, &host.ListOptions{})

		fmt.Printf("Restore completed.")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)

	restoreCmd.Flags().String("from", "", "The file to restore from")
}
