package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	version   = "dev"
	snapBuild string
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
		err := checkSnapRestrictions(cmd, args)
		if err != nil {
			return err
		}
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

func getDefaultHostFile() string {
	// Snap confinement doesn't allow to read other than
	if runtime.GOOS == "linux" && snapBuild == "yes" {
		return "/etc/hosts"
	}

	envHostFile := os.Getenv("HOSTCTL_FILE")
	if envHostFile != "" {
		return envHostFile
	}

	if runtime.GOOS == "windows" {
		return `C:/Windows/System32/Drivers/etc/hosts`
	}

	return "/etc/hosts"
}

func checkSnapRestrictions(cmd *cobra.Command, _ []string) error {
	from, _ := cmd.Flags().GetString("from")
	src, _ := cmd.Flags().GetString("host-file")

	defaultSrc := getDefaultHostFile()

	if snapBuild != "yes" {
		return nil
	}
	if from != "" || src != defaultSrc {
		return fmt.Errorf("can't use --from or --host-file. Snap confinement restristrions doesn't allow to read other than /etc/hosts file")
	}
	return nil
}
