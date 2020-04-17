package host

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/spf13/afero"
)

// NewFile creates a new File from the given src on default OS filesystem
func NewFile(src string) (*File, error) {
	return NewWithFs(src, afero.NewOsFs())
}

// NewWithFs creates a new File with src and an existing filesystem
func NewWithFs(src string, fs afero.Fs) (*File, error) {
	if fs == nil {
		fs = afero.NewOsFs()
	}

	s, err := fs.Open(src)
	if err != nil {
		return nil, err
	}

	f := &File{src: s, fs: fs}

	_, _ = f.src.Seek(0, io.SeekStart)
	data, err := Parse(f.src)
	f.data = data

	if err != nil {
		return nil, err
	}

	return f, nil
}

// GetStatus returns a map with the status of the given profiles
func (f *File) GetStatus(profiles []string) map[string]ProfileStatus {
	st := map[string]ProfileStatus{}

	for _, name := range profiles {
		profile, ok := f.data.Profiles[name]
		if !ok {
			continue
		}

		st[name] = profile.Status
	}

	return st
}

// GetEnabled returns a list of profiles that are Enabled
func (f *File) GetEnabled() []string {
	enabled := []string{}

	for _, name := range f.data.ProfileNames {
		if f.data.Profiles[name].Status == Enabled {
			enabled = append(enabled, name)
		}
	}

	return enabled
}

// GetDisabled returns a list of profiles that are Enabled
func (f *File) GetDisabled() []string {
	disabled := []string{}

	for _, name := range f.data.ProfileNames {
		if f.data.Profiles[name].Status == Disabled {
			disabled = append(disabled, name)
		}
	}

	return disabled
}

// GetProfile return a Profile from the list
func (f *File) GetProfile(name string) (*Profile, error) {
	profile, ok := f.data.Profiles[name]
	if !ok {
		return nil, ErrUnknownProfile
	}

	return profile, nil
}

// GetProfileNames return a list of all profile names
func (f *File) GetProfileNames() []string {
	return f.data.ProfileNames
}

// AddRoutes add route information to a given profile
func (f *File) AddRoutes(name, ip string, hostnames []string) error {
	profile, err := f.GetProfile(name)
	if err != nil && !errors.Is(err, ErrUnknownProfile) {
		return err
	}

	if profile == nil {
		p := Profile{
			Name:   name,
			Status: Enabled,
			Routes: map[string]*Route{},
		}

		p.AddRoutes(ip, hostnames)

		return f.AddProfile(p)
	}

	profile.AddRoutes(ip, hostnames)

	return nil
}

// RemoveRoutes removes route information from a given profile.
// also removes the profile if gets empty.
func (f *File) RemoveRoutes(name string, routes []string) (bool, error) {
	p, err := f.GetProfile(name)
	if err != nil {
		return false, err
	}

	p.RemoveRoutes(routes)

	if len(p.Routes) == 0 {
		err := f.RemoveProfile(p.Name)
		if err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}

// WriteTo overwrite file with hosts info
func (f *File) WriteTo(src string) error {
	h, err := f.fs.OpenFile(src, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	return f.writeToFile(h)
}

// Flush overwrite file with hosts info
func (f *File) Flush() error {
	h, err := f.fs.OpenFile(f.src.Name(), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer h.Close()

	return f.writeToFile(h)
}

// writeToFile overwrite file with hosts info
func (f *File) writeToFile(dst afero.File) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if f.data == nil {
		return ErrNoContent
	}

	err := dst.Truncate(0)
	if err != nil {
		return err
	}

	err = f.data.DefaultProfile.Render(dst)
	if err != nil {
		return err
	}

	f.writeBanner(dst)

	for _, name := range f.data.ProfileNames {
		if name == Default {
			continue
		}

		profile := f.data.Profiles[name]

		err := profile.Render(dst)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *File) writeBanner(w io.StringWriter) {
	if f.hasBanner {
		return
	}

	_, _ = w.WriteString(fmt.Sprintf("%s\n", banner))
	f.hasBanner = true
}

// Close closes the underlying file
func (f *File) Close() {
	f.src.Close()
}

func contains(s []string, n string) bool {
	for _, x := range s {
		if x == n {
			return true
		}
	}

	return false
}

func remove(s []string, n string) []string {
	list := []string{}

	for _, x := range s {
		if x != n {
			list = append(list, x)
		}
	}

	return list
}
