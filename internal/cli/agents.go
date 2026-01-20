package cli

import (
	"github.com/spf13/cobra"
)

var agentsCmd = &cobra.Command{
	Use:     "agents",
	Aliases: []string{"a"},
	Short:   "Manage Claude Code agents",
	Long:    `Manage Claude Code agents in ~/.claude/agents/ directory.`,
}

func init() {
	rootCmd.AddCommand(agentsCmd)
}
