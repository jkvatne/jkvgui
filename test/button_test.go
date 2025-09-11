package test

import (
	"flag"
	"log/slog"
	"testing"
	"time"

	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
)

var saveScreen bool

func init() {
	flag.BoolVar(&saveScreen, "test.save", false, "Save the captured screen to ./test-assets, making this the reference image.")
}

func TestButtons(t *testing.T) {
	sys.Init()
	defer sys.Shutdown()
	sys.NoScaling = true
	slog.SetLogLoggerLevel(slog.LevelError)
	sys.CreateWindow(0, 0, 400, 150, "Test", 1, 1.0)
	sys.CurrentWindow.SetSize(400, 150)
	sys.StartFrame(theme.Canvas.Bg())

	// Draw buttons
	wid.Show(10, 10, 400, wid.Btn("Primary", gpu.Home, nil, wid.Filled, ""))
	wid.Show(150, 10, 400, wid.Btn("Secondary", gpu.Home, nil, wid.Filled.Role(theme.Secondary), ""))
	wid.Show(300, 10, 400, wid.Btn("", gpu.Home, nil, wid.Round, ""))
	wid.Show(10, 50, 400, wid.Btn("Outline", nil, nil, wid.Outline, ""))
	wid.Show(150, 50, 400, wid.Btn("Text", nil, nil, wid.Text, ""))
	wid.Show(300, 50, 400, wid.Btn("", gpu.Home, nil, wid.Round.Role(theme.Secondary), ""))
	wid.Show(10, 100, 400, wid.Btn("Size 12", gpu.Home, nil, wid.Filled.Font(gpu.Normal12), ""))
	wid.Show(150, 100, 400, wid.Btn("Size 20", gpu.Home, nil, wid.Filled.Role(theme.Secondary).Font(gpu.Normal20), ""))
	wid.Show(300, 100, 400, wid.Btn("Surface", nil, nil, wid.Filled.Role(theme.Surface), ""))
	// Verify resulting image
	VerifyScreen(t, "TestButtons", 400, 150, saveScreen)
	// Place breakpoint here in order to look at the screen output.
	time.Sleep(time.Millisecond)

}
