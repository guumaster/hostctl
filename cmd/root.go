package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hostctl",
	Short: "Manage your hosts file like a pro",
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
			fmt.Printf("Using hosts file: %s\n", host)
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

func getDefaultHostFile() string {
	envHostFile := os.Getenv("HOSTCTL_FILE")
	if envHostFile != "" {
		return envHostFile
	}

	if runtime.GOOS == "windows" {
		return `C:/Windows/System32/Drivers/etc/hosts`
	}

	return "/etc/hosts"
}

func init() {
	rootCmd.PersistentFlags().StringP("profile", "p", "", "Choose a profile")
	rootCmd.PersistentFlags().String("host-file", getDefaultHostFile(), "Hosts file path")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Run command without output")

	rootCmd.PersistentFlags().StringSliceP("column", "c", nil, "Columns to show on lists")
	rootCmd.PersistentFlags().Bool("raw", false, "Output without table borders")
}

// isPiped detect if there is any input through STDIN
func isPiped() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	notPipe := info.Mode()&os.ModeNamedPipe == 0
	return !notPipe || info.Size() > 0
}
