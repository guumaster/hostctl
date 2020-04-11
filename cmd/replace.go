package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/guumaster/hostctl/pkg/host"
)

// replaceCmd represents the setFromFile command
var replaceCmd = &cobra.Command{
	Use:     "replace [profile] [domains] [flags]",
	Aliases: []string{"set"},
	Short:   "Replace content to a profile in your hosts file.",
	Long: `
Reads from a file and set content to a profile in your hosts file.
If the profile already exists it will be overwritten.
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return host.MissingProfileError
		} else if len(args) > 1 {
			return fmt.Errorf("specify only one profile")
		}
		if err := containsDefault(args); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, profiles []string) error {
		src, _ := cmd.Flags().GetString("host-file")
		from, _ := cmd.Flags().GetString("from")

		in := cmd.InOrStdin()

		var r io.Reader
		var err error
		if isPiped() || in != os.Stdin {
			r = in
		} else {
			r, err = os.Open(from)
			if err != nil {
				return err
			}
		}

		p, err := host.NewProfileFromReader(r)
		if err != nil {
			return err
		}

		h, err := host.NewFile(src)
		if err != nil {
			return err
		}

		p.Name = profiles[0]
		p.Status = host.Enabled

		err = h.ReplaceProfile(*p)
		if err != nil {
			return err
		}

		return h.WriteTo(src)
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		return postActionCmd(cmd, args, nil, true)
	},
}

func init() {
	rootCmd.AddCommand(replaceCmd)

	replaceCmd.Flags().StringP("from", "f", "", "file to read")
}
