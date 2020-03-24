package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/guumaster/hostctl/pkg/host"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a profile from your hosts file.",
	Long: `
Completely remove a profile content from your hosts file.
It cannot be undone unless you have a backup and restore it.

If you want to remove a profile but would like to use it later, 
use 'hosts disable' instead.
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		profile, _ := cmd.Flags().GetString("profile")
		all, _ := cmd.Flags().GetBool("all")

		if !all && profile == "" {
			return host.MissingProfileError
		}

		if profile == "default" {
			return host.DefaultProfileError
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		profile, _ := cmd.Flags().GetString("profile")
		dst, _ := cmd.Flags().GetString("host-file")

		err := host.RemoveProfile(&host.RemoveProfileOptions{
			Dst:     dst,
			Profile: profile,
		})
		if err != nil {
			return err
		}

		fmt.Printf("Profile '%s' removed.", profile)

		return nil
	},
}

// removeDomainsCmd represents the remove command
var removeDomainsCmd = &cobra.Command{
	Use:   "domains",
	Short: "Remove domains from your hosts file.",
	Long: `
Completely remove domains from your hosts file.
It cannot be undone unless you have a backup and restore it.
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		profile, _ := cmd.Flags().GetString("profile")
		all, _ := cmd.Flags().GetBool("all")

		if !all && profile == "" {
			return host.MissingProfileError
		}

		if profile == "default" {
			return host.DefaultProfileError
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		profile, _ := cmd.Flags().GetString("profile")
		dst, _ := cmd.Flags().GetString("host-file")

		err := host.RemoveDomains(&host.RemoveDomainsOptions{
			Dst:     dst,
			Profile: profile,
			Domains: args,
		})
		if err != nil {
			return err
		}

		return host.ListProfiles(dst, &host.ListOptions{
			Profile: profile,
		})
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
	removeCmd.Flags().Bool("all", false, "Remove all profiles")

	removeCmd.AddCommand(removeDomainsCmd)
	removeDomainsCmd.Flags().Bool("all", false, "Remove all domains")
}
