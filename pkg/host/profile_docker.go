package host

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// DockerOptions contains parameters to sync with docker and docker-compose
type DockerOptions struct {
	Domain      string
	Network     string
	ComposeFile string
	ProjectName string
	KeepPrefix  bool
	Cli         *client.Client
}

func NewProfileFromDocker(ctx context.Context, opts *DockerOptions) (*Profile, error) {
	cli := opts.Cli
	if opts.Cli == nil {
		cli, err := client.NewEnvClient()
		if err != nil {
			return nil, err
		}
		defer cli.Close()
	}

	f := filters.NewArgs()
	f.Add("status", "running")

	networkID, err := getNetworkID(ctx, cli, opts)
	if err != nil {
		return nil, err
	}
	if networkID != "" {
		f.Add("network", networkID)
	}

	var composeServices []string
	if opts.ComposeFile != "" {
		composeServices, err = parseComposeFile(opts.ComposeFile, opts.ProjectName)
		if err != nil {
			return nil, err
		}
	}

	list, err := cli.ContainerList(ctx, types.ContainerListOptions{
		Filters: f,
	})
	if err != nil {
		return nil, err
	}

	profile := &Profile{
		Routes: map[string]*Route{},
	}

	for _, c := range list {
		for _, n := range c.NetworkSettings.Networks {
			if networkID != "" && n.NetworkID != networkID {
				continue
			}
			name := strings.Replace(c.Names[0], "/", "", -1)

			if len(composeServices) == 0 {
				name := fmt.Sprintf("%s.%s", name, opts.Domain)
				profile.AddRoute(n.IPAddress, name)
				continue
			}

			for _, c := range composeServices {
				match, err := regexp.MatchString(fmt.Sprintf("^%s(_[0-9]+)?", c), name)
				if err != nil {
					return nil, err
				}
				if match {
					name = fmt.Sprintf("%s.%s", name, opts.Domain)
					if !opts.KeepPrefix {
						name = strings.Replace(name, opts.ProjectName+"_", "", 1)
					}
					name := fmt.Sprintf("%s.%s", name, opts.Domain)
					profile.AddRoute(n.IPAddress, name)
				}
			}
		}
	}
	return profile, nil
}
