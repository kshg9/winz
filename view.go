package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	appName  = "Winz Lab Generator"
	subtitle = "Scaffold college lab boilerplate instantly"
)

var (
	titleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FAFAFA")).Background(lipgloss.Color("#7B2FBE")).Padding(0, 2)
	subStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#888")).MarginBottom(1)
	selected    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7B2FBE"))
	normal      = lipgloss.NewStyle().Foreground(lipgloss.Color("#BDBDBD"))
	searchStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Bold(true)
	hint        = lipgloss.NewStyle().Foreground(lipgloss.Color("#666")).MarginTop(1)
	success     = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).MarginTop(1)
	errorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5F57")).MarginTop(1)
)

func (m model) View() string {
	if m.quitting {
		return "\n"
	}

	var b strings.Builder
	b.WriteString(titleStyle.Render(" " + appName + " "))
	b.WriteString("\n")
	b.WriteString(subStyle.Render(subtitle))
	b.WriteString("\n")

	// 1. Render Search / Filter Header
	if m.mode == modeSearch {
		b.WriteString(searchStyle.Render(fmt.Sprintf("Search: %s█", m.query)))
		b.WriteString("\n\n")
	} else if m.query != "" {
		b.WriteString(hint.Render(fmt.Sprintf("Filtered by: '%s'", m.query)))
		b.WriteString("\n\n")
	} else {
		b.WriteString("\n\n")
	}

	// 2. Render the Paginated List
	if len(m.filtered) == 0 {
		b.WriteString(normal.Render("  no matches"))
		b.WriteString("\n")
	} else {
		// Calculate available vertical space (Subtract ~10 lines for headers/footers)
		listHeight := m.windowHeight - 10
		if listHeight < 5 {
			listHeight = 5 // Fallback for very small terminals or before first WindowSizeMsg
		}

		// Calculate current page bounds based on cursor position
		page := m.cursor / listHeight
		start := page * listHeight
		end := start + listHeight
		if end > len(m.filtered) {
			end = len(m.filtered)
		}

		// Render only the items on the current page
		for i := start; i < end; i++ {
			item := m.filtered[i]
			if i == m.cursor {
				b.WriteString(selected.Render("▶ " + item))
			} else {
				b.WriteString(normal.Render("  " + item))
			}
			b.WriteString("\n")
		}

		// Show a pagination indicator if there are hidden items
		if len(m.filtered) > listHeight {
			totalPages := ((len(m.filtered) - 1) / listHeight) + 1
			progress := fmt.Sprintf("  --- Page %d of %d ---", page+1, totalPages)
			b.WriteString(hint.Render(progress))
			b.WriteString("\n")
		} else {
			b.WriteString("\n") // Maintain consistent spacing
		}
	}

	// 3. Render Status
	if m.status != "" {
		if m.statusErr {
			b.WriteString(errorStyle.Render("⚠ " + m.status))
		} else {
			b.WriteString(success.Render("✓ " + m.status))
		}
		b.WriteString("\n")
	} else {
		b.WriteString("\n")
	}

	// 4. Render Footer Hints
	if m.mode == modeNormal {
		b.WriteString(hint.Render("f or / search • j/k move • enter init • q quit"))
	} else {
		b.WriteString(hint.Render("esc stop searching • enter init"))
	}

	return b.String()
}
