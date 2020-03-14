package host

import (
	"errors"
	"os"
)

// AddFromFile reads content from a file and adds it as a profile into your hosts file.
// If you pass reset=true it will delete all previous content of the profile.
func AddFromFile(from, dst, profile string, reset bool) error {
	if from == "" {
		return errors.New("missing source file")
	}
	if dst == "" {
		return errors.New("missing destination file")
	}
	if profile == "" {
		profile = "default"
	}

	currData, err := ReadHostFile(dst)
	if err != nil {
		return err
	}
	newData, _ := ReadHostFileStrict(from)

	if reset {
		currData.profiles[profile] = hostLines{}
	}
	currData.profiles[profile] = append(currData.profiles[profile], newData.profiles["default"]...)

	dstFile, err := os.OpenFile(dst, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	return WriteToFile(dstFile, currData)
}
