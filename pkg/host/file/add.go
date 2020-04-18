package file

import (
	"github.com/guumaster/hostctl/pkg/host/errors"
	"github.com/guumaster/hostctl/pkg/host/types"
)

// AddProfile adds a profile to the list
func (f *File) AddProfile(p *types.Profile) error {
	if p.Name == types.Default {
		return errors.ErrDefaultProfileError
	}

	f.MergeProfiles(&types.Content{
		ProfileNames: []string{p.Name},
		Profiles: map[string]*types.Profile{
			p.Name: p,
		},
	})

	return nil
}
