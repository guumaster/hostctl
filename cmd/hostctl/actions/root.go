package actions

import (
	"os"

	"github.com/guumaster/cligger"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
)

// NewRootCmd creates the base command for hostctl.
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "hostctl",
		Short:   "Your dev tool to manage /etc/hosts like a pro",
		Version: version,
		Long: `
hostctl is a CLI tool to manage your hosts file with ease.
You can have multiple profiles, enable/disable exactly what
you need each time with a simple interface.
`,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			noColor, _ := cmd.Flags().GetBool("no-color")

			if noColor {
				cligger.DisableColor()
			}
			cligger.SetWriter(cmd.OutOrStdout())
			host, _ := cmd.Flags().GetString("host-file")

			showHostFile := host != getDefaultHostFile() || os.Getenv("HOSTCTL_FILE") != ""

			quiet := needsQuietOutput(cmd)
			isHelper := isHelperCmd(cmd)

			if showHostFile && !quiet && !isHelper {
				cligger.Info("Using hosts file: %s\n", host)
			}

			return nil
		},
	}

	// set CLI version
	rootCmd.Version = version

	// rootCmd
	rootCmd.PersistentFlags().String("host-file", getDefaultHostFile(), "Hosts file path")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Run command without output")
	rootCmd.PersistentFlags().Bool("raw", false, "Output without borders (same as -o raw)")
	rootCmd.PersistentFlags().StringP("out", "o", "table", "Output type (table|raw|markdown|json)")
	rootCmd.PersistentFlags().StringSliceP("column", "c", nil, "Column names to show on lists. comma separated")

	rootCmd.PersistentFlags().Bool("no-color", false, "force colorless output")

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
	addCmd.PersistentFlags().
		DurationP("wait", "w", -1, "Enables a profile for a specific amount of time. (example: 5m, 1h)")

	// remove
	removeCmd.Flags().Bool("all", false, "Remove all profiles")

	// add domains
	addDomainsCmd, removeDomainsCmd := newAddRemoveDomainsCmd()
	addDomainsCmd.Flags().String("ip", "127.0.0.1", "domains ip")

	// replace
	replaceCmd := newReplaceCmd()
	replaceCmd.Flags().StringP("from", "f", "", "file to read")

	// toggle
	toggleCmd := newToggleCmd()
	toggleCmd.Flags().DurationP("wait", "w", -1, "Toggles a profile for a specific amount of time")

	// enable / disable
	enableCmd, disableCmd := newEnableDisableCmd()

	waitUsage := "Enables a profile for a specific amount of time"

	enableCmd.Flags().BoolP("all", "", false, "Enable all profiles")
	enableCmd.Flags().Bool("only", false, "Disable all other profiles")
	enableCmd.Flags().DurationP("wait", "w", -1, waitUsage)

	disableCmd.Flags().BoolP("all", "", false, "Disable all profiles")
	disableCmd.Flags().Bool("only", false, "Enable all other profiles")
	disableCmd.Flags().DurationP("wait", "w", -1, waitUsage)

	// backup
	backupCmd := newBackupCmd()
	backupCmd.Flags().String("path", cwd, "A path to save the backup")

	// restore
	restoreCmd := newRestoreCmd()
	restoreCmd.Flags().String("from", "", "The file to restore from")

	// sync
	syncCmd := newSyncCmd()
	syncCmd.PersistentFlags().DurationP("wait", "w", -1, waitUsage)

	// sync docker
	syncDockerCmd := newSyncDockerCmd(removeCmd, nil)
	syncDockerCmd.Flags().String("network", "", "Filter containers from a specific network")
	syncDockerCmd.Flags().StringP("domain", "d", "loc", "domain where your docker containers will be added")

	// sync docker compose
	syncDockerComposeCmd := newSyncDockerComposeCmd(removeCmd, nil)
	syncDockerComposeCmd.Flags().String("network", "", "Filter containers from a specific network")
	syncDockerComposeCmd.Flags().StringP("domain", "d", "loc", "domain where your docker containers will be added")
	syncDockerComposeCmd.Flags().String("compose-file", "", "path to docker-compose.yml")
	syncDockerComposeCmd.Flags().String("project-name", "", "docker compose project name")
	syncDockerComposeCmd.Flags().Bool("prefix", false, "keep project name prefix from domain name")

	// list
	listCmd := newListCmd()

	// status
	statusCmd := newStatusCmd()

	// register sub-commands
	addCmd.AddCommand(addDomainsCmd)
	removeCmd.AddCommand(removeDomainsCmd)
	syncCmd.AddCommand(syncDockerCmd)
	syncCmd.AddCommand(syncDockerComposeCmd)

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
