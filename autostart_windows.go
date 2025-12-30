//go:build windows

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

const (
	AppName        = "TunF"
	RunRegistryKey = `Software\Microsoft\Windows\CurrentVersion\Run`
)

// SetAutoStart enables or disables the application auto-start on Windows
func SetAutoStart(enable bool) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, RunRegistryKey, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()

	if enable {
		execPath, err := os.Executable()
		if err != nil {
			return err
		}
		// Convert to absolute path just in case
		absPath, err := filepath.Abs(execPath)
		if err != nil {
			return err
		}
		// Add --hidden flag for autostart
		cmd := fmt.Sprintf("\"%s\" --hidden", absPath)
		return k.SetStringValue(AppName, cmd)
	} else {
		// Try to delete, ignore error if it doesn't exist
		_ = k.DeleteValue(AppName)
		return nil
	}
}

// IsAutoStartEnabled checks if the application is set to auto-start
func IsAutoStartEnabled() bool {
	k, err := registry.OpenKey(registry.CURRENT_USER, RunRegistryKey, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer k.Close()

	_, _, err = k.GetStringValue(AppName)
	return err == nil
}
