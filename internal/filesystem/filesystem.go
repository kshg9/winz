package filesystem

import (
	"os"
	"path/filepath"
)

// EnsureDir creates a directory tree if it does not already exist.
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0o755)
}

// WriteFileSafe writes content atomically after ensuring parent directories exist.
func WriteFileSafe(path string, content []byte, perm os.FileMode) error {
	if err := EnsureDir(filepath.Dir(path)); err != nil {
		return err
	}

	tmp := path + ".tmp"

	// Ensure we don't leave orphaned .tmp files behind if WriteFile or Rename fails.
	// If Rename succeeds, os.Remove will quietly fail (which is fine) because the file is gone.
	defer os.Remove(tmp)

	if err := os.WriteFile(tmp, content, perm); err != nil {
		return err
	}

	return os.Rename(tmp, path)
}
