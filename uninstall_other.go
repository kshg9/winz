//go:build !windows

package main

import "fmt"

func runUninstall() {
	fmt.Println("boom, not implemented on this platform")
}
