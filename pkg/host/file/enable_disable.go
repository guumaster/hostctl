package file

import (
	"github.com/guumaster/hostctl/pkg/host"
	"github.com/guumaster/hostctl/pkg/host/errors"
)

// Enable marks profiles as enable by uncommenting all hosts lines
// making the routing work again.
func (f *File) Enable(profiles []string) error {
	return f.changeTo(profiles, host.Enabled)
}

// Disable marks profiles as disable by commenting all hosts lines.
// The content remains on the file and can be enabled later.
func (f *File) Disable(profiles []string) error {
	return f.changeTo(profiles, host.Disabled)
}

// EnableAll marks all profiles as enable by uncommenting all hosts lines
// making the routing work again.
func (f *File) EnableAll() error {
	return f.changeTo(f.GetProfileNames(), host.Enabled)
}

// DisableAll marks all profiles as disable by commenting all hosts lines.
// The content remains on the file and can be enabled later.
func (f *File) DisableAll() error {
	return f.changeTo(f.GetProfileNames(), host.Disabled)
}

// DisableOnly marks profiles as disable and enable all other profiles
func (f *File) DisableOnly(profiles []string) error {
	return f.changeToSplitted(profiles, host.Disabled)
}

// EnableOnly marks profiles as enable and disable all other profiles
func (f *File) EnableOnly(profiles []string) error {
	return f.changeToSplitted(profiles, host.Enabled)
}

func (f *File) changeToSplitted(profiles []string, status host.Status) error {
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

func invertStatus(s host.Status) host.Status {
	if s == host.Enabled {
		return host.Disabled
	}

	return host.Enabled
}
func (f *File) changeTo(profiles []string, status host.Status) error {
	for _, name := range profiles {
		if name == host.Default {
			continue
		}

		p, ok := f.data.Profiles[name]
		if !ok {
			return errors.ErrUnknownProfile
		}

		p.Status = status
		f.data.Profiles[name] = p
	}

	return nil
}
