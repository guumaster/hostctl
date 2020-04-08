package cmd

import (
	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync some system IPs with a profile.",
	Long: `
Reads IPs and names from some local system and sync it with a profile in your hosts file.
`,
}

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.AddCommand(syncDockerCmd)
	syncCmd.AddCommand(syncDockerComposeCmd)

	syncDockerComposeCmd.Flags().String("compose-file", "", "path to docker-compose.yml")
	syncDockerComposeCmd.Flags().String("project-name", "", "docker compose project name")
	syncDockerComposeCmd.Flags().Bool("prefix", false, "keep project name prefix from domain name")

	syncCmd.PersistentFlags().String("network", "", "Filter containers from a specific network")
	syncCmd.PersistentFlags().StringP("domain", "d", "loc", "domain where your docker containers will be added")
}
