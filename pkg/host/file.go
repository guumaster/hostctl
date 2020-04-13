package host

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/afero"
)

func NewFile(src string) (*File, error) {
	return NewWithFs(src, afero.NewOsFs())
}

func NewWithFs(src string, fs afero.Fs) (*File, error) {
	if fs == nil {
		fs = afero.NewOsFs()
	}
	s, err := fs.Open(src)
	if err != nil {
		return nil, err
	}
	f := &File{
		src: s,
		fs:  fs,
	}
	err = f.read()
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (f *File) read() error {
	_, _ = f.src.Seek(0, io.SeekStart)
	data, err := Parse(f.src)
	f.data = data
	return err
}

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

func (f *File) GetEnabled() []string {
	enabled := []string{}
	for _, name := range f.data.ProfileNames {
		if f.data.Profiles[name].Status == Enabled {
			enabled = append(enabled, name)
		}
	}
	return enabled
}

func (f *File) GetProfile(name string) (*Profile, error) {
	profile, ok := f.data.Profiles[name]
	if !ok {
		return nil, UnknownProfileError
	}
	return &profile, nil
}

func (f *File) GetProfileNames() []string {
	return f.data.ProfileNames
}

func (f *File) GetDisabled() []string {
	disabled := []string{}
	for _, name := range f.data.ProfileNames {
		if f.data.Profiles[name].Status == Disabled {
			disabled = append(disabled, name)
		}
	}
	return disabled
}

func (f *File) AddRoutes(name, ip string, hostnames []string) error {
	profile, err := f.GetProfile(name)
	if err != nil {
		return err
	}

	profile.AddRoutes(ip, hostnames)
	return nil
}

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
		return fmt.Errorf("no content to write")
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
		if name == "default" {
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
