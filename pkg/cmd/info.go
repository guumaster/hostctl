package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:    "info",
		Hidden: true,
		Run: func(cmd *cobra.Command, profiles []string) {
			_, _ = fmt.Fprint(cmd.OutOrStdout(), `
        __                       __             __     __
       / /_    ____     _____   / /_   _____   / /_   / /
      / __ \  / __ \   / ___/  / __/  / ___/  / __/  / /
     / / / / / /_/ /  (__  )  / /_   / /__   / /_   / /
    /_/ /_/  \____/  /____/   \__/   \___/   \__/  /_/  Version: `+version+`

     Your dev tool to manage /etc/hosts like a pro
`)
		},
	}
}
