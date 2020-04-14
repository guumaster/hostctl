package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hostctl",
	Short: "Your dev tool to manage /etc/hosts like a pro",
	Long: `
    __                    __           __     __
   / /_   ____    _____  / /_  _____  / /_   / /
  / __ \ / __ \  / ___/ / __/ / ___/ / __/  / /
 / / / // /_/ / (__  ) / /_  / /__  / /_   / /
/_/ /_/ \____/ /____/  \__/  \___/  \__/  /_/


hostctl is a CLI tool to manage your hosts file with ease.
You can have multiple profiles, enable/disable exactly what
you need each time with a simple interface.
`,
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		host, _ := cmd.Flags().GetString("host-file")
		quiet, _ := cmd.Flags().GetBool("quiet")

		defaultHostsFile := getDefaultHostFile()
		if (host != defaultHostsFile || os.Getenv("HOSTCTL_FILE") != "") && !quiet {
			fmt.Fprintf(cmd.OutOrStdout(), "Using hosts file: %s\n", host)
		}

		return nil
	},
}

// Execute is the main entrypoint for CLI usage
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = version

	rootCmd.PersistentFlags().String("host-file", getDefaultHostFile(), "Hosts file path")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Run command without output")
	rootCmd.PersistentFlags().Bool("raw", false, "Output without table borders")
	rootCmd.PersistentFlags().StringSliceP("column", "c", nil, "Columns to show on lists")
}
