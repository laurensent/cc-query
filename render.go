package main

import (
	"strings"

	"github.com/charmbracelet/glamour"
)

// renderMarkdown renders markdown content to ANSI-styled terminal output.
func renderMarkdown(content string) (string, error) {
	width := getTermWidth()

	var styleOpt glamour.TermRendererOption
	cfg := loadConfig()
	switch cfg.Theme {
	case "dark", "light", "dracula", "pink", "ascii", "notty":
		styleOpt = glamour.WithStandardStyle(cfg.Theme)
	default:
		styleOpt = glamour.WithAutoStyle()
	}

	r, err := glamour.NewTermRenderer(
		styleOpt,
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return "", err
	}
	out, err := r.Render(content)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out) + "\n", nil
}
