package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// isPiped returns true if stdin is not a terminal (i.e., data is being piped in).
func isPiped() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) == 0
}

// isStdoutTerminal returns true if stdout is connected to a terminal.
func isStdoutTerminal() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

// readPipe reads all data from stdin and returns it as a string.
func readPipe() (string, error) {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", fmt.Errorf("failed to read stdin: %w", err)
	}
	return string(data), nil
}

// promptModel is a bubbletea model for interactive single-line input.
type promptModel struct {
	input    textinput.Model
	result   string
	canceled bool
}

func (m promptModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m promptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.result = m.input.Value()
			return m, tea.Quit
		case tea.KeyEsc, tea.KeyCtrlC:
			m.canceled = true
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m promptModel) View() string {
	return "> " + m.input.View()
}

// errCanceled is returned when the user cancels interactive input (ESC / Ctrl+C).
var errCanceled = fmt.Errorf("canceled")

// readInteractivePrompt reads a single line using bubbletea with full
// emacs keybindings, ESC to cancel, and proper terminal handling.
// Returns (input, nil) on success, or ("", errCanceled) when canceled.
func readInteractivePrompt() (string, error) {
	ti := textinput.New()
	ti.Focus()
	ti.Prompt = ""

	p := tea.NewProgram(promptModel{input: ti}, tea.WithOutput(os.Stderr))
	result, err := p.Run()
	if err != nil {
		return "", err
	}

	final := result.(promptModel)
	if final.canceled {
		return "", errCanceled
	}
	return strings.TrimSpace(final.result), nil
}
