package types

import (
	"errors"
)

var (
	// ErrMissingProfile when the profile is mandatory
	ErrMissingProfile = errors.New("missing profile name")

	// ErrUnknownProfile when the profile is not present
	ErrUnknownProfile = errors.New("unknown profile name")

	// ErrDefaultProfile when trying to edit default content
	ErrDefaultProfile = errors.New("'default' profile should not be handled by hostctl")

	// ErrNoContent when data to write is empty
	ErrNoContent = errors.New("no content to write")

	// ErrNotPresentIP when looking for an IP not contained in profile
	ErrNotPresentIP = errors.New("ip not present")

	// ErrUnknownNetworkID when you pass an invalid network ID to sync docker
	ErrUnknownNetworkID = errors.New("unknown network ID")

	// ErrInvalidIP when the IP is malformed
	ErrInvalidIP = errors.New("invalid ip")

	// ErrInvalidProfileHeader when the profile header is invalid
	ErrInvalidProfileHeader = errors.New("invalid format for profile header")

	// ErrMinikubeStatus when minikube profile is not running
	ErrMinikubeStatus = errors.New("minikube profile has to be running")

	// ErrMinikubeIngress when minikube doesn't have ingress addon enabled
	ErrMinikubeIngress = errors.New("minikube profile doesn't have ingress addon enabled")

	// ErrKubernetesNamespace when no namespace is given
	ErrKubernetesNamespace = errors.New("namespace parameter is required or use --all-namespaces")
)
