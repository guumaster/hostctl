package cmd

import (
	"github.com/spf13/cobra"

	"github.com/guumaster/hostctl/pkg/host"
)

type enableDisableFn func(h *host.File, profiles []string, only, all bool) error

var enableActionFn = func(h *host.File, profiles []string, only, all bool) error {
	switch {
	case only:
		return h.EnableOnly(profiles)
	case all:
		return h.EnableAll()
	default:
		return h.Enable(profiles)
	}
}

var disableActionFn = func(h *host.File, profiles []string, only, all bool) error {
	switch {
	case only:
		return h.DisableOnly(profiles)
	case all:
		return h.DisableAll()
	default:
		return h.Disable(profiles)
	}
}

func newEnableDisableCmd() (*cobra.Command, *cobra.Command) {
	enableCmd := &cobra.Command{
		Use:   "enable [profiles] [flags]",
		Short: "Enable a profile on your hosts file.",
		Long: `
Enables an existing profile.
It will be listed as "on" while it is enabled.
`,
		Args: commonCheckArgsWithAll,
		RunE: makeEnableDisable(enableActionFn),
	}

	disableCmd := &cobra.Command{
		Use:   "disable [profiles] [flags]",
		Short: "Disable a profile from your hosts file.",
		Long: `
Disable a profile from your hosts file without removing it.
It will be listed as "off" while it is disabled.
`,
		Args: commonCheckArgsWithAll,
		RunE: makeEnableDisable(disableActionFn),
	}

	enableCmd.PostRunE = func(cmd *cobra.Command, args []string) error {
		return postActionCmd(cmd, args, disableCmd, true)
	}

	disableCmd.PostRunE = func(cmd *cobra.Command, args []string) error {
		return postActionCmd(cmd, args, enableCmd, true)
	}

	return enableCmd, disableCmd
}

func makeEnableDisable(actionFn enableDisableFn) func(cmd *cobra.Command, profiles []string) error {
	return func(cmd *cobra.Command, profiles []string) error {
		src, _ := cmd.Flags().GetString("host-file")
		only, _ := cmd.Flags().GetBool("only")
		all, _ := cmd.Flags().GetBool("all")

		h, err := host.NewFile(src)
		if err != nil {
			return err
		}

		err = actionFn(h, profiles, only, all)
		if err != nil {
			return err
		}

		return h.Flush()
	}
}
