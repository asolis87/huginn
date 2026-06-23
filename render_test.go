package main

import (
	"strings"
	"testing"
)

// These tests assert PROPERTIES of the rendered output (content present, right
// number of separators, correct shell escaping) rather than exact byte strings.
// Property tests survive cosmetic tweaks and only fail on real regressions.

// --- Plain renderer ----------------------------------------------------------

func TestPlainRenderShowsAllSegmentText(t *testing.T) {
	segs := []Segment{
		{text: "~/dev", fg: "blue", bg: "blue"},
		{text: "main", fg: "yellow", bg: "yellow"},
	}
	out := plainRenderer{}.Render(ShellZsh, segs)

	for _, seg := range segs {
		if !strings.Contains(out, seg.text) {
			t.Errorf("plain output missing segment text %q\ngot: %q", seg.text, out)
		}
	}
}

func TestPlainRenderUncoloredSegmentHasNoAnsi(t *testing.T) {
	// A segment with no fg should render as bare text — no escape sequences.
	out := plainRenderer{}.Render(ShellZsh, []Segment{{text: "plain"}})
	if strings.Contains(out, "\x1b[") {
		t.Errorf("uncolored segment emitted ANSI: %q", out)
	}
}

func TestPlainDefaultIsSingleSpaceJoin(t *testing.T) {
	// The cardinal invariant: with default plain config the output must be the
	// original look — segments joined by a single bare space, nothing else.
	segs := []Segment{{text: "a"}, {text: "b"}, {text: "c"}}
	out := plainRenderer{cfg: defaultConfig().Plain}.Render(ShellFish, segs)

	if out != "a b c" {
		t.Errorf("default plain join changed: got %q, want %q", out, "a b c")
	}
}

func TestPlainCustomSeparator(t *testing.T) {
	segs := []Segment{{text: "a"}, {text: "b"}, {text: "c"}}
	out := plainRenderer{cfg: PlainConfig{Separator: " | "}}.Render(ShellFish, segs)

	if out != "a | b | c" {
		t.Errorf("custom separator: got %q, want %q", out, "a | b | c")
	}
}

func TestPlainWrapEachSegment(t *testing.T) {
	segs := []Segment{{text: "a"}, {text: "b"}}
	out := plainRenderer{cfg: PlainConfig{Separator: " ", WrapLeft: "[", WrapRight: "]"}}.
		Render(ShellFish, segs)

	if out != "[a] [b]" {
		t.Errorf("wrap: got %q, want %q", out, "[a] [b]")
	}
}

func TestPlainEmptySeparatorFallsBackToSpace(t *testing.T) {
	// An empty separator must not glue segments together; it falls back to a space.
	segs := []Segment{{text: "a"}, {text: "b"}}
	out := plainRenderer{cfg: PlainConfig{Separator: ""}}.Render(ShellFish, segs)

	if out != "a b" {
		t.Errorf("empty separator should fall back to space: got %q", out)
	}
}

func TestPlainSeparatorColorWithoutWrapEmitsNoStrayEscapes(t *testing.T) {
	// separator_color set, no wrap, segments have no fg: only the separator should
	// carry color. Wrap is empty, so NO escape pairs may wrap the empty wrap slots.
	// Expect exactly one colored separator and no other escape sequences.
	segs := []Segment{{text: "a"}, {text: "b"}}
	out := plainRenderer{cfg: PlainConfig{Separator: " | ", SeparatorColor: "red"}}.
		Render(ShellFish, segs)

	// fish emits raw escapes (no %{ %}): the only \x1b runs belong to the one
	// separator (a set + a reset = two \x1b). Empty wraps must add none.
	if got := strings.Count(out, "\x1b["); got != 2 {
		t.Errorf("expected exactly 2 escapes (one colored separator), got %d: %q", got, out)
	}
}

// --- Shell escaping (the easiest thing to break silently) --------------------

func TestColorWrappingPerShell(t *testing.T) {
	segs := []Segment{{text: "x", fg: "red", bg: "red"}}

	zsh := plainRenderer{}.Render(ShellZsh, segs)
	if !strings.Contains(zsh, "%{") || !strings.Contains(zsh, "%}") {
		t.Errorf("zsh render must wrap escapes in %%{ %%}: %q", zsh)
	}

	fish := plainRenderer{}.Render(ShellFish, segs)
	if strings.Contains(fish, "%{") {
		t.Errorf("fish render must NOT use zsh %%{ %%} delimiters: %q", fish)
	}
}

// --- Powerline renderer ------------------------------------------------------

func TestPowerlineHasSeparatorBetweenBlocks(t *testing.T) {
	// N blocks → N separators: one between each pair PLUS the closing cap on the
	// last block. (Closing arrow is part of the powerline look.)
	segs := []Segment{
		{text: "a", fg: "white", bg: "blue"},
		{text: "b", fg: "white", bg: "green"},
		{text: "c", fg: "white", bg: "red"},
	}
	out := powerlineRenderer{}.Render(ShellFish, segs)

	got := strings.Count(out, powerlineRight)
	if got != len(segs) {
		t.Errorf("powerline with %d blocks: got %d separators, want %d\n%q",
			len(segs), got, len(segs), out)
	}
}

func TestPowerlineShowsAllSegmentText(t *testing.T) {
	segs := []Segment{
		{text: "~/dev", fg: "white", bg: "blue"},
		{text: "main", fg: "white", bg: "green"},
	}
	out := powerlineRenderer{}.Render(ShellFish, segs)

	for _, seg := range segs {
		if !strings.Contains(out, seg.text) {
			t.Errorf("powerline output missing segment text %q\ngot: %q", seg.text, out)
		}
	}
}

func TestPowerlineSkipsBackgroundlessSegments(t *testing.T) {
	// Segments without a bg are not part of the powerline chain (the prompt
	// symbol is handled separately). They contribute no separator.
	segs := []Segment{
		{text: "block", fg: "white", bg: "blue"},
		{text: "loose", fg: "green"}, // no bg
	}
	out := powerlineRenderer{}.Render(ShellFish, segs)

	// One block → one separator (its closing cap), not two.
	if got := strings.Count(out, powerlineRight); got != 1 {
		t.Errorf("expected 1 separator for 1 block, got %d: %q", got, out)
	}
}

// --- Symbol layout (renderPrompt orchestration) ------------------------------

func TestSymbolLayoutNewLine(t *testing.T) {
	cfg := defaultConfig()
	cfg.SymbolOnNewLine = true
	out := renderPrompt(ShellFish, 0, 0, expensiveNone, "", cfg)

	if !strings.Contains(out, "\n") {
		t.Errorf("symbol_on_new_line should put the symbol on its own line: %q", out)
	}
	if !strings.Contains(out, cfg.Symbol.Char) {
		t.Errorf("symbol char %q missing: %q", cfg.Symbol.Char, out)
	}
}

func TestSymbolOmittedInPowerlineInline(t *testing.T) {
	cfg := defaultConfig()
	cfg.Style = "powerline"
	cfg.SymbolOnNewLine = false
	cfg.Symbol.Char = "❯"
	out := renderPrompt(ShellFish, 0, 0, expensiveNone, "", cfg)

	// Powerline inline ends on the closing arrow; the symbol is omitted.
	if strings.Contains(out, "❯") {
		t.Errorf("powerline inline should omit the prompt symbol: %q", out)
	}
}
