package host

import (
	"os"

	"github.com/guumaster/tablewriter"
)

func createTableWriter(opts *ListOptions) *tablewriter.Table {
	out := opts.Writer
	if out == nil {
		out = os.Stdout
	}

	table := tablewriter.NewWriter(out)
	table.SetHeader(opts.Columns)

	if opts.RawTable {
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
	}

	return table
}

func getRow(line *tableRow, columns []string) []string {
	row := []string{}

	if line.Comment != "" {
		return row
	}

	for _, c := range columns {
		switch c {
		case "profile":
			row = append(row, line.Profile)
		case "status":
			row = append(row, line.Status)
		case "ip", "ips":
			row = append(row, line.IP)
		case "domain", "domains":
			row = append(row, line.Host)
		}
	}

	return row
}
