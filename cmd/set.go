package cmd

import (
	"github.com/spf13/cobra"

	"github.com/guumaster/hostctl/pkg/host"
)

// setFromFileCmd represents the setFromFile command
var setFromFileCmd = &cobra.Command{
	Use:   "set",
	Short: "Set content to a profile in your hosts file.",
	Long: `
Reads from a file and set content to a profile in your hosts file.
If the profile already exists it will be overwritten.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		from, _ := cmd.Flags().GetString("from")
		profile, _ := cmd.Flags().GetString("profile")

		h, _ := cmd.Flags().GetString("host-file")

		return host.AddFromFile(from, h, profile, true)
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		profile, _ := cmd.Flags().GetString("profile")
		return host.CheckProfile(profile)
	},
}

func init() {
	rootCmd.AddCommand(setFromFileCmd)

	setFromFileCmd.Flags().StringP("from", "f", "", "file to read")
}
