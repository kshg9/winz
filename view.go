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
	modeStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#9D7CD8")).Bold(true)
)

func (m model) View() string {
	if m.quitting {
		return "bye\n"
	}

	modeText := "NORMAL"
	if m.mode == modeSearch {
		modeText = "SEARCH"
	}

	var b strings.Builder
	b.WriteString(titleStyle.Render(" " + appName + " "))
	b.WriteString("\n")
	b.WriteString(subStyle.Render(subtitle + "  •  " + modeStyle.Render(modeText)))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("Search: %s\n\n", m.query))

	if len(m.filtered) == 0 {
		b.WriteString(normal.Render("  no matches"))
		b.WriteString("\n")
	} else {
		start, end := m.visibleRange()
		for i := start; i < end; i++ {
			item := m.filtered[i]
			line := "  " + item
			if i == m.cursor {
				b.WriteString(selected.Render("▶ " + item))
			} else {
				b.WriteString(normal.Render(line))
			}
			b.WriteString("\n")
		}
		if len(m.filtered) > (end - start) {
			b.WriteString(subStyle.Render(fmt.Sprintf("showing %d-%d of %d", start+1, end, len(m.filtered))))
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

	b.WriteString(hint.Render("normal: j/k navigate • enter init • f search • q quit | search: type • esc normal"))
	return b.String()
}
