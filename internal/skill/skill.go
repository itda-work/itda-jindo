package skill

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// Skill represents a Claude Code skill
type Skill struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	AllowedTools []string `json:"allowed_tools"`
	Path         string   `json:"path"`
}

// ParseSkillFile parses a SKILL.md or skill.md file and returns a Skill
func ParseSkillFile(path string) (result *Skill, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	skill := &Skill{
		Path: path,
	}

	scanner := bufio.NewScanner(file)
	inFrontmatter := false
	lineCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineCount++

		// Check for frontmatter delimiter
		if strings.TrimSpace(line) == "---" {
			if !inFrontmatter {
				inFrontmatter = true
				continue
			} else {
				// End of frontmatter
				break
			}
		}

		if !inFrontmatter {
			continue
		}

		// Parse YAML-like frontmatter (simple key: value)
		if idx := strings.Index(line, ":"); idx > 0 {
			key := strings.TrimSpace(line[:idx])
			value := strings.TrimSpace(line[idx+1:])

			switch key {
			case "name":
				skill.Name = value
			case "description":
				skill.Description = value
			case "allowed-tools":
				// Parse comma-separated tools
				tools := strings.Split(value, ",")
				for _, tool := range tools {
					tool = strings.TrimSpace(tool)
					if tool != "" {
						skill.AllowedTools = append(skill.AllowedTools, tool)
					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return skill, nil
}

// Store manages skills in a directory
type Store struct {
	baseDir string
}

// NewStore creates a new skill store
func NewStore(baseDir string) *Store {
	return &Store{baseDir: baseDir}
}

// List returns all skills in the store
func (s *Store) List() ([]*Skill, error) {
	var skills []*Skill

	// Expand ~ to home directory
	dir := s.baseDir
	if strings.HasPrefix(dir, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		dir = filepath.Join(home, dir[2:])
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return skills, nil
		}
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		skillDir := filepath.Join(dir, entry.Name())

		// Try SKILL.md first, then skill.md
		skillFile := filepath.Join(skillDir, "SKILL.md")
		if _, err := os.Stat(skillFile); os.IsNotExist(err) {
			skillFile = filepath.Join(skillDir, "skill.md")
			if _, err := os.Stat(skillFile); os.IsNotExist(err) {
				continue
			}
		}

		skill, err := ParseSkillFile(skillFile)
		if err != nil {
			continue
		}

		// Use directory name if name is empty
		if skill.Name == "" {
			skill.Name = entry.Name()
		}

		skills = append(skills, skill)
	}

	return skills, nil
}
