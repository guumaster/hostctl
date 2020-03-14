package cmd

import (
	"github.com/spf13/cobra"

	"github.com/guumaster/hostctl/pkg/host"
)

// addFromFileCmd represents the fromFile command
var addFromFileCmd = &cobra.Command{
	Use:   "add",
	Short: "Add content to a profile in your hosts file.",
	Long: `
Reads from a file and set content to a profile in your hosts file.
If the profile already exists it will be added to it.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		src, _ := cmd.Flags().GetString("host-file")
		from, _ := cmd.Flags().GetString("from")
		profile, _ := cmd.Flags().GetString("profile")

		h, _ := cmd.Flags().GetString("host-file")

		err := host.AddFromFile(from, h, profile, false)
		if err != nil {
			return err
		}

		return host.ListProfiles(src, &host.ListOptions{
			Profile: profile,
		})
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		profile, _ := cmd.Flags().GetString("profile")

		err := host.NotEmptyProfile(profile)
		if err != nil {
			return err
		}

		return host.ValidProfile(profile)
	},
}

func init() {
	rootCmd.AddCommand(addFromFileCmd)

	addFromFileCmd.Flags().StringP("from", "f", "", "file to read")
}
