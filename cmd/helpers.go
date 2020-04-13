package cmd

import (
	"os"
	"runtime"

	"github.com/guumaster/hostctl/pkg/host"
)

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
