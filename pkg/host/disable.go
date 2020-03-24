package host

import (
	"fmt"
	"os"
)

// Disable marks a profile as disable by commenting all hosts lines.
// The content remains on the file and can be enabled later.
func Disable(dst, profile string) error {
	if dst == "" {
		return MissingDestError
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
				disableProfile(h, p)
			}
		}
	} else {
		disableProfile(h, profile)
	}

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	return WriteToFile(dstFile, h)
}

func disableProfile(h *hostFile, profile string) {
	for i, r := range h.profiles[profile] {
		if !IsDisabled(r) {
			h.profiles[profile][i] = fmt.Sprintf("# %s", r)
		}
	}
}
