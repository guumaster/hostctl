package file

import (
	"github.com/guumaster/hostctl/pkg/host"
	"github.com/guumaster/hostctl/pkg/host/errors"
)

// Toggle alternates between enable and disable status of a profile.
func (f *File) Toggle(profiles []string) error {
	for _, name := range profiles {
		if name == host.Default {
			continue
		}

		p, ok := f.data.Profiles[name]
		if !ok {
			return errors.ErrUnknownProfile
		}

		if p.Status == host.Enabled {
			p.Status = host.Disabled
		} else {
			p.Status = host.Enabled
		}

		f.data.Profiles[name] = p
	}

	return nil
}
