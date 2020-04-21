package file

import (
	"github.com/guumaster/hostctl/pkg/types"
)

// Toggle alternates between enable and disable status of a types.
func (f *File) Toggle(profiles []string) error {
	for _, name := range profiles {
		if name == types.Default {
			continue
		}

		p, ok := f.data.Profiles[name]
		if !ok {
			return types.ErrUnknownProfile
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
