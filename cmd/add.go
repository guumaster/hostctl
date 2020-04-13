package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/guumaster/hostctl/pkg/host"
)

// addCmd represents the fromFile command
var addCmd = &cobra.Command{
	Use:     "add-to [profile] [flags]",
	Aliases: []string{"add"},
	Short:   "Add content to a profile in your hosts file.",
	Long: `
Reads from a file and set content to a profile in your hosts file.
If the profile already exists it will be added to it.`,
	Args: commonCheckArgs,
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

		err = h.AddProfile(*p)
		if err != nil {
			return err
		}

		return h.Flush()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		return postActionCmd(cmd, args, removeCmd, true)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.AddCommand(addDomainsCmd)

	addCmd.Flags().StringP("from", "f", "", "file to read")
	addCmd.PersistentFlags().DurationP("wait", "w", -1, "Enables a profile for a specific amount of time")

	addDomainsCmd.Flags().String("ip", "127.0.0.1", "domains ip")
}
