package host

import (
	"os"
	"strings"

	"github.com/guumaster/tablewriter"
)

// ListOptions contains available options for listing.
type ListOptions struct {
	Profile string
}

// ListProfiles shows a table with profile names status and routing information
func ListProfiles(src string, opts *ListOptions) error {
	var profile = ""
	if opts.Profile != "" {
		profile = opts.Profile
	}
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	h, err := Read(f, true)
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Profile", "Status", "IP", "Domain"})

	if profile == "default" || profile == "" {
		appendProfile("default", table, h.profiles["default"])

		if len(h.profiles) > 1 {
			table.AddSeparator()
		}
	}

	i := 0
	for p, data := range h.profiles {
		i++
		if profile != "" && p != profile {
			continue
		}
		if p == "default" {
			continue
		}

		appendProfile(p, table, data)

		if i < len(h.profiles) {
			table.AddSeparator()
		}
	}
	table.Render()
	return nil
}

func appendProfile(profile string, table *tablewriter.Table, data hostLines) {
	for _, r := range data {
		if r == "" {
			continue
		}
		if !IsHostLine(r) {
			continue
		}
		rs := strings.Split(cleanLine(r), " ")

		status := "on"
		ip, domain := rs[0], rs[1]
		if IsDisabled(r) {
			// skip empty comments lines
			if rs[1] == "" {
				continue
			}
			status = "off"
			ip, domain = rs[1], rs[2]
		}
		table.Append([]string{
			profile,
			status,
			ip,
			domain,
		})
	}
}
