package file

import (
	"errors"

	errors2 "github.com/guumaster/hostctl/pkg/host/errors"
	"github.com/guumaster/hostctl/pkg/host/types"
)

// ReplaceProfile removes previous profile with same name and add new profile to the list
func (f *File) ReplaceProfile(p *types.Profile) error {
	err := f.RemoveProfile(p.Name)
	if err != nil && !errors.Is(err, errors2.ErrUnknownProfile) {
		return err
	}

	return f.AddProfile(p)
}
