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
 🄷🄾🅂🅃🄲🅃🄻

hostctl is a CLI tool to manage your hosts file with ease.
You can have multiple profiles, enable/disable exactly what
you need each time with a simple interface.
`,
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		host, _ := cmd.Flags().GetString("host-file")

		fmt.Printf("Using hosts file: %s\n", host)

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	defaultHostsFile := "/etc/hosts"

	if runtime.GOOS == "windows" {
		defaultHostsFile = `C:/Windows/System32/Drivers/etc/hosts`
	}

	envHostFile := os.Getenv("HOSTCTL_FILE")
	if envHostFile != "" {
		defaultHostsFile = envHostFile
	}

	rootCmd.PersistentFlags().StringP("profile", "p", "", "Choose a profile")
	rootCmd.PersistentFlags().String("host-file", defaultHostsFile, "Hosts file path")
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
