package repo

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

// setupTestRepo creates a temporary directory structure for testing.
func setupTestRepo(t *testing.T) string {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "repo-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	return tmpDir
}

// createFile creates a file with the given content.
func createFile(t *testing.T, path, content string) {
	t.Helper()
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("failed to create directory %s: %v", dir, err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create file %s: %v", path, err)
	}
}

// createDir creates a directory.
func createDir(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatalf("failed to create directory %s: %v", path, err)
	}
}

func TestScanSkills(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T, repoPath string)
		expected []BrowseItem
	}{
		{
			name: "empty repo",
			setup: func(t *testing.T, repoPath string) {
				// No setup needed
			},
			expected: nil,
		},
		{
			name: "skills in root skills/ directory",
			setup: func(t *testing.T, repoPath string) {
				createFile(t, filepath.Join(repoPath, "skills", "web-fetch", "SKILL.md"), "# Web Fetch")
				createFile(t, filepath.Join(repoPath, "skills", "code-review", "skill.md"), "# Code Review")
			},
			expected: []BrowseItem{
				{Name: "code-review", Path: "skills/code-review", Type: TypeSkill},
				{Name: "web-fetch", Path: "skills/web-fetch", Type: TypeSkill},
			},
		},
		{
			name: "skills in .claude/skills/ directory",
			setup: func(t *testing.T, repoPath string) {
				createFile(t, filepath.Join(repoPath, ".claude", "skills", "my-skill", "SKILL.md"), "# My Skill")
			},
			expected: []BrowseItem{
				{Name: "my-skill", Path: ".claude/skills/my-skill", Type: TypeSkill},
			},
		},
		{
			name: "skills in both directories",
			setup: func(t *testing.T, repoPath string) {
				createFile(t, filepath.Join(repoPath, "skills", "root-skill", "SKILL.md"), "# Root Skill")
				createFile(t, filepath.Join(repoPath, ".claude", "skills", "claude-skill", "SKILL.md"), "# Claude Skill")
			},
			expected: []BrowseItem{
				{Name: "root-skill", Path: "skills/root-skill", Type: TypeSkill},
				{Name: "claude-skill", Path: ".claude/skills/claude-skill", Type: TypeSkill},
			},
		},
		{
			name: "same name in both directories (both shown)",
			setup: func(t *testing.T, repoPath string) {
				createFile(t, filepath.Join(repoPath, "skills", "common", "SKILL.md"), "# Common Root")
				createFile(t, filepath.Join(repoPath, ".claude", "skills", "common", "SKILL.md"), "# Common Claude")
			},
			expected: []BrowseItem{
				{Name: "common", Path: "skills/common", Type: TypeSkill},
				{Name: "common", Path: ".claude/skills/common", Type: TypeSkill},
			},
		},
		{
			name: "directory without SKILL.md is ignored",
			setup: func(t *testing.T, repoPath string) {
				createDir(t, filepath.Join(repoPath, "skills", "no-skill-md"))
				createFile(t, filepath.Join(repoPath, "skills", "no-skill-md", "README.md"), "# README")
				createFile(t, filepath.Join(repoPath, "skills", "valid-skill", "SKILL.md"), "# Valid")
			},
			expected: []BrowseItem{
				{Name: "valid-skill", Path: "skills/valid-skill", Type: TypeSkill},
			},
		},
		{
			name: "files in skills/ directory are ignored (only directories)",
			setup: func(t *testing.T, repoPath string) {
				createFile(t, filepath.Join(repoPath, "skills", "not-a-skill.md"), "# Not a skill")
				createFile(t, filepath.Join(repoPath, "skills", "real-skill", "SKILL.md"), "# Real")
			},
			expected: []BrowseItem{
				{Name: "real-skill", Path: "skills/real-skill", Type: TypeSkill},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoPath := setupTestRepo(t)
			defer func() { _ = os.RemoveAll(repoPath) }()

			tt.setup(t, repoPath)

			store := NewStore(repoPath)
			items, err := store.scanSkills(repoPath)
			if err != nil {
				t.Fatalf("scanSkills failed: %v", err)
			}

			// Sort for consistent comparison
			sort.Slice(items, func(i, j int) bool {
				return items[i].Path < items[j].Path
			})
			sort.Slice(tt.expected, func(i, j int) bool {
				return tt.expected[i].Path < tt.expected[j].Path
			})

			if len(items) != len(tt.expected) {
				t.Errorf("got %d items, want %d", len(items), len(tt.expected))
				t.Errorf("got: %+v", items)
				t.Errorf("want: %+v", tt.expected)
				return
			}

			for i := range items {
				if items[i].Name != tt.expected[i].Name ||
					items[i].Path != tt.expected[i].Path ||
					items[i].Type != tt.expected[i].Type {
					t.Errorf("item %d: got %+v, want %+v", i, items[i], tt.expected[i])
				}
			}
		})
	}
}

