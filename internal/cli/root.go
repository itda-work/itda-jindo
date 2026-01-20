package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "jindo",
	Short:   "Claude Code configuration manager",
	Version: Version,
	Long: `jindo is a CLI tool for managing Claude Code configurations
including skills, commands, and agents.

Use 'jindo --help' for all available commands.`,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}
