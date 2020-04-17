package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:    "info",
		Hidden: true,
		Run: func(cmd *cobra.Command, profiles []string) {
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Your dev tool to manage /etc/hosts like a pro!")
		},
	}
}
