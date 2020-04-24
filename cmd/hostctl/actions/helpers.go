package actions

import (
	"bytes"
	"io"
	"io/ioutil"
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

// isPiped detect if there is any input through STDIN
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

func getDefaultHostFile(snapBuild bool) string {
	// Snap confinement doesn't allow to read other than
	if runtime.GOOS == "linux" && snapBuild {
		return "/etc/hosts"
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

func checkSnapRestrictions(cmd *cobra.Command, isSnap bool) error {
	from, _ := cmd.Flags().GetString("from")
	src, _ := cmd.Flags().GetString("host-file")

	defaultSrc := getDefaultHostFile(isSnap)

	if !isSnap {
		return nil
	}

	if from != "" || src != defaultSrc && !isValidURL(from) {
		return types.ErrSnapConfinement
	}

	return nil
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

func readerFromURL(url string) (io.Reader, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)

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
