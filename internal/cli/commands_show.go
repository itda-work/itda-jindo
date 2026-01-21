package cli

import (
	"fmt"
	"os"

	"github.com/itda-work/jindo/internal/command"
	"github.com/spf13/cobra"
)

var (
	commandsShowBrief  bool
	commandsShowGlobal bool
	commandsShowLocal  bool
)

var commandsShowCmd = &cobra.Command{
	Use:     "show <command-name>",
	Aliases: []string{"s"},
	Short:   "Show command details",
	Long: `Show the full content of a specific command from ~/.claude/commands/ (global) or .claude/commands/ (local) directory.

Use --local to show from the current directory's .claude/commands/.`,
	Args: cobra.ExactArgs(1),
	RunE: runCommandsShow,
}

func init() {
	commandsCmd.AddCommand(commandsShowCmd)
	commandsShowCmd.Flags().BoolVar(&commandsShowBrief, "brief", false, "Show only metadata (name, description)")
	commandsShowCmd.Flags().BoolVarP(&commandsShowGlobal, "global", "g", false, "Show from global ~/.claude/commands/ (default)")
	commandsShowCmd.Flags().BoolVarP(&commandsShowLocal, "local", "l", false, "Show from local .claude/commands/")
}

func runCommandsShow(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true
	name := args[0]

	// Determine scope (default: global)
	scope := ScopeGlobal
	if commandsShowLocal {
		scope = ScopeLocal
	}

	store := command.NewStore(GetPathByScope(scope, "commands"))

	if commandsShowBrief {
		return showCommandBrief(store, name)
	}

	return showCommandFull(store, name)
}

func showCommandBrief(store *command.Store, name string) error {
	cmd, err := store.Get(name)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("command not found: %s", name)
		}
		return fmt.Errorf("failed to get command: %w", err)
	}

	fmt.Printf("Name:        %s\n", cmd.Name)
	fmt.Printf("Description: %s\n", cmd.Description)
	fmt.Printf("Path:        %s\n", cmd.Path)

	return nil
}

func showCommandFull(store *command.Store, name string) error {
	content, err := store.GetContent(name)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("command not found: %s", name)
		}
		return fmt.Errorf("failed to get command content: %w", err)
	}

	fmt.Print(content)
	return nil
}
