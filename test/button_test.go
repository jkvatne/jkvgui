package test

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
	sys.Init()
	defer sys.Shutdown()
	slog.SetLogLoggerLevel(slog.LevelError)
	sys.CreateWindow(0, 0, 800, 300, "Test", 1, 1.0)
	sys.StartFrame(theme.Canvas.Bg())
	gpu.RoundedRect(gpu.CurrentInfo.WindowRect.Reduce(1), 0, 1, f32.Transparent, f32.Red)

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
	VerifyScreen(t, "TestButtons", 800, 300, saveScreen)
	// Place breakpoint here in order to look at the screen output.
	time.Sleep(time.Millisecond)

}
