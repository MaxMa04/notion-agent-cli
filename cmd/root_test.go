package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/MaxMa04/notion-agent-cli/internal/config"
)

func setupTestConfig(t *testing.T) {
	t.Helper()
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(tmpDir, ".config"))
	t.Setenv("HOME", tmpDir)

	// Clear env vars that could interfere
	t.Setenv("NOTION_TOKEN", "")
	t.Setenv("NOTION_PROFILE", "")

	// Reset global flag
	profileFlag = ""

	cfg := &config.Config{}
	cfg.SetProfile("default", &config.Profile{
		Token:         "default-token",
		WorkspaceName: "Default",
	})
	cfg.SetProfile("work", &config.Profile{
		Token:         "work-token",
		WorkspaceName: "Work",
	})
	cfg.SetProfile("alfred", &config.Profile{
		Token:         "alfred-token",
		WorkspaceName: "Alfred Bot",
	})
	cfg.CurrentProfile = "default"

	if err := config.Save(cfg); err != nil {
		t.Fatalf("Save config: %v", err)
	}
}

func TestGetToken_EnvVar(t *testing.T) {
	setupTestConfig(t)
	t.Setenv("NOTION_TOKEN", "env-token-value")

	token, err := getToken()
	if err != nil {
		t.Fatalf("getToken() error = %v", err)
	}
	if token != "env-token-value" {
		t.Errorf("token = %q, want %q", token, "env-token-value")
	}
}

func TestGetToken_EnvOverridesProfileFlag(t *testing.T) {
	setupTestConfig(t)
	t.Setenv("NOTION_TOKEN", "env-wins")
	profileFlag = "work"

	token, err := getToken()
	if err != nil {
		t.Fatalf("getToken() error = %v", err)
	}
	if token != "env-wins" {
		t.Errorf("NOTION_TOKEN should take priority, got %q", token)
	}
}

func TestGetToken_ProfileFlag(t *testing.T) {
	setupTestConfig(t)
	profileFlag = "work"

	token, err := getToken()
	if err != nil {
		t.Fatalf("getToken() error = %v", err)
	}
	if token != "work-token" {
		t.Errorf("token = %q, want %q", token, "work-token")
	}
}

func TestGetToken_ProfileFlagOverridesNotionProfile(t *testing.T) {
	setupTestConfig(t)
	profileFlag = "work"
	t.Setenv("NOTION_PROFILE", "alfred")

	token, err := getToken()
	if err != nil {
		t.Fatalf("getToken() error = %v", err)
	}
	if token != "work-token" {
		t.Errorf("--profile flag should beat NOTION_PROFILE, got %q", token)
	}
}

func TestGetToken_NotionProfileEnv(t *testing.T) {
	setupTestConfig(t)
	t.Setenv("NOTION_PROFILE", "alfred")

	token, err := getToken()
	if err != nil {
		t.Fatalf("getToken() error = %v", err)
	}
	if token != "alfred-token" {
		t.Errorf("token = %q, want %q", token, "alfred-token")
	}
}

func TestGetToken_CurrentProfile(t *testing.T) {
	setupTestConfig(t)

	token, err := getToken()
	if err != nil {
		t.Fatalf("getToken() error = %v", err)
	}
	if token != "default-token" {
		t.Errorf("token = %q, want %q", token, "default-token")
	}
}

func TestGetToken_ProfileNotFound(t *testing.T) {
	setupTestConfig(t)
	profileFlag = "nonexistent"

	_, err := getToken()
	if err == nil {
		t.Fatal("getToken() should error for unknown profile")
	}
}

func TestGetToken_NotionProfileNotFound(t *testing.T) {
	setupTestConfig(t)
	t.Setenv("NOTION_PROFILE", "nonexistent")

	_, err := getToken()
	if err == nil {
		t.Fatal("getToken() should error for unknown NOTION_PROFILE")
	}
}

func TestGetToken_NoConfig(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(tmpDir, ".config"))
	t.Setenv("HOME", tmpDir)
	t.Setenv("NOTION_TOKEN", "")
	t.Setenv("NOTION_PROFILE", "")
	profileFlag = ""

	_, err := getToken()
	if err == nil {
		t.Fatal("getToken() should error when no config exists")
	}
}

func TestGetToken_LegacyConfigPath(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(tmpDir, ".config"))
	t.Setenv("HOME", tmpDir)
	t.Setenv("NOTION_TOKEN", "")
	t.Setenv("NOTION_PROFILE", "")
	profileFlag = ""

	// Write config to legacy notion-cli path
	legacyDir := filepath.Join(tmpDir, ".config", "notion-cli")
	if err := os.MkdirAll(legacyDir, 0700); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}
	cfg := `{"current_profile":"default","profiles":{"default":{"token":"legacy-works"}}}`
	if err := os.WriteFile(filepath.Join(legacyDir, "config.json"), []byte(cfg), 0600); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	token, err := getToken()
	if err != nil {
		t.Fatalf("getToken() error = %v (should fall back to legacy path)", err)
	}
	if token != "legacy-works" {
		t.Errorf("token = %q, want %q", token, "legacy-works")
	}
}
