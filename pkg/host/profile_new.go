package host

import (
	"bufio"
	"io"
)

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
