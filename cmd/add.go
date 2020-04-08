package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/guumaster/hostctl/pkg/host"
)

// addCmd represents the fromFile command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add content to a profile in your hosts file.",
	Long: `
Reads from a file and set content to a profile in your hosts file.
If the profile already exists it will be added to it.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		profile, _ := cmd.Flags().GetString("profile")

		if profile == "" {
			return host.MissingProfileError
		}

		if profile == "default" {
			return host.DefaultProfileError
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		src, _ := cmd.Flags().GetString("host-file")
		from, _ := cmd.Flags().GetString("from")
		profile, _ := cmd.Flags().GetString("profile")

		if isPiped() {
			return host.AddFromReader(os.Stdin, &host.AddFromFileOptions{
				Dst:     src,
				Profile: profile,
				Reset:   false,
			})
		}

		return host.AddFromFile(&host.AddFromFileOptions{
			From:    from,
			Dst:     src,
			Profile: profile,
			Reset:   false,
		})
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		return postActionCmd(cmd, args, removeCmd)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.AddCommand(addDomainsCmd)

	addCmd.Flags().StringP("from", "f", "", "file to read")
	addCmd.PersistentFlags().DurationP("wait", "w", -1, "Enables a profile for a specific amount of time")

	addDomainsCmd.Flags().String("ip", "127.0.0.1", "domains ip")
}
