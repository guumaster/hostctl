package file

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/spf13/afero"

	"github.com/guumaster/hostctl/pkg/profile"
	"github.com/guumaster/hostctl/pkg/types"
)

// File container to handle a hosts file
type File struct {
	fs        afero.Fs
	src       afero.File
	data      *types.Content
	hasBanner bool
	mutex     sync.Mutex
}

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
	data, err := profile.Parse(f.src)
	f.data = data

	if err != nil {
		return nil, err
	}

	return f, nil
}

// GetStatus returns a map with the status of the given profiles
func (f *File) GetStatus(profiles []string) map[string]types.Status {
	st := map[string]types.Status{}

	for _, name := range profiles {
		p, ok := f.data.Profiles[name]
		if !ok {
			continue
		}

		st[name] = p.Status
	}

	return st
}

// GetEnabled returns a list of profiles that are Enabled
func (f *File) GetEnabled() []string {
	enabled := []string{}

	for _, name := range f.data.ProfileNames {
		if f.data.Profiles[name].Status == types.Enabled {
			enabled = append(enabled, name)
		}
	}

	return enabled
}

// GetDisabled returns a list of profiles that are Enabled
func (f *File) GetDisabled() []string {
	disabled := []string{}

	for _, name := range f.data.ProfileNames {
		if f.data.Profiles[name].Status == types.Disabled {
			disabled = append(disabled, name)
		}
	}

	return disabled
}

// GetProfile return a Profile from the list
func (f *File) GetProfile(name string) (*types.Profile, error) {
	p, ok := f.data.Profiles[name]
	if !ok {
		return nil, types.ErrUnknownProfile
	}

	return p, nil
}

// GetProfileNames return a list of all profile names
func (f *File) GetProfileNames() []string {
	return f.data.ProfileNames
}

// AddRoute adds a single route information to a given profile
func (f *File) AddRoute(name string, route *types.Route) error {
	return f.AddRoutes(name, []*types.Route{route})
}

// AddRoutes adds routes information to a given profile
func (f *File) AddRoutes(name string, routes []*types.Route) error {
	p, err := f.GetProfile(name)
	if err != nil && !errors.Is(err, types.ErrUnknownProfile) {
		return err
	}

	if p == nil {
		p = &types.Profile{
			Name:   name,
			Status: types.Enabled,
			Routes: map[string]*types.Route{},
		}

		p.AddRoutes(routes)

		return f.AddProfile(p)
	}

	p.AddRoutes(routes)

	return nil
}

// RemoveHostnames removes route information from a given types.
// also removes the profile if gets empty.
func (f *File) RemoveHostnames(name string, routes []string) (bool, error) {
	p, err := f.GetProfile(name)
	if err != nil {
		return false, err
	}

	p.RemoveHostnames(routes)

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
		return types.ErrNoContent
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
		if name == types.Default {
			continue
		}

		p := f.data.Profiles[name]

		err := p.Render(dst)
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

	_, _ = w.WriteString(fmt.Sprintf("%s\n", Banner))
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
