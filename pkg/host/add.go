package host

import (
	"os"
)

// commonAddOptions contains common options for adding.
type commonAddOptions struct {
	Dst     string
	Profile string
	Reset   bool
}

// AddFromFileOptions contains available options for adding from file.
type AddFromFileOptions struct {
	Dst     string
	Profile string
	Reset   bool
	From    string
}

// AddFromArgsOptions contains available options for adding from arguments.
type AddFromArgsOptions struct {
	Dst     string
	Profile string
	Reset   bool
	Domains []string
	IP      string
}

// AddFromFile reads content from a file and adds it as a profile into your hosts file.
// If you pass reset=true it will delete all previous content of the profile.
func AddFromFile(opts *AddFromFileOptions) error {
	if opts.From == "" {
		return MissingSourceError
	}
	newData, _ := ReadHostFileStrict(opts.From)

	return add(newData, &commonAddOptions{
		opts.Dst,
		opts.Profile,
		opts.Reset,
	})
}

func AddFromArgs(opts *AddFromArgsOptions) error {
	if len(opts.Domains) == 0 {
		return MissingDomainsError
	}
	newData := ReadFromArgs(opts.Domains, opts.IP)

	return add(newData, &commonAddOptions{
		opts.Dst,
		opts.Profile,
		opts.Reset,
	})
}

func add(n *hostFile, opts *commonAddOptions) error {
	if opts.Dst == "" {
		return MissingDestError
	}
	if opts.Profile == "" {
		opts.Profile = "default"
	}

	currData, err := ReadHostFile(opts.Dst)
	if err != nil {
		return err
	}

	if opts.Reset {
		currData.profiles[opts.Profile] = hostLines{}
	}
	currData.profiles[opts.Profile] = append(currData.profiles[opts.Profile], n.profiles["default"]...)

	dstFile, err := os.OpenFile(opts.Dst, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	return WriteToFile(dstFile, currData)
}
