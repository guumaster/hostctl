package host

// Enable marks profiles as enable by uncommenting all hosts lines
// making the routing work again.
func (f *File) Enable(profiles []string) error {
	for _, p := range profiles {
		if p == "default" {
			continue
		}
		profile, ok := f.data.Profiles[p]
		if !ok {
			return UnknownProfileError
		}
		profile.Status = Enabled
		f.data.Profiles[p] = profile
	}
	return nil
}

// EnableAll marks all profiles as enable by uncommenting all hosts lines
// making the routing work again.
func (f *File) EnableAll() error {
	return f.Enable(f.data.ProfileNames)
}

// EnableOnly marks profiles as enable and disable all other profiles
func (f *File) EnableOnly(profiles []string) error {
	for _, name := range f.data.ProfileNames {
		if name == "default" {
			continue
		}
		profile, ok := f.data.Profiles[name]
		if !ok {
			return UnknownProfileError
		}

		if contains(profiles, name) {
			profile.Status = Enabled
		} else {
			profile.Status = Disabled
		}
		f.data.Profiles[name] = profile
	}
	return nil
}
