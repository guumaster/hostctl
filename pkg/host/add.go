package host

// AddProfile adds a profile to the list
func (f *File) AddProfile(profile Profile) error {
	if profile.Name == Default {
		return ErrDefaultProfileError
	}

	f.MergeProfiles(&Content{
		ProfileNames: []string{profile.Name},
		Profiles: map[string]*Profile{
			profile.Name: &profile,
		},
	})

	return nil
}
