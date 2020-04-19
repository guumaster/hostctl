package file

import (
	"github.com/guumaster/hostctl/pkg/host"
	"github.com/guumaster/hostctl/pkg/host/errors"
)

// AddProfile adds a profile to the list
func (f *File) AddProfile(p *host.Profile) error {
	if p.Name == host.Default {
		return errors.ErrDefaultProfile
	}

	f.MergeProfiles([]*host.Profile{p})

	return nil
}
