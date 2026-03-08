package filesystem

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// EnsureDir creates a directory tree if it does not already exist.
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0o755)
}

// NormalizeNewlines converts LF line endings to CRLF on Windows.
func NormalizeNewlines(content []byte) []byte {
	if runtime.GOOS != "windows" {
		return content
	}

	s := strings.ReplaceAll(string(content), "\r\n", "\n")
	s = strings.ReplaceAll(s, "\n", "\r\n")
	return []byte(s)
}

// WriteFileSafe writes content atomically after ensuring parent directories exist.
func WriteFileSafe(path string, content []byte, perm os.FileMode) error {
	if err := EnsureDir(filepath.Dir(path)); err != nil {
		return err
	}

	normalized := NormalizeNewlines(content)
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, normalized, perm); err != nil {
		return err
	}

	return os.Rename(tmp, path)
}
