package parser

import (
	"fmt"
	"strings"

	dtypes "github.com/docker/docker/api/types"

	"github.com/guumaster/hostctl/pkg/docker"
	"github.com/guumaster/hostctl/pkg/types"
)

// NewProfileFromDocker creates a new profile from docker info.
func NewProfileFromDocker(opts *docker.Options) (*types.Profile, error) {
	p := &types.Profile{}

	containers, err := docker.GetContainerList(opts)
	if err != nil {
		return nil, err
	}

	addFromContainer(p, containers, opts)

	return p, err
}

// NewProfileFromDockerCompose creates a new profile from docker info.
func NewProfileFromDockerCompose(opts *docker.Options) (*types.Profile, error) {
	p := &types.Profile{}

	containers, err := docker.GetContainerList(opts)
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

func addFromContainer(profile *types.Profile, containers []dtypes.Container, opts *docker.Options) {
	var routes []*types.Route

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

func addFromComposeService(p *types.Profile, containers []dtypes.Container, srv []string, opts *docker.Options) {
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

func routeFromService(srv []string, name string, ip string, opts *docker.Options) []*types.Route {
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
