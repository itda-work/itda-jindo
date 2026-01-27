package config

import (
	"testing"
)

func TestParseDotKey(t *testing.T) {
	t.Helper()

	tests := []struct {
		name    string
		key     string
		want    []string
		wantErr error
	}{
		{
			name: "single key",
			key:  "common",
			want: []string{"common"},
		},
		{
			name: "multi level key",
			key:  "common.api_keys.tiingo",
			want: []string{"common", "api_keys", "tiingo"},
		},
		{
			name:    "empty string",
			key:     "",
			wantErr: ErrInvalidKey,
		},
		{
			name:    "whitespace only",
			key:     "   ",
			wantErr: ErrInvalidKey,
		},
		{
			name:    "consecutive dots",
			key:     "a..b",
			wantErr: ErrInvalidKey,
		},
		{
			name: "trimmed whitespace",
			key:  "  common.api_keys  ",
			want: []string{"common", "api_keys"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDotKey(tt.key)
			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("parseDotKey(%q) error = %v, want %v", tt.key, err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseDotKey(%q) unexpected error: %v", tt.key, err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("parseDotKey(%q) = %v, want %v", tt.key, got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("parseDotKey(%q)[%d] = %q, want %q", tt.key, i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestGetNestedValue(t *testing.T) {
	t.Helper()

	data := map[string]any{
		"level1": "value1",
		"nested": map[string]any{
			"level2": "value2",
			"deep": map[string]any{
				"level3": "value3",
			},
		},
		"notmap": "string_value",
	}

	tests := []struct {
		name    string
		keys    []string
		want    any
		wantErr error
	}{
		{
			name: "single level",
			keys: []string{"level1"},
			want: "value1",
		},
		{
			name: "multi level",
			keys: []string{"nested", "deep", "level3"},
			want: "value3",
		},
		{
			name:    "key not found",
			keys:    []string{"nonexistent"},
			wantErr: ErrKeyNotFound,
		},
		{
			name:    "intermediate not a map",
			keys:    []string{"notmap", "child"},
			wantErr: ErrNotAMap,
		},
		{
			name: "empty keys returns data",
			keys: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getNestedValue(data, tt.keys)
			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("getNestedValue() error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("getNestedValue() unexpected error: %v", err)
			}
			if tt.want != nil && got != tt.want {
				t.Errorf("getNestedValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetNestedValue(t *testing.T) {
	t.Helper()

	tests := []struct {
		name    string
		initial map[string]any
		keys    []string
		value   any
		wantErr error
		check   func(t *testing.T, data map[string]any)
	}{
		{
			name:    "single level",
			initial: map[string]any{},
			keys:    []string{"key"},
			value:   "value",
			check: func(t *testing.T, data map[string]any) {
				if data["key"] != "value" {
					t.Errorf("expected data[key] = value, got %v", data["key"])
				}
			},
		},
		{
			name:    "multi level with auto-created maps",
			initial: map[string]any{},
			keys:    []string{"a", "b", "c"},
			value:   "deep",
			check: func(t *testing.T, data map[string]any) {
				a := data["a"].(map[string]any)
				b := a["b"].(map[string]any)
				if b["c"] != "deep" {
					t.Errorf("expected data[a][b][c] = deep, got %v", b["c"])
				}
			},
		},
		{
			name:    "overwrite existing value",
			initial: map[string]any{"key": "old"},
			keys:    []string{"key"},
			value:   "new",
			check: func(t *testing.T, data map[string]any) {
				if data["key"] != "new" {
					t.Errorf("expected data[key] = new, got %v", data["key"])
				}
			},
		},
		{
			name:    "intermediate not a map",
			initial: map[string]any{"key": "string"},
			keys:    []string{"key", "child"},
			value:   "val",
			wantErr: ErrNotAMap,
		},
		{
			name:    "empty keys",
			initial: map[string]any{},
			keys:    []string{},
			value:   "val",
			wantErr: ErrInvalidKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := setNestedValue(tt.initial, tt.keys, tt.value)
			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("setNestedValue() error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("setNestedValue() unexpected error: %v", err)
			}
			if tt.check != nil {
				tt.check(t, tt.initial)
			}
		})
	}
}

func TestDeleteNestedValue(t *testing.T) {
	t.Helper()

	tests := []struct {
		name    string
		initial map[string]any
		keys    []string
		wantErr error
		check   func(t *testing.T, data map[string]any)
	}{
		{
			name:    "single level",
			initial: map[string]any{"key": "value"},
			keys:    []string{"key"},
			check: func(t *testing.T, data map[string]any) {
				if _, ok := data["key"]; ok {
					t.Error("expected key to be deleted")
				}
			},
		},
		{
			name: "nested key",
			initial: map[string]any{
				"parent": map[string]any{
					"child": "value",
					"other": "keep",
				},
			},
			keys: []string{"parent", "child"},
			check: func(t *testing.T, data map[string]any) {
				parent := data["parent"].(map[string]any)
				if _, ok := parent["child"]; ok {
					t.Error("expected child to be deleted")
				}
				if parent["other"] != "keep" {
					t.Error("expected other to remain")
				}
			},
		},
		{
			name:    "nonexistent key in path",
			initial: map[string]any{},
			keys:    []string{"a", "b"},
			wantErr: ErrKeyNotFound,
		},
		{
			name:    "intermediate not a map",
			initial: map[string]any{"key": "string"},
			keys:    []string{"key", "child"},
			wantErr: ErrNotAMap,
		},
		{
			name:    "empty keys",
			initial: map[string]any{},
			keys:    []string{},
			wantErr: ErrInvalidKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := deleteNestedValue(tt.initial, tt.keys)
			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("deleteNestedValue() error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("deleteNestedValue() unexpected error: %v", err)
			}
			if tt.check != nil {
				tt.check(t, tt.initial)
			}
		})
	}
}

func TestParseValue(t *testing.T) {
	t.Helper()

	tests := []struct {
		name  string
		input string
		want  any
	}{
		{name: "true", input: "true", want: true},
		{name: "false", input: "false", want: false},
		{name: "positive integer", input: "123", want: int64(123)},
		{name: "negative integer", input: "-42", want: int64(-42)},
		{name: "float", input: "3.14", want: float64(3.14)},
		{name: "string", input: "hello", want: "hello"},
		{name: "empty string", input: "", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseValue(tt.input)
			if got != tt.want {
				t.Errorf("ParseValue(%q) = %v (%T), want %v (%T)", tt.input, got, got, tt.want, tt.want)
			}
		})
	}
}

func TestToEnvKey(t *testing.T) {
	t.Helper()

	tests := []struct {
		name   string
		dotKey string
		want   string
	}{
		{
			name:   "single key",
			dotKey: "market",
			want:   "ITDA_MARKET",
		},
		{
			name:   "multi level",
			dotKey: "common.api_keys.tiingo",
			want:   "ITDA_COMMON_API_KEYS_TIINGO",
		},
		{
			name:   "key with hyphen",
			dotKey: "skills.quant-data.format",
			want:   "ITDA_SKILLS_QUANT-DATA_FORMAT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toEnvKey(tt.dotKey)
			if got != tt.want {
				t.Errorf("toEnvKey(%q) = %q, want %q", tt.dotKey, got, tt.want)
			}
		})
	}
}
