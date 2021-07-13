package actions

import (
	"context"
	"io"
	"os"

	"github.com/guumaster/hostctl/pkg/file"
	"github.com/guumaster/hostctl/pkg/parser"
	"github.com/guumaster/hostctl/pkg/types"
	"github.com/spf13/cobra"
)

type addRemoveFn func(h *file.File, p *types.Profile) error

func newAddRemoveCmd() (*cobra.Command, *cobra.Command) {
	addCmd := newAddCmd()
	removeCmd := newRemoveCmd()

	addCmd.PostRunE = func(cmd *cobra.Command, args []string) error {
		return postActionCmd(cmd, args, removeCmd, true)
	}

	return addCmd, removeCmd
}

func newAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add [profiles] [flags]",
		Short: "Add content to a profile in your hosts file.",
		Long: `
Reads from a file and set content to a profile in your hosts file.
If the profile already exists it will be added to it.`,
		Args: commonCheckArgs,
		RunE: makeAddReplace(func(h *file.File, p *types.Profile) error {
			return h.AddProfile(p)
		}),
	}
}

func newReplaceCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "replace [profile] [flags]",
		Short: "Replace content to a profile in your hosts file.",
		Long: `
Reads from a file and set content to a profile in your hosts file.
If the profile already exists it will be overwritten.
`,
		Args: commonCheckArgs,
		RunE: makeAddReplace(func(h *file.File, p *types.Profile) error {
			return h.ReplaceProfile(p)
		}),
		PostRunE: postRunListOnly,
	}
}

func makeAddReplace(actionFn addRemoveFn) func(cmd *cobra.Command, profiles []string) error {
	return func(cmd *cobra.Command, profiles []string) error {
		src, _ := cmd.Flags().GetString("host-file")
		from, _ := cmd.Flags().GetString("from")
		in := cmd.InOrStdin()

		p, err := getProfileFromInput(in, from)
		if err != nil {
			return err
		}

		h, err := file.NewFile(src)
		if err != nil {
			return err
		}

		p.Name = profiles[0]
		p.Status = types.Enabled

		err = actionFn(h, p)
		if err != nil {
			return err
		}

		return h.Flush()
	}
}

func getProfileFromInput(in io.Reader, from string) (*types.Profile, error) {
	var (
		r   io.Reader
		err error
	)

	switch {
	case isPiped() || in != os.Stdin:
		r = in

	case isValidURL(from):
		ctx := context.Background()
		r, err = readerFromURL(ctx, from)

	default:
		r, err = os.Open(from)
	}

	if err != nil {
		return nil, err
	}

	return parser.ParseProfile(r)
}
