package host

import (
	"bufio"
	"io"
)

func NewProfileFromReader(r io.Reader) (*Profile, error) {
	p := &Profile{}
	s := bufio.NewScanner(r)
	for s.Scan() {
		appendLine(p, string(s.Bytes()))

		if err := s.Err(); err != nil {
			return nil, err
		}
	}
	return p, nil
}
