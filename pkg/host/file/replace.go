package file

import (
	"errors"

	"github.com/guumaster/hostctl/pkg/host"
	errors2 "github.com/guumaster/hostctl/pkg/host/errors"
)

// ReplaceProfile removes previous profile with same name and add new profile to the list
func (f *File) ReplaceProfile(p *host.Profile) error {
	err := f.RemoveProfile(p.Name)
	if err != nil && !errors.Is(err, errors2.ErrUnknownProfile) {
		return err
	}

	return f.AddProfile(p)
}
