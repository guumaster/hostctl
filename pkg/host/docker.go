package host

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"gopkg.in/yaml.v2"
)

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
