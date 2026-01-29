package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigDir(t *testing.T) {
	t.Run("uses XDG_CONFIG_HOME when set", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", "/tmp/xdg-config")
		got := configDir()
		want := "/tmp/xdg-config/ask"
		if got != want {
			t.Errorf("configDir() = %q, want %q", got, want)
		}
	})

	t.Run("falls back to ~/.config", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", "")
		got := configDir()
		home, _ := os.UserHomeDir()
		want := filepath.Join(home, ".config", "ask")
		if got != want {
			t.Errorf("configDir() = %q, want %q", got, want)
		}
	})
}

func TestDataDir(t *testing.T) {
	t.Run("uses XDG_DATA_HOME when set", func(t *testing.T) {
		t.Setenv("XDG_DATA_HOME", "/tmp/xdg-data")
		got := dataDir()
		want := "/tmp/xdg-data/ask"
		if got != want {
			t.Errorf("dataDir() = %q, want %q", got, want)
		}
	})

	t.Run("falls back to ~/.local/share", func(t *testing.T) {
		t.Setenv("XDG_DATA_HOME", "")
		got := dataDir()
		home, _ := os.UserHomeDir()
		want := filepath.Join(home, ".local", "share", "ask")
		if got != want {
			t.Errorf("dataDir() = %q, want %q", got, want)
		}
	})
}

func TestLoadConfig(t *testing.T) {
	t.Run("returns zero value for missing file", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", "/tmp/nonexistent-ccq-config-dir")
		cfg := loadConfig()
		if cfg.DefaultModel != "" || cfg.RawOutput != false {
			t.Errorf("loadConfig() = %+v, want zero value", cfg)
		}
	})

	t.Run("loads valid config", func(t *testing.T) {
		dir := t.TempDir()
		t.Setenv("XDG_CONFIG_HOME", dir)
		cqDir := filepath.Join(dir, "ask")
		os.MkdirAll(cqDir, 0755)
		os.WriteFile(filepath.Join(cqDir, "config.json"), []byte(`{"default_model":"opus","raw_output":true}`), 0644)

		cfg := loadConfig()
		if cfg.DefaultModel != "opus" {
			t.Errorf("DefaultModel = %q, want %q", cfg.DefaultModel, "opus")
		}
		if !cfg.RawOutput {
			t.Error("RawOutput = false, want true")
		}
	})

	t.Run("returns zero value for invalid json", func(t *testing.T) {
		dir := t.TempDir()
		t.Setenv("XDG_CONFIG_HOME", dir)
		cqDir := filepath.Join(dir, "ask")
		os.MkdirAll(cqDir, 0755)
		os.WriteFile(filepath.Join(cqDir, "config.json"), []byte(`not json`), 0644)

		cfg := loadConfig()
		if cfg.DefaultModel != "" {
			t.Errorf("DefaultModel = %q, want empty for invalid json", cfg.DefaultModel)
		}
	})
}

func TestDefaultConfigJSON(t *testing.T) {
	data := defaultConfigJSON()
	if len(data) == 0 {
		t.Error("defaultConfigJSON() returned empty data")
	}
	if data[len(data)-1] != '\n' {
		t.Error("defaultConfigJSON() should end with newline")
	}
}
