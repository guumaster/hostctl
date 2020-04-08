package cmd

import (
	"context"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"github.com/guumaster/hostctl/pkg/host"
)

// syncDockerComposeCmd represents the sync docker command
var syncDockerComposeCmd = &cobra.Command{
	Use:   "docker-compose",
	Short: "Sync your docker-compose containers IPs with a profile.",
	Long: `
Reads from a docker-compose.yml file  the list of containers and add names and IPs to a profile in your hosts file.
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		profile, _ := cmd.Flags().GetString("profile")

		if profile == "default" {
			return host.DefaultProfileError
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		hostFile, _ := cmd.Flags().GetString("host-file")
		profile, _ := cmd.Flags().GetString("profile")
		domain, _ := cmd.Flags().GetString("domain")
		network, _ := cmd.Flags().GetString("network")
		composeFile, _ := cmd.Flags().GetString("compose-file")
		projectName, _ := cmd.Flags().GetString("project-name")
		prefix, _ := cmd.Flags().GetBool("prefix")
		quiet, _ := cmd.Flags().GetBool("quiet")

		if composeFile == "" {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			composeFile = path.Join(cwd, "docker-compose.yml")
		}

		if projectName == "" {
			projectName = guessProjectName(composeFile)
		}

		if profile == "" && projectName == "" {
			return host.MissingProfileError
		}

		if profile == "" {
			profile = projectName
		}

		if domain == "" {
			domain = fmt.Sprintf("%s.loc", profile)
		}

		ctx := context.Background()

		err := host.AddFromDocker(ctx, &host.AddFromDockerOptions{
			Dst:     hostFile,
			Domain:  domain,
			Profile: profile,
			Watch:   false,
			Docker: &host.DockerOptions{
				ComposeFile: composeFile,
				ProjectName: projectName,
				Network:     network,
				KeepPrefix:  prefix,
			},
		})
		if err != nil {
			return err
		}

		if quiet {
			return nil
		}
		return host.ListProfiles(hostFile, &host.ListOptions{
			Profile: profile,
		})
	},
}

func guessProjectName(composeFile string) string {
	reg := regexp.MustCompile("[^a-z0-9-]+")
	base := path.Base(path.Dir(composeFile))
	base = strings.ToLower(base)
	base = reg.ReplaceAllString(base, "")
	return base
}
