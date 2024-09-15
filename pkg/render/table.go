package render

import (
	"io"
	"os"
	"strings"

	"github.com/guumaster/tablewriter"

	"github.com/guumaster/hostctl/pkg/types"
)

// TableRendererOptions contains options to render a table.
type TableRendererOptions struct {
	Writer  io.Writer
	Columns []string
}

// RendererType represents all the existing renderers.
type RendererType string

//nolint:gochecknoglobals
var (
	Markdown RendererType = "markdown"
	Table    RendererType = "table"
	Raw      RendererType = "raw"
	JSON     RendererType = "json"
)

type meta struct {
	Rows int
	Raw  bool
}

// TableRenderer is the Renderer used to output tables.
type TableRenderer struct {
	Type    RendererType
	Columns []string
	table   *tablewriter.Table
	opts    *TableRendererOptions
	meta    *meta
}

func createTableWriter(opts *TableRendererOptions) *tablewriter.Table {
	if len(opts.Columns) == 0 {
		opts.Columns = types.DefaultColumns
	}

	out := opts.Writer
	if out == nil {
		out = os.Stdout
	}

	table := tablewriter.NewWriter(out)
	table.SetHeader(opts.Columns)

	return table
}

// NewTableRenderer creates an instance of TableRenderer.
func NewTableRenderer(opts *TableRendererOptions) TableRenderer {
	table := createTableWriter(opts)

	return TableRenderer{
		Type:    Table,
		Columns: opts.Columns,
		table:   table,
		opts:    opts,
		meta: &meta{
			Rows: 0,
		},
	}
}

// AppendRow adds a new row to the list.
func (t TableRenderer) AppendRow(row *types.Row) {
	r := []string{}

	if row.Comment != "" {
		return
	}

	for _, c := range t.Columns {
		switch c {
		case "profile":
			r = append(r, row.Profile)
		case "status":
			r = append(r, row.Status)
		case "ip", "ips":
			r = append(r, row.IP)
		case "domain", "domains":
			hosts := strings.Join(row.Hosts, " ")
			r = append(r, hosts)
		}
	}

	if len(r) > 0 {
		t.meta.Rows++
		t.table.Append(r)
	}
}

// AddSeparator adds a separator line to the list.
func (t TableRenderer) AddSeparator() {
	if !t.meta.Raw && t.meta.Rows > 0 {
		t.table.AddSeparator()
	}
}

// Render prints a table representation of row content.
func (t TableRenderer) Render() error { //nolint: unparam
	if t.meta.Rows > 0 {
		t.table.Render()
	}

	return nil
}
