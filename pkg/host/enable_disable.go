package host

// Enable marks profiles as enable by uncommenting all hosts lines
// making the routing work again.
func (f *File) Enable(profiles []string) error {
	return f.changeTo(profiles, Enabled)
}

// Disable marks profiles as disable by commenting all hosts lines.
// The content remains on the file and can be enabled later.
func (f *File) Disable(profiles []string) error {
	return f.changeTo(profiles, Disabled)
}

// EnableAll marks all profiles as enable by uncommenting all hosts lines
// making the routing work again.
func (f *File) EnableAll() error {
	return f.changeTo(f.GetProfileNames(), Enabled)
}

// DisableAll marks all profiles as disable by commenting all hosts lines.
// The content remains on the file and can be enabled later.
func (f *File) DisableAll() error {
	return f.changeTo(f.GetProfileNames(), Disabled)
}

// DisableOnly marks profiles as disable and enable all other profiles
func (f *File) DisableOnly(profiles []string) error {
	return f.changeToSplitted(profiles, Disabled)
}

// EnableOnly marks profiles as enable and disable all other profiles
func (f *File) EnableOnly(profiles []string) error {
	return f.changeToSplitted(profiles, Enabled)
}

func (f *File) changeToSplitted(profiles []string, status ProfileStatus) error {
	for _, name := range f.data.ProfileNames {
		if name == Default {
			continue
		}

		profile, ok := f.data.Profiles[name]
		if !ok {
			return ErrUnknownProfile
		}

		if contains(profiles, name) {
			profile.Status = status
		} else {
			profile.Status = invertStatus(status)
		}

		f.data.Profiles[name] = profile
	}

	return nil
}

func invertStatus(s ProfileStatus) ProfileStatus {
	if s == Enabled {
		return Disabled
	}

	return Enabled
}
func (f *File) changeTo(profiles []string, status ProfileStatus) error {
	for _, p := range profiles {
		if p == Default {
			continue
		}

		profile, ok := f.data.Profiles[p]
		if !ok {
			return ErrUnknownProfile
		}

		profile.Status = status
		f.data.Profiles[p] = profile
	}

	return nil
}
