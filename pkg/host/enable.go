package host

// Enable marks a profile as enable by uncommenting all hosts lines
// making the routing work again.
func Enable(dst, profile string) error {
	h, err := getHostData(dst, profile)
	if err != nil {
		return err
	}

	if profile == "" {
		for p := range h.profiles {
			if p != "default" {
				enableProfile(h, p)
			}
		}
	} else {
		enableProfile(h, profile)
	}

	return writeHostData(dst, h)
}

func enableProfile(h *hostFile, profile string) {
	for i, r := range h.profiles[profile] {
		if IsDisabled(r) {
			h.profiles[profile][i] = EnableLine(r)
		}
	}

}
