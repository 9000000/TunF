package main

import (
	"embed"
	"os"
	"syscall"
	"unsafe"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

const mutexName = "TunF_SingleInstance_Mutex"

var (
	kernel32      = syscall.NewLazyDLL("kernel32.dll")
	user32        = syscall.NewLazyDLL("user32.dll")
	createMutex   = kernel32.NewProc("CreateMutexW")
	findWindow    = user32.NewProc("FindWindowW")
	showWindow    = user32.NewProc("ShowWindow")
	setForeground = user32.NewProc("SetForegroundWindow")
)

const (
	SW_RESTORE = 9
)

func isAlreadyRunning() bool {
	mutexNamePtr, _ := syscall.UTF16PtrFromString(mutexName)
	handle, _, err := createMutex.Call(0, 0, uintptr(unsafe.Pointer(mutexNamePtr)))

	if handle == 0 {
		return true
	}

	// ERROR_ALREADY_EXISTS = 183
	if err.(syscall.Errno) == 183 {
		return true
	}

	return false
}

func showExistingWindow() {
	windowName, _ := syscall.UTF16PtrFromString("TunF")
	hwnd, _, _ := findWindow.Call(0, uintptr(unsafe.Pointer(windowName)))

	if hwnd != 0 {
		showWindow.Call(hwnd, SW_RESTORE)
		setForeground.Call(hwnd)
	}
}

func main() {
	// Check if already running
	if isAlreadyRunning() {
		showExistingWindow()
		return
	}

	// Create an instance of the app structure
	app := NewApp()

	// Parse command line arguments
	startHidden := false
	for _, arg := range os.Args {
		if arg == "--hidden" {
			startHidden = true
			break
		}
	}

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "TunF",
		Width:  420,
		Height: 780,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 15, G: 23, B: 42, A: 1},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		Bind: []interface{}{
			app,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			BackdropType:         windows.Mica,
		},
		HideWindowOnClose: true,
		StartHidden:       startHidden,
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
