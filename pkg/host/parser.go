package host

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"regexp"
	"strings"
)

var (
	profileNameRe = regexp.MustCompile(`# profile(?:.(on|off))?\s+([a-z0-9-_.\s]+)`)
	profileEnd    = regexp.MustCompile(`(?i)# end\s*`)
	//	disableRe     = regexp.MustCompile(`^#\s*`)
	spaceRemover = regexp.MustCompile(`\s+`)
	tabReplacer  = regexp.MustCompile(`\t+`)
)

// Parse reads content from reader into Data struct.
func Parse(r io.Reader) (*Content, error) {
	data := &Content{
		ProfileNames: []string{},
		Profiles:     map[string]Profile{},
	}

	currProfile := ""

	s := bufio.NewScanner(r)
	for s.Scan() {
		b := s.Bytes()

		switch {
		case profileNameRe.Match(b):
			p, _ := parseProfileHeader(b)

			currProfile = p.Name
			data.ProfileNames = append(data.ProfileNames, currProfile)
			data.Profiles[currProfile] = Profile{
				Name:   currProfile,
				Status: p.Status,
			}

		case profileEnd.Match(b):
			currProfile = ""

		case currProfile != "":
			profile := data.Profiles[currProfile]
			p := appendLine(&profile, string(b))
			data.Profiles[currProfile] = *p

		default:
			row := parseToDefault(b, currProfile)
			data.DefaultProfile = append(data.DefaultProfile, row)

		}

		if err := s.Err(); err != nil {
			return nil, err
		}
	}
	return data, nil
}

func parseToDefault(b []byte, currProfile string) *tableRow {
	var row *tableRow
	if len(b) == 0 {
		if currProfile == "" {
			row = &tableRow{Comment: ""}
		}
		return row
	}

	line, ok := parseLine(string(b))
	if !ok {
		row = &tableRow{
			Comment: string(b),
		}
	} else {
		status := Enabled
		if off, _ := regexp.Match("^#", b); off {
			status = Disabled
		}
		row = &tableRow{
			Profile: "default",
			Status:  string(status),
			IP:      line.IP.String(),
			Host:    line.HostNames[0],
		}
	}
	return row
}

func parseProfileHeader(b []byte) (*Profile, error) {
	rs := profileNameRe.FindSubmatch(b)
	if len(rs) != 3 || string(rs[2]) == "" {
		return nil, fmt.Errorf("invalid format for profile header")
	}

	status := Enabled
	if string(rs[1]) == string(Disabled) {
		status = Disabled
	}
	return &Profile{
		Name:   strings.TrimSpace(string(rs[2])),
		Status: status,
	}, nil
}

func uniqueStrings(xs []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range xs {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func appendLine(p *Profile, line string) *Profile {
	if line == "" {
		return p
	}
	route, ok := parseLine(line)
	if !ok {
		return p
	}
	ip := route.IP.String()
	if p.Routes == nil {
		p.Routes = map[string]*Route{}
		p.Routes[ip] = route
	} else if p.Routes[ip] == nil {
		p.Routes[ip] = route
	} else {
		p.Routes[ip].HostNames = append(p.Routes[ip].HostNames, route.HostNames...)
		p.Routes[ip].HostNames = uniqueStrings(p.Routes[ip].HostNames)
	}
	return p
}

// parseLine checks if a line is a host line or a comment line.
func parseLine(str string) (*Route, bool) {
	p := strings.Split(cleanLine(str), " ")
	i := 0
	if p[0] == "#" && len(p) > 1 {
		i = 1
	}
	ip := net.ParseIP(p[i])

	if ip == nil {
		return nil, false
	}

	return &Route{IP: ip, HostNames: p[i+1:]}, true
}

func cleanLine(line string) string {
	clean := spaceRemover.ReplaceAllString(line, " ")
	clean = tabReplacer.ReplaceAllString(clean, " ")
	clean = strings.TrimSpace(clean)
	return clean
}