func TestScanCommands(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T, repoPath string)
		expected []BrowseItem
	}{
		{
			name: "empty repo",
			setup: func(t *testing.T, repoPath string) {
				// No setup needed
			},
			expected: nil,
		},
		{
			name: "commands in root commands/ directory",
			setup: func(t *testing.T, repoPath string) {
				createFile(t, filepath.Join(repoPath, "commands", "commit.md"), "# Commit")
				createFile(t, filepath.Join(repoPath, "commands", "review.md"), "# Review")
			},
			expected: []BrowseItem{
				{Name: "commit", Path: "commands/commit.md", Type: TypeCommand},
				{Name: "review", Path: "commands/review.md", Type: TypeCommand},
			},
		},
		{
			name: "commands in .claude/commands/ directory",
			setup: func(t *testing.T, repoPath string) {
				createFile(t, filepath.Join(repoPath, ".claude", "commands", "my-cmd.md"), "# My Cmd")
			},
			expected: []BrowseItem{
				{Name: "my-cmd", Path: ".claude/commands/my-cmd.md", Type: TypeCommand},
			},
		},
		{
			name: "commands in nested subdirectory (game/)",
			setup: func(t *testing.T, repoPath string) {
				createFile(t, filepath.Join(repoPath, "commands", "game", "init.md"), "# Game Init")
				createFile(t, filepath.Join(repoPath, "commands", "game", "publish.md"), "# Game Publish")
				createFile(t, filepath.Join(repoPath, "commands", "commit.md"), "# Commit")
			},
			expected: []BrowseItem{
				{Name: "commit", Path: "commands/commit.md", Type: TypeCommand},
				{Name: "game:init", Path: "commands/game/init.md", Type: TypeCommand},
				{Name: "game:publish", Path: "commands/game/publish.md", Type: TypeCommand},
			},
		},
		{
			name: "commands in nested subdirectory within .claude/",
			setup: func(t *testing.T, repoPath string) {
				createFile(t, filepath.Join(repoPath, ".claude", "commands", "dev", "test.md"), "# Test")
				createFile(t, filepath.Join(repoPath, ".claude", "commands", "deploy.md"), "# Deploy")
			},
			expected: []BrowseItem{
				{Name: "deploy", Path: ".claude/commands/deploy.md", Type: TypeCommand},
				{Name: "dev:test", Path: ".claude/commands/dev/test.md", Type: TypeCommand},
			},
		},
		{
			name: "deeply nested commands (only one level supported)",
			setup: func(t *testing.T, repoPath string) {
				createFile(t, filepath.Join(repoPath, "commands", "a", "b", "deep.md"), "# Deep")
				createFile(t, filepath.Join(repoPath, "commands", "a", "shallow.md"), "# Shallow")
			},
			expected: []BrowseItem{
				{Name: "a:shallow", Path: "commands/a/shallow.md", Type: TypeCommand},
			},
		},
		{
			name: "non-md files are ignored",
			setup: func(t *testing.T, repoPath string) {
				createFile(t, filepath.Join(repoPath, "commands", "valid.md"), "# Valid")
				createFile(t, filepath.Join(repoPath, "commands", "invalid.txt"), "Not a command")
				createFile(t, filepath.Join(repoPath, "commands", "script.sh"), "#!/bin/bash")
			},
			expected: []BrowseItem{
				{Name: "valid", Path: "commands/valid.md", Type: TypeCommand},
			},
		},
		{
			name: "same name in both directories",
			setup: func(t *testing.T, repoPath string) {
				createFile(t, filepath.Join(repoPath, "commands", "common.md"), "# Common Root")
				createFile(t, filepath.Join(repoPath, ".claude", "commands", "common.md"), "# Common Claude")
			},
			expected: []BrowseItem{
				{Name: "common", Path: "commands/common.md", Type: TypeCommand},
				{Name: "common", Path: ".claude/commands/common.md", Type: TypeCommand},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoPath := setupTestRepo(t)
			defer func() { _ = os.RemoveAll(repoPath) }()

			tt.setup(t, repoPath)

			store := NewStore(repoPath)
			items, err := store.scanCommands(repoPath)
			if err != nil {
				t.Fatalf("scanCommands failed: %v", err)
			}

			// Sort for consistent comparison
			sort.Slice(items, func(i, j int) bool {
				return items[i].Path < items[j].Path
			})
			sort.Slice(tt.expected, func(i, j int) bool {
				return tt.expected[i].Path < tt.expected[j].Path
			})

			if len(items) != len(tt.expected) {
				t.Errorf("got %d items, want %d", len(items), len(tt.expected))
				t.Errorf("got: %+v", items)
				t.Errorf("want: %+v", tt.expected)
				return
			}

			for i := range items {
				if items[i].Name != tt.expected[i].Name ||
					items[i].Path != tt.expected[i].Path ||
					items[i].Type != tt.expected[i].Type {
					t.Errorf("item %d: got %+v, want %+v", i, items[i], tt.expected[i])
				}
			}
		})
	}
}

