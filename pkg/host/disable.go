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
