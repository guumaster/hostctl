package file

import (
	"errors"

	"github.com/guumaster/hostctl/pkg/types"
)

// ReplaceProfile removes previous profile with same name and add new profile to the list
func (f *File) ReplaceProfile(p *types.Profile) error {
	err := f.RemoveProfile(p.Name)
	if err != nil && !errors.Is(err, types.ErrUnknownProfile) {
		return err
	}

	return f.AddProfile(p)
}
