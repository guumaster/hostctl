package host

import (
	"fmt"
	"io"
	"net"
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

// AddRoute adds multiple routes to the profile
func (p *Profile) AddRoutes(ip string, hostnames []string) {
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
		return nil, fmt.Errorf("invalid ip '%s'", ip)
	}

	hosts, ok := p.Routes[key.String()]
	if !ok {
		return nil, fmt.Errorf("ip '%s' not present in profile '%s' ", key, p.Name)
	}

	return hosts.HostNames, nil
}

// GetHostNames returns all hostnames of the profile.
func (p *Profile) GetAllHostNames() ([]string, error) {
	list := []string{}

	for _, r := range p.Routes {
		list = append(list, r.HostNames...)
	}

	return list, nil
}

// Render writes the profile content to the given StringWriter
func (p *Profile) Render(w io.StringWriter) error {
	_, err := w.WriteString(fmt.Sprintf("\n# profile.%s %s\n", p.Status, p.Name))
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

			_, err = w.WriteString(fmt.Sprintf("%s%s %s\n", prefix, ip, host))
			if err != nil {
				return err
			}
		}
	}

	_, err = w.WriteString("# end\n")
	if err != nil {
		return err
	}

	return nil
}

// Render writes the default profile content to the given StringWriter
func (d DefaultProfile) Render(w io.StringWriter) error {
	for _, row := range d {
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

		_, err := w.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}
