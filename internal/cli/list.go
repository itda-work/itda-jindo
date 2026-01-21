package cli

import (
	"encoding/json"
	"fmt"

	"github.com/itda-work/jindo/internal/agent"
	"github.com/itda-work/jindo/internal/command"
	"github.com/itda-work/jindo/internal/skill"
	"github.com/spf13/cobra"
)

var listJSON bool

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "ls"},
	Short:   "List all skills, agents, and commands",
	Long:    `List all configured skills, agents, and commands from ~/.claude/ and .claude/ directories.`,
	RunE:    runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVar(&listJSON, "json", false, "Output in JSON format")
}

type listItem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type scopedListOutput struct {
	Skills   []listItem `json:"skills"`
	Agents   []listItem `json:"agents"`
	Commands []listItem `json:"commands"`
}

type listOutput struct {
	Global scopedListOutput `json:"global"`
	Local  scopedListOutput `json:"local,omitempty"`
}

func runList(cmd *cobra.Command, _ []string) error {
	cmd.SilenceUsage = true

	// Get global items
	globalSkillStore := skill.NewStore(GetGlobalPath("skills"))
	globalAgentStore := agent.NewStore(GetGlobalPath("agents"))
	globalCommandStore := command.NewStore(GetGlobalPath("commands"))

	globalSkills, _ := globalSkillStore.List()
	globalAgents, _ := globalAgentStore.List()
	globalCommands, _ := globalCommandStore.List()

	// Get local items (if .claude exists)
	var localSkills []*skill.Skill
	var localAgents []*agent.Agent
	var localCommands []*command.Command

	if localPath := GetLocalPath("skills"); localPath != "" {
		localSkillStore := skill.NewStore(localPath)
		localSkills, _ = localSkillStore.List()
	}
	if localPath := GetLocalPath("agents"); localPath != "" {
		localAgentStore := agent.NewStore(localPath)
		localAgents, _ = localAgentStore.List()
	}
	if localPath := GetLocalPath("commands"); localPath != "" {
		localCommandStore := command.NewStore(localPath)
		localCommands, _ = localCommandStore.List()
	}

	hasLocal := len(localSkills) > 0 || len(localAgents) > 0 || len(localCommands) > 0

	if listJSON {
		return printListJSON(globalSkills, globalAgents, globalCommands, localSkills, localAgents, localCommands)
	}

	// Print Global section
	fmt.Println("=== Global (~/.claude/) ===")
	fmt.Println()

	fmt.Println("Skills:")
	if len(globalSkills) == 0 {
		fmt.Println("  No skills found.")
	} else {
		printSkillsTable(globalSkills)
	}
	fmt.Println()

	fmt.Println("Agents:")
	if len(globalAgents) == 0 {
		fmt.Println("  No agents found.")
	} else {
		printAgentsTable(globalAgents)
	}
	fmt.Println()

	fmt.Println("Commands:")
	if len(globalCommands) == 0 {
		fmt.Println("  No commands found.")
	} else {
		printCommandsTable(globalCommands)
	}

	// Print Local section only if has items
	if hasLocal {
		fmt.Println()
		fmt.Println("=== Local (.claude/) ===")
		fmt.Println()

		if len(localSkills) > 0 {
			fmt.Println("Skills:")
			printSkillsTable(localSkills)
			fmt.Println()
		}

		if len(localAgents) > 0 {
			fmt.Println("Agents:")
			printAgentsTable(localAgents)
			fmt.Println()
		}

		if len(localCommands) > 0 {
			fmt.Println("Commands:")
			printCommandsTable(localCommands)
		}
	}

	return nil
}

func printListJSON(globalSkills []*skill.Skill, globalAgents []*agent.Agent, globalCommands []*command.Command,
	localSkills []*skill.Skill, localAgents []*agent.Agent, localCommands []*command.Command) error {

	toListItems := func(skills []*skill.Skill, agents []*agent.Agent, commands []*command.Command) scopedListOutput {
		output := scopedListOutput{
			Skills:   make([]listItem, 0, len(skills)),
			Agents:   make([]listItem, 0, len(agents)),
			Commands: make([]listItem, 0, len(commands)),
		}
		for _, s := range skills {
			output.Skills = append(output.Skills, listItem{Name: s.Name, Description: s.Description})
		}
		for _, a := range agents {
			output.Agents = append(output.Agents, listItem{Name: a.Name, Description: a.Description})
		}
		for _, c := range commands {
			output.Commands = append(output.Commands, listItem{Name: c.Name, Description: c.Description})
		}
		return output
	}

	output := listOutput{
		Global: toListItems(globalSkills, globalAgents, globalCommands),
		Local:  toListItems(localSkills, localAgents, localCommands),
	}

	jsonOutput, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(jsonOutput))
	return nil
}
