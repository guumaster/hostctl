package profile

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	dtypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"

	"github.com/guumaster/hostctl/pkg/docker"
	"github.com/guumaster/hostctl/pkg/types"
)

// DockerOptions contains parameters to sync with docker and docker-compose
type DockerOptions struct {
	Domain      string
	Network     string
	NetworkID   string
	ComposeFile string
	ProjectName string
	KeepPrefix  bool
	Cli         *client.Client
}

// NewProfileFromDocker creates a new profile from docker info
func NewProfileFromDocker(ctx context.Context, opts *DockerOptions) (*types.Profile, error) {
	var (
		networkID string
		err       error
	)

	cli := opts.Cli
	if cli == nil {
		cli, err = client.NewEnvClient()
		if err != nil {
			return nil, err
		}
	}
	defer cli.Close()

	if opts.NetworkID == "" && opts.Network != "" {
		networkID, err = docker.GetNetworkID(ctx, cli, opts.Network)
		if err != nil {
			return nil, err
		}

		opts.NetworkID = networkID
	}

	p := &types.Profile{
		Routes: map[string]*types.Route{},
	}

	containers, err := docker.GetContainerList(ctx, cli, networkID)
	if err != nil {
		return nil, err
	}

	composeServices, err := getComposeServices(opts)
	if err != nil {
		return nil, err
	}

	if len(composeServices) == 0 {
		addFromContainer(p, containers, opts)
		return p, nil
	}

	err = addFromComposeService(p, containers, composeServices, opts)

	return p, err
}

func getComposeServices(opts *DockerOptions) ([]string, error) {
	if opts.ComposeFile == "" {
		return nil, nil
	}

	f, err := os.Open(opts.ComposeFile)

	if err != nil {
		return nil, err
	}

	return docker.ParseComposeFile(f, opts.ProjectName)
}

func addFromContainer(profile *types.Profile, containers []dtypes.Container, opts *DockerOptions) {
	for _, c := range containers {
		for _, n := range c.NetworkSettings.Networks {
			if opts.NetworkID != "" && n.NetworkID != opts.NetworkID {
				continue
			}

			name := strings.Replace(c.Names[0], "/", "", -1)

			name = fmt.Sprintf("%s.%s", name, opts.Domain)
			profile.AddRoute(n.IPAddress, name)
		}
	}
}

func addFromComposeService(p *types.Profile, containers []dtypes.Container, srv []string, opts *DockerOptions) error {
	for _, c := range containers {
		for _, n := range c.NetworkSettings.Networks {
			if opts.NetworkID != "" && n.NetworkID != opts.NetworkID {
				continue
			}

			name := strings.Replace(c.Names[0], "/", "", -1)

			for _, c := range srv {
				match, err := regexp.MatchString(fmt.Sprintf("^%s(_[0-9]+)?", c), name)
				if err != nil {
					return err
				}

				if match {
					name = fmt.Sprintf("%s.%s", name, opts.Domain)
					if !opts.KeepPrefix {
						name = strings.Replace(name, opts.ProjectName+"_", "", 1)
					}

					name := fmt.Sprintf("%s.%s", name, opts.Domain)

					p.AddRoute(n.IPAddress, name)
				}
			}
		}
	}

	return nil
}
