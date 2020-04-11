package host

import (
	"fmt"
	"io"
	"net"
)

func (p *Profile) String() string {
	return fmt.Sprintf("[%s]%s", p.Status, p.Name)
}

func (p *Profile) GetStatus() string {
	return string(p.Status)
}
func (p *Profile) AddRoute(ip, hostname string) {
	if p.Routes[ip] == nil {
		p.Routes[ip] = &Route{
			IP:        net.ParseIP(ip),
			HostNames: []string{hostname},
		}
	} else {
		p.Routes[ip].HostNames = append(p.Routes[ip].HostNames, hostname)
	}
}

func (p *Profile) AddRoutes(ip string, hostnames []string) {
	if p.Routes[ip] == nil {
		p.Routes[ip] = &Route{
			IP:        net.ParseIP(ip),
			HostNames: hostnames,
		}
	} else {
		p.Routes[ip].HostNames = append(p.Routes[ip].HostNames, hostnames...)
	}
}

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

func (p *Profile) GetHostNames(ip string) ([]string, error) {
	key := net.ParseIP(ip)
	if key == nil {
		return nil, fmt.Errorf("invalid ip '%s'", key)
	}
	hosts, ok := p.Routes[key.String()]
	if !ok {
		return nil, fmt.Errorf("ip '%s' not present in profile '%s' ", key, p.Name)
	}

	return hosts.HostNames, nil
}

func (p *Profile) Render(w io.StringWriter) error {
	_, err := w.WriteString(fmt.Sprintf("\n# profile.%s %s\n", p.Status, p.Name))
	if err != nil {
		return err
	}

	for ip, routes := range p.Routes {
		for _, host := range routes.HostNames {
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
