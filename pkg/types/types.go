package types

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

// Content contains complete data of all profiles.
type Content struct {
	DefaultProfile DefaultProfile
	ProfileNames   []string
	Profiles       map[string]*Profile
}

// Status represents the status of a Profile.
type Status string

const (
	// Enabled marks a profile active on your hosts file.
	Enabled Status = "on"
	// Disabled marks a profile not active on your hosts file.
	Disabled Status = "off"

	// Default is the name of the default profile.
	Default = "default"
)

// DefaultProfile contains data for the default profile.
type DefaultProfile []*Row

// Render writes the default profile content to the given StringWriter.
func (d DefaultProfile) Render(w io.StringWriter) error {
	tmp := bytes.NewBufferString("")

	for i, row := range d {
		line := getLine(row)
		nextLine := ""

		if i+1 < len(d) {
			nextLine = getLine(d[i+1])
		}

		// skips two consecutive empty lines
		if line == "" && nextLine == "" {
			continue
		}

		_, err := tmp.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	// Write to input writer after knowing the profile is well formed
	_, err := w.WriteString(tmp.String())

	return err
}

func getLine(row *Row) string {
	line := ""
	if row.Comment != "" {
		line = row.Comment
	} else {
		prefix := ""
		if row.Status == string(Disabled) {
			prefix = "# "
		}

		hosts := strings.Join(row.Hosts, " ")
		line = fmt.Sprintf("%s%s %s", prefix, row.IP, hosts)
	}

	return strings.TrimSpace(line)
}
