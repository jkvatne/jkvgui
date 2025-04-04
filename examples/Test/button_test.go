package main

import (
	"flag"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
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
	_ = gpu.InitWindow(400, 200, "Test", 2)
	defer gpu.Shutdown()
	sys.Initialize(gpu.Window, 12)
	gpu.BackgroundColor(f32.White)

	// Draw buttons
	wid.Show(10, 10, wid.Btn("Primary", gpu.Home, nil, wid.Filled, ""))
	wid.Show(150, 10, wid.Btn("Secondary", gpu.Home, nil, wid.Filled.Role(theme.Secondary), ""))
	wid.Show(300, 10, wid.Btn("", gpu.Home, nil, wid.Round, ""))
	wid.Show(10, 50, wid.Btn("Outline", nil, nil, wid.Outline, ""))
	wid.Show(150, 50, wid.Btn("Text", nil, nil, wid.Text, ""))
	wid.Show(300, 50, wid.Btn("", gpu.Home, nil, wid.Round.Role(theme.Secondary), ""))
	wid.Show(10, 100, wid.Btn("Size 1.0", gpu.Home, nil, wid.Filled.Size(1.0), ""))
	wid.Show(150, 100, wid.Btn("Size 2.0", gpu.Home, nil, wid.Filled.Role(theme.Secondary).Size(2.0), ""))
	wid.Show(300, 100, wid.Btn("Surface", nil, nil, wid.Filled.Role(theme.Surface), ""))

	// Verify resulting image
	VerifyScreen(t, "TestButtons", 400, 200, saveScreen)
	// Place breakpoint here in order to look at the screen output.
	time.Sleep(time.Millisecond)

}
