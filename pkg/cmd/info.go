package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// infoCmd represents the enable command
var infoCmd = &cobra.Command{
	Use:    "info",
	Hidden: true,
	Run: func(cmd *cobra.Command, profiles []string) {
		fmt.Fprintln(cmd.OutOrStdout(), "Your dev tool to manage /etc/hosts like a pro!")
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
