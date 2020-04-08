package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/guumaster/hostctl/pkg/host"
)

var postActionCmd = func(cmd *cobra.Command, args []string, postCmd *cobra.Command) error {
	src, _ := cmd.Flags().GetString("host-file")
	quiet, _ := cmd.Flags().GetBool("quiet")
	duration, _ := cmd.Flags().GetDuration("wait")
	profile, _ := cmd.Flags().GetString("profile")
	raw, _ := cmd.Flags().GetBool("raw")
	cols, _ := cmd.Flags().GetStringSlice("column")

	var err error
	if !quiet {
		err = host.ListProfiles(src, &host.ListOptions{
			Profile:  profile,
			RawTable: raw,
			Columns:  cols,
		})
		if err != nil {
			return err
		}
	}

	action := postCmd.Name()
	if action == "domains" {
		action = "remove domains"
	}

	if duration >= 0 && !quiet {
		fmt.Printf("\nWaiting for %s or ctrl+c to %s from profile '%s'\n\n", duration, action, profile)
	}

	if duration >= 0 {
		doneCh := waitSignalOrDuration(duration)
		<-doneCh

		return postCmd.RunE(cmd, args)
	}
	return nil
}

func waitSignalOrDuration(d time.Duration) <-chan struct{} {
	done := make(chan struct{}, 0)
	sig := make(chan os.Signal)

	if d < 0 {
		d = -d
	}

	if d == 0 {
		// NOTE: It's large enough, practically it will never timeout.
		d = 999999 * time.Hour
	}

	signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)

	go func() {
		for {
			select {
			case <-time.After(d):
				done <- struct{}{}
				return
			case <-sig:
				done <- struct{}{}
				return
			}
		}
	}()

	return done
}
