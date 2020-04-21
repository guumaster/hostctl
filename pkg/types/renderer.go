package types

// DefaultColumns is the list of default columns to use when showing table list
var DefaultColumns = []string{"profile", "status", "ip", "domain"}

// ProfilesOnlyColumns are the columns used for profile status list
var ProfilesOnlyColumns = []string{"profile", "status"}

// Renderer is the interface to render hosts file content
type Renderer interface {
	AppendRow(row *Row)
	AddSeparator()
	Render() error
}

// Row represents a line for all output types
type Row struct {
	Comment string
	Profile string
	Status  string
	IP      string
	Host    string
}

type Meta struct {
	Rows int
	Raw  bool
}
