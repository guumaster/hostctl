package host

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/guumaster/hostctl/pkg/host/errors"
	"github.com/guumaster/hostctl/pkg/host/render"
)

// Profile contains all data of a single profile
type Profile struct {
	Name   string
	Status Status
	IPList []string
	Routes map[string]*Route
}

// DefaultProfile contains data for the default profile
type DefaultProfile []*render.Row

// Route contains hostnames of all routes with the same IP
type Route struct {
	IP        net.IP
	HostNames []string
}

// ProfileStatus represents the status of a Profile
type Status string

const (
	// Enabled marks a profile active on your hosts file.
	Enabled Status = "on"
	// Disabled marks a profile not active on your hosts file.
	Disabled Status = "off"

	// Default is the name of the default profile
	Default = "default"
)

// String returns a string representation of the profile
func (p *Profile) String() string {
	return fmt.Sprintf("[%s]%s", p.Status, p.Name)
}

// GetStatus returns a string value of ProfileStatus
func (p *Profile) GetStatus() string {
	return string(p.Status)
}

func (p *Profile) appendIP(n string) {
	for _, c := range p.IPList {
		if c == n {
			return
		}
	}

	p.IPList = append(p.IPList, n)
}

// AddRoute adds a single route to the profile
func (p *Profile) AddRoute(ip, hostname string) {
	if p.Routes[ip] == nil {
		p.appendIP(ip)
		p.Routes[ip] = &Route{
			IP:        net.ParseIP(ip),
			HostNames: []string{hostname},
		}
	} else {
		p.Routes[ip].HostNames = append(p.Routes[ip].HostNames, hostname)
	}
}

// AddRoutes adds multiple routes to the profile
func (p *Profile) AddRoutes(ip string, hostnames []string) {
	if p.Routes == nil {
		p.Routes = map[string]*Route{}
	}

	if p.Routes[ip] == nil {
		p.appendIP(ip)
		p.Routes[ip] = &Route{
			IP:        net.ParseIP(ip),
			HostNames: hostnames,
		}
	} else {
		p.Routes[ip].HostNames = append(p.Routes[ip].HostNames, hostnames...)
	}
}

// RemoveRoutes removes multiple hostnames of a profile
func (p *Profile) RemoveRoutes(hostnames []string) {
	for _, h := range hostnames {
		for ip, r := range p.Routes {
			r.HostNames = remove(r.HostNames, h)
			if len(r.HostNames) == 0 {
				delete(p.Routes, ip)
			}
		}
	}
}

// GetHostNames returns a list of all hostnames of the given ip.
func (p *Profile) GetHostNames(ip string) ([]string, error) {
	key := net.ParseIP(ip)
	if key == nil {
		return nil, fmt.Errorf("%w '%s'", errors.ErrInvalidIP, ip)
	}

	hosts, ok := p.Routes[key.String()]
	if !ok {
		return nil, fmt.Errorf("%w: %s[%s] ", errors.ErrNotPresentIP, key, p.Name)
	}

	return hosts.HostNames, nil
}

// GetAllHostNames returns all hostnames of the profile.
func (p *Profile) GetAllHostNames() ([]string, error) {
	list := []string{}

	for _, r := range p.Routes {
		list = append(list, r.HostNames...)
	}

	return list, nil
}

// Render writes the profile content to the given StringWriter
func (p *Profile) Render(w io.StringWriter) error {
	tmp := bytes.NewBufferString("")

	_, err := tmp.WriteString(fmt.Sprintf("\n# profile.%s %s\n", p.Status, p.Name))
	if err != nil {
		return err
	}

	for _, ip := range p.IPList {
		route := p.Routes[ip]
		for _, host := range route.HostNames {
			prefix := ""
			if p.Status == Disabled {
				prefix = "# "
			}

			_, err = tmp.WriteString(fmt.Sprintf("%s%s %s\n", prefix, ip, host))
			if err != nil {
				return err
			}
		}
	}

	_, err = tmp.WriteString("# end\n")
	if err != nil {
		return err
	}

	// Write to input writer after knowing the profile is well formed
	_, err = w.WriteString(tmp.String())

	return err
}

// Render writes the default profile content to the given StringWriter
func (d DefaultProfile) Render(w io.StringWriter) error {
	tmp := bytes.NewBufferString("")

	for i, row := range d {
		line := getLine(row)
		nextLine := ""

		if i+1 < len(d) {
			nextLine = getLine(d[i+1])
		}

		// skips two consecutive empty lines
		if line == "" && nextLine == "" {
			continue
		}

		_, err := tmp.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	// Write to input writer after knowing the profile is well formed
	_, err := w.WriteString(tmp.String())

	return err
}

func getLine(row *render.Row) string {
	line := ""
	if row.Comment != "" {
		line = row.Comment
	} else {
		prefix := ""
		if row.Status == string(Disabled) {
			prefix = "# "
		}

		line = fmt.Sprintf("%s%s %s", prefix, row.IP, row.Host)
	}

	return strings.TrimSpace(line)
}

func remove(s []string, n string) []string {
	list := []string{}

	for _, x := range s {
		if x != n {
			list = append(list, x)
		}
	}

	return list
}
