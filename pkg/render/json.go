package render

import (
	"encoding/json"
	"io"

	"github.com/guumaster/hostctl/pkg/types"
)

// JSONRendererOptions contains options to render JSON content.
type JSONRendererOptions struct {
	Writer      io.Writer
	Columns     []string
	OnlyEnabled bool
}

// JSONRenderer is the Renderer used to output JSON.
type JSONRenderer struct {
	Type    RendererType
	Columns []string
	w       io.Writer
	data    *data
}

type data struct {
	lines []line
}

// NewJSONRenderer creates an instance of JSONRenderer.
func NewJSONRenderer(opts *JSONRendererOptions) JSONRenderer {
	if len(opts.Columns) == 0 {
		opts.Columns = types.DefaultColumns
	}

	return JSONRenderer{
		Type:    JSON,
		Columns: opts.Columns,
		w:       opts.Writer,
		data:    &data{},
	}
}

// AddSeparator not used on JSONRenderer.
func (j JSONRenderer) AddSeparator() {
	// not used
}

type line struct {
	Profile string
	Status  string
	IP      string
	Hosts   []string
}

// AppendRow adds a new row to the list.
func (j JSONRenderer) AppendRow(row *types.Row) {
	if row.Comment != "" {
		return
	}

	l := line{
		Profile: row.Profile,
		Status:  row.Status,
		IP:      row.IP,
		Hosts:   row.Hosts,
	}
	j.data.lines = append(j.data.lines, l)
}

// Render returns a JSON representation of the list content.
func (j JSONRenderer) Render() error {
	enc := json.NewEncoder(j.w)

	return enc.Encode(j.data.lines)
}
