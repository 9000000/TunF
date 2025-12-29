//go:build !windows

package main

import "fmt"

func SetAutoStart(enable bool) error {
	// TODO: Implement for Linux/macOS
	return fmt.Errorf("auto-start is not supported on this platform yet")
}

func IsAutoStartEnabled() bool {
	return false
}
