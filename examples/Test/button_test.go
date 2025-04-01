package main

import (
	"flag"
	"github.com/jkvatne/jkvgui/btn"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/icon"
	"github.com/jkvatne/jkvgui/sys"
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
	sys.Initialize(gpu.Window, 12)
	gpu.BackgroundColor(f32.White)

	// Draw buttons
	wid.Show(10, 10, btn.Filled("Primary", icon.Home, nil, btn.Role(theme.Primary), ""))
	wid.Show(150, 10, btn.Filled("Secondary", icon.Home, nil, btn.Role(theme.Secondary), ""))
	wid.Show(300, 10, btn.Round(icon.Home, nil, btn.Default.Role(theme.Secondary), ""))
	wid.Show(10, 50, btn.Outline("Outline", nil, nil, btn.Role(theme.Surface), ""))
	wid.Show(150, 50, btn.Text("Text", nil, nil, btn.Role(theme.SurfaceContainer), ""))
	wid.Show(300, 50, btn.Filled("", icon.Home, nil, btn.Role(theme.Secondary), ""))
	wid.Show(10, 100, btn.Filled("Size 1.0", icon.Home, nil, btn.Role(theme.Primary).Size(1.0), ""))
	wid.Show(150, 100, btn.Filled("Size 2.0", icon.Home, nil, btn.Role(theme.Secondary).Size(2.0), ""))
	wid.Show(300, 100, btn.Filled("Surface", nil, nil, btn.Role(theme.Surface), ""))

	// Verify resulting image
	VerifyScreen(t, "TestButtons", 400, 200, saveScreen)
	// Place breakpoint here in order to look at the screen output.
	time.Sleep(100 * time.Millisecond)

}
