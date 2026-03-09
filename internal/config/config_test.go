package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func setupTestHome(t *testing.T) {
	t.Helper()
	tmpDir := t.TempDir()

	// Use t.Setenv so env is automatically restored after each test.
	// Set XDG_CONFIG_HOME to a subdir of tmpDir — configDir() checks this first,
	// so we bypass os.UserHomeDir() entirely (avoids caching issues on CI).
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(tmpDir, ".config"))
	t.Setenv("HOME", tmpDir)
}

func TestSaveAndLoad(t *testing.T) {
	setupTestHome(t)

	cfg := &Config{
		Token:         "test-token-value",
		WorkspaceName: "Test Workspace",
		WorkspaceID:   "ws-123",
		BotID:         "bot-456",
	}

	if err := Save(cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if loaded.Token != cfg.Token {
		t.Errorf("Token = %q, want %q", loaded.Token, cfg.Token)
	}
	if loaded.WorkspaceName != cfg.WorkspaceName {
		t.Errorf("WorkspaceName = %q, want %q", loaded.WorkspaceName, cfg.WorkspaceName)
	}
	if loaded.WorkspaceID != cfg.WorkspaceID {
		t.Errorf("WorkspaceID = %q, want %q", loaded.WorkspaceID, cfg.WorkspaceID)
	}
	if loaded.BotID != cfg.BotID {
		t.Errorf("BotID = %q, want %q", loaded.BotID, cfg.BotID)
	}
}

func TestLoadMissing(t *testing.T) {
	setupTestHome(t)

	_, err := Load()
	if err == nil {
		t.Error("Load() should error when config file doesn't exist")
	}
}

func TestConfigFilePermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix file permissions not applicable on Windows")
	}
	setupTestHome(t)

	cfg := &Config{Token: "secret-token"}
	if err := Save(cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	configFile := configPath()
	info, err := os.Stat(configFile)
	if err != nil {
		t.Fatalf("Stat() error = %v", err)
	}

	perm := info.Mode().Perm()
	if perm&0077 != 0 {
		t.Errorf("Config file permissions = %o, want no group/other access", perm)
	}
}

