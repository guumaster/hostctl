package actions

import (
	"os"

	"github.com/spf13/cobra"
)

func newCompletionCmd(rootCmd *cobra.Command) *cobra.Command {
	completionCmd := &cobra.Command{
		Use:    "completion <bash|zsh|fish>",
		Short:  "Generate bash zsh or fish completion script",
		Hidden: true,
	}

	bashCompletionCmd := newBashCompletionCmd(rootCmd)
	zshCompletionCmd := newZshCompletionCmd(rootCmd)
	fishCompletionCmd := newFishCompletionCmd(rootCmd)
	completionCmd.AddCommand(bashCompletionCmd, zshCompletionCmd, fishCompletionCmd)

	return completionCmd
}

func newBashCompletionCmd(rootCmd *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:   "bash",
		Short: "Generate bash completion script",
		RunE: func(cmd *cobra.Command, args []string) error {
			return rootCmd.GenBashCompletion(os.Stdout)
		},
	}
}

func newZshCompletionCmd(rootCmd *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:   "zsh",
		Short: "Generate zsh completion script",
		RunE: func(cmd *cobra.Command, args []string) error {
			return rootCmd.GenZshCompletion(os.Stdout)
		},
	}
}

func newFishCompletionCmd(rootCmd *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:   "fish",
		Short: "Generate fish completion script",
		RunE: func(cmd *cobra.Command, args []string) error {
			return rootCmd.GenFishCompletion(os.Stdout, true)
		},
	}
}
