package test

import (
	"log/slog"
	"testing"
	"time"

	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
)

var text = "abcdefg hijklmn opqrst"

func TestEditCursor(t *testing.T) {
	sys.Init()
	defer sys.Shutdown()
	sys.NoScaling = true
	slog.SetLogLoggerLevel(slog.LevelError)
	sys.CreateWindow(0, 0, 600, 70, "Test", 1, 1.0)
	sys.StartFrame(theme.Canvas.Bg())
	// Simulate click between j and k
	sys.SimPos(420, 30)
	sys.SimLeftBtnPress()
	// Draw buttons
	wid.Display(10, 10, 570, wid.Edit(&text, "Test", nil, nil))
	sys.SimLeftBtnRelease()
	wid.Display(10, 10, 570, wid.Edit(&text, "Test", nil, nil))
	// Verify resulting image
	VerifyScreen(t, "TestEditCursor", 600, 70, saveScreen)
	sys.EndFrame()
	// Place breakpoint here in order to look at the screen output.
	time.Sleep(time.Microsecond)
}

func TestEdit(t *testing.T) {
	sys.Init()
	defer sys.Shutdown()
	sys.CreateWindow(0, 0, 600, 70, "Test", 1, 1.0)
	sys.StartFrame(theme.Canvas.Bg())
	// Simulate doubleclick between j and k
	sys.SimPos(420, 30)
	sys.SimLeftBtnPress()
	// Draw buttons'
	wid.Display(10, 10, 570, wid.Edit(&text, "Test", nil, nil))
	sys.SimLeftBtnRelease()
	wid.Display(10, 10, 570, wid.Edit(&text, "Test", nil, nil))
	sys.SimLeftBtnPress()
	wid.Display(10, 10, 570, wid.Edit(&text, "Test", nil, nil))
	sys.SimLeftBtnRelease()
	wid.Display(10, 10, 570, wid.Edit(&text, "Test", nil, nil))
	// Verify resulting image
	VerifyScreen(t, "TestEdit", 600, 70, saveScreen)
	sys.WindowList[0].SwapBuffers()
	// Place breakpoint here in order to look at the screen output.
	time.Sleep(time.Millisecond)
	sys.Shutdown()
}
