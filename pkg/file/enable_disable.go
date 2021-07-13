package file

import (
	"github.com/guumaster/hostctl/pkg/types"
)

// Enable marks profiles as enable by uncommenting all hosts lines
// making the routing work again.
func (f *File) Enable(profiles []string) error {
	return f.changeTo(profiles, types.Enabled)
}

// Disable marks profiles as disable by commenting all hosts lines.
// The content remains on the file and can be enabled later.
func (f *File) Disable(profiles []string) error {
	return f.changeTo(profiles, types.Disabled)
}

// EnableAll marks all profiles as enable by uncommenting all hosts lines
// making the routing work again.
func (f *File) EnableAll() error {
	return f.changeTo(f.GetProfileNames(), types.Enabled)
}

// DisableAll marks all profiles as disable by commenting all hosts lines.
// The content remains on the file and can be enabled later.
func (f *File) DisableAll() error {
	return f.changeTo(f.GetProfileNames(), types.Disabled)
}

// DisableOnly marks profiles as disable and enable all other profiles.
func (f *File) DisableOnly(profiles []string) error {
	return f.changeToSplitted(profiles, types.Disabled)
}

// EnableOnly marks profiles as enable and disable all other profiles.
func (f *File) EnableOnly(profiles []string) error {
	return f.changeToSplitted(profiles, types.Enabled)
}

func (f *File) changeToSplitted(profiles []string, status types.Status) error {
	for _, name := range f.data.ProfileNames {
		p := f.data.Profiles[name]

		if contains(profiles, name) {
			p.Status = status
		} else {
			p.Status = invertStatus(status)
		}

		f.data.Profiles[name] = p
	}

	return nil
}

func invertStatus(s types.Status) types.Status {
	if s == types.Enabled {
		return types.Disabled
	}

	return types.Enabled
}

func (f *File) changeTo(profiles []string, status types.Status) error {
	for _, name := range profiles {
		if name == types.Default {
			continue
		}

		p, ok := f.data.Profiles[name]
		if !ok {
			return types.ErrUnknownProfile
		}

		p.Status = status
		f.data.Profiles[name] = p
	}

	return nil
}
