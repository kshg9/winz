package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/axiahq/winz/internal/generator"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	g := generator.New(embeddedTemplates)
	if err := runCLI(g, os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func runCLI(g *generator.Generator, args []string) error {
	if len(args) == 0 || args[0] == "tui" {
		p := tea.NewProgram(initialModel(g), tea.WithAltScreen())
		_, err := p.Run()
		return err
	}

	switch args[0] {
	case "list":
		templates, err := g.ListTemplates()
		if err != nil {
			return err
		}
		for _, t := range templates {
			fmt.Println(t)
		}
		return nil
	case "init":
		if len(args) < 2 {
			return fmt.Errorf("usage: %s init <template> [target]", filepath.Base(os.Args[0]))
		}
		templateName := args[1]
		target := filepath.Base(templateName)
		if len(args) >= 3 {
			target = args[2]
		}
		if err := g.Generate(templateName, target); err != nil {
			return err
		}
		fmt.Printf("initialized %s in %s\n", templateName, target)
		return nil
	case "uninstall":
		runUninstall()
		return nil
	default:
		return fmt.Errorf("unknown command %q (available: list, init, tui, uninstall)", args[0])
	}
}
