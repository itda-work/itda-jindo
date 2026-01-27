package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	c := New()
	if c == nil {
		t.Fatal("New() returned nil")
	}
	if !c.IsEmpty() {
		t.Error("New() config should be empty")
	}
}

func TestLoadFromPath(t *testing.T) {
	t.Run("valid TOML file", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "config.toml")
		content := []byte("[common]\ndefault_market = \"kr\"\n")
		if err := os.WriteFile(path, content, 0644); err != nil {
			t.Fatal(err)
		}

		c, err := LoadFromPath(path)
		if err != nil {
			t.Fatalf("LoadFromPath() error: %v", err)
		}

		val, err := c.Get("common.default_market")
		if err != nil {
			t.Fatalf("Get() error: %v", err)
		}
		if val != "kr" {
			t.Errorf("Get(common.default_market) = %v, want kr", val)
		}
	})

	t.Run("nonexistent file returns empty config", func(t *testing.T) {
		c, err := LoadFromPath("/nonexistent/path/config.toml")
		if err != nil {
			t.Fatalf("LoadFromPath() error: %v", err)
		}
		if !c.IsEmpty() {
			t.Error("expected empty config for nonexistent file")
		}
	})

	t.Run("invalid TOML returns error", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "config.toml")
		content := []byte("[invalid\n")
		if err := os.WriteFile(path, content, 0644); err != nil {
			t.Fatal(err)
		}

		_, err := LoadFromPath(path)
		if err == nil {
			t.Error("LoadFromPath() expected error for invalid TOML")
		}
	})
}

func TestSaveAndLoadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")

	// Create and save
	c := New()
	if err := c.Set("common.default_market", "kr"); err != nil {
		t.Fatal(err)
	}
	if err := c.Set("common.api_keys.tiingo", "test-key"); err != nil {
		t.Fatal(err)
	}

	// SaveToPath calls EnsureConfigDir which uses GetConfigDir.
	// Write directly to avoid that dependency.
	tomlStr, err := c.ToTOML()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(tomlStr), 0644); err != nil {
		t.Fatal(err)
	}

	// Load and verify
	loaded, err := LoadFromPath(path)
	if err != nil {
		t.Fatalf("LoadFromPath() error: %v", err)
	}

	val, err := loaded.Get("common.default_market")
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if val != "kr" {
		t.Errorf("round-trip: common.default_market = %v, want kr", val)
	}

	val, err = loaded.Get("common.api_keys.tiingo")
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if val != "test-key" {
		t.Errorf("round-trip: common.api_keys.tiingo = %v, want test-key", val)
	}
}

