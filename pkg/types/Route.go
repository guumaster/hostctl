package types

import (
	"net"
)

// Route contains hostnames of all routes with the same IP
type Route struct {
	IP        net.IP
	HostNames []string
}

func NewRoute(ip string, hostnames ...string) *Route {
	return &Route{
		IP:        net.ParseIP(ip),
		HostNames: hostnames,
	}
}