func TestScanAgents(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T, repoPath string)
		expected []BrowseItem
	}{
		{
			name: "empty repo",
			setup: func(t *testing.T, repoPath string) {
				// No setup needed
			},
			expected: nil,
		},
		{
			name: "agents in root agents/ directory",
			setup: func(t *testing.T, repoPath string) {
				createFile(t, filepath.Join(repoPath, "agents", "researcher.md"), "# Researcher")
				createFile(t, filepath.Join(repoPath, "agents", "coder.md"), "# Coder")
			},
			expected: []BrowseItem{
				{Name: "coder", Path: "agents/coder.md", Type: TypeAgent},
				{Name: "researcher", Path: "agents/researcher.md", Type: TypeAgent},
			},
		},
		{
			name: "agents in .claude/agents/ directory",
			setup: func(t *testing.T, repoPath string) {
				createFile(t, filepath.Join(repoPath, ".claude", "agents", "my-agent.md"), "# My Agent")
			},
			expected: []BrowseItem{
				{Name: "my-agent", Path: ".claude/agents/my-agent.md", Type: TypeAgent},
			},
		},
		{
			name: "agents in nested subdirectory",
			setup: func(t *testing.T, repoPath string) {
				createFile(t, filepath.Join(repoPath, "agents", "dev", "tester.md"), "# Tester")
				createFile(t, filepath.Join(repoPath, "agents", "dev", "reviewer.md"), "# Reviewer")
				createFile(t, filepath.Join(repoPath, "agents", "general.md"), "# General")
			},
			expected: []BrowseItem{
				{Name: "general", Path: "agents/general.md", Type: TypeAgent},
				{Name: "dev:reviewer", Path: "agents/dev/reviewer.md", Type: TypeAgent},
				{Name: "dev:tester", Path: "agents/dev/tester.md", Type: TypeAgent},
			},
		},
		{
			name: "agents in both directories with nested",
			setup: func(t *testing.T, repoPath string) {
				createFile(t, filepath.Join(repoPath, "agents", "root.md"), "# Root")
				createFile(t, filepath.Join(repoPath, "agents", "team", "lead.md"), "# Lead")
				createFile(t, filepath.Join(repoPath, ".claude", "agents", "claude.md"), "# Claude")
				createFile(t, filepath.Join(repoPath, ".claude", "agents", "team", "member.md"), "# Member")
			},
			expected: []BrowseItem{
				{Name: "root", Path: "agents/root.md", Type: TypeAgent},
				{Name: "team:lead", Path: "agents/team/lead.md", Type: TypeAgent},
				{Name: "claude", Path: ".claude/agents/claude.md", Type: TypeAgent},
				{Name: "team:member", Path: ".claude/agents/team/member.md", Type: TypeAgent},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoPath := setupTestRepo(t)
			defer func() { _ = os.RemoveAll(repoPath) }()

			tt.setup(t, repoPath)

			store := NewStore(repoPath)
			items, err := store.scanAgents(repoPath)
			if err != nil {
				t.Fatalf("scanAgents failed: %v", err)
			}

			// Sort for consistent comparison
			sort.Slice(items, func(i, j int) bool {
				return items[i].Path < items[j].Path
			})
			sort.Slice(tt.expected, func(i, j int) bool {
				return tt.expected[i].Path < tt.expected[j].Path
			})

			if len(items) != len(tt.expected) {
				t.Errorf("got %d items, want %d", len(items), len(tt.expected))
				t.Errorf("got: %+v", items)
				t.Errorf("want: %+v", tt.expected)
				return
			}

			for i := range items {
				if items[i].Name != tt.expected[i].Name ||
					items[i].Path != tt.expected[i].Path ||
					items[i].Type != tt.expected[i].Type {
					t.Errorf("item %d: got %+v, want %+v", i, items[i], tt.expected[i])
				}
			}
		})
	}
}

