package main

import (
	"os"
	"strings"
	"time"

	"github.com/energye/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type trayHistoryItem struct {
	menuItem *systray.MenuItem
	port     string
	target   string
}

func (a *App) setupTray() {
	systray.Run(a.onTrayReady, a.onTrayExit)
}

func (a *App) onTrayReady() {
	if len(iconData) > 0 {
		systray.SetIcon(iconData)
	}

	systray.SetTitle("TunF")
	systray.SetTooltip("TunF - Port Forwarding")

	mToggle := systray.AddMenuItem("Start Proxy", "Toggle Proxy ON/OFF")
	mShow := systray.AddMenuItem("Show Window", "Show the main window")

	systray.AddSeparator()

	mRecent := systray.AddMenuItem("Recent Connections", "Quickly fill from history")
	historyItems := make([]*trayHistoryItem, 5)
	for i := 0; i < 5; i++ {
		item := mRecent.AddSubMenuItem("-", "")
		item.Hide()
		historyItems[i] = &trayHistoryItem{menuItem: item}
	}

	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit the application")

	// Update Tray Menu Status Helper
	updateStatus := func() {
		if a.proxyService.IsRunning() {
			mToggle.SetTitle("Stop Proxy")
			if len(iconActiveData) > 0 {
				systray.SetIcon(iconActiveData)
			}
		} else {
			mToggle.SetTitle("Start Proxy")
			if len(iconData) > 0 {
				systray.SetIcon(iconData)
			}
		}

		// Update History Submenu
		config := LoadConfig()
		for i := 0; i < 5; i++ {
			if i < len(config.History) {
				parts := strings.Split(config.History[i], "|")
				if len(parts) == 2 {
					port, target := parts[0], parts[1]
					historyItems[i].port = port
					historyItems[i].target = target
					historyItems[i].menuItem.SetTitle(port + " -> " + target)
					historyItems[i].menuItem.Show()
				}
			} else {
				historyItems[i].menuItem.Hide()
			}
		}
	}

	// Toggle proxy on left-click
	systray.SetOnClick(func(menu systray.IMenu) {
		if a.proxyService.IsRunning() {
			a.StopProxy()
		} else {
			config := LoadConfig()
			a.StartProxy(config.LastListenPort, config.LastTargetAddr, config.AutoOpenFirewall)
		}
		updateStatus()
	})

	// Menu Item Click Handlers
	mToggle.Click(func() {
		if a.proxyService.IsRunning() {
			a.StopProxy()
		} else {
			config := LoadConfig()
			a.StartProxy(config.LastListenPort, config.LastTargetAddr, config.AutoOpenFirewall)
		}
		updateStatus()
	})

	mShow.Click(func() {
		runtime.WindowShow(a.ctx)
	})

	mQuit.Click(func() {
		systray.Quit()
		runtime.Quit(a.ctx)
		// Force exit if it doesn't close in 1 second
		go func() {
			time.Sleep(1 * time.Second)
			os.Exit(0)
		}()
	})

	// History items click handlers
	for i := 0; i < 5; i++ {
		idx := i
		historyItems[idx].menuItem.Click(func() {
			a.quickStart(historyItems[idx], updateStatus)
		})
	}

	// Handle refresh channel in background
	go func() {
		for range a.refreshTrayChan {
			updateStatus()
		}
	}()

	updateStatus()
}

func (a *App) quickStart(item *trayHistoryItem, refresh func()) {
	if a.proxyService.IsRunning() {
		a.StopProxy()
	}
	config := LoadConfig()
	a.StartProxy(item.port, item.target, config.AutoOpenFirewall)
	refresh()
}

func (a *App) onTrayExit() {}
