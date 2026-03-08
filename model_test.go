package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestSearchModeKeyHandling(t *testing.T) {
	m := model{all: []string{"sql_lab/exp1", "web_lab/exp1"}, filtered: []string{"sql_lab/exp1", "web_lab/exp1"}, mode: modeNormal}

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
	m1 := updated.(model)
	if m1.mode != modeSearch {
		t.Fatalf("expected search mode after f")
	}

	updated, _ = m1.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	m2 := updated.(model)
	if m2.quitting {
		t.Fatalf("q in search mode should not quit")
	}
	if m2.query != "q" {
		t.Fatalf("expected query to collect runes in search mode, got %q", m2.query)
	}

	updated, _ = m2.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m3 := updated.(model)
	if m3.mode != modeNormal {
		t.Fatalf("expected esc to return to normal mode")
	}
}
