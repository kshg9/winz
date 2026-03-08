package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FAFAFA")).Background(lipgloss.Color("#7B2FBE")).Padding(0, 2)
	subStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#888")).MarginBottom(1)
	selected   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7B2FBE"))
	normal     = lipgloss.NewStyle().Foreground(lipgloss.Color("#BDBDBD"))
	hint       = lipgloss.NewStyle().Foreground(lipgloss.Color("#666")).MarginTop(1)
	success    = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).MarginTop(1)
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5F57")).MarginTop(1)
)

func (m model) View() string {
	if m.quitting {
		return "bye\n"
	}

	var b strings.Builder
	b.WriteString(titleStyle.Render(" " + appName + " "))
	b.WriteString("\n")
	b.WriteString(subStyle.Render(subtitle))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("Search: %s\n\n", m.query))

	if len(m.filtered) == 0 {
		b.WriteString(normal.Render("  no matches"))
		b.WriteString("\n")
	} else {
		for i, item := range m.filtered {
			line := "  " + item
			if i == m.cursor {
				b.WriteString(selected.Render("▶ " + item))
			} else {
				b.WriteString(normal.Render(line))
			}
			b.WriteString("\n")
		}
	}

	if m.status != "" {
		if m.statusErr {
			b.WriteString(errorStyle.Render("⚠ " + m.status))
		} else {
			b.WriteString(success.Render("✓ " + m.status))
		}
		b.WriteString("\n")
	}

	b.WriteString(hint.Render("type to fuzzy-search • j/k or ↑/↓ move • enter init • q quit"))
	return b.String()
}