func TestScanHooks(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T, repoPath string)
		expected []BrowseItem
	}{
		{
			name: "empty repo",
			setup: func(t *testing.T, repoPath string) {
				// No setup needed
			},
			expected: nil,
		},
		{
			name: "hooks in root hooks/ directory",
			setup: func(t *testing.T, repoPath string) {
				createFile(t, filepath.Join(repoPath, "hooks", "pre-commit.sh"), "#!/bin/bash")
				createFile(t, filepath.Join(repoPath, "hooks", "post-build"), "#!/bin/bash")
			},
			expected: []BrowseItem{
				{Name: "post-build", Path: "hooks/post-build", Type: TypeHook},
				{Name: "pre-commit.sh", Path: "hooks/pre-commit.sh", Type: TypeHook},
			},
		},
		{
			name: "hooks in .claude/hooks/ directory",
			setup: func(t *testing.T, repoPath string) {
				createFile(t, filepath.Join(repoPath, ".claude", "hooks", "my-hook.sh"), "#!/bin/bash")
			},
			expected: []BrowseItem{
				{Name: "my-hook.sh", Path: ".claude/hooks/my-hook.sh", Type: TypeHook},
			},
		},
		{
			name: "hooks in both directories",
			setup: func(t *testing.T, repoPath string) {
				createFile(t, filepath.Join(repoPath, "hooks", "root-hook"), "#!/bin/bash")
				createFile(t, filepath.Join(repoPath, ".claude", "hooks", "claude-hook"), "#!/bin/bash")
			},
			expected: []BrowseItem{
				{Name: "root-hook", Path: "hooks/root-hook", Type: TypeHook},
				{Name: "claude-hook", Path: ".claude/hooks/claude-hook", Type: TypeHook},
			},
		},
		{
			name: "directories in hooks/ are ignored",
			setup: func(t *testing.T, repoPath string) {
				createDir(t, filepath.Join(repoPath, "hooks", "subdir"))
				createFile(t, filepath.Join(repoPath, "hooks", "subdir", "nested.sh"), "#!/bin/bash")
				createFile(t, filepath.Join(repoPath, "hooks", "valid-hook"), "#!/bin/bash")
			},
			expected: []BrowseItem{
				{Name: "valid-hook", Path: "hooks/valid-hook", Type: TypeHook},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoPath := setupTestRepo(t)
			defer func() { _ = os.RemoveAll(repoPath) }()

			tt.setup(t, repoPath)

			store := NewStore(repoPath)
			items, err := store.scanHooks(repoPath)
			if err != nil {
				t.Fatalf("scanHooks failed: %v", err)
			}

			// Sort for consistent comparison
			sort.Slice(items, func(i, j int) bool {
				return items[i].Path < items[j].Path
			})
			sort.Slice(tt.expected, func(i, j int) bool {
				return tt.expected[i].Path < tt.expected[j].Path
			})

			if len(items) != len(tt.expected) {
				t.Errorf("got %d items, want %d", len(items), len(tt.expected))
				t.Errorf("got: %+v", items)
				t.Errorf("want: %+v", tt.expected)
				return
			}

			for i := range items {
				if items[i].Name != tt.expected[i].Name ||
					items[i].Path != tt.expected[i].Path ||
					items[i].Type != tt.expected[i].Type {
					t.Errorf("item %d: got %+v, want %+v", i, items[i], tt.expected[i])
				}
			}
		})
	}
}

