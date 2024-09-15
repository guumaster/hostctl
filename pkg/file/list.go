package file

import (
	"github.com/guumaster/hostctl/pkg/types"
)

// ListOptions contains available options for listing.
type ListOptions struct {
	Renderer     types.Renderer
	Profiles     []string
	ProfilesOnly bool
	StatusFilter types.Status
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

// ProfileStatus shows a table only with profile names status.
func (f *File) ProfileStatus(r types.Renderer, profiles []string) {
	for _, name := range f.data.ProfileNames {
		currProfile := f.data.Profiles[name]

		if profiles != nil && !includeProfile(name, profiles) {
			continue
		}

		r.AppendRow(&types.Row{
			Profile: currProfile.Name,
			Status:  currProfile.GetStatus(),
		})
	}

	_ = r.Render()
}

// List shows a table with profile names status and routing information.
func (f *File) List(r types.Renderer, opts *ListOptions) {
	addDefault(f, r, opts)

	for _, name := range f.data.ProfileNames {
		addProfiles(f.data.Profiles[name], r, opts)
	}

	_ = r.Render()
}

func addDefault(f *File, r types.Renderer, opts *ListOptions) {
	// First check if default should be shown
	if !includeProfile(types.Default, opts.Profiles) {
		return
	}

	for _, row := range f.data.DefaultProfile {
		if row.Comment == "" && row.Profile != "" {
			r.AppendRow(row)
		}
	}

	if len(f.data.DefaultProfile) > 0 && len(f.data.Profiles) > 0 {
		r.AddSeparator()
	}
}

func addProfiles(p *types.Profile, r types.Renderer, opts *ListOptions) {
	if !includeProfile(p.Name, opts.Profiles) {
		return
	}

	if opts.StatusFilter != "" && p.Status != opts.StatusFilter {
		return
	}

	for _, ip := range p.IPList {
		route := p.Routes[ip]
		r.AppendRow(&types.Row{
			Profile: p.Name,
			Status:  p.GetStatus(),
			IP:      route.IP.String(),
			Hosts:    route.HostNames,
		})
	}

	if len(p.IPList) > 0 {
		r.AddSeparator()
	}
}
