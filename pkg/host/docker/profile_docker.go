package docker

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"

	types2 "github.com/guumaster/hostctl/pkg/host/types"
)

// DockerOptions contains parameters to sync with docker and docker-compose
type Options struct {
	Domain      string
	Network     string
	ComposeFile string
	ProjectName string
	KeepPrefix  bool
	Cli         *client.Client
}

func containerList(ctx context.Context, opts *Options) ([]types.Container, error) {
	var err error

	cli := opts.Cli
	if opts.Cli == nil {
		cli, err = client.NewEnvClient()
		if err != nil {
			return nil, err
		}

		opts.Cli = cli
	}
	defer cli.Close()

	f := filters.NewArgs()
	f.Add("status", "running")

	networkID, err := getNetworkID(ctx, opts)
	if err != nil {
		return nil, err
	}

	if networkID != "" {
		f.Add("network", networkID)
	}

	return cli.ContainerList(ctx, types.ContainerListOptions{Filters: f})
}

// NewProfileFromDocker creates a new profile from docker info
func NewProfileFromDocker(ctx context.Context, opts *Options) (*types2.Profile, error) {
	containers, err := containerList(ctx, opts)
	if err != nil {
		return nil, err
	}

	var composeServices []string
	if opts.ComposeFile != "" {
		composeServices, err = parseComposeFile(opts.ComposeFile, opts.ProjectName)
		if err != nil {
			return nil, err
		}
	}

	p := &types2.Profile{
		Routes: map[string]*types2.Route{},
	}

	return addToProfile(ctx, p, containers, composeServices, opts)
}

func addToProfile(
	ctx context.Context,
	profile *types2.Profile,
	containers []types.Container,
	composeServices []string,
	opts *Options) (*types2.Profile, error) {
	networkID, err := getNetworkID(ctx, opts)
	if err != nil {
		return nil, err
	}

	for _, c := range containers {
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
