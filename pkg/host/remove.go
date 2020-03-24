package host

import (
	"strings"
)

type RemoveProfileOptions struct {
	Dst     string
	Profile string
}

type RemoveDomainsOptions struct {
	Dst     string
	Profile string
	Domains []string
}

// RemoveProfile removes a profile from a hosts file.
func RemoveProfile(opts *RemoveProfileOptions) error {
	h, err := getHostData(opts.Dst, opts.Profile)
	if err != nil {
		return err
	}

	if opts.Profile == "" {
		for p := range h.profiles {
			if p != "default" {
				delete(h.profiles, p)
			}
		}
	} else {
		delete(h.profiles, opts.Profile)
	}

	return writeHostData(opts.Dst, h)
}

// RemoveDomains removes domains from a hosts file.
func RemoveDomains(opts *RemoveDomainsOptions) error {
	h, err := getHostData(opts.Dst, opts.Profile)
	if err != nil {
		return err
	}

	if opts.Profile == "" {
		for p := range h.profiles {
			if p != "default" {
				domains := h.profiles[p]
				h.profiles[p] = removeFromProfile(domains, opts.Domains)
			}
		}
	} else {
		domains := h.profiles[opts.Profile]
		h.profiles[opts.Profile] = removeFromProfile(domains, opts.Domains)
	}

	return writeHostData(opts.Dst, h)
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
