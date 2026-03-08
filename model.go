package main

import tea "github.com/charmbracelet/bubbletea"

type model struct {
	cursor   int
	result   string
	isDanger bool
}

func initialModel() model {
	return model{cursor: 0}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "j", "down":
			if m.cursor < len(options)-1 {
				m.cursor++
			}

		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}

		case "g":
			m.cursor = 0

		case "G":
			m.cursor = len(options) - 1

		case "enter", " ":
			opt := options[m.cursor]
			m.result = opt.result
			m.isDanger = opt.danger
			if opt.danger {
				runUninstall()
				// return m, tea.Quit
			}

		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}
