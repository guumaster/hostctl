package profile

import (
	"github.com/guumaster/hostctl/pkg/k8s"
	"github.com/guumaster/hostctl/pkg/k8s/minikube"
	"github.com/guumaster/hostctl/pkg/types"
)

// NewProfileFromMinikube creates a new profile from a minikube profile and k8s namespace
func NewProfileFromMinikube(mini *minikube.Profile, ns string) (*types.Profile, error) {
	if mini.Status != minikube.Running {
		return nil, types.ErrMinikubeStatus
	}

	if !mini.IngressEnabled {
		return nil, types.ErrMinikubeIngress
	}

	cli, err := k8s.NewClientset()
	if err != nil {
		return nil, err
	}

	list, err := k8s.GetIngresses(cli, ns)
	if err != nil {
		return nil, err
	}

	p := &types.Profile{}

	for _, in := range list {
		p.AddRoute(types.NewRoute(in.IP.String(), in.Hostname))
	}

	return p, nil
}
