package host

import (
	"os"
	"strings"

	"github.com/guumaster/tablewriter"
)

// DefaultColumns is the list of default columns to use when showing table list
var DefaultColumns = []string{"profile", "status", "ip", "domain"}

// ListOptions contains available options for listing.
type ListOptions struct {
	Profile  string
	RawTable bool
	Columns  []string
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

	cols := opts.Columns
	if len(cols) == 0 {
		cols = DefaultColumns
	}

	table.SetHeader(cols)

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

	if profile == "default" || profile == "" {
		appendProfile("default", table, cols, h.profiles["default"])

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

		appendProfile(p, table, cols, data)

		if i < len(h.profiles) && !opts.RawTable {
			table.AddSeparator()
		}
	}
	table.Render()
	return nil
}

func appendProfile(profile string, table *tablewriter.Table, cols []string, data hostLines) {
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
		var row []string
		for _, c := range cols {
			switch c {
			case "profile":
				row = append(row, profile)
			case "status":
				row = append(row, status)
			case "ip", "ips":
				row = append(row, ip)
			case "domain", "domains":
				row = append(row, domain)
			}
		}
		table.Append(row)
	}
}
