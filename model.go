package main

import (
	"fmt"
	"path"
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
		// 1. Global shortcuts (always active)
		switch msg.Type {
		case tea.KeyCtrlC:
			m.quitting = true
			return m, tea.Quit
		case tea.KeyEnter:
			if len(m.filtered) == 0 {
				return m, nil
			}
			selected := m.filtered[m.cursor]
			// Replaced custom function with standard library
			target := path.Base(selected)

			if err := m.generator.Generate(selected, target); err != nil {
				m.status = err.Error()
				m.statusErr = true
			} else {
				m.status = fmt.Sprintf("initialized %s in ./%s", selected, target)
				m.statusErr = false
				// Optional: auto-quit after successful generation
				m.quitting = true
				return m, tea.Quit
			}
			return m, nil
		}

		// 2. Mode-specific handling
		switch m.mode {
		case modeNormal:
			switch msg.String() {
			case "q":
				m.quitting = true
				return m, tea.Quit
			case "f", "/": // '/' is a very common standard for search
				m.mode = modeSearch
				m.query = ""
				m.refilter()
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
			}
		case modeSearch:
			// Check specific types first
			switch msg.Type {
			case tea.KeyEsc:
				m.mode = modeNormal
				m.query = ""
				m.refilter()
			case tea.KeyUp:
				if m.cursor > 0 {
					m.cursor--
				}
			case tea.KeyDown:
				if m.cursor < len(m.filtered)-1 {
					m.cursor++
				}
			case tea.KeyRunes:
				m.query += msg.String()
				m.refilter()
			default:
				// Fallback to string checking for backspace to handle terminal quirks
				// TIL interesting story of backspace quirks
				if msg.String() == "backspace" {
					if len(m.query) > 0 {
						m.query = m.query[:len(m.query)-1]
						m.refilter()
					}
				}
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

	// Tie breaker
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

	// Ensure cursor stays within bounds after filtering
	if m.cursor >= len(m.filtered) {
		m.cursor = max(0, len(m.filtered)-1)
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
