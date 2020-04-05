package host

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"gopkg.in/yaml.v2"
)

// AddFromFileOptions contains available options for adding from file.
type AddFromDockerOptions struct {
	Dst     string
	Domain  string
	Profile string
	Watch   bool
	Docker  *DockerOptions
}

// DockerOptions contains parameters to sync with docker and docker-compose
type DockerOptions struct {
	Network     string
	ComposeFile string
	ProjectName string
	KeepPrefix  bool
}

func AddFromDocker(ctx context.Context, opts *AddFromDockerOptions) error {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	f := filters.NewArgs()
	f.Add("status", "running")

	networkID, err := getNetworkID(ctx, cli, opts.Docker)
	if err != nil {
		return err
	}
	if networkID != "" {
		f.Add("network", networkID)
	}

	var containers []string
	if opts.Docker.ComposeFile != "" {
		containers, err = parseComposeFile(opts.Docker.ComposeFile, opts.Docker.ProjectName)
		if err != nil {
			return err
		}
	}

	list, err := cli.ContainerList(ctx, types.ContainerListOptions{
		Filters: f,
	})
	if err != nil {
		return err
	}

	var lines []string
	for _, c := range list {
		for _, n := range c.NetworkSettings.Networks {
			if networkID != "" && n.NetworkID != networkID {
				continue
			}
			name := strings.Replace(c.Names[0], "/", "", -1)

			if len(containers) == 0 {
				name = fmt.Sprintf("%s.%s", name, opts.Domain)
				lines = append(lines, fmt.Sprintf("%s %s", n.IPAddress, name))
			} else {
				for _, c := range containers {
					match, err := regexp.MatchString(fmt.Sprintf("^%s(_[0-9]+)?", c), name)
					if err != nil {
						return err
					}
					if match {
						name = fmt.Sprintf("%s.%s", name, opts.Domain)
						if !opts.Docker.KeepPrefix {
							name = strings.Replace(name, opts.Docker.ProjectName+"_", "", 1)
						}
						lines = append(lines, fmt.Sprintf("%s %s", n.IPAddress, name))
					}
				}
			}
		}
	}

	newData := &hostFile{
		profiles: profileMap{
			"default": lines,
		},
	}
	return add(newData, &commonAddOptions{
		opts.Dst,
		opts.Profile,
		true,
	})
}

type composeData struct {
	Services map[string]composeService `yaml:"services"`
}

type composeService struct {
	ContainerName string `yaml:"container_name"`
}

func parseComposeFile(file, projectName string) ([]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	data := &composeData{}
	err = yaml.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}

	var containers []string
	for serv, data := range data.Services {
		name := data.ContainerName
		if data.ContainerName == "" {
			name = fmt.Sprintf("%s_%s", projectName, serv)
		}
		containers = append(containers, name)
	}
	return containers, nil
}

func getNetworkID(ctx context.Context, cli *client.Client, opts *DockerOptions) (string, error) {
	if opts == nil || opts.Network == "" {
		return "", nil
	}

	var networkID string
	nets, err := cli.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		return "", err
	}

	for _, net := range nets {
		if net.Name == opts.Network || net.ID == opts.Network {
			networkID = net.ID
			break
		}
	}
	if networkID == "" {
		return "", fmt.Errorf("unknown network name or ID: '%s'", opts.Network)
	}
	return networkID, nil
}
