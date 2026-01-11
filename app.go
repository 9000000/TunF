package main

import (
	"context"
	"fmt"

	"github.com/energye/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx             context.Context
	proxyService    *ProxyService
	refreshTrayChan chan bool
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		proxyService:    NewProxyService(),
		refreshTrayChan: make(chan bool, 1),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	go a.setupTray()
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	systray.Quit()
}

// GetConfig returns the saved configuration
func (a *App) GetConfig() Config {
	config := LoadConfig()
	// Sync with actual registry state
	config.AutoStart = IsAutoStartEnabled()
	return config
}

// SetAutoStart enables or disables auto-start with Windows
func (a *App) SetAutoStart(enable bool) string {
	err := SetAutoStart(enable)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}

	config := LoadConfig()
	config.AutoStart = enable
	SaveConfig(config)

	return "Success"
}

// StartProxy starts the TCP proxy and saves the successful config
func (a *App) StartProxy(listenPort string, targetAddr string, manageFirewall bool) string {
	listenAddr := fmt.Sprintf(":%s", listenPort)
	err := a.proxyService.Start(listenAddr, targetAddr, manageFirewall)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}

	// Persist successful config
	config := LoadConfig()
	config.LastListenPort = listenPort
	config.LastTargetAddr = targetAddr
	config.AutoOpenFirewall = manageFirewall
	config.History = AddToHistory(config.History, listenPort, targetAddr)
	config.TargetHistory = AddValueToHistory(config.TargetHistory, targetAddr)
	config.ProxyPortHistory = AddValueToHistory(config.ProxyPortHistory, listenPort)
	SaveConfig(config)

	// Notify tray to refresh history
	select {
	case a.refreshTrayChan <- true:
	default:
	}

	runtime.EventsEmit(a.ctx, "proxy-state-change", true)
	return "Success"
}

// StopProxy stops the TCP proxy
func (a *App) StopProxy() {
	a.proxyService.Stop()
	// Notify tray to update status
	select {
	case a.refreshTrayChan <- true:
	default:
	}
	runtime.EventsEmit(a.ctx, "proxy-state-change", false)
}

// IsProxyRunning returns true if the proxy is running
func (a *App) IsProxyRunning() bool {
	return a.proxyService.IsRunning()
}
