package file

import (
	"github.com/guumaster/hostctl/pkg/host/errors"
	"github.com/guumaster/hostctl/pkg/host/types"
)

// Toggle alternates between enable and disable status of a profile.
func (f *File) Toggle(profiles []string) error {
	for _, name := range profiles {
		if name == types.Default {
			continue
		}

		p, ok := f.data.Profiles[name]
		if !ok {
			return errors.ErrUnknownProfile
		}

		if p.Status == types.Enabled {
			p.Status = types.Disabled
		} else {
			p.Status = types.Enabled
		}

		f.data.Profiles[name] = p
	}

	return nil
}
