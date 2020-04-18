package profile

import (
	"bufio"
	"io"

	"github.com/guumaster/hostctl/pkg/host/parser"
	"github.com/guumaster/hostctl/pkg/host/types"
)

// NewProfileFromReader creates a new profile reading lines from a reader
func NewProfileFromReader(r io.Reader, uniq bool) (*types.Profile, error) {
	p := &types.Profile{}
	s := bufio.NewScanner(r)

	for s.Scan() {
		appendLine(p, string(s.Bytes()))

		if err := s.Err(); err != nil {
			return nil, err
		}
	}

	if uniq {
		for _, r := range p.Routes {
			r.HostNames = uniqueStrings(r.HostNames)
		}
	}

	return p, nil
}

func appendLine(p *types.Profile, line string) {
	if line == "" {
		return
	}

	route, ok := parser.ParseLine(line)
	if !ok {
		return
	}

	ip := route.IP.String()
	p.AddRoutes(ip, route.HostNames)
}

func uniqueStrings(xs []string) []string {
	var list []string

	keys := make(map[string]bool)

	for _, entry := range xs {
		if _, value := keys[entry]; !value {
			keys[entry] = true

			list = append(list, entry)
		}
	}

	return list
}
