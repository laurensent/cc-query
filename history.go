package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const timeFormat = "2006-01-02 15:04:05"

func historyPath() string {
	return filepath.Join(dataDir(), "history")
}

// saveHistory appends a query to the history file.
func saveHistory(query string) {
	if query == "" {
		return
	}
	p := historyPath()
	_ = os.MkdirAll(filepath.Dir(p), 0755)
	f, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	escaped := strings.ReplaceAll(query, "\\", "\\\\")
	escaped = strings.ReplaceAll(escaped, "\n", "\\n")
	escaped = strings.ReplaceAll(escaped, "\t", "\\t")
	fmt.Fprintf(f, "%s\t%s\n", time.Now().Format(timeFormat), escaped)
}

type historyEntry struct {
	Time  string
	Query string
}

func loadHistory() ([]historyEntry, error) {
	f, err := os.Open(historyPath())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	var entries []historyEntry
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) != 2 {
			continue
		}
		query := parts[1]
		query = strings.ReplaceAll(query, "\\t", "\t")
		query = strings.ReplaceAll(query, "\\n", "\n")
		query = strings.ReplaceAll(query, "\\\\", "\\")
		entries = append(entries, historyEntry{Time: parts[0], Query: query})
	}
	return entries, scanner.Err()
}

// --- bubbletea history browser ---

// historyItem implements list.Item for the bubbles/list component.
type historyItem struct {
	entry historyEntry
}

func (i historyItem) Title() string       { return firstLine(i.entry.Query, 80) }
func (i historyItem) Description() string { return i.entry.Time }
func (i historyItem) FilterValue() string { return i.entry.Query }

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#E36C38")).
			MarginLeft(1)
)

type historyBrowser struct {
	list     list.Model
	selected *historyEntry
}

func (m historyBrowser) Init() tea.Cmd {
	return nil
}

func (m historyBrowser) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Don't intercept keys when filtering is active
		if m.list.FilterState() == list.Filtering {
			break
		}
		switch msg.Type {
		case tea.KeyEnter:
			if item, ok := m.list.SelectedItem().(historyItem); ok {
				m.selected = &item.entry
			}
			return m, tea.Quit
		case tea.KeyEsc, tea.KeyCtrlC:
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m historyBrowser) View() string {
	return m.list.View()
}

// interactiveHistory opens an interactive list for browsing and re-running past queries.
func interactiveHistory() error {
	entries, err := loadHistory()
	if err != nil {
		return fmt.Errorf("failed to read history: %w", err)
	}
	if len(entries) == 0 {
		fmt.Println("No history yet.")
		return nil
	}

	// Most recent first
	items := make([]list.Item, len(entries))
	for i, e := range entries {
		items[len(entries)-1-i] = historyItem{entry: e}
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color("#E36C38"))
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(lipgloss.Color("#E36C38"))

	l := list.New(items, delegate, 0, 0)
	l.Title = "Query History"
	l.Styles.Title = titleStyle
	l.SetFilteringEnabled(true)
	l.SetShowStatusBar(true)
	l.SetStatusBarItemName("query", "queries")

	m := historyBrowser{list: l}
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithOutput(os.Stderr))
	result, err := p.Run()
	if err != nil {
		return err
	}

	final := result.(historyBrowser)
	if final.selected == nil {
		return nil
	}

	saveHistory(final.selected.Query)
	if cfg.Mode == "api" {
		return runAPI(final.selected.Query, model, cfg)
	}
	return runClaude(final.selected.Query, model)
}

// clearHistory removes the history file.
func clearHistory() error {
	p := historyPath()
	if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to clear history: %w", err)
	}
	fmt.Println("History cleared.")
	return nil
}

func firstLine(s string, maxLen int) string {
	line := strings.SplitN(s, "\n", 2)[0]
	if len(line) > maxLen {
		return line[:maxLen-3] + "..."
	}
	return line
}
