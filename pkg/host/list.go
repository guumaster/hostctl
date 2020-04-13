package host

import (
	"io"

	"github.com/guumaster/tablewriter"
)

// DefaultColumns is the list of default columns to use when showing table list
var DefaultColumns = []string{"profile", "status", "ip", "domain"}

// ProfilesOnlyColumns are the columns used for profile status list
var ProfilesOnlyColumns = []string{"profile", "status"}

// ListOptions contains available options for listing.
type ListOptions struct {
	Profiles     []string
	RawTable     bool
	Columns      []string
	ProfilesOnly bool
	StatusFilter ProfileStatus
	Writer       io.Writer
}

func includeProfile(needle string, stack []string) bool {
	if len(stack) == 0 {
		return true
	}
	for _, s := range stack {
		if s == needle {
			return true
		}
	}
	return false
}

// ProfileStatus shows a table only with profile names status
func (f *File) ProfileStatus(opts *ListOptions) {
	opts.Columns = ProfilesOnlyColumns

	table := createTableWriter(opts)

	for _, name := range f.data.ProfileNames {
		currProfile := f.data.Profiles[name]
		if !includeProfile(name, opts.Profiles) {
			continue
		}

		table.Append([]string{currProfile.Name, currProfile.GetStatus()})
	}

	table.Render()
}

// List shows a table with profile names status and routing information
func (f *File) List(opts *ListOptions) {
	if len(opts.Columns) == 0 {
		opts.Columns = DefaultColumns
	}

	table := createTableWriter(opts)

	added := addDefault(f, table, opts)
	if added && len(f.data.Profiles) > 0 && !opts.RawTable {
		table.AddSeparator()
	}
	for _, name := range f.data.ProfileNames {
		added := addProfiles(f.data.Profiles[name], table, opts)
		if added && !opts.RawTable {
			table.AddSeparator()
		}
	}

	table.Render()
}

func addDefault(f *File, table *tablewriter.Table, opts *ListOptions) bool {
	// First check if default should be shown
	if !includeProfile("default", opts.Profiles) {
		return false
	}

	i := 0
	for _, line := range f.data.DefaultProfile {
		i++
		if line.Comment == "" && line.Profile != "" {
			row := getRow(line, opts.Columns)
			if len(row) > 0 {
				table.Append(row)
			}
		}
	}
	return i > 0
}

func addProfiles(p Profile, table *tablewriter.Table, opts *ListOptions) bool {
	if !includeProfile(p.Name, opts.Profiles) {
		return false
	}

	if opts.StatusFilter != "" && p.Status != opts.StatusFilter {
		return false
	}

	for _, route := range p.Routes {
		for _, h := range route.HostNames {
			line := &tableRow{
				Profile: p.Name,
				Status:  p.GetStatus(),
				IP:      route.IP.String(),
				Host:    h,
			}
			row := getRow(line, opts.Columns)
			if len(row) > 0 {
				table.Append(row)
			}
		}
	}

	return true
}
