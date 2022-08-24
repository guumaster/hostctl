package actions

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"os"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/guumaster/hostctl/pkg/render"
	"github.com/guumaster/hostctl/pkg/types"
)

func commonCheckProfileOnly(_ *cobra.Command, args []string) error {
	if len(args) == 0 {
		return types.ErrMissingProfile
	}

	if err := containsDefault(args); err != nil {
		return err
	}

	return nil
}

func commonCheckArgsWithAll(cmd *cobra.Command, args []string) error {
	all, _ := cmd.Flags().GetBool("all")
	if all && len(args) > 0 {
		return ErrIncompatibleAllFlag
	}

	if !all && len(args) == 0 {
		return types.ErrMissingProfile
	}

	if err := containsDefault(args); err != nil {
		return err
	}

	return nil
}

func commonCheckArgs(_ *cobra.Command, args []string) error {
	if len(args) == 0 {
		return types.ErrMissingProfile
	} else if len(args) > 1 {
		return ErrMultipleProfiles
	}

	if err := containsDefault(args); err != nil {
		return err
	}

	return nil
}

// isPiped detect if there is any input through STDIN.
func isPiped() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	notPipe := info.Mode()&os.ModeNamedPipe == 0

	return !notPipe || info.Size() > 0
}

func containsDefault(args []string) error {
	for _, p := range args {
		if p == types.Default {
			return types.ErrDefaultProfile
		}
	}

	return nil
}

func getDefaultHostFile() string {
	if runtime.GOOS == "linux" {
		return "/etc/hosts" //nolint: goconst
	}

	envHostFile := os.Getenv("HOSTCTL_FILE")
	if envHostFile != "" {
		return envHostFile
	}

	if runtime.GOOS == "windows" {
		return `C:/Windows/System32/Drivers/etc/hosts`
	}

	return "/etc/hosts"
}

// isValidURL tests a string to determine if it is a well-structured url or not.
func isValidURL(s string) bool {
	_, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	}

	u, err := url.Parse(s)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

func readerFromURL(ctx context.Context, url string) (io.Reader, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)

	return bytes.NewReader(b), err
}

func getRenderer(cmd *cobra.Command, opts *render.TableRendererOptions) types.Renderer {
	raw, _ := cmd.Flags().GetBool("raw")
	out, _ := cmd.Flags().GetString("out")
	cols, _ := cmd.Flags().GetStringSlice("column")

	if opts == nil {
		opts = &render.TableRendererOptions{}
	}

	if len(opts.Columns) == 0 {
		opts.Columns = cols
	}

	if opts.Writer == nil {
		opts.Writer = cmd.OutOrStdout()
	}

	//nolint: goconst
	switch {
	case raw || out == "raw":
		return render.NewRawRenderer(opts)

	case out == "md" || out == "markdown":
		return render.NewMarkdownRenderer(opts)

	case out == "json":
		return render.NewJSONRenderer(&render.JSONRendererOptions{
			Writer:  cmd.OutOrStdout(),
			Columns: cols,
		})

	default:
		return render.NewTableRenderer(opts)
	}
}

func isHelperCmd(cmd *cobra.Command) bool {
	executed := cmd.Name()
	helpers := []string{"info", "help", "completion", "gen-md-docs"}

	for _, name := range helpers {
		if executed == name {
			return true
		}
	}

	return false
}
