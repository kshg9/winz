package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/axiahq/winz/internal/generator"
	tea "github.com/charmbracelet/bubbletea"
)

type inputMode int

const (
	modeNormal inputMode = iota
	modeSearch
)

type model struct {
	generator    *generator.Generator
	all          []string
	filtered     []string
	cursor       int
	query        string
	status       string
	statusErr    bool
	quitting     bool
	mode         inputMode
	windowHeight int
}

func initialModel(g *generator.Generator) model {
	items, err := g.ListTemplates()
	m := model{generator: g, all: items, filtered: items, mode: modeNormal}
	if err != nil {
		m.status = err.Error()
		m.statusErr = true
	}
	return m
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowHeight = msg.Height
		return m, nil
	case tea.KeyMsg:
		if m.mode == modeSearch {
			return m.updateSearch(msg)
		}
		return m.updateNormal(msg)
	}

	return m, nil
}

func (m model) updateNormal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		m.quitting = true
		return m, tea.Quit
	case "j", "down":
		if m.cursor < len(m.filtered)-1 {
			m.cursor++
		}
	case "k", "up":
		if m.cursor > 0 {
			m.cursor--
		}
	case "g", "home":
		m.cursor = 0
	case "G", "end":
		if len(m.filtered) > 0 {
			m.cursor = len(m.filtered) - 1
		}
	case "f":
		m.mode = modeSearch
	case "enter":
		if len(m.filtered) == 0 {
			return m, nil
		}
		selected := m.filtered[m.cursor]
		target := lastPathSegment(selected)
		if err := m.generator.Generate(selected, target); err != nil {
			m.status = err.Error()
			m.statusErr = true
		} else {
			m.status = fmt.Sprintf("initialized %s in %s", selected, target)
			m.statusErr = false
		}
	}
	return m, nil
}

func (m model) updateSearch(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = modeNormal
	case "backspace":
		if len(m.query) > 0 {
			m.query = m.query[:len(m.query)-1]
			m.refilter()
		}
	case "ctrl+c":
		m.quitting = true
		return m, tea.Quit
	default:
		if msg.Type == tea.KeyRunes {
			m.query += msg.String()
			m.refilter()
		}
	}
	return m, nil
}

func (m *model) refilter() {
	if strings.TrimSpace(m.query) == "" {
		m.filtered = append([]string(nil), m.all...)
		m.cursor = 0
		return
	}

	type scored struct {
		name  string
		score int
	}
	var matches []scored
	for _, item := range m.all {
		s := fuzzyScore(strings.ToLower(item), strings.ToLower(m.query))
		if s >= 0 {
			matches = append(matches, scored{name: item, score: s})
		}
	}

	sort.Slice(matches, func(i, j int) bool {
		if matches[i].score == matches[j].score {
			return matches[i].name < matches[j].name
		}
		return matches[i].score < matches[j].score
	})

	m.filtered = m.filtered[:0]
	for _, match := range matches {
		m.filtered = append(m.filtered, match.name)
	}
	if m.cursor >= len(m.filtered) {
		m.cursor = len(m.filtered) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
}

func (m model) visibleRange() (start, end int) {
	maxRows := 12
	if m.windowHeight > 0 {
		usable := m.windowHeight - 8
		if usable > 0 {
			maxRows = usable
		}
	}
	if maxRows < 5 {
		maxRows = 5
	}

	if len(m.filtered) <= maxRows {
		return 0, len(m.filtered)
	}

	start = m.cursor - maxRows/2
	if start < 0 {
		start = 0
	}
	end = start + maxRows
	if end > len(m.filtered) {
		end = len(m.filtered)
		start = end - maxRows
	}
	return start, end
}

func fuzzyScore(item, query string) int {
	idx := 0
	score := 0
	for _, ch := range query {
		next := strings.IndexRune(item[idx:], ch)
		if next < 0 {
			return -1
		}
		score += idx + next
		idx += next + 1
	}
	return score
}

func lastPathSegment(s string) string {
	parts := strings.Split(s, "/")
	return parts[len(parts)-1]
}
