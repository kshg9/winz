package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	if len(args) == 0 {
		p := tea.NewProgram(initialModel(g), tea.WithAltScreen())
		_, err := p.Run()
		return err
	}

	switch args[0] {

	/*
		List all the templates
	*/
	case "list":
		templates, err := g.ListTemplates()
		if err != nil {
			return err
		}
		for _, t := range templates {
			fmt.Println(t)
		}
		return nil

	/*
		Initialization function for templates
	*/
	case "init":
		if len(args) < 2 {
			return fmt.Errorf("usage: %s init <query> [target]", filepath.Base(os.Args[0]))
		}

		// The user's search query (e.g., "senso interf" or "typescript")
		query := args[1]

		templates, err := g.ListTemplates()
		if err != nil {
			return err
		}

		matches := findMatches(query, templates)

		// Handle the 3 possible outcomes: Zero, Multiple, or Exact match
		if len(matches) == 0 {
			return fmt.Errorf("no templates matched the query %q", query)
		}

		if len(matches) > 1 {
			// Collision detected! Print the helpful error.
			errMsg := fmt.Sprintf("Ambiguous query %q. Did you mean:\n", query)
			for _, m := range matches {
				errMsg += fmt.Sprintf("  - %s\n", m)
			}
			errMsg += "Please add more keywords to narrow it down."
			return fmt.Errorf("%s", errMsg)
		}

		// Exactly 1 match found! Proceed.
		templateName := matches[0]

		// Set the target directory (defaulting to the final folder name)
		target := filepath.Base(templateName)
		if len(args) >= 3 {
			target = args[2]
		}

		if err := g.Generate(templateName, target); err != nil {
			return err
		}
		fmt.Printf("Initialized %s in %s\n", templateName, target)
		return nil

	/*
		Handle Uninstall case
	*/
	case "uninstall":
		runUninstall()
		return nil

	default:
		return fmt.Errorf("unknown command %q (available: list, init, tui, uninstall)", args[0])
	}
}

// findMatches searches templates using space-separated keywords.
func findMatches(query string, templates []string) []string {
	queryTerms := strings.Fields(strings.ToLower(query))
	var matches []string

	for _, t := range templates {
		normalizedPath := strings.ToLower(t)
		normalizedPath = strings.ReplaceAll(normalizedPath, "/", " ")
		normalizedPath = strings.ReplaceAll(normalizedPath, "_", " ")
		normalizedPath = strings.ReplaceAll(normalizedPath, "-", " ")

		// Check if ALL typed terms exist in the normalized path
		isMatch := true
		for _, term := range queryTerms {
			if !strings.Contains(normalizedPath, term) {
				isMatch = false
				break
			}
		}

		if isMatch {
			matches = append(matches, t)
		}
	}
	return matches
}
