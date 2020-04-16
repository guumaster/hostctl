package host

import (
	"errors"
)

func (f *File) ReplaceProfile(p Profile) error {
	err := f.RemoveProfile(p.Name)
	if err != nil && !errors.Is(err, ErrUnknownProfile) {
		return err
	}

	return f.AddProfile(p)
}
