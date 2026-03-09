package generator

import (
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/axiahq/winz/internal/filesystem"
)

const templatesRoot = "internal/templates"

type Generator struct {
	templates fs.FS
}

func New(templates fs.FS) *Generator {
	return &Generator{templates: templates}
}

// ListTemplates returns all available template paths that contain a README.md.
func (g *Generator) ListTemplates() ([]string, error) {
	candidates := map[string]bool{}

	err := fs.WalkDir(g.templates, templatesRoot, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// The README.md acts as our marker for a valid selectable project
		if d.IsDir() || path.Base(p) != "README.md" {
			return nil
		}
		dir := path.Dir(p)
		rel := strings.TrimPrefix(dir, templatesRoot+"/")
		if rel != "" {
			candidates[rel] = true
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	out := make([]string, 0, len(candidates))
	for name := range candidates {
		out = append(out, name)
	}
	sort.Strings(out)
	return out, nil
}

// Generate safely copies the selected lab/project into the targetPath.
func (g *Generator) Generate(templateName string, targetPath string) error {
	// Clean the incoming name to ensure it plays nicely with io/fs
	templateName = path.Clean(filepath.ToSlash(templateName))
	srcRoot := path.Join(templatesRoot, templateName)

	if _, err := fs.Stat(g.templates, srcRoot); err != nil {
		return fmt.Errorf("template %q not found", templateName)
	}

	if err := filesystem.EnsureDir(targetPath); err != nil {
		return err
	}

	return fs.WalkDir(g.templates, srcRoot, func(srcPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory itself to prevent pathing issues
		if srcPath == srcRoot {
			return nil
		}

		// Calculate a clean relative path based on the root of the selected template
		relPath := srcPath[len(srcRoot)+1:]

		// Safely translate the fs.FS path to the local OS path
		destPath := filepath.Join(targetPath, filepath.FromSlash(relPath))

		if d.IsDir() {
			return filesystem.EnsureDir(destPath)
		}

		// Read and copy verbatim (no text/template parsing)
		content, err := fs.ReadFile(g.templates, srcPath)
		if err != nil {
			return err
		}

		return filesystem.WriteFileSafe(destPath, content, 0o644)
	})
}
