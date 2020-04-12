package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var postActionCmd = func(cmd *cobra.Command, args []string, postCmd *cobra.Command, list bool) error {
	quiet, _ := cmd.Flags().GetBool("quiet")
	duration, _ := cmd.Flags().GetDuration("wait")

	if !quiet && list {
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
		p := strings.Join(args, ", ")
		if duration == 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "\nWaiting until ctrl+c to %s from profile '%s'\n\n", action, p)
		} else if duration > 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "\nWaiting for %s or ctrl+c to %s from profile '%s'\n\n", duration, action, p)
		}
	}

	if duration >= 0 {
		doneCh := waitSignalOrDuration(duration)
		<-doneCh

		// Add new line to separate from "^C" output
		fmt.Fprintln(cmd.OutOrStdout())

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
