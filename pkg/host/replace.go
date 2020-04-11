package host

import (
	"errors"
)

func (f *File) ReplaceProfile(p Profile) error {
	err := f.RemoveProfile(p.Name)
	if err != nil && !errors.Is(err, UnknownProfileError) {
		return err
	}

	return f.AddProfile(p)
}
