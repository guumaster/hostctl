package file

import (
	"github.com/guumaster/hostctl/pkg/types"
)

// RemoveProfiles removes given profiles from the list
func (f *File) RemoveProfiles(profiles []string) error {
	for _, p := range profiles {
		err := f.RemoveProfile(p)
		if err != nil {
			return err
		}
	}

	return nil
}

// RemoveProfile removes given profile from the list
func (f *File) RemoveProfile(name string) error {
	var names []string

	if name == types.Default {
		return types.ErrDefaultProfile
	}

	_, ok := f.data.Profiles[name]
	if !ok {
		return types.ErrUnknownProfile
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
