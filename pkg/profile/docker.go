package profile

import (
	"context"
	"fmt"
	"io"
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
	ComposeFile io.Reader
	ProjectName string
	KeepPrefix  bool
	Cli         *client.Client
}

// NewProfileFromDockerCompose creates a new profile from docker info
func NewProfileFromDockerCompose(opts *DockerOptions) (*types.Profile, error) {
	p := &types.Profile{}

	containers, err := getContainerList(opts)
	if err != nil {
		return nil, err
	}

	composeServices, err := docker.ParseComposeFile(opts.ComposeFile, opts.ProjectName)
	if err != nil {
		return nil, err
	}

	addFromComposeService(p, containers, composeServices, opts)

	return p, nil
}

// NewProfileFromDocker creates a new profile from docker info
func NewProfileFromDocker(opts *DockerOptions) (*types.Profile, error) {
	p := &types.Profile{}

	err := checkCli(opts)
	if err != nil {
		return nil, err
	}

	containers, err := getContainerList(opts)
	if err != nil {
		return nil, err
	}

	addFromContainer(p, containers, opts)

	return p, err
}

func checkCli(opts *DockerOptions) error {
	cli := opts.Cli
	if cli == nil {
		cli, err := client.NewEnvClient()
		if err != nil {
			return err
		}

		opts.Cli = cli
	}

	return nil
}

func getContainerList(opts *DockerOptions) ([]dtypes.Container, error) {
	var (
		networkID string
		err       error
	)

	ctx := context.Background()

	if opts.NetworkID == "" && opts.Network != "" {
		networkID, err = docker.GetNetworkID(ctx, opts.Cli, opts.Network)
		if err != nil {
			return nil, err
		}

		opts.NetworkID = networkID
	}

	return docker.GetContainerList(ctx, opts.Cli, networkID)
}

func addFromContainer(profile *types.Profile, containers []dtypes.Container, opts *DockerOptions) {
	routes := []*types.Route{}

	for _, c := range containers {
		for _, n := range c.NetworkSettings.Networks {
			if opts.NetworkID != "" && n.NetworkID != opts.NetworkID {
				continue
			}

			name := strings.Replace(c.Names[0], "/", "", -1)

			name = fmt.Sprintf("%s.%s", name, opts.Domain)
			routes = append(routes, types.NewRoute(n.IPAddress, name))
		}
	}

	profile.AddRoutes(routes)
}

func addFromComposeService(p *types.Profile, containers []dtypes.Container, srv []string, opts *DockerOptions) {
	for _, c := range containers {
		for _, n := range c.NetworkSettings.Networks {
			if opts.NetworkID != "" && n.NetworkID != opts.NetworkID {
				continue
			}

			name := strings.Replace(c.Names[0], "/", "", -1)

			routes := routeFromService(srv, name, n.IPAddress, opts)
			if len(routes) > 0 {
				p.AddRoutes(routes)
			}
		}
	}
}

func routeFromService(srv []string, name string, ip string, opts *DockerOptions) []*types.Route {
	// removeSuffix := regexp.MustCompile("(_[0-9]+)$")
	routes := []*types.Route{}

	for _, s := range srv {
		if !strings.Contains(name, s) {
			continue
		}

		// name := removeSuffix.ReplaceAllString(name, "")
		if !opts.KeepPrefix {
			name = strings.Replace(name, opts.ProjectName+"_", "", 1)
		}

		name = fmt.Sprintf("%s.%s", name, opts.Domain)
		routes = append(routes, types.NewRoute(ip, name))
	}

	return routes
}
