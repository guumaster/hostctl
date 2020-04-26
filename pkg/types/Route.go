package types

import (
	"net"
)

// Route contains hostnames of all routes with the same IP
type Route struct {
	IP        net.IP
	HostNames []string
}

// NewRoute creates an new Route
func NewRoute(ip string, hostnames ...string) *Route {
	return &Route{
		IP:        net.ParseIP(ip),
		HostNames: hostnames,
	}
}
