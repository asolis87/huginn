package main

import "testing"

// TestApplyThemeKnown verifies the huginn theme seeds segment colors.
func TestApplyThemeKnown(t *testing.T) {
	cfg := defaultConfig()
	applyTheme(&cfg, "huginn")

	if cfg.Git.Bg != "#5e81ac" {
		t.Errorf("git bg = %q, want huginn steel blue #5e81ac", cfg.Git.Bg)
	}
	if cfg.Cwd.Bg != "#4c566a" {
		t.Errorf("cwd bg = %q, want huginn graphite #4c566a", cfg.Cwd.Bg)
	}
}

// TestApplyThemeUnknown verifies an unknown theme leaves config untouched, so
// the base defaults survive (never crash, never blank the prompt).
func TestApplyThemeUnknown(t *testing.T) {
	cfg := defaultConfig()
	before := cfg.Cwd.Color
	applyTheme(&cfg, "does-not-exist")
	if cfg.Cwd.Color != before {
		t.Errorf("unknown theme changed cwd color: %q -> %q", before, cfg.Cwd.Color)
	}
}