func TestProfileSupport(t *testing.T) {
	setupTestHome(t)

	cfg := &Config{}

	// Add profiles
	cfg.SetProfile("default", &Profile{
		Token:         "default-token",
		WorkspaceName: "Default Workspace",
	})
	cfg.SetProfile("work", &Profile{
		Token:         "work-token",
		WorkspaceName: "Work Workspace",
	})
	cfg.CurrentProfile = "default"

	if err := Save(cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Check profiles exist
	if len(loaded.Profiles) != 2 {
		t.Errorf("got %d profiles, want 2", len(loaded.Profiles))
	}

	// Check current profile
	profile := loaded.GetCurrentProfile()
	if profile == nil {
		t.Fatal("GetCurrentProfile() returned nil")
	}
	if profile.Token != "default-token" {
		t.Errorf("Token = %q, want %q", profile.Token, "default-token")
	}

	// Switch profile
	loaded.CurrentProfile = "work"
	profile = loaded.GetCurrentProfile()
	if profile.Token != "work-token" {
		t.Errorf("Token = %q, want %q", profile.Token, "work-token")
	}
}

func TestMigrateToProfiles(t *testing.T) {
	setupTestHome(t)

	// Create legacy config
	cfg := &Config{
		Token:         "legacy-token",
		WorkspaceName: "Legacy Workspace",
		WorkspaceID:   "ws-legacy",
		BotID:         "bot-legacy",
	}

	cfg.MigrateToProfiles()

	// Check migration
	if len(cfg.Profiles) != 1 {
		t.Errorf("got %d profiles, want 1", len(cfg.Profiles))
	}
	if cfg.CurrentProfile != "default" {
		t.Errorf("CurrentProfile = %q, want %q", cfg.CurrentProfile, "default")
	}

	profile := cfg.GetCurrentProfile()
	if profile == nil {
		t.Fatal("GetCurrentProfile() returned nil")
	}
	if profile.Token != "legacy-token" {
		t.Errorf("Token = %q, want %q", profile.Token, "legacy-token")
	}

	// Legacy fields should be cleared
	if cfg.Token != "" {
		t.Errorf("legacy Token should be cleared, got %q", cfg.Token)
	}
}

func TestListProfiles(t *testing.T) {
	cfg := &Config{
		Profiles: map[string]*Profile{
			"work":     {Token: "t1"},
			"personal": {Token: "t2"},
			"default":  {Token: "t3"},
		},
	}

	profiles := cfg.ListProfiles()
	if len(profiles) != 3 {
		t.Errorf("got %d profiles, want 3", len(profiles))
	}

	// Should be sorted alphabetically
	if profiles[0] != "default" {
		t.Errorf("profiles[0] = %q, want %q", profiles[0], "default")
	}
	if profiles[1] != "personal" {
		t.Errorf("profiles[1] = %q, want %q", profiles[1], "personal")
	}
	if profiles[2] != "work" {
		t.Errorf("profiles[2] = %q, want %q", profiles[2], "work")
	}
}

func TestGetProfile(t *testing.T) {
	cfg := &Config{
		Profiles: map[string]*Profile{
			"default": {Token: "t1"},
			"work":    {Token: "t2"},
		},
	}

	p := cfg.GetProfile("work")
	if p == nil {
		t.Fatal("GetProfile(\"work\") returned nil")
	}
	if p.Token != "t2" {
		t.Errorf("Token = %q, want %q", p.Token, "t2")
	}
}

func TestGetProfile_NotFound(t *testing.T) {
	cfg := &Config{
		Profiles: map[string]*Profile{
			"default": {Token: "t1"},
		},
	}

	p := cfg.GetProfile("nonexistent")
	if p != nil {
		t.Errorf("GetProfile(\"nonexistent\") should return nil, got %+v", p)
	}
}

func TestGetProfile_NilProfiles(t *testing.T) {
	cfg := &Config{}
	p := cfg.GetProfile("anything")
	if p != nil {
		t.Errorf("GetProfile on empty config should return nil, got %+v", p)
	}
}

func TestLegacyConfigDirFallback(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(tmpDir, ".config"))
	t.Setenv("HOME", tmpDir)

	// Write config to legacy path
	legacyDir := filepath.Join(tmpDir, ".config", "notion-cli")
	if err := os.MkdirAll(legacyDir, 0700); err != nil {
		t.Fatal(err)
	}
	legacyCfg := &Config{Token: "legacy-token", WorkspaceName: "Legacy"}
	data, _ := json.Marshal(legacyCfg)
	if err := os.WriteFile(filepath.Join(legacyDir, "config.json"), data, 0600); err != nil {
		t.Fatal(err)
	}

	// Load should find it via legacy path
	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if loaded.Token != "legacy-token" {
		t.Errorf("Token = %q, want %q (from legacy path)", loaded.Token, "legacy-token")
	}

	// Save should write to new path
	if err := Save(loaded); err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	newPath := filepath.Join(tmpDir, ".config", "notion-agent-cli", "config.json")
	if _, err := os.Stat(newPath); os.IsNotExist(err) {
		t.Error("Save() should write to new config path notion-agent-cli/")
	}
}

func TestGetCurrentProfileLegacyFallback(t *testing.T) {
	// Test that GetCurrentProfile falls back to legacy fields
	cfg := &Config{
		Token:         "legacy-token",
		WorkspaceName: "Legacy",
	}

	profile := cfg.GetCurrentProfile()
	if profile == nil {
		t.Fatal("GetCurrentProfile() returned nil")
	}
	if profile.Token != "legacy-token" {
		t.Errorf("Token = %q, want %q", profile.Token, "legacy-token")
	}
}
