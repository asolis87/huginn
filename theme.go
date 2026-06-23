package main

// Theme is a named color preset. It seeds the config's colors before the user's
// own config.toml is layered on top, so a theme gives a coherent look out of the
// box while still being fully overridable per field.
//
// Each segment carries an explicit foreground AND background: in powerline mode
// a good theme uses a dark block background with a light, contrasting text/icon
// color — they are not the same value. In plain mode only the foreground shows.
type Theme struct {
	CwdFg, CwdBg           string
	GitFg, GitBg           string
	GitDirtyFg, GitDirtyBg string
	NodeFg, NodeBg         string
	NodeMismatch           string
	SymbolSuccess          string
	SymbolError            string
}

// themes is the registry of built-in presets, selected by name via the config
// `theme` key. "huginn" is the default.
var themes = map[string]Theme{
	// huginn — Odin's raven: dark graphite plumage, icy "frost" blue highlights
	// (the iridescent sheen of raven feathers), warm rune-gold accents (Odin).
	// Palette derived from Nord, which is itself Norse-themed — fitting.
	"huginn": {
		CwdFg: "#eceff4", CwdBg: "#4c566a", // light text on graphite
		GitFg: "#eceff4", GitBg: "#5e81ac", // light text on steel blue
		GitDirtyFg: "#2e3440", GitDirtyBg: "#bf616a", // dark text on muted red
		NodeFg: "#2e3440", NodeBg: "#a3be8c", // dark text on moss green
		NodeMismatch:  "#ebcb8b", // rune gold
		SymbolSuccess: "#a3be8c", // moss green
		SymbolError:   "#bf616a", // muted red
	},
}

// applyTheme seeds cfg with the named theme's colors. It is applied BEFORE the
// user's config file, so user overrides still win. An unknown name is a no-op
// (cfg keeps the built-in defaults).
func applyTheme(cfg *Config, name string) {
	t, ok := themes[name]
	if !ok {
		return
	}
	cfg.Cwd.Color = t.CwdFg
	cfg.Cwd.Bg = t.CwdBg
	cfg.Git.Color = t.GitFg
	cfg.Git.Bg = t.GitBg
	cfg.Git.DirtyColor = t.GitDirtyFg
	cfg.Git.DirtyBg = t.GitDirtyBg
	cfg.Node.Color = t.NodeFg
	cfg.Node.Bg = t.NodeBg
	cfg.Node.MismatchColor = t.NodeMismatch
	cfg.Symbol.SuccessColor = t.SymbolSuccess
	cfg.Symbol.ErrorColor = t.SymbolError
}
