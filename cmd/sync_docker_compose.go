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
		prefix, _ := cmd.Flags().GetBool("prefix")

		compose, err := getComposeInfo(cmd)
		if err != nil {
			return err
		}

		if profile == "" && compose.ProjectName == "" {
			return host.MissingProfileError
		}

		if profile == "" {
			profile = compose.ProjectName
			_ = cmd.Flags().Set("profile", profile)
		}

		if domain == "" {
			domain = fmt.Sprintf("%s.loc", profile)
		}

		ctx := context.Background()

		return host.AddFromDocker(ctx, &host.AddFromDockerOptions{
			Dst:     hostFile,
			Domain:  domain,
			Profile: profile,
			Watch:   false,
			Docker: &host.DockerOptions{
				ComposeFile: compose.File,
				ProjectName: compose.ProjectName,
				Network:     network,
				KeepPrefix:  prefix,
			},
		})
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		return postActionCmd(cmd, args, removeCmd)
	},
}

type ComposeInfo struct {
	ProjectName string
	File        string
}

func getComposeInfo(cmd *cobra.Command) (*ComposeInfo, error) {
	name, _ := cmd.Flags().GetString("project-name")
	f, _ := cmd.Flags().GetString("compose-file")

	if f == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		f = path.Join(cwd, "docker-compose.yml")
	}

	if name == "" {
		reg := regexp.MustCompile("[^a-z0-9-]+")
		name = path.Base(path.Dir(f))
		name = strings.ToLower(name)
		name = reg.ReplaceAllString(name, "")
	}
	return &ComposeInfo{
		ProjectName: name,
		File:        f,
	}, nil
}
