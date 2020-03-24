package host

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"strings"
)

var (
	profileStart = regexp.MustCompile(`(?i)# profile\s+([a-z0-9-_]+)\s*`)
	profileEnd   = regexp.MustCompile(`(?i)# end\s*`)
	disableRe    = regexp.MustCompile(`^#\s*`)
	spaceRemover = regexp.MustCompile(`\s+`)
	tabReplacer  = regexp.MustCompile(`\t+`)
)

type hostFile struct {
	profiles profileMap
}

type profileMap map[string]hostLines

type hostLines []string

// ReadHostFile open a file an read content into a hostFile struct
func ReadHostFile(file string) (*hostFile, error) {
	fromFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	return Read(fromFile, false)
}

// ReadHostFileStrict open a file an read content into a hostFile struct. removes all comments.
func ReadHostFileStrict(file string) (*hostFile, error) {
	fromFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	return Read(fromFile, true)
}

// Read returns hosts file content grouped by profiles.
// If you pass strict=true it would remove all comments.
func Read(r io.Reader, strict bool) (*hostFile, error) {
	h := &hostFile{
		profiles: profileMap{},
	}

	ln := 0
	s := bufio.NewScanner(r)
	open := ""
	for s.Scan() {
		ln++
		b := s.Bytes()

		switch {

		case profileStart.Match(b):
			open = strings.TrimSpace(strings.Split(string(b), "# profile")[1])
			h.profiles[open] = []string{}

		case profileEnd.Match(b):
			open = ""

		case open != "":
			if strict && !IsHostLine(string(b)) {
				// skip
			} else {
				h.profiles[open] = append(h.profiles[open], string(b))
			}

		default:
			h.profiles["default"] = append(h.profiles["default"], string(b))
		}

		if err := s.Err(); err != nil {
			return nil, err
		}
	}
	return h, nil
}

// ReadFromArgs read arguments into a hostFile struct
func ReadFromArgs(domains []string, ip string) *hostFile {
	dom := make([]string, len(domains))
	for k, d := range domains {
		dom[k] = fmt.Sprintf("%s %s", ip, d)
	}
	newData := &hostFile{
		profiles: profileMap{
			"default": dom,
		},
	}
	return newData
}

// IsHostLine checks if a line is a host line or a comment line.
func IsHostLine(line string) bool {
	p := strings.Split(cleanLine(line), " ")
	i := 0
	if p[0] == "#" && len(p) > 1 {
		i = 1
	}
	ip := net.ParseIP(p[i])

	return ip != nil
}

func cleanLine(line string) string {
	return tabReplacer.ReplaceAllString(spaceRemover.ReplaceAllString(line, " "), " ")
}

// IsDisabled check if a line starts with a # comment marker.
func IsDisabled(line string) bool {
	return disableRe.MatchString(line)
}

// EnableLine removes the # comment marker of the line.
func EnableLine(line string) string {
	return disableRe.ReplaceAllString(line, "")
}

func banner() string {
	return `
##################################################################
# Content under this line is handled by hostctl. DO NOT EDIT. 
##################################################################`
}

func containsBanner(str string) bool {
	b := banner()
	m, _ := regexp.MatchString(b, str)
	return m
}

// WriteToFile write hosts content to file
func WriteToFile(f *os.File, h *hostFile) error {
	err := f.Truncate(0)
	if err != nil {
		return err
	}

	for _, l := range h.profiles["default"] {
		_, err = f.WriteString(l + "\n")
		if err != nil {
			return err
		}
	}

	content := strings.Join(h.profiles["default"], "\n")
	if !containsBanner(content) {
		f.WriteString("\n" + banner() + "\n")
	}

	for n, p := range h.profiles {
		if n == "default" {
			continue
		}
		err = addProfile(f, n, p)
		if err != nil {
			return err
		}
	}
	return nil
}

func addProfile(f *os.File, profile string, hl hostLines) error {
	_, err := f.WriteString(fmt.Sprintf("# profile %s\n", profile))
	if err != nil {
		return err
	}
	for _, l := range hl {
		_, err = f.WriteString(l + "\n")
		if err != nil {
			return err
		}
	}
	_, err = f.WriteString("# end\n")
	return err
}

func getHostData(dst, profile string) (*hostFile, error) {
	h := &hostFile{
		profiles: profileMap{},
	}

	if dst == "" {
		return h, MissingDestError
	}

	h, err := ReadHostFile(dst)
	if err != nil {
		return h, err
	}
	_, ok := h.profiles[profile]
	if profile != "" && !ok {
		return h, fmt.Errorf("profile '%s' doesn't exists in file", profile)
	}

	return h, nil
}

func writeHostData(dst string, h *hostFile) error {
	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	return WriteToFile(dstFile, h)
}
