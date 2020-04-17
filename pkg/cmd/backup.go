package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/guumaster/hostctl/pkg/host"
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup [flags]",
	Short: "Creates a backup copy of your hosts file",
	Long: `
Creates a backup copy of your hosts file with the date in .YYYYMMDD
as extension.
`,
	Args: cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		src, _ := cmd.Flags().GetString("host-file")
		dst, _ := cmd.Flags().GetString("path")
		quiet, _ := cmd.Flags().GetBool("quiet")

		h, err := host.NewFile(src)
		if err != nil {
			return err
		}

		fname, err := h.Backup(dst)
		if err != nil {
			return err
		}

		if !quiet {
			fmt.Fprintf(cmd.OutOrStdout(), "Backup '%s' created.\n", fname)
		}

		return nil
	},
	PostRunE: postRunListOnly,
}

func init() {
	rootCmd.AddCommand(backupCmd)

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	backupCmd.Flags().String("path", cwd, "A path to save the backup")
}
