package host

import (
	"strings"
)

// RemoveProfile removes a profile from a hosts file.
func RemoveProfile(dst, profile string) error {
	h, err := getHostData(dst, profile)
	if err != nil {
		return err
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

	return writeHostData(dst, h)
}

// RemoveDomains removes domains from a hosts file.
func RemoveDomains(dst, profile string, domains []string) error {
	if len(domains) == 0 {
		return MissingDomainsError
	}
	h, err := getHostData(dst, profile)
	if err != nil {
		return err
	}

	lines := h.profiles[profile]
	h.profiles[profile] = removeFromProfile(lines, domains)
	if len(h.profiles[profile]) == 0 {
		delete(h.profiles, profile)
	}

	return writeHostData(dst, h)
}

func removeFromProfile(lines hostLines, remove []string) hostLines {
	newProfile := make([]string, 0)
	for _, l := range lines {
		rs := strings.Split(cleanLine(l), " ")
		domain := rs[1]
		if IsDisabled(l) {
			// skip empty comments lines
			if rs[1] == "" {
				continue
			}
			domain = rs[2]
		}
		canAdd := true
		for _, r := range remove {
			if r == domain {
				canAdd = false
				break
			}
		}
		if canAdd {
			newProfile = append(newProfile, l)
		}
	}

	return newProfile
}
