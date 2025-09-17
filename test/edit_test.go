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
	w := sys.CreateWindow(0, 0, 600, 70, "Test", 1, 1.0)
	sys.LoadOpenGl(w)
	w.StartFrame(theme.Canvas.Bg())
	// Simulate click between j and k
	w.SimPos(420, 30)
	w.SimLeftBtnPress()
	// Draw buttons
	wid.Display(w, 10, 10, 570, wid.Edit(&text, "Test", nil, nil))
	w.SimLeftBtnRelease()
	wid.Display(w, 10, 10, 570, wid.Edit(&text, "Test", nil, nil))
	// Verify resulting image
	VerifyScreen(t, w, "TestEditCursor", 600, 70, saveScreen)
	w.EndFrame()
	// Place breakpoint here in order to look at the screen output.
	time.Sleep(time.Microsecond)
}

func TestEdit(t *testing.T) {
	sys.Init()
	defer sys.Shutdown()
	w := sys.CreateWindow(0, 0, 600, 70, "Test", 1, 1.0)
	sys.LoadOpenGl(w)
	w.StartFrame(theme.Canvas.Bg())
	// Simulate doubleclick between j and k
	w.SimPos(420, 30)
	w.SimLeftBtnPress()
	// Draw buttons'
	wid.Display(w, 10, 10, 570, wid.Edit(&text, "Test", nil, nil))
	w.SimLeftBtnRelease()
	wid.Display(w, 10, 10, 570, wid.Edit(&text, "Test", nil, nil))
	w.SimLeftBtnPress()
	wid.Display(w, 10, 10, 570, wid.Edit(&text, "Test", nil, nil))
	w.SimLeftBtnRelease()
	wid.Display(w, 10, 10, 570, wid.Edit(&text, "Test", nil, nil))
	// Verify resulting image
	VerifyScreen(t, w, "TestEdit", 600, 70, saveScreen)
	w.EndFrame()
	// Place breakpoint here in order to look at the screen output.
	time.Sleep(time.Millisecond)
	sys.Shutdown()
}
