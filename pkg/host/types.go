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

// File container to handle a hosts file
type File struct {
	fs        afero.Fs
	src       afero.File
	data      *Content
	hasBanner bool
	mutex     sync.Mutex
}

// Content contains complete data of all profiles
type Content struct {
	DefaultProfile DefaultProfile
	ProfileNames   []string
	Profiles       map[string]*Profile
}

// Profile contains all data of a single profile
type Profile struct {
	Name   string
	Status ProfileStatus
	IPList []string
	Routes map[string]*Route
}

//  DefaultProfile contains data for the default profile
type DefaultProfile []*tableRow

type tableRow struct {
	Comment string
	Profile string
	Status  string
	IP      string
	Host    string
}

// Route contains hostnames of all routes with the same IP
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
