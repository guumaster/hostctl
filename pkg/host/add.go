package host

func (f *File) AddProfile(profile Profile) error {
	if profile.Name == "default" {
		return DefaultProfileError
	}
	f.MergeProfiles(&Content{
		ProfileNames: []string{profile.Name},
		Profiles: map[string]Profile{
			profile.Name: profile,
		},
	})
	return nil
}
