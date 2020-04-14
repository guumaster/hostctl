package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// completionCmd represents a completion command.
var completionCmd = &cobra.Command{
	Use:    "completion <bash|zsh>",
	Short:  "Generate bash or zsh completion script",
	Hidden: true,
}

// bashCompletionCmd represents a bash completion command.
var bashCompletionCmd = &cobra.Command{
	Use:   "bash",
	Short: "Generate bash completion script",
	RunE: func(cmd *cobra.Command, args []string) error {
		return rootCmd.GenBashCompletion(os.Stdout)
	},
}

// zshCompletionCmd represents a  zsh completion command.
var zshCompletionCmd = &cobra.Command{
	Use:   "zsh",
	Short: "Generate zsh completion script",
	RunE: func(cmd *cobra.Command, args []string) error {
		return rootCmd.GenZshCompletion(os.Stdout)
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
	completionCmd.AddCommand(bashCompletionCmd, zshCompletionCmd)
}
