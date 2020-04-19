package render

import (
	"io"
	"os"

	"github.com/guumaster/tablewriter"
)

type TableRendererOptions struct {
	Writer  io.Writer
	Columns []string
}

type RendererType string

var (
	Markdown RendererType = "markdown"
	Table    RendererType = "table"
	Raw      RendererType = "raw"
	JSON     RendererType = "json"
)

type TableRenderer struct {
	Type    RendererType
	Columns []string
	table   *tablewriter.Table
	opts    *TableRendererOptions
	meta    *meta
}

func createTableWriter(opts *TableRendererOptions) *tablewriter.Table {
	if len(opts.Columns) == 0 {
		opts.Columns = DefaultColumns
	}

	out := opts.Writer
	if out == nil {
		out = os.Stdout
	}

	table := tablewriter.NewWriter(out)
	table.SetHeader(opts.Columns)

	return table
}

func NewTableRenderer(opts *TableRendererOptions) TableRenderer {
	table := createTableWriter(opts)

	return TableRenderer{
		Type:    Table,
		Columns: opts.Columns,
		table:   table,
		opts:    opts,
		meta: &meta{
			rows: 0,
		},
	}
}

func (t TableRenderer) AppendRow(row *Row) {
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
			r = append(r, row.Host)
		}
	}

	if len(r) > 0 {
		t.meta.rows++
		t.table.Append(r)
	}
}

func (t TableRenderer) AddSeparator() {
	if !t.meta.raw && t.meta.rows > 0 {
		t.table.AddSeparator()
	}
}

func (t TableRenderer) Render() error {
	if t.meta.rows > 0 {
		t.table.Render()
	}

	return nil
}
