package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/guumaster/hostctl/pkg/host"
)

func postRunListOnly(cmd *cobra.Command, args []string) error {
	return postActionCmd(cmd, args, nil, true)
}

func commonCheckProfileOnly(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return host.MissingProfileError
	}
	if err := containsDefault(args); err != nil {
		return err
	}
	return nil
}

func commonCheckArgsWithAll(cmd *cobra.Command, args []string) error {
	all, _ := cmd.Flags().GetBool("all")
	if all && len(args) > 0 {
		return fmt.Errorf("args must be empty with --all flag")
	}
	if !all && len(args) == 0 {
		return host.MissingProfileError
	}
	if err := containsDefault(args); err != nil {
		return err
	}
	return nil
}

func commonCheckArgs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return host.MissingProfileError
	} else if len(args) > 1 {
		return fmt.Errorf("specify only one profile")
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
		if p == "default" {
			return host.DefaultProfileError
		}
	}
	return nil
}

func getDefaultHostFile() string {
	envHostFile := os.Getenv("HOSTCTL_FILE")
	if envHostFile != "" {
		return envHostFile
	}

	if runtime.GOOS == "windows" {
		return `C:/Windows/System32/Drivers/etc/hosts`
	}

	return "/etc/hosts"
}
