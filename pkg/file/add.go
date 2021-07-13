package file

import (
	"github.com/guumaster/hostctl/pkg/types"
)

// AddProfile adds a profile to the list.
func (f *File) AddProfile(p *types.Profile) error {
	if p.Name == types.Default {
		return types.ErrDefaultProfile
	}

	f.MergeProfiles([]*types.Profile{p})

	return nil
}
