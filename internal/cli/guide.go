package cli

import (
	"github.com/spf13/cobra"
)

var guideCmd = &cobra.Command{
	Use:     "guide",
	Aliases: []string{"g"},
	Short:   "AI-powered usage guide for skills, hooks, agents, and commands",
	Long: `Get AI-powered usage guides for Claude Code configurations.

The guide command analyzes your skills, hooks, agents, or commands and provides:
- When it gets triggered (trigger conditions)
- How to use it effectively (usage scenarios)
- Practical examples
- Customization suggestions and improvements

By default, it provides a one-shot explanation. Use -i for interactive mode
where AI asks about your context and provides personalized guidance.`,
	Example: `  # Get usage guide for a skill
  jd guide skills my-skill

  # Interactive mode - AI asks questions for personalized guidance
  jd guide skills my-skill -i

  # Get usage guide for a hook
  jd guide hooks pre-commit`,
}

func init() {
	rootCmd.AddCommand(guideCmd)
}
