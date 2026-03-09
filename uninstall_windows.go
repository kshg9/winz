//go:build windows

package main

import (
	"fmt"
	"os"
	"os/exec"
)

func runUninstall() {
	exe, err := os.Executable()
	if err != nil {
		fmt.Println("could not find own path:", err)
		return
	}

	// Create the self-destruct batch file
	batchPath := os.TempDir() + `\winz_uninstall.bat`
	batch := fmt.Sprintf(`@echo off
:: Wait 2 seconds for the Go program to fully exit
ping -n 2 127.0.0.1 > nul
:: Delete the executable
del /f /q "%s"
:: Delete this batch file itself
del /f /q "%%~f0"
`, exe)

	if err := os.WriteFile(batchPath, []byte(batch), 0644); err != nil {
		fmt.Println("could not write uninstall batch:", err)
		return
	}

	// Actually LAUNCH the batch file detached so it runs in the background
	cmd := exec.Command("cmd.exe", "/C", "start", "/b", batchPath)
	if err := cmd.Start(); err != nil {
		fmt.Println("failed to start uninstaller:", err)
		return
	}

	fmt.Println("Winz has been uninstalled. The binary will be deleted in a few seconds.")
	// Exit immediately so the batch file can delete the unlocked .exe
	os.Exit(0)
}
