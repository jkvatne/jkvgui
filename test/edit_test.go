package test

import (
	"log/slog"
	"testing"
	"time"

	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
)

var text = "Some text with ÅgÅgØøÆæ"

func TestEdit(t *testing.T) {
	slog.Info("TestEdit")
	sys.Init()
	defer sys.Shutdown()
	sys.NoScaling = true
	slog.SetLogLoggerLevel(slog.LevelError)
	w := sys.CreateWindow(0, 0, 600, 70, "Test", 1, 1.0)
	w.StartFrame(theme.Canvas.Bg())
	// Draw Edit widget
	wid.Display(w, 10, 10, 570, wid.Edit(&text, "Test", nil, nil))
	// Verify resulting image
	VerifyScreen(t, w, "TestEdit", 600, 70, saveScreen)
	w.EndFrame()
	// Place breakpoint here in order to look at the screen output.
	time.Sleep(time.Microsecond)
	sys.Shutdown()
}

func TestEditCursor(t *testing.T) {
	slog.Info("TestEditCursor")
	sys.Init()
	defer sys.Shutdown()
	sys.NoScaling = true
	slog.SetLogLoggerLevel(slog.LevelError)
	w := sys.CreateWindow(0, 0, 600, 70, "Test", 1, 1.0)
	w.StartFrame(theme.Canvas.Bg())
	// Draw Edit widget
	// wid.Display(w, 10, 10, 570, wid.Edit(&text, "Test", nil, nil))
	// Simulate single click between g and q
	sys.BlinkState.Store(true)
	w.SimLeftBtnPress(420, 30)
	wid.Display(w, 10, 10, 570, wid.Edit(&text, "Test", nil, nil))
	w.SimLeftBtnRelease(420, 30)
	wid.Display(w, 10, 10, 570, wid.Edit(&text, "Test", nil, nil))
	// Verify resulting image
	VerifyScreen(t, w, "TestEditCursor", 600, 70, saveScreen)
	w.EndFrame()
	// Place breakpoint here in order to look at the screen output.
	time.Sleep(time.Microsecond)
	sys.Shutdown()
}

func TestEditSelect(t *testing.T) {
	slog.Info("TestEditSelect")
	sys.Init()
	defer sys.Shutdown()
	sys.NoScaling = true
	slog.SetLogLoggerLevel(slog.LevelError)
	w := sys.CreateWindow(0, 0, 600, 70, "Test", 1, 1.0)
	w.StartFrame(theme.Canvas.Bg())
	// Draw Edit widget
	wid.Display(w, 10, 10, 570, wid.Edit(&text, "Test", nil, nil))
	// Simulate double-click between j and k
	sys.BlinkState.Store(true)
	w.SimLeftBtnPress(420, 30)
	wid.Display(w, 10, 10, 570, wid.Edit(&text, "Test", nil, nil))
	w.SimLeftBtnRelease(420, 30)
	wid.Display(w, 10, 10, 570, wid.Edit(&text, "Test", nil, nil))
	w.SimLeftBtnPress(420, 30)
	wid.Display(w, 10, 10, 570, wid.Edit(&text, "Test", nil, nil))
	w.SimLeftBtnRelease(420, 30)
	wid.Display(w, 10, 10, 570, wid.Edit(&text, "Test", nil, nil))
	// Verify resulting image
	VerifyScreen(t, w, "TestEditSelect", 600, 70, saveScreen)
	w.EndFrame()
	// Place breakpoint here in order to look at the screen output.
	time.Sleep(time.Millisecond)
	sys.Shutdown()
}
