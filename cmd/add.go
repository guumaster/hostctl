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

		err := host.AddFromFile(&host.AddFromFileOptions{
			From:    from,
			Dst:     h,
			Profile: profile,
			Reset:   false,
		})
		if err != nil {
			return err
		}

		return host.ListProfiles(src, &host.ListOptions{
			Profile: profile,
		})
	},
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
}

// addDomainsCmd represents the fromFile command
var addDomainsCmd = &cobra.Command{
	Use:   "domains",
	Short: "Add content in your hosts file.",
	Long: `
Set content in your hosts file.
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
		ip, _ := cmd.Flags().GetString("ip")
		profile, _ := cmd.Flags().GetString("profile")
		h, _ := cmd.Flags().GetString("host-file")

		err := host.AddFromArgs(&host.AddFromArgsOptions{
			Domains: args,
			IP:      ip,
			Dst:     h,
			Profile: profile,
			Reset:   false,
		})
		if err != nil {
			return err
		}

		err = host.Enable(src, profile)
		if err != nil {
			return err
		}

		return host.ListProfiles(src, &host.ListOptions{
			Profile: profile,
		})
	},
}

func init() {
	rootCmd.AddCommand(addFromFileCmd)

	addFromFileCmd.Flags().StringP("from", "f", "", "file to read")

	addFromFileCmd.AddCommand(addDomainsCmd)

	addDomainsCmd.Flags().String("ip", "127.0.0.1", "domains ip")
}
