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
	w := sys.CreateWindow(0, 0, 400, 150, "Test", 1, 1.0)
	sys.LoadOpenGl(w)
	// Draw buttons
	w.StartFrame(theme.Canvas.Bg())
	wid.Display(w, 10, 10, 400, wid.Btn("Primary", gpu.Home, nil, wid.Filled, ""))
	wid.Display(w, 150, 10, 400, wid.Btn("Secondary", gpu.Home, nil, wid.Filled.Role(theme.Secondary), ""))
	wid.Display(w, 300, 10, 400, wid.Btn("", gpu.Home, nil, wid.Round, ""))
	wid.Display(w, 10, 50, 400, wid.Btn("Outline", nil, nil, wid.Outline, ""))
	wid.Display(w, 150, 50, 400, wid.Btn("Text", nil, nil, wid.Text, ""))
	wid.Display(w, 300, 50, 400, wid.Btn("", gpu.Home, nil, wid.Round.Role(theme.Secondary), ""))
	wid.Display(w, 10, 100, 400, wid.Btn("Size 12", gpu.Home, nil, wid.Filled.Font(gpu.Normal12), ""))
	wid.Display(w, 150, 100, 400, wid.Btn("Size 20", gpu.Home, nil, wid.Filled.Role(theme.Secondary).Font(gpu.Normal20), ""))
	wid.Display(w, 300, 100, 400, wid.Btn("Surface", nil, nil, wid.Filled.Role(theme.Surface), ""))
	// Verify resulting image
	VerifyScreen(t, w, "TestButtons", 400, 150, saveScreen)
	w.EndFrame()
	// Place breakpoint here in order to look at the screen output.
	time.Sleep(time.Millisecond)

}
