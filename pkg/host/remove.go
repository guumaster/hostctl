package host

import (
	"errors"
	"fmt"
	"os"
)

// Remove removes a profile from a hosts file.
func Remove(dst, profile string) error {
	if dst == "" {
		return errors.New("missing destination file")
	}

	h, err := ReadHostFile(dst)
	if err != nil {
		return err
	}
	_, ok := h.profiles[profile]
	if profile != "" && !ok {
		return fmt.Errorf("profile '%s' doesn't exists in file", profile)
	}

	if profile == "" {
		for p := range h.profiles {
			if p != "default" {
				delete(h.profiles, p)
			}
		}
	} else {
		delete(h.profiles, profile)
	}

	dstFile, err := os.OpenFile(dst, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	return WriteToFile(dstFile, h)
}
