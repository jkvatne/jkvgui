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
	wid.Show(10, 10, btn.Btn("Primary", icon.Home, nil, btn.Filled, ""))
	wid.Show(150, 10, btn.Btn("Secondary", icon.Home, nil, btn.Filled.Role(theme.Secondary), ""))
	wid.Show(300, 10, btn.Btn("", icon.Home, nil, btn.Round, ""))
	wid.Show(10, 50, btn.Btn("Outline", nil, nil, btn.Outline, ""))
	wid.Show(150, 50, btn.Btn("Text", nil, nil, btn.Text, ""))
	wid.Show(300, 50, btn.Btn("", icon.Home, nil, btn.Round.Role(theme.Secondary), ""))
	wid.Show(10, 100, btn.Btn("Size 1.0", icon.Home, nil, btn.Filled.Size(1.0), ""))
	wid.Show(150, 100, btn.Btn("Size 2.0", icon.Home, nil, btn.Filled.Role(theme.Secondary).Size(2.0), ""))
	wid.Show(300, 100, btn.Btn("Surface", nil, nil, btn.Filled.Role(theme.Surface), ""))

	// Verify resulting image
	VerifyScreen(t, "TestButtons", 400, 200, saveScreen)
	// Place breakpoint here in order to look at the screen output.
	time.Sleep(100 * time.Millisecond)

}
