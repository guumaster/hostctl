package host

func (f *File) RemoveProfiles(profiles []string) error {
	for _, p := range profiles {
		err := f.RemoveProfile(p)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *File) RemoveProfile(name string) error {
	var names []string

	if name == Default {
		return ErrDefaultProfileError
	}

	_, ok := f.data.Profiles[name]
	if !ok {
		return ErrUnknownProfile
	}

	delete(f.data.Profiles, name)

	for _, n := range f.data.ProfileNames {
		if n != name {
			names = append(names, n)
		}
	}

	f.data.ProfileNames = names

	return nil
}
