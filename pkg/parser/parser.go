package parser

import (
	"bufio"
	"io"
	"net"
	"regexp"
	"strings"

	"github.com/guumaster/hostctl/pkg/types"
)

var (
	profileNameRe = regexp.MustCompile(`# profile(?:.(on|off))?\s+([a-z0-9-_.\s]+)`)
	profileEnd    = regexp.MustCompile(`(?i)# end\s*`)
	spaceRemover  = regexp.MustCompile(`\s+`)
	tabReplacer   = regexp.MustCompile(`\t+`)
	endingComment = regexp.MustCompile(`(.[^#]*).*`)
)

// Parser is the interface for content parsers.
type Parser interface {
	Parse(reader io.Reader) types.Content
}

// Parse reads content from reader into Data struct.
func Parse(r io.Reader) (*types.Content, error) {
	data := &types.Content{
		ProfileNames: []string{},
		Profiles:     map[string]*types.Profile{},
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
			data.Profiles[currProfile] = &types.Profile{
				Name:   currProfile,
				Status: p.Status,
			}

		case profileEnd.Match(b):
			currProfile = ""

		case currProfile != "":
			line := string(b)
			if line == "" {
				continue
			}

			route, ok := parseRouteLine(line)
			if !ok {
				continue
			}

			data.Profiles[currProfile].AddRoute(route)

		default:
			row := parseToDefault(b, currProfile)

			// When hostctl banner line is detected, remove previous and next line
			if isBannerLine(row) {
				data.DefaultProfile = data.DefaultProfile[0 : len(data.DefaultProfile)-1]

				s.Scan() // skip next line
			} else {
				data.DefaultProfile = append(data.DefaultProfile, row)
			}
		}

		if err := s.Err(); err != nil {
			return nil, err
		}
	}

	return data, nil
}

func isBannerLine(r *types.Row) bool {
	return strings.Contains(r.Comment, "# Content under this line is handled by hostctl. DO NOT EDIT.")
}

func parseToDefault(b []byte, currProfile string) *types.Row {
	var row *types.Row

	if len(b) == 0 {
		if currProfile == "" {
			row = &types.Row{Comment: ""}
		}

		return row
	}

	line, ok := parseRouteLine(string(b))
	if !ok {
		row = &types.Row{
			Comment: string(b),
		}
	} else {
		status := types.Enabled
		if off, _ := regexp.Match("^#", b); off {
			status = types.Disabled
		}
		row = &types.Row{
			Profile: types.Default,
			Status:  string(status),
			IP:      line.IP.String(),
			Hosts:   line.HostNames,
		}
	}

	return row
}

func parseProfileHeader(b []byte) (*types.Profile, error) {
	rs := profileNameRe.FindSubmatch(b)
	if len(rs) != 3 || string(rs[2]) == "" {
		return nil, types.ErrInvalidProfileHeader
	}

	status := types.Enabled
	if string(rs[1]) == string(types.Disabled) {
		status = types.Disabled
	}

	return &types.Profile{
		Name:   strings.TrimSpace(string(rs[2])),
		Status: status,
	}, nil
}

// parseRouteLine checks if a line is a host line or a comment line.
func parseRouteLine(str string) (*types.Route, bool) {
	clean := spaceRemover.ReplaceAllString(str, " ")
	clean = tabReplacer.ReplaceAllString(clean, " ")
	clean = strings.TrimSpace(clean)

	if len(clean) == 0 {
		return nil, false
	}

	result := endingComment.FindStringSubmatch(clean)
	tResult := strings.TrimSpace(result[1])
	p := strings.Split(tResult, " ")

	i := 0
	if p[0] == "#" && len(p) > 1 {
		i = 1
	}

	ip := net.ParseIP(p[i])
	hostnames := p[i+1:]

	if ip == nil || len(hostnames) == 0 {
		return nil, false
	}

	return &types.Route{IP: ip, HostNames: hostnames}, true
}

// ParseProfile creates a new profile reading lines from a reader.
func ParseProfile(r io.Reader) (*types.Profile, error) {
	p := &types.Profile{}
	s := bufio.NewScanner(r)

	var routes []*types.Route

	for s.Scan() {
		line := string(s.Bytes())
		if line == "" {
			continue
		}

		route, ok := parseRouteLine(line)
		if !ok {
			continue
		}

		routes = append(routes, route)

		if err := s.Err(); err != nil {
			return nil, err
		}
	}

	p.AddRoutes(routes)

	return p, nil
}
