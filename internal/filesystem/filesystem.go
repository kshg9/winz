package filesystem

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// EnsureDir creates a directory tree if it does not already exist.
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0o755)
}

// NormalizeNewlines converts LF line endings to CRLF on Windows for text-like files.
func NormalizeNewlines(path string, content []byte) []byte {
	if runtime.GOOS != "windows" || !shouldNormalize(path, content) {
		return content
	}

	s := strings.ReplaceAll(string(content), "\r\n", "\n")
	s = strings.ReplaceAll(s, "\n", "\r\n")
	return []byte(s)
}

func shouldNormalize(path string, content []byte) bool {
	if bytes.IndexByte(content, 0x00) >= 0 {
		return false
	}

	switch strings.ToLower(filepath.Ext(path)) {
	case ".png", ".jpg", ".jpeg", ".gif", ".bmp", ".webp", ".ico", ".pdf", ".zip", ".gz", ".jar", ".exe", ".dll", ".so", ".ttf", ".woff", ".woff2", ".mp3", ".mp4":
		return false
	default:
		return true
	}
}

// WriteFileSafe writes content atomically after ensuring parent directories exist.
func WriteFileSafe(path string, content []byte, perm os.FileMode) error {
	if err := EnsureDir(filepath.Dir(path)); err != nil {
		return err
	}

	normalized := NormalizeNewlines(path, content)
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, normalized, perm); err != nil {
		return err
	}

	return os.Rename(tmp, path)
}
