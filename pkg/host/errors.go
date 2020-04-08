package host

import (
	"errors"
)

// MissingProfileError when the profile is mandatory
var MissingProfileError = errors.New("missing profile name")

// UnknownProfileError when the profile is not present
var UnknownProfileError = errors.New("unknown profile name")

// DefaultProfileError when trying to edit default content
var DefaultProfileError = errors.New("'default' profile should not be handled by hostctl")

// MissingDomainsError when trying to set/add domains and none were given
var MissingDomainsError = errors.New("no domains provided")

// MissingSourceError when trying to write to a file
var MissingDestError = errors.New("missing destination file")

// MissingSourceError when trying to read from a file
var MissingSourceError = errors.New("missing source file")
