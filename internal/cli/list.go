package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/itda-work/jindo/internal/agent"
	"github.com/itda-work/jindo/internal/command"
	"github.com/itda-work/jindo/internal/skill"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var listJSON bool

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "ls"},
	Short:   "List all skills, agents, and commands",
	Long:    `List all configured skills, agents, and commands from ~/.claude/ directory.`,
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

type listOutput struct {
	Skills   []listItem `json:"skills"`
	Agents   []listItem `json:"agents"`
	Commands []listItem `json:"commands"`
}

func runList(_ *cobra.Command, _ []string) error {
	skillStore := skill.NewStore("~/.claude/skills")
	agentStore := agent.NewStore("~/.claude/agents")
	commandStore := command.NewStore("~/.claude/commands")

	skills, err := skillStore.List()
	if err != nil {
		skills = nil
	}

	agents, err := agentStore.List()
	if err != nil {
		agents = nil
	}

	commands, err := commandStore.List()
	if err != nil {
		commands = nil
	}

	if listJSON {
		return printListJSON(skills, agents, commands)
	}

	printListGrouped(skills, agents, commands)
	return nil
}

func printListJSON(skills []*skill.Skill, agents []*agent.Agent, commands []*command.Command) error {
	output := listOutput{
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

	jsonOutput, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(jsonOutput))
	return nil
}

func printListGrouped(skills []*skill.Skill, agents []*agent.Agent, commands []*command.Command) {
	total := 0

	// Get terminal width
	termWidth := 80
	if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && w > 0 {
		termWidth = w
	}

	// Calculate max name width for alignment
	maxNameWidth := 0
	for _, s := range skills {
		if len(s.Name) > maxNameWidth {
			maxNameWidth = len(s.Name)
		}
	}
	for _, a := range agents {
		if len(a.Name) > maxNameWidth {
			maxNameWidth = len(a.Name)
		}
	}
	for _, c := range commands {
		if len(c.Name) > maxNameWidth {
			maxNameWidth = len(c.Name)
		}
	}
	if maxNameWidth > 30 {
		maxNameWidth = 30
	}

	// indent(2) + name + gap(2) + desc
	indent := 2
	gap := 2
	descWidth := termWidth - indent - maxNameWidth - gap
	if descWidth < 20 {
		descWidth = 20
	}

	// Skills
	fmt.Println("Skills:")
	if len(skills) == 0 {
		fmt.Println("  (none)")
	} else {
		for _, s := range skills {
			printItem(s.Name, firstSentence(s.Description), maxNameWidth, descWidth, indent)
		}
		total += len(skills)
	}
	fmt.Println()

	// Agents
	fmt.Println("Agents:")
	if len(agents) == 0 {
		fmt.Println("  (none)")
	} else {
		for _, a := range agents {
			printItem(a.Name, firstSentence(a.Description), maxNameWidth, descWidth, indent)
		}
		total += len(agents)
	}
	fmt.Println()

	// Commands
	fmt.Println("Commands:")
	if len(commands) == 0 {
		fmt.Println("  (none)")
	} else {
		for _, c := range commands {
			printItem(c.Name, firstSentence(c.Description), maxNameWidth, descWidth, indent)
		}
		total += len(commands)
	}
	fmt.Println()

	fmt.Printf("Total: %d items (%d skills, %d agents, %d commands)\n",
		total, len(skills), len(agents), len(commands))
}

func printItem(name, desc string, nameWidth, descWidth, indent int) {
	lines := wrapText(desc, descWidth)
	if len(lines) == 0 {
		lines = []string{""}
	}

	// First line with name
	fmt.Printf("%*s%-*s  %s\n", indent, "", nameWidth, name, lines[0])

	// Remaining lines with indent
	padding := strings.Repeat(" ", indent+nameWidth+2)
	for i := 1; i < len(lines); i++ {
		fmt.Printf("%s%s\n", padding, lines[i])
	}
}

func wrapText(text string, width int) []string {
	if width <= 0 || len(text) == 0 {
		return []string{text}
	}

	var lines []string
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{}
	}

	currentLine := words[0]
	for _, word := range words[1:] {
		if len(currentLine)+1+len(word) <= width {
			currentLine += " " + word
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}
	lines = append(lines, currentLine)

	return lines
}

func firstSentence(s string) string {
	// Find first sentence ending with . ! or ?
	for i, r := range s {
		if r == '.' || r == '!' || r == '?' {
			// Check if next char is space or end of string (to avoid cutting "e.g." or "v1.0")
			if i+1 >= len(s) || s[i+1] == ' ' || s[i+1] == '\n' {
				return strings.TrimSpace(s[:i+1])
			}
		}
	}
	return strings.TrimSpace(s)
}
