package actions

import (
	"github.com/spf13/cobra"
)

func newSyncCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sync",
		Short: "Sync some system IPs with a profile.",
		Long: `
Reads IPs and names from some local system and sync it with a profile in your hosts file.
`,
	}
}
