package generator

import (
	"bytes"
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/axiahq/winz/internal/filesystem"
)

const templatesRoot = "internal/templates"

// Generator copies and renders embedded templates into a target directory.
type Generator struct {
	templates fs.FS
}

func New(templates fs.FS) *Generator {
	return &Generator{templates: templates}
}

// ListTemplates returns all available template paths.
func (g *Generator) ListTemplates() ([]string, error) {
	candidates := map[string]bool{}

	err := fs.WalkDir(g.templates, templatesRoot, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
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

// Generate copies the selected template into targetPath.
func (g *Generator) Generate(templateName string, targetPath string) error {
	templateName = strings.TrimPrefix(filepath.ToSlash(templateName), "/")
	srcRoot := path.Join(templatesRoot, templateName)

	if _, err := fs.Stat(g.templates, srcRoot); err != nil {
		return fmt.Errorf("template %q not found", templateName)
	}

	if err := filesystem.EnsureDir(targetPath); err != nil {
		return err
	}

	data := map[string]any{
		"Template": templateName,
		"Year":     time.Now().Year(),
	}

	return fs.WalkDir(g.templates, srcRoot, func(srcPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel := strings.TrimPrefix(srcPath, srcRoot+"/")
		if rel == srcPath || rel == "" {
			return nil
		}

		destRel := rel
		if strings.HasSuffix(destRel, ".tmpl") {
			destRel = strings.TrimSuffix(destRel, ".tmpl")
		}
		destPath := filepath.Join(targetPath, filepath.FromSlash(destRel))

		if d.IsDir() {
			return filesystem.EnsureDir(destPath)
		}

		content, err := fs.ReadFile(g.templates, srcPath)
		if err != nil {
			return err
		}

		if strings.HasSuffix(srcPath, ".tmpl") {
			tpl, err := template.New(path.Base(srcPath)).Parse(string(content))
			if err != nil {
				return err
			}
			var buf bytes.Buffer
			if err := tpl.Execute(&buf, data); err != nil {
				return err
			}
			content = buf.Bytes()
		}

		return filesystem.WriteFileSafe(destPath, content, 0o644)
	})
}
