package cmd

import (
	"github.com/spf13/cobra"

	"github.com/guumaster/cligger"

	"github.com/guumaster/hostctl/pkg/file"
)

func newRestoreCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "restore [flags]",
		Short: "Restore hosts file content from a backup file.",
		Long: `
Reads from a file and replace the content of your hosts file.

WARNING: the complete hosts file will be overwritten with the backup data.
`,
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			dst, _ := cmd.Flags().GetString("host-file")
			from, _ := cmd.Flags().GetString("from")
			quiet, _ := cmd.Flags().GetBool("quiet")

			h, err := file.NewFile(dst)
			if err != nil {
				return err
			}

			err = h.Restore(from)
			if err != nil {
				return err
			}

			if !quiet {
				cligger.Success("File '%s' restored.\n\n", from)
			}

			return nil
		},
		PostRunE: postRunListOnly,
	}
}
