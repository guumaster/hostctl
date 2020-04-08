package cmd

import (
	"os"

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
		quiet, _ := cmd.Flags().GetBool("quiet")

		h, _ := cmd.Flags().GetString("host-file")

		var err error
		if isPiped() {
			err = host.AddFromReader(os.Stdin, &host.AddFromFileOptions{
				Dst:     h,
				Profile: profile,
				Reset:   true,
			})
		} else {
			err = host.AddFromFile(&host.AddFromFileOptions{
				From:    from,
				Dst:     h,
				Profile: profile,
				Reset:   true,
			})
		}
		if err != nil {
			return err
		}

		if quiet {
			return nil
		}
		return host.ListProfiles(src, &host.ListOptions{
			Profile: profile,
		})
	},
}

func init() {
	rootCmd.AddCommand(setFromFileCmd)

	setFromFileCmd.Flags().StringP("from", "f", "", "file to read")

	setFromFileCmd.AddCommand(setDomainsCmd)

	setDomainsCmd.Flags().String("ip", "127.0.0.1", "domains ip")
}
