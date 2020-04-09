package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var postActionCmd = func(cmd *cobra.Command, args []string, postCmd *cobra.Command) error {
	quiet, _ := cmd.Flags().GetBool("quiet")
	duration, _ := cmd.Flags().GetDuration("wait")
	profile, _ := cmd.Flags().GetString("profile")

	if !quiet {
		err := listCmd.RunE(cmd, args)
		if err != nil {
			return err
		}
	}

	if postCmd == nil {
		return nil
	}

	action := postCmd.Name()
	if action == "domains" {
		action = "remove domains"
	}

	if !quiet {
		if duration == 0 {
			fmt.Printf("\nWaiting until ctrl+c to %s from profile '%s'\n\n", action, profile)
		} else if duration > 0 {
			fmt.Printf("\nWaiting for %s or ctrl+c to %s from profile '%s'\n\n", duration, action, profile)
		}
	}

	if duration >= 0 {
		doneCh := waitSignalOrDuration(duration)
		<-doneCh

		// Add new line to separate from "^C" output
		fmt.Println()

		err := postCmd.RunE(cmd, args)
		if err != nil {
			return err
		}
		if quiet {
			return nil
		}
		return listCmd.RunE(cmd, args)
	}
	return nil
}

func waitSignalOrDuration(d time.Duration) <-chan struct{} {
	done := make(chan struct{})
	sig := make(chan os.Signal)

	if d < 0 {
		d = -d
	}

	if d == 0 {
		// NOTE: It's large enough, practically it will never timeout.
		d = 999999 * time.Hour
	}

	signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

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
