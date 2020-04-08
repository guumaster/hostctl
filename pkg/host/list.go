package host

import (
	"os"
	"strings"

	"github.com/guumaster/tablewriter"
)

// DefaultColumns is the list of default columns to use when showing table list
var DefaultColumns = []string{"profile", "status", "ip", "domain"}

// ProfilesOnlyColumns are the columns used for profile status list
var ProfilesOnlyColumns = []string{"profile", "status"}

type ProfileStatus string

const (
	// Enabled marks a profile active on your hosts file.
	Enabled ProfileStatus = "on"
	// Disabled marks a profile not active on your hosts file.
	Disabled ProfileStatus = "off"
)

// ListOptions contains available options for listing.
type ListOptions struct {
	Profile      string
	RawTable     bool
	Columns      []string
	ProfilesOnly bool
	StatusFilter ProfileStatus
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

	if len(opts.Columns) == 0 {
		opts.Columns = DefaultColumns
	}
	if opts.ProfilesOnly {
		opts.Columns = ProfilesOnlyColumns
	}

	table.SetHeader(opts.Columns)

	if opts.RawTable {
		table.SetAutoWrapText(false)
		table.SetAutoFormatHeaders(true)
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetCenterSeparator("")
		table.SetColumnSeparator("\t")
		table.SetRowSeparator("")
		table.SetHeaderLine(false)
		table.SetBorder(false)
		table.SetTablePadding("\t") // pad with tabs
		table.SetNoWhiteSpace(true)
	}

	// First check if default should be shown
	if (profile == "default" || profile == "") && !opts.ProfilesOnly {
		appendProfile("default", table, h.profiles["default"], opts)

		if len(h.profiles) > 1 && !opts.RawTable {
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

		appendProfile(p, table, data, opts)

		if i < len(h.profiles) && !opts.RawTable {
			table.AddSeparator()
		}
	}
	table.Render()
	return nil
}

func appendProfile(profile string, table *tablewriter.Table, data hostLines, opts *ListOptions) {
	for _, r := range data {
		if r == "" {
			continue
		}
		if !IsHostLine(r) {
			continue
		}
		rs := strings.Split(cleanLine(r), " ")

		status := Enabled
		ip, domain := rs[0], rs[1]
		if IsDisabled(r) {
			// skip empty comments lines
			if rs[1] == "" {
				continue
			}
			status = Disabled
			ip, domain = rs[1], rs[2]
		}
		if opts.StatusFilter != "" && status != opts.StatusFilter {
			continue
		}
		if opts.ProfilesOnly {
			table.Append([]string{profile, string(status)})
			return
		}
		var row []string
		for _, c := range opts.Columns {
			switch c {
			case "profile":
				row = append(row, profile)
			case "status":
				row = append(row, string(status))
			case "ip", "ips":
				row = append(row, ip)
			case "domain", "domains":
				row = append(row, domain)
			}
		}
		table.Append(row)
	}
}
