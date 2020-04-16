package host

func (f *File) MergeProfiles(content *Content) {
	for _, newName := range content.ProfileNames {
		newP := content.Profiles[newName]

		_, ok := f.data.Profiles[newName]
		if !ok {
			f.data.ProfileNames = append(f.data.ProfileNames, newName)
			f.data.Profiles[newName] = newP

			continue
		}

		baseP := f.data.Profiles[newName]
		if baseP.Routes == nil {
			baseP.Routes = map[string]*Route{}
		}

		for _, r := range newP.Routes {
			ip := r.IP.String()
			if _, ok := baseP.Routes[ip]; ok {
				baseP.Routes[ip].HostNames = append(baseP.Routes[ip].HostNames, r.HostNames...)
			} else {
				baseP.appendIP(ip)
				baseP.Routes[ip] = r
			}
		}

		f.data.Profiles[newName] = baseP
	}
}
