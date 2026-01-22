package guide

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// GuideType represents the type of guide
type GuideType string

const (
	TypeSkill    GuideType = "skills"
	TypeHook     GuideType = "hooks"
	TypeAgent    GuideType = "agents"
	TypeCommand  GuideType = "commands"
	TypeClaudemd GuideType = "claudemd"
)

// Guide represents a cached guide
type Guide struct {
	Type      GuideType
	ID        string
	Content   string
	CreatedAt time.Time
	Path      string
}

// Store manages cached guides
type Store struct {
	baseDir string
}

// NewStore creates a new guide store
func NewStore() (*Store, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	baseDir := filepath.Join(home, ".claude", "jindo", "guides")
	return &Store{baseDir: baseDir}, nil
}

// GetDir returns the directory for a guide type
func (s *Store) GetDir(guideType GuideType) string {
	return filepath.Join(s.baseDir, string(guideType))
}

// GetPath returns the path for a specific guide
func (s *Store) GetPath(guideType GuideType, id string) string {
	// Sanitize ID for filename
	safeID := sanitizeFilename(id)
	return filepath.Join(s.GetDir(guideType), safeID+".md")
}

// Exists checks if a guide exists
func (s *Store) Exists(guideType GuideType, id string) bool {
	path := s.GetPath(guideType, id)
	_, err := os.Stat(path)
	return err == nil
}

// Get retrieves a cached guide
func (s *Store) Get(guideType GuideType, id string) (*Guide, error) {
	path := s.GetPath(guideType, id)

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	// Parse frontmatter to get created_at if present
	createdAt := info.ModTime()
	contentStr := string(content)

	if parsedTime, body, ok := parseFrontmatter(contentStr); ok {
		createdAt = parsedTime
		contentStr = body
	}

	return &Guide{
		Type:      guideType,
		ID:        id,
		Content:   contentStr,
		CreatedAt: createdAt,
		Path:      path,
	}, nil
}

// Save saves a guide to cache
func (s *Store) Save(guideType GuideType, id string, content string) (*Guide, error) {
	dir := s.GetDir(guideType)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	path := s.GetPath(guideType, id)
	now := time.Now()

	// Add frontmatter with timestamp
	fullContent := fmt.Sprintf(`---
type: %s
id: %s
created_at: %s
---

%s`, guideType, id, now.Format(time.RFC3339), content)

	if err := os.WriteFile(path, []byte(fullContent), 0644); err != nil {
		return nil, err
	}

	return &Guide{
		Type:      guideType,
		ID:        id,
		Content:   content,
		CreatedAt: now,
		Path:      path,
	}, nil
}

// Delete removes a cached guide
func (s *Store) Delete(guideType GuideType, id string) error {
	path := s.GetPath(guideType, id)
	return os.Remove(path)
}

// List lists all cached guides of a type
func (s *Store) List(guideType GuideType) ([]*Guide, error) {
	dir := s.GetDir(guideType)

	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var guides []*Guide
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		id := strings.TrimSuffix(entry.Name(), ".md")
		guide, err := s.Get(guideType, id)
		if err != nil {
			continue
		}
		guides = append(guides, guide)
	}

	return guides, nil
}

// FormatAge returns a human-readable age string
func FormatAge(t time.Time) string {
	duration := time.Since(t)

	if duration < time.Minute {
		return "방금 전"
	} else if duration < time.Hour {
		mins := int(duration.Minutes())
		return fmt.Sprintf("%d분 전", mins)
	} else if duration < 24*time.Hour {
		hours := int(duration.Hours())
		return fmt.Sprintf("%d시간 전", hours)
	} else if duration < 7*24*time.Hour {
		days := int(duration.Hours() / 24)
		return fmt.Sprintf("%d일 전", days)
	} else if duration < 30*24*time.Hour {
		weeks := int(duration.Hours() / 24 / 7)
		return fmt.Sprintf("%d주 전", weeks)
	} else {
		return t.Format("2006-01-02")
	}
}

// parseFrontmatter extracts created_at from frontmatter
func parseFrontmatter(content string) (time.Time, string, bool) {
	if !strings.HasPrefix(content, "---\n") {
		return time.Time{}, content, false
	}

	endIdx := strings.Index(content[4:], "\n---")
	if endIdx == -1 {
		return time.Time{}, content, false
	}

	frontmatter := content[4 : 4+endIdx]
	body := strings.TrimPrefix(content[4+endIdx+4:], "\n")

	// Parse created_at
	re := regexp.MustCompile(`created_at:\s*(.+)`)
	matches := re.FindStringSubmatch(frontmatter)
	if len(matches) < 2 {
		return time.Time{}, body, false
	}

	t, err := time.Parse(time.RFC3339, strings.TrimSpace(matches[1]))
	if err != nil {
		return time.Time{}, body, false
	}

	return t, body, true
}

// sanitizeFilename makes a string safe for use as a filename
func sanitizeFilename(s string) string {
	// Replace problematic characters
	s = strings.ReplaceAll(s, "/", "-")
	s = strings.ReplaceAll(s, "\\", "-")
	s = strings.ReplaceAll(s, ":", "-")
	s = strings.ReplaceAll(s, "*", "-")
	s = strings.ReplaceAll(s, "?", "-")
	s = strings.ReplaceAll(s, "\"", "-")
	s = strings.ReplaceAll(s, "<", "-")
	s = strings.ReplaceAll(s, ">", "-")
	s = strings.ReplaceAll(s, "|", "-")
	return s
}
