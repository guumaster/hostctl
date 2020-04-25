package file

import (
	"github.com/guumaster/hostctl/pkg/types"
)

// MergeFile joins new content with existing content
func (f *File) MergeFile(from *File) {
	var ps []*types.Profile
	for _, p := range from.data.Profiles {
		ps = append(ps, p)
	}

	f.MergeProfiles(ps)
}

// MergeProfiles joins new profiles with existing content
func (f *File) MergeProfiles(profiles []*types.Profile) {
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
			baseP.Routes = map[string]*types.Route{}
		}

		routes := []*types.Route{}
		for _, r := range newP.Routes {
			routes = append(routes, r)
		}

		baseP.AddRoutes(routes)

		f.data.Profiles[newName] = baseP
	}
}
