package cli

import (
	"github.com/spf13/cobra"
)

var claudemdCmd = &cobra.Command{
	Use:     "claudemd",
	Aliases: []string{"cm"},
	Short:   "Manage CLAUDE.md configuration files",
	Long: `Manage CLAUDE.md configuration files for Claude Code.

CLAUDE.md files provide instructions and context to Claude Code sessions.
This command helps analyze and optimize these configuration files.`,
}

func init() {
	rootCmd.AddCommand(claudemdCmd)
}
