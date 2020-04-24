package actions

import (
	"errors"
)

var (
	// ErrIncompatibleAllFlag when you can't use --all flag
	ErrIncompatibleAllFlag = errors.New("args must be empty with --all flag")

	// ErrMultipleProfiles when you can use only a single profile
	ErrMultipleProfiles = errors.New("specify only one profile")

	// ErrEmptyProfiles when trying to update empty profile list
	ErrEmptyProfiles = errors.New("there are no profiles")

	// ErrReadingFile when a file can't be read
	ErrReadingFile = errors.New("error reading data from file")
)
