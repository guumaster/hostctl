package file

import (
	"github.com/guumaster/hostctl/pkg/host"
)

// MergeProfiles joins new content with existing content
func (f *File) MergeProfiles(profiles []*host.Profile) {
	for _, newP := range profiles {
		newName := newP.Name

		_, ok := f.data.Profiles[newName]
		if !ok {
			f.data.ProfileNames = append(f.data.ProfileNames, newName)
			f.data.Profiles[newName] = newP

			continue
		}

		baseP := f.data.Profiles[newName]
		if baseP.Routes == nil {
			baseP.Routes = map[string]*host.Route{}
		}

		for _, r := range newP.Routes {
			ip := r.IP.String()
			baseP.AddRoutes(ip, r.HostNames)
		}

		f.data.Profiles[newName] = baseP
	}
}
