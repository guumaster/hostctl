package host

import (
	"errors"
)

// ValidProfile controls that you don't ask to handle 'default' profile.
func ValidProfile(p string) error {
	if p == "default" {
		return errors.New("'default' profile should not be handled by hostctl")
	}

	return nil
}

// NotEmptyProfile controls the profile name.
func NotEmptyProfile(p string) error {
	if p == "" {
		return errors.New("profile can't be empty")
	}

	return nil
}
