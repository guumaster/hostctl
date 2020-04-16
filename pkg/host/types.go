package host

import (
	"net"
	"sync"

	"github.com/spf13/afero"
)

const banner = `
##################################################################
# Content under this line is handled by hostctl. DO NOT EDIT.
##################################################################`

type File struct {
	fs        afero.Fs
	src       afero.File
	data      *Content
	hasBanner bool
	mutex     sync.Mutex
}

type Content struct {
	DefaultProfile DefaultProfile
	ProfileNames   []string
	Profiles       map[string]*Profile
}

type Profile struct {
	Name   string
	Status ProfileStatus
	IPList []string
	Routes map[string]*Route
}

type DefaultProfile []*tableRow

type tableRow struct {
	Comment string
	Profile string
	Status  string
	IP      string
	Host    string
}

type Route struct {
	IP        net.IP
	HostNames []string
}

type ProfileStatus string

const (
	// Enabled marks a profile active on your hosts file.
	Enabled ProfileStatus = "on"
	// Disabled marks a profile not active on your hosts file.
	Disabled ProfileStatus = "off"

	// Default is the name of the default profile
	Default = "default"
)
