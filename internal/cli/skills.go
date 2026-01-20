package cli

import (
	"github.com/spf13/cobra"
)

var skillsCmd = &cobra.Command{
	Use:   "skills",
	Short: "Manage Claude Code skills",
	Long:  `Manage Claude Code skills in ~/.claude/skills/ directory.`,
}

func init() {
	rootCmd.AddCommand(skillsCmd)
}
