// package main contains CLI entrypoint
package main

import (
	"log"
	"os"

	"github.com/guumaster/hostctl/pkg/cmd"
)

func main() {
	_, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	rootCmd := cmd.NewRootCmd()

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
