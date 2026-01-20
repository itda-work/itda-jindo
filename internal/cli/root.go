package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "jd",
	Short: "Claude Code configuration manager",
	Version: Version,
	Long: `jd is a CLI tool for managing Claude Code configurations
including skills, commands, and agents.

Subcommand aliases: skills(s), commands(c), agents(a)

Use 'jd --help' for all available commands.`,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}
