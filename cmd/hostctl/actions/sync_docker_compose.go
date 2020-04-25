package actions

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"github.com/guumaster/hostctl/pkg/file"
	"github.com/guumaster/hostctl/pkg/profile"
	"github.com/guumaster/hostctl/pkg/types"
)

type composeInfo struct {
	ProjectName string
	File        string
}

type getOptionsFn func(cmd *cobra.Command, profiles []string) (*profile.DockerOptions, error)

func newSyncDockerComposeCmd(removeCmd *cobra.Command, getOptionsFn getOptionsFn) *cobra.Command {
	if getOptionsFn == nil {
		getOptionsFn = defaultGetOptions
	}

	return &cobra.Command{
		Use:   "docker-compose [profile] [flags]",
		Short: "Sync your docker-compose containers IPs with a profile.",
		Long: `
Reads from a docker-compose.yml file  the list of containers and add names and IPs to a profile in your hosts file.
`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("profile")

			if name == "default" {
				return types.ErrDefaultProfile
			}
			return nil
		},
		Args: commonCheckArgs,
		RunE: func(cmd *cobra.Command, profiles []string) error {
			src, _ := cmd.Flags().GetString("host-file")
			name := profiles[0]

			opts, err := getOptionsFn(cmd, profiles)
			if err != nil {
				return err
			}

			p, err := profile.NewProfileFromDockerCompose(opts)
			if err != nil {
				return err
			}

			h, err := file.NewFile(src)
			if err != nil {
				return err
			}

			p.Name = name
			p.Status = types.Enabled

			err = h.ReplaceProfile(p)
			if err != nil {
				return err
			}

			return h.Flush()
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return postActionCmd(cmd, args, removeCmd, false)
		},
	}
}

func defaultGetOptions(cmd *cobra.Command, profiles []string) (*profile.DockerOptions, error) {
	domain, _ := cmd.Flags().GetString("domain")
	network, _ := cmd.Flags().GetString("network")
	prefix, _ := cmd.Flags().GetBool("prefix")

	name := profiles[0]

	compose, err := getComposeInfo(cmd)
	if err != nil {
		return nil, err
	}

	if name == "" && compose.ProjectName == "" {
		return nil, types.ErrMissingProfile
	}

	if name == "" {
		name = compose.ProjectName
		profiles = append(profiles, name)
		cmd.SetArgs(profiles)
	}

	if domain == "" {
		domain = fmt.Sprintf("%s.loc", name)
	}

	f, err := os.Open(compose.File)
	if err != nil {
		return nil, err
	}

	return &profile.DockerOptions{
		Domain:      domain,
		Network:     network,
		ComposeFile: f,
		ProjectName: compose.ProjectName,
		KeepPrefix:  prefix,
		Cli:         nil,
	}, nil
}

func getComposeInfo(cmd *cobra.Command) (*composeInfo, error) {
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

	return &composeInfo{
		ProjectName: name,
		File:        f,
	}, nil
}
