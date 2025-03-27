package main

import (
	"flag"
	"github.com/jkvatne/jkvgui/button"
	"github.com/jkvatne/jkvgui/callback"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/icon"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
	"log/slog"
	"testing"
	"time"
)

var saveScreen bool

func init() {
	flag.BoolVar(&saveScreen, "test.save", false, "Save the captured screen to ./test-assets, making this the reference image.")
}

func TestButtons(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelError)
	_ = gpu.InitWindow(400, 200, "Test", 0)
	defer gpu.Shutdown()
	callback.Initialize(gpu.Window)
	gpu.BackgroundColor(f32.White)

	// Draw buttons
	wid.Show(10, 10, button.Filled("Primary", icon.Home, nil, button.Role(theme.Primary), ""))
	wid.Show(150, 10, button.Filled("Secondary", icon.Home, nil, button.Role(theme.Secondary), ""))
	wid.Show(300, 10, button.Round(icon.Home, nil, button.DefaultButtonStyle.Role(theme.Secondary), ""))
	wid.Show(10, 50, button.Outline("Outline", nil, nil, button.Role(theme.Surface), ""))
	wid.Show(150, 50, button.Text("Text", nil, nil, button.Role(theme.SurfaceContainer), ""))
	wid.Show(300, 50, button.Filled("", icon.Home, nil, button.Role(theme.Secondary), ""))
	wid.Show(10, 100, button.Filled("Size 1.0", icon.Home, nil, button.Role(theme.Primary).Size(1.0), ""))
	wid.Show(150, 100, button.Filled("Size 2.0", icon.Home, nil, button.Role(theme.Secondary).Size(2.0), ""))
	wid.Show(300, 100, button.Filled("Surface", nil, nil, button.Role(theme.Surface), ""))

	// Verify resulting image
	VerifyScreen(t, "TestButtons", 400, 200, saveScreen)
	// Place breakpoint here in order to look at the screen output.
	time.Sleep(100 * time.Millisecond)

}