func TestGetSet(t *testing.T) {
	t.Run("single key", func(t *testing.T) {
		c := New()
		if err := c.Set("key", "value"); err != nil {
			t.Fatal(err)
		}
		val, err := c.Get("key")
		if err != nil {
			t.Fatal(err)
		}
		if val != "value" {
			t.Errorf("Get(key) = %v, want value", val)
		}
	})

	t.Run("nested key", func(t *testing.T) {
		c := New()
		if err := c.Set("a.b.c", "deep"); err != nil {
			t.Fatal(err)
		}
		val, err := c.Get("a.b.c")
		if err != nil {
			t.Fatal(err)
		}
		if val != "deep" {
			t.Errorf("Get(a.b.c) = %v, want deep", val)
		}
	})

	t.Run("intermediate maps auto-created", func(t *testing.T) {
		c := New()
		if err := c.Set("x.y.z", 42); err != nil {
			t.Fatal(err)
		}
		// Access intermediate map
		val, err := c.Get("x.y")
		if err != nil {
			t.Fatal(err)
		}
		m, ok := val.(map[string]any)
		if !ok {
			t.Fatalf("expected map, got %T", val)
		}
		if m["z"] != 42 {
			t.Errorf("expected x.y.z = 42, got %v", m["z"])
		}
	})

	t.Run("nonexistent key", func(t *testing.T) {
		c := New()
		_, err := c.Get("nonexistent")
		if err != ErrKeyNotFound {
			t.Errorf("Get(nonexistent) error = %v, want ErrKeyNotFound", err)
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("existing key", func(t *testing.T) {
		c := New()
		_ = c.Set("key", "value")
		if err := c.Delete("key"); err != nil {
			t.Fatal(err)
		}
		_, err := c.Get("key")
		if err != ErrKeyNotFound {
			t.Error("expected key to be deleted")
		}
	})

	t.Run("nested key", func(t *testing.T) {
		c := New()
		_ = c.Set("a.b", "val")
		_ = c.Set("a.c", "keep")
		if err := c.Delete("a.b"); err != nil {
			t.Fatal(err)
		}
		_, err := c.Get("a.b")
		if err != ErrKeyNotFound {
			t.Error("expected a.b to be deleted")
		}
		val, err := c.Get("a.c")
		if err != nil {
			t.Fatal(err)
		}
		if val != "keep" {
			t.Errorf("a.c = %v, want keep", val)
		}
	})
}

func TestGetWithEnv(t *testing.T) {
	t.Run("env var takes priority", func(t *testing.T) {
		c := New()
		_ = c.Set("common.market", "kr")
		t.Setenv("ITDA_COMMON_MARKET", "us")

		val, ok := c.GetWithEnv("common.market")
		if !ok {
			t.Fatal("GetWithEnv() returned false")
		}
		if val != "us" {
			t.Errorf("GetWithEnv() = %v, want us", val)
		}
	})

	t.Run("config fallback", func(t *testing.T) {
		c := New()
		_ = c.Set("common.market", "kr")

		val, ok := c.GetWithEnv("common.market")
		if !ok {
			t.Fatal("GetWithEnv() returned false")
		}
		if val != "kr" {
			t.Errorf("GetWithEnv() = %v, want kr", val)
		}
	})

	t.Run("neither env nor config", func(t *testing.T) {
		c := New()

		_, ok := c.GetWithEnv("nonexistent.key")
		if ok {
			t.Error("GetWithEnv() returned true for nonexistent key")
		}
	})
}

func TestGetWithVendorEnv(t *testing.T) {
	const vendorEnv = "TEST_VENDOR_API_KEY_FOR_JINDO"

	t.Run("vendor env var takes priority", func(t *testing.T) {
		c := New()
		_ = c.Set("common.api_keys.testvendor", "config-key")
		t.Setenv("ITDA_COMMON_API_KEYS_TESTVENDOR", "itda-key")
		t.Setenv(vendorEnv, "vendor-key")

		val, ok := c.GetWithVendorEnv("common.api_keys.testvendor", vendorEnv)
		if !ok {
			t.Fatal("GetWithVendorEnv() returned false")
		}
		if val != "vendor-key" {
			t.Errorf("GetWithVendorEnv() = %v, want vendor-key", val)
		}
	})

	t.Run("ITDA env fallback", func(t *testing.T) {
		c := New()
		_ = c.Set("common.api_keys.testvendor", "config-key")
		t.Setenv("ITDA_COMMON_API_KEYS_TESTVENDOR", "itda-key")

		val, ok := c.GetWithVendorEnv("common.api_keys.testvendor", vendorEnv)
		if !ok {
			t.Fatal("GetWithVendorEnv() returned false")
		}
		if val != "itda-key" {
			t.Errorf("GetWithVendorEnv() = %v, want itda-key", val)
		}
	})

	t.Run("config fallback", func(t *testing.T) {
		c := New()
		_ = c.Set("common.api_keys.testvendor", "config-key")

		val, ok := c.GetWithVendorEnv("common.api_keys.testvendor", vendorEnv)
		if !ok {
			t.Fatal("GetWithVendorEnv() returned false")
		}
		if val != "config-key" {
			t.Errorf("GetWithVendorEnv() = %v, want config-key", val)
		}
	})

	t.Run("none found", func(t *testing.T) {
		c := New()

		_, ok := c.GetWithVendorEnv("common.api_keys.testvendor", vendorEnv)
		if ok {
			t.Error("GetWithVendorEnv() returned true when nothing is set")
		}
	})
}

func TestIsEmpty(t *testing.T) {
	t.Run("empty config", func(t *testing.T) {
		c := New()
		if !c.IsEmpty() {
			t.Error("IsEmpty() = false for new config")
		}
	})

	t.Run("non-empty config", func(t *testing.T) {
		c := New()
		_ = c.Set("key", "value")
		if c.IsEmpty() {
			t.Error("IsEmpty() = true for config with data")
		}
	})
}

func TestToMap(t *testing.T) {
	c := New()
	_ = c.Set("key", "value")

	m := c.ToMap()
	if m["key"] != "value" {
		t.Errorf("ToMap()[key] = %v, want value", m["key"])
	}
}

func TestToTOML(t *testing.T) {
	c := New()
	_ = c.Set("common.market", "kr")

	tomlStr, err := c.ToTOML()
	if err != nil {
		t.Fatalf("ToTOML() error: %v", err)
	}
	if tomlStr == "" {
		t.Error("ToTOML() returned empty string")
	}

	// Verify it can be loaded back
	dir := t.TempDir()
	path := filepath.Join(dir, "test.toml")
	if err := os.WriteFile(path, []byte(tomlStr), 0644); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadFromPath(path)
	if err != nil {
		t.Fatalf("LoadFromPath() error: %v", err)
	}
	val, err := loaded.Get("common.market")
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if val != "kr" {
		t.Errorf("ToTOML round-trip: common.market = %v, want kr", val)
	}
}
