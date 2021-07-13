// package main contains CLI entrypoint
package main

import (
	"os"

	"github.com/guumaster/cligger"
	"github.com/guumaster/hostctl/cmd/hostctl/actions"
)

func main() {
	_, err := os.Getwd()
	if err != nil {
		cligger.Fatal("error: %w\n", err)
	}

	rootCmd := actions.NewRootCmd()

	if err := rootCmd.Execute(); err != nil {
		cligger.Fatal("error: %s\n", err)
	}
}
