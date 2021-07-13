package render

import (
	"github.com/guumaster/tablewriter"
)

// NewRawRenderer creates an instance of TableRenderer without borders.
func NewRawRenderer(opts *TableRendererOptions) TableRenderer {
	table := createTableWriter(opts)

	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("\t")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)

	return TableRenderer{
		Columns: opts.Columns,
		table:   table,
		opts:    opts,
		meta: &meta{
			Rows: 0,
			Raw:  true,
		},
	}
}
