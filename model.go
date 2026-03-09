package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/axiahq/winz/internal/generator"
	tea "github.com/charmbracelet/bubbletea"
)

type mode int

const (
	modeNormal mode = iota
	modeSearch
)

// model represents the state of the TUI.  We track a handful of fields
// used by the bubbletea engine, plus the list of available templates and
// what subset of them the user can currently see.
type model struct {
	generator    *generator.Generator
	all          []string
	filtered     []string
	cursor       int
	query        string
	status       string
	statusErr    bool
	quitting     bool
	windowHeight int
	mode         mode
}

func initialModel(g *generator.Generator) model {
	items, err := g.ListTemplates()
	m := model{generator: g, all: items, filtered: items}
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
		// handle search mode transitions first
		if m.mode == modeNormal {
			switch msg.String() {
			case "f":
				m.mode = modeSearch
				m.query = ""
				m.filtered = append([]string(nil), m.all...)
				m.cursor = 0
				return m, nil
			}
		} else if m.mode == modeSearch {
			if msg.Type == tea.KeyEsc {
				m.mode = modeNormal
				m.query = ""
				m.filtered = append([]string(nil), m.all...)
				m.cursor = 0
				return m, nil
			}
		}

		switch msg.String() {
		case "ctrl+c":
			// always quit regardless of mode
			m.quitting = true
			return m, tea.Quit
		case "q":
			// only quit in normal mode
			if m.mode == modeNormal {
				m.quitting = true
				return m, tea.Quit
			}
			if m.mode == modeSearch && msg.Type == tea.KeyRunes {
				m.query += msg.String()
				m.refilter()
			}
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
		case "backspace":
			if len(m.query) > 0 {
				m.query = m.query[:len(m.query)-1]
				m.refilter()
			}
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
		default:
			if m.mode == modeSearch && msg.Type == tea.KeyRunes {
				m.query += msg.String()
				m.refilter()
			}
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
