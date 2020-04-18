package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version   = "dev"
	snapBuild string
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
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
			isSnapBuild := snapBuild == "yes"
			err := checkSnapRestrictions(cmd, isSnapBuild)
			if err != nil {
				return err
			}
			host, _ := cmd.Flags().GetString("host-file")

			showHostFile := host != getDefaultHostFile(isSnapBuild) || os.Getenv("HOSTCTL_FILE") != ""

			quiet := needsQuietOutput(cmd)

			if showHostFile && !quiet {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Using hosts file: %s\n", host)
			}

			return nil
		},
	}

	// set CLI version
	rootCmd.Version = version
	isSnapBuild := snapBuild == "yes"

	// rootCmd
	rootCmd.PersistentFlags().String("host-file", getDefaultHostFile(isSnapBuild), "Hosts file path")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Run command without output")
	rootCmd.PersistentFlags().Bool("raw", false, "Output without borders (same as -o raw)")
	rootCmd.PersistentFlags().StringP("out", "o", "table", "Output type (table|raw|markdown|json)")
	rootCmd.PersistentFlags().StringSliceP("column", "c", nil, "Columns to show on lists")

	registerCommands(rootCmd)

	return rootCmd
}

func needsQuietOutput(cmd *cobra.Command) bool {
	quiet, _ := cmd.Flags().GetBool("quiet")
	out, _ := cmd.Flags().GetString("out")

	if quiet {
		return true
	}

	switch {
	case quiet:
		return true
	case out == "json":
		return true
	case out == "md" || out == "markdown":
		return false

	default:
		return false
	}
}

func registerCommands(rootCmd *cobra.Command) {
	cwd, _ := os.Getwd()

	// Helper commands
	infoCmd := newInfoCmd()
	completionCmd := newCompletionCmd(rootCmd)
	genMdDocsCmd := newGenMdDocsCmd(rootCmd)
	genMdDocsCmd.Flags().String("path", cwd, "Path to save the docs files")

	// add
	addCmd, removeCmd := newAddRemoveCmd()

	addCmd.Flags().StringP("from", "f", "", "file to read")
	addCmd.PersistentFlags().DurationP("wait", "w", -1, "Enables a profile for a specific amount of time")
	addCmd.PersistentFlags().BoolP("uniq", "u", false, "only keep uniq domains per IP")

	// remove
	removeCmd.Flags().Bool("all", false, "Remove all profiles")

	// add domains

	addDomainsCmd, removeDomainsCmd := newAddRemoveDomainsCmd()
	addDomainsCmd.Flags().String("ip", "127.0.0.1", "domains ip")

	// replace
	replaceCmd := newReplaceCmd()
	replaceCmd.Flags().StringP("from", "f", "", "file to read")
	replaceCmd.Flags().BoolP("uniq", "u", false, "only keep uniq domains per IP")

	// toggle
	toggleCmd := newToggleCmd()
	toggleCmd.Flags().DurationP("wait", "w", -1, "Toggles a profile for a specific amount of time")

	// enable / disable
	enableCmd, disableCmd := newEnableDisableCmd()

	enableCmd.Flags().BoolP("all", "", false, "Enable all profiles")
	enableCmd.Flags().Bool("only", false, "Disable all other profiles")
	enableCmd.Flags().DurationP("wait", "w", -1, "Enables a profile for a specific amount of time")

	disableCmd.Flags().BoolP("all", "", false, "Disable all profiles")
	disableCmd.Flags().Bool("only", false, "Enable all other profiles")
	disableCmd.Flags().DurationP("wait", "w", -1, "Enables a profile for a specific amount of time")

	// backup
	backupCmd := newBackupCmd()
	backupCmd.Flags().String("path", cwd, "A path to save the backup")

	// restore
	restoreCmd := newRestoreCmd()
	restoreCmd.Flags().String("from", "", "The file to restore from")

	// sync
	syncCmd := newSyncCmd()

	syncCmd.PersistentFlags().String("network", "", "Filter containers from a specific network")
	syncCmd.PersistentFlags().StringP("domain", "d", "loc", "domain where your docker containers will be added")
	syncCmd.PersistentFlags().DurationP("wait", "w", -1, "Enables a profile for a specific amount of time")

	// sync docker
	syncDockerCmd := newSyncDockerCmd(removeCmd)

	// sync docker compose
	syncDockerComposeCmd := newSyncDockerComposeCmd(removeCmd)
	syncDockerComposeCmd.Flags().String("compose-file", "", "path to docker-compose.yml")
	syncDockerComposeCmd.Flags().String("project-name", "", "docker compose project name")
	syncDockerComposeCmd.Flags().Bool("prefix", false, "keep project name prefix from domain name")

	syncCmd.AddCommand(syncDockerCmd)
	syncCmd.AddCommand(syncDockerComposeCmd)

	// list
	listCmd := newListCmd()

	// status
	statusCmd := newStatusCmd()

	// register sub-commands
	addCmd.AddCommand(addDomainsCmd)
	removeCmd.AddCommand(removeDomainsCmd)

	// register all commands
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(backupCmd)
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(disableCmd)
	rootCmd.AddCommand(enableCmd)
	rootCmd.AddCommand(genMdDocsCmd)
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(replaceCmd)
	rootCmd.AddCommand(restoreCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(toggleCmd)
}
