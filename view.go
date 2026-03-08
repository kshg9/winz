package main

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7B2FBE")).
			Padding(0, 2).
			MarginBottom(1)

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7B2FBE")).
			PaddingLeft(2)

	descStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			PaddingLeft(6)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#AAAAAA")).
			PaddingLeft(4)

	resultStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#04B575")).
			PaddingLeft(2).
			MarginTop(1)

	dangerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF5F57")).
			PaddingLeft(2).
			MarginTop(1)

	hintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#555555")).
			MarginTop(1).
			PaddingLeft(2)
)

func (m model) View() string {
	s := titleStyle.Render("  demotool  ") + "\n\n"

	for i, opt := range options {
		if m.cursor == i {
			s += selectedStyle.Render("▶ "+opt.label) + "\n"
			s += descStyle.Render(opt.desc) + "\n"
		} else {
			s += normalStyle.Render("  "+opt.label) + "\n"
		}
	}

	if m.result != "" {
		if m.isDanger {
			s += dangerStyle.Render("⚠  "+m.result) + "\n"
		} else {
			s += resultStyle.Render("→  "+m.result) + "\n"
		}
	}

	s += hintStyle.Render("j/k  ↑/↓  move   •   g/G  top/bottom   •   enter  select   •   q  quit")

	return s
}
