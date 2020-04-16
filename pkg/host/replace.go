package host

import (
	"errors"
)

// ReplaceProfile removes previous profile with same name and add new profile to the list
func (f *File) ReplaceProfile(p Profile) error {
	err := f.RemoveProfile(p.Name)
	if err != nil && !errors.Is(err, ErrUnknownProfile) {
		return err
	}

	return f.AddProfile(p)
}
