package main

import (
	"strings"

	"github.com/getlantern/systray"
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
		} else {
			mToggle.SetTitle("Start Proxy")
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

	// Main Tray Loop
	go func() {
		for {
			select {
			case <-mToggle.ClickedCh:
				if a.proxyService.IsRunning() {
					a.StopProxy()
				} else {
					config := LoadConfig()
					a.StartProxy(config.LastListenPort, config.LastTargetAddr, config.AutoOpenFirewall)
				}
				updateStatus()

			case <-mShow.ClickedCh:
				runtime.WindowShow(a.ctx)

			case <-mQuit.ClickedCh:
				systray.Quit()
				runtime.Quit(a.ctx)

			case <-a.refreshTrayChan:
				updateStatus()

			case <-historyItems[0].menuItem.ClickedCh:
				a.quickStart(historyItems[0], updateStatus)
			case <-historyItems[1].menuItem.ClickedCh:
				a.quickStart(historyItems[1], updateStatus)
			case <-historyItems[2].menuItem.ClickedCh:
				a.quickStart(historyItems[2], updateStatus)
			case <-historyItems[3].menuItem.ClickedCh:
				a.quickStart(historyItems[3], updateStatus)
			case <-historyItems[4].menuItem.ClickedCh:
				a.quickStart(historyItems[4], updateStatus)
			}
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
