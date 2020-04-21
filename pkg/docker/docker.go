package docker

import (
	"context"
	"fmt"

	dtypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"

	"github.com/guumaster/hostctl/pkg/types"
)

// GetContainerList returns a list of running docker containers, filter by network if networkID passed
func GetContainerList(ctx context.Context, cli *client.Client, networkID string) ([]dtypes.Container, error) {
	f := filters.NewArgs()
	f.Add("status", "running")

	if networkID != "" {
		f.Add("network", networkID)
	}

	return cli.ContainerList(ctx, dtypes.ContainerListOptions{Filters: f})
}

// GetNetworkID returns the an ID that match a network name
func GetNetworkID(ctx context.Context, cli *client.Client, network string) (string, error) {
	var networkID string

	if network == "" {
		return "", nil
	}

	nets, err := cli.NetworkList(ctx, dtypes.NetworkListOptions{})
	if err != nil {
		return "", err
	}

	for _, net := range nets {
		if net.Name == network || net.ID == network {
			networkID = net.ID
			break
		}
	}

	if networkID == "" {
		return "", fmt.Errorf("%w: '%s'", types.ErrUnknownNetworkID, network)
	}

	return networkID, nil
}
