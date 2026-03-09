package main

import "embed"

// embeddedTemplates stores all scaffold templates directly inside the binary.
//
//go:embed internal/templates
var embeddedTemplates embed.FS