func TestBrowse(t *testing.T) {
	tests := []struct {
		name       string
		typeFilter PackageType
		setup      func(t *testing.T, repoPath string)
		expected   []BrowseItem
	}{
		{
			name:       "browse all types",
			typeFilter: "",
			setup: func(t *testing.T, repoPath string) {
				createFile(t, filepath.Join(repoPath, "skills", "my-skill", "SKILL.md"), "# Skill")
				createFile(t, filepath.Join(repoPath, "commands", "my-cmd.md"), "# Cmd")
				createFile(t, filepath.Join(repoPath, "agents", "my-agent.md"), "# Agent")
				createFile(t, filepath.Join(repoPath, "hooks", "my-hook.sh"), "#!/bin/bash")
			},
			expected: []BrowseItem{
				{Name: "my-skill", Path: "skills/my-skill", Type: TypeSkill},
				{Name: "my-cmd", Path: "commands/my-cmd.md", Type: TypeCommand},
				{Name: "my-agent", Path: "agents/my-agent.md", Type: TypeAgent},
				{Name: "my-hook.sh", Path: "hooks/my-hook.sh", Type: TypeHook},
			},
		},
		{
			name:       "browse only skills",
			typeFilter: TypeSkill,
			setup: func(t *testing.T, repoPath string) {
				createFile(t, filepath.Join(repoPath, "skills", "my-skill", "SKILL.md"), "# Skill")
				createFile(t, filepath.Join(repoPath, "commands", "my-cmd.md"), "# Cmd")
			},
			expected: []BrowseItem{
				{Name: "my-skill", Path: "skills/my-skill", Type: TypeSkill},
			},
		},
		{
			name:       "browse only commands",
			typeFilter: TypeCommand,
			setup: func(t *testing.T, repoPath string) {
				createFile(t, filepath.Join(repoPath, "skills", "my-skill", "SKILL.md"), "# Skill")
				createFile(t, filepath.Join(repoPath, "commands", "my-cmd.md"), "# Cmd")
			},
			expected: []BrowseItem{
				{Name: "my-cmd", Path: "commands/my-cmd.md", Type: TypeCommand},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoPath := setupTestRepo(t)
			defer func() { _ = os.RemoveAll(repoPath) }()

			tt.setup(t, repoPath)

			// Create a store and manually set up the test
			store := NewStore(repoPath)

			// Directly call the internal browse logic
			var items []BrowseItem

			if tt.typeFilter == "" || tt.typeFilter == TypeSkill {
				skillItems, _ := store.scanSkills(repoPath)
				items = append(items, skillItems...)
			}
			if tt.typeFilter == "" || tt.typeFilter == TypeCommand {
				cmdItems, _ := store.scanCommands(repoPath)
				items = append(items, cmdItems...)
			}
			if tt.typeFilter == "" || tt.typeFilter == TypeAgent {
				agentItems, _ := store.scanAgents(repoPath)
				items = append(items, agentItems...)
			}
			if tt.typeFilter == "" || tt.typeFilter == TypeHook {
				hookItems, _ := store.scanHooks(repoPath)
				items = append(items, hookItems...)
			}

			if len(items) != len(tt.expected) {
				t.Errorf("got %d items, want %d", len(items), len(tt.expected))
				return
			}
		})
	}
}
