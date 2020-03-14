package host

import (
	"errors"
)

// CheckProfile controls that you don't ask to handle 'default' profile.
func CheckProfile(p string) error {
	if p == "default" {
		return errors.New("'default' profile should not be handled by hostctl")
	}

	return nil
}
