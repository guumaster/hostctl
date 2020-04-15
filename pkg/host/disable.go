package host

// Disable marks profiles as disable by commenting all hosts lines.
// The content remains on the file and can be enabled later.
func (f *File) Disable(profiles []string) error {
	for _, p := range profiles {
		if p == "default" {
			continue
		}
		profile, ok := f.data.Profiles[p]
		if !ok {
			return UnknownProfileError
		}
		profile.Status = Disabled
		f.data.Profiles[p] = profile
	}
	return nil
}

// DisableAll marks all profiles as disable by commenting all hosts lines.
// The content remains on the file and can be enabled later.
func (f *File) DisableAll() error {
	return f.Disable(f.data.ProfileNames)
}

// DisableOnly marks profiles as enable and disable all other profiles
func (f *File) DisableOnly(profiles []string) error {
	for _, name := range f.data.ProfileNames {
		if name == "default" {
			continue
		}
		profile, ok := f.data.Profiles[name]
		if !ok {
			return UnknownProfileError
		}

		if contains(profiles, name) {
			profile.Status = Disabled
		} else {
			profile.Status = Enabled
		}
		f.data.Profiles[name] = profile
	}
	return nil
}
