package cli

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/itda-work/jindo/internal/skill"
	"github.com/spf13/cobra"
)

var skillsListJSON bool

var skillsListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "ls"},
	Short:   "List all skills",
	Long:    `List all skills from ~/.claude/skills/ and .claude/skills/ directories.`,
	RunE:    runSkillsList,
}

func init() {
	skillsCmd.AddCommand(skillsListCmd)
	skillsListCmd.Flags().BoolVar(&skillsListJSON, "json", false, "Output in JSON format")
}

// skillsListOutput represents JSON output for skills list with scope
type skillsListOutput struct {
	Global []*skill.Skill `json:"global"`
	Local  []*skill.Skill `json:"local,omitempty"`
}

func runSkillsList(cmd *cobra.Command, _ []string) error {
	cmd.SilenceUsage = true

	// Get global skills
	globalStore := skill.NewStore(GetGlobalPath("skills"))
	globalSkills, err := globalStore.List()
	if err != nil {
		globalSkills = nil
	}

	// Get local skills (if .claude/skills exists)
	var localSkills []*skill.Skill
	if localPath := GetLocalPath("skills"); localPath != "" {
		localStore := skill.NewStore(localPath)
		localSkills, _ = localStore.List()
	}

	if skillsListJSON {
		output := skillsListOutput{
			Global: globalSkills,
			Local:  localSkills,
		}
		jsonOutput, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(jsonOutput))
		return nil
	}

	// Print global section
	fmt.Println("=== Global (~/.claude/skills/) ===")
	if len(globalSkills) == 0 {
		fmt.Println("No skills found.")
	} else {
		printSkillsTable(globalSkills)
	}

	// Print local section only if exists and has items
	if len(localSkills) > 0 {
		fmt.Println()
		fmt.Println("=== Local (.claude/skills/) ===")
		printSkillsTable(localSkills)
	}

	return nil
}

func printSkillsJSON(skills []*skill.Skill) error {
	output, err := json.MarshalIndent(skills, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}

func printSkillsTable(skills []*skill.Skill) {
	// Calculate column widths
	idWidth := len("ID")
	descWidth := len("DESCRIPTION")
	toolsWidth := len("ALLOWED-TOOLS")

	for _, s := range skills {
		// Use directory name as skill ID (used in commands)
		skillID := filepath.Base(filepath.Dir(s.Path))
		if len(skillID) > idWidth {
			idWidth = len(skillID)
		}
		// Truncate description for display
		desc := s.Description
		if len(desc) > 50 {
			desc = desc[:47] + "..."
		}
		if len(desc) > descWidth {
			descWidth = len(desc)
		}
		tools := strings.Join(s.AllowedTools, ", ")
		if len(tools) > toolsWidth {
			toolsWidth = len(tools)
		}
	}

	// Cap widths
	if idWidth > 30 {
		idWidth = 30
	}
	if descWidth > 50 {
		descWidth = 50
	}
	if toolsWidth > 30 {
		toolsWidth = 30
	}

	// Print header
	fmt.Printf("%-*s  %-*s  %-*s\n",
		idWidth, "ID",
		descWidth, "DESCRIPTION",
		toolsWidth, "ALLOWED-TOOLS")
	fmt.Printf("%s  %s  %s\n",
		strings.Repeat("-", idWidth),
		strings.Repeat("-", descWidth),
		strings.Repeat("-", toolsWidth))

	// Print rows
	for _, s := range skills {
		// Use directory name as skill ID
		skillID := filepath.Base(filepath.Dir(s.Path))
		if len(skillID) > idWidth {
			skillID = skillID[:idWidth-3] + "..."
		}

		desc := s.Description
		if len(desc) > descWidth {
			desc = desc[:descWidth-3] + "..."
		}

		tools := strings.Join(s.AllowedTools, ", ")
		if len(tools) > toolsWidth {
			tools = tools[:toolsWidth-3] + "..."
		}

		fmt.Printf("%-*s  %-*s  %-*s\n",
			idWidth, skillID,
			descWidth, desc,
			toolsWidth, tools)
	}

	fmt.Printf("\nTotal: %d skills\n", len(skills))
}
