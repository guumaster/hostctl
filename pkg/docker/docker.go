package docker

import (
	"context"
	"fmt"
	"io"

	dtypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"

	"github.com/guumaster/hostctl/pkg/types"
)

// Options contains parameters to sync with docker and docker-compose.
type Options struct {
	Domain      string
	Network     string
	NetworkID   string
	ComposeFile io.Reader
	ProjectName string
	KeepPrefix  bool
	Cli         *client.Client
}

// GetContainerList returns a list of running docker containers, filter by network if networkID passed.
func GetContainerList(opts *Options) ([]dtypes.Container, error) {
	var (
		networkID string
		err       error
	)

	ctx := context.Background()

	err = checkCli(opts)
	if err != nil {
		return nil, err
	}

	if opts.NetworkID == "" && opts.Network != "" {
		networkID, err = GetNetworkID(ctx, opts.Cli, opts.Network)
		if err != nil {
			return nil, err
		}

		opts.NetworkID = networkID
	}

	f := filters.NewArgs()
	f.Add("status", "running")

	if networkID != "" {
		f.Add("network", networkID)
	}

	return opts.Cli.ContainerList(ctx, dtypes.ContainerListOptions{Filters: f})
}

// GetNetworkID returns the an ID that match a network name.
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

func checkCli(opts *Options) error {
	cli := opts.Cli
	if cli == nil {
		cli, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			return err
		}

		opts.Cli = cli
	}

	return nil
}
