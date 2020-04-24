package actions

import (
	"os"

	"github.com/spf13/cobra"
)

func newCompletionCmd(rootCmd *cobra.Command) *cobra.Command {
	var completionCmd = &cobra.Command{
		Use:    "completion <bash|zsh>",
		Short:  "Generate bash or zsh completion script",
		Hidden: true,
	}

	bashCompletionCmd := newBashCompletionCmd(rootCmd)
	zshCompletionCmd := newZshCompletionCmd(rootCmd)
	completionCmd.AddCommand(bashCompletionCmd, zshCompletionCmd)

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
