package actions

import (
	"github.com/guumaster/cligger"
	"github.com/guumaster/hostctl/pkg/file"
	"github.com/spf13/cobra"
)

func newBackupCmd() *cobra.Command {
	return &cobra.Command{
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

			h, err := file.NewFile(src)
			if err != nil {
				return err
			}

			fname, err := h.Backup(dst)
			if err != nil {
				return err
			}

			if !quiet {
				cligger.Success("Backup '%s' created.\n", fname)
			}

			return nil
		},
		PostRunE: postRunListOnly,
	}
}
