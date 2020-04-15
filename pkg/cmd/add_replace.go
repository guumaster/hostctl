package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/guumaster/hostctl/pkg/host"
)

// addCmd represents the fromFile command
var addCmd = &cobra.Command{
	Use:   "add [profiles] [flags]",
	Short: "Add content to a profile in your hosts file.",
	Long: `
Reads from a file and set content to a profile in your hosts file.
If the profile already exists it will be added to it.`,
	Args: commonCheckArgs,
	RunE: makeAddReplace("add"),
	PostRunE: func(cmd *cobra.Command, args []string) error {
		return postActionCmd(cmd, args, removeCmd, true)
	},
}

// replaceCmd represents the setFromFile command
var replaceCmd = &cobra.Command{
	Use:   "replace [profile] [flags]",
	Short: "Replace content to a profile in your hosts file.",
	Long: `
Reads from a file and set content to a profile in your hosts file.
If the profile already exists it will be overwritten.
`,
	Args:     commonCheckArgs,
	RunE:     makeAddReplace("replace"),
	PostRunE: postRunListOnly,
}

func init() {

	addCmd.Flags().StringP("from", "f", "", "file to read")
	addCmd.PersistentFlags().DurationP("wait", "w", -1, "Enables a profile for a specific amount of time")
	addCmd.PersistentFlags().BoolP("uniq", "u", false, "only keep uniq domains per IP")

	rootCmd.AddCommand(addCmd)

	replaceCmd.Flags().StringP("from", "f", "", "file to read")
	replaceCmd.Flags().BoolP("uniq", "u", false, "only keep uniq domains per IP")
	rootCmd.AddCommand(replaceCmd)

	addDomainsCmd.Flags().String("ip", "127.0.0.1", "domains ip")
	addCmd.AddCommand(addDomainsCmd)

}

func makeAddReplace(action string) func(cmd *cobra.Command, profiles []string) error {
	actionFn := func(h *host.File, p *host.Profile) error {
		return h.AddProfile(*p)
	}
	if action == "replace" {
		actionFn = func(h *host.File, p *host.Profile) error {
			return h.ReplaceProfile(*p)
		}
	}

	return func(cmd *cobra.Command, profiles []string) error {
		src, _ := cmd.Flags().GetString("host-file")
		from, _ := cmd.Flags().GetString("from")
		uniq, _ := cmd.Flags().GetBool("uniq")

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

		p, err := host.NewProfileFromReader(r, uniq)
		if err != nil {
			return err
		}

		h, err := host.NewFile(src)
		if err != nil {
			return err
		}

		p.Name = profiles[0]
		p.Status = host.Enabled

		err = actionFn(h, p)
		if err != nil {
			return err
		}

		return h.Flush()
	}
}
