package profile

import (
	"bufio"
	"io"

	"github.com/guumaster/hostctl/pkg/types"
)

// NewProfileFromReader creates a new profile reading lines from a reader
func NewProfileFromReader(r io.Reader, uniq bool) (*types.Profile, error) {
	p := &types.Profile{}
	s := bufio.NewScanner(r)

	routes := []*types.Route{}

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

	if uniq {
		p.AddRoutesUniq(routes)
	} else {
		p.AddRoutes(routes)
	}

	return p, nil
}
