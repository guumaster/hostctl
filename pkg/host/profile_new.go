package host

import (
	"bufio"
	"io"
)

// NewProfileFromReader creates a new profile reading lines from a reader
func NewProfileFromReader(r io.Reader, uniq bool) (*Profile, error) {
	p := &Profile{}
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
