package cmd

import (
	"os"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/guumaster/hostctl/pkg/host"
)

func commonCheckProfileOnly(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return host.ErrMissingProfile
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
		return host.ErrMissingProfile
	}

	if err := containsDefault(args); err != nil {
		return err
	}

	return nil
}

func commonCheckArgs(_ *cobra.Command, args []string) error {
	if len(args) == 0 {
		return host.ErrMissingProfile
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
		if p == "default" {
			return host.ErrDefaultProfileError
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

	if from != "" || src != defaultSrc {
		return host.ErrSnapConfinement
	}

	return nil
}
