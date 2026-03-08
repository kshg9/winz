//go:build windows

package main

import (
	"fmt"
	"os"
	"strings"
)

func runUninstall() {
	exe, err := os.Executable()
	if err != nil {
		fmt.Println("could not find own path:", err)
		return
	}

	// Strip install dir from user PATH
	installDir := exe[:strings.LastIndex(exe, `\`)]
	removeFromPath(installDir)

	// Remove the binary itself (Windows won't let you delete a running exe,
	// so we schedule deletion on next reboot via a temp batch file)
	batchPath := os.TempDir() + `\demotool_uninstall.bat`
	batch := fmt.Sprintf(`@echo off
ping -n 2 127.0.0.1 > nul
del /f /q "%s"
del /f /q "%%~f0"
`, exe)

	if err := os.WriteFile(batchPath, []byte(batch), 0644); err != nil {
		fmt.Println("could not write uninstall batch:", err)
		return
	}

	// Launch the batch detached so it runs after we exit
	// (uses cmd /c start so it doesn't block)
	fmt.Println("demotool removed from PATH.")
	fmt.Println("Binary will be deleted on next login.")
}

func removeFromPath(dir string) {
	key := `HKCU\Environment`
	// Read current user PATH via reg query
	out, err := runCmd("reg", "query", key, "/v", "PATH")
	if err != nil {
		return
	}

	// Parse the value out of reg query output
	lines := strings.Split(out, "\n")
	var currentPath string
	for _, l := range lines {
		if strings.Contains(l, "PATH") && strings.Contains(l, `\`) {
			parts := strings.SplitN(strings.TrimSpace(l), "    ", 3)
			if len(parts) == 3 {
				currentPath = strings.TrimSpace(parts[2])
			}
		}
	}

	if currentPath == "" {
		return
	}

	parts := strings.Split(currentPath, ";")
	filtered := parts[:0]
	for _, p := range parts {
		if strings.TrimSpace(p) != dir {
			filtered = append(filtered, p)
		}
	}

	newPath := strings.Join(filtered, ";")
	runCmd("reg", "add", key, "/v", "PATH", "/t", "REG_EXPAND_SZ", "/d", newPath, "/f") //nolint
}
