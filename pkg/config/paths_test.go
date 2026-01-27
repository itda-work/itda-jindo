package config

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestConstants(t *testing.T) {
	if AppName != "itda-skills" {
		t.Errorf("AppName = %q, want %q", AppName, "itda-skills")
	}
	if ConfigFileName != "config.toml" {
		t.Errorf("ConfigFileName = %q, want %q", ConfigFileName, "config.toml")
	}
}

func TestGetConfigDir(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping XDG tests on Windows")
	}

	t.Run("XDG_CONFIG_HOME set", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", "/tmp/test-xdg")

		dir, err := GetConfigDir()
		if err != nil {
			t.Fatalf("GetConfigDir() error: %v", err)
		}
		want := filepath.Join("/tmp/test-xdg", AppName)
		if dir != want {
			t.Errorf("GetConfigDir() = %q, want %q", dir, want)
		}
	})

	t.Run("XDG_CONFIG_HOME unset uses default", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", "")

		dir, err := GetConfigDir()
		if err != nil {
			t.Fatalf("GetConfigDir() error: %v", err)
		}
		// Should end with .config/itda-skills
		if !strings.HasSuffix(dir, filepath.Join(".config", AppName)) {
			t.Errorf("GetConfigDir() = %q, expected suffix %q", dir, filepath.Join(".config", AppName))
		}
	})
}

func TestGetConfigPath(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping XDG tests on Windows")
	}

	t.Setenv("XDG_CONFIG_HOME", "/tmp/test-xdg")

	path, err := GetConfigPath()
	if err != nil {
		t.Fatalf("GetConfigPath() error: %v", err)
	}
	if !strings.HasSuffix(path, ConfigFileName) {
		t.Errorf("GetConfigPath() = %q, expected to end with %q", path, ConfigFileName)
	}
	want := filepath.Join("/tmp/test-xdg", AppName, ConfigFileName)
	if path != want {
		t.Errorf("GetConfigPath() = %q, want %q", path, want)
	}
}
