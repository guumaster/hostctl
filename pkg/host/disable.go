package host

import (
	"fmt"
)

// Disable marks a profile as disable by commenting all hosts lines.
// The content remains on the file and can be enabled later.
func Disable(dst, profile string) error {
	h, err := getHostData(dst, profile)
	if err != nil {
		return err
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

	return writeHostData(dst, h)
}

func disableProfile(h *hostFile, profile string) {
	for i, r := range h.profiles[profile] {
		if !IsDisabled(r) {
			h.profiles[profile][i] = fmt.Sprintf("# %s", r)
		}
	}
}
