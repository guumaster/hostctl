package cmd

import (
	"errors"
	"github.com/guumaster/tablewriter"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/guumaster/hostctl/pkg/host"
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Creates a backup copy of your hosts file",
	Long: `
Creates a backup copy of your hosts file with the date in .YYYYMMDD 
as extension. 
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		profile, _ := cmd.Flags().GetString("profile")

		if profile != "" {
			return errors.New("backup can only be done to whole file. remove profile flag")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		src, _ := cmd.Flags().GetString("host-file")
		dst, _ := cmd.Flags().GetString("path")

		backupFile, err := host.BackupFile(src, dst)

		table := tablewriter.NewWriter(os.Stdout)
		table.Append([]string{backupFile})
		table.Render()

		return err
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	backupCmd.Flags().String("path", cwd, "A path to save the backup")
}
