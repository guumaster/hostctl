package host

// Toggle alternates between enable and disable status of a profile.
func (f *File) Toggle(profiles []string) error {
	for _, p := range profiles {
		if p == "default" {
			continue
		}
		profile, ok := f.data.Profiles[p]
		if !ok {
			return UnknownProfileError
		}
		if profile.Status == Enabled {
			profile.Status = Disabled
		} else {
			profile.Status = Enabled

		}
		f.data.Profiles[p] = profile
	}
	return nil
}
