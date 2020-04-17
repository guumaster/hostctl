package host

import (
	"errors"
)

var (
	// ErrMissingProfile when the profile is mandatory
	ErrMissingProfile = errors.New("missing profile name")

	// ErrUnknownProfile when the profile is not present
	ErrUnknownProfile = errors.New("unknown profile name")

	// ErrDefaultProfileError when trying to edit default content
	ErrDefaultProfileError = errors.New("'default' profile should not be handled by hostctl")

	// ErrMissingDomainsError when trying to set/add domains and none were given
	ErrMissingDomainsError = errors.New("no domains provided")

	// ErrMissingDestError when trying to write to a file
	ErrMissingDestError = errors.New("missing destination file")

	// ErrMissingSourceError when trying to read from a file
	ErrMissingSourceError = errors.New("missing source file")

	// ErrSnapConfinement when trying to read files on snap installation
	ErrSnapConfinement = errors.New("can't use --from or --host-file. " +
		"Snap confinement restrictions doesn't allow to read other than /etc/hosts file")
)
