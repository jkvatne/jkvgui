package main

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/mouse"
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/wid"
	"log/slog"
	"testing"
	"time"
)

var text = "abcdefg hijklmn opqrst"

func TestEditCursor(t *testing.T) {
	_ = gpu.InitWindow(600, 70, "Test", 2)
	defer gpu.Shutdown()
	sys.Initialize(gpu.Window, 24)
	// Simulate click between j and k
	mouse.SimPos(420, 30)
	mouse.SimLeftBtnPress()
	// Draw buttons'
	wid.Show(10, 10, 570, wid.Edit(&text, "Test", nil, nil))
	mouse.SimLeftBtnRelease()
	wid.Show(10, 10, 570, wid.Edit(&text, "Test", nil, nil))
	// Verify resulting image
	VerifyScreen(t, "TestEditCursor", 600, 70, saveScreen)
	// Place breakpoint here in order to look at the screen output.
	time.Sleep(time.Microsecond)
}

func TestEdit(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelError)
	_ = gpu.InitWindow(600, 70, "Test", 2)
	defer gpu.Shutdown()
	sys.Initialize(gpu.Window, 24)
	gpu.BackgroundColor(f32.White)
	// Simulate doubleclick between j and k
	mouse.SimPos(420, 30)
	mouse.SimLeftBtnPress()
	// Draw buttons'
	wid.Show(10, 10, 570, wid.Edit(&text, "Test", nil, nil))
	mouse.SimLeftBtnRelease()
	wid.Show(10, 10, 570, wid.Edit(&text, "Test", nil, nil))
	mouse.SimLeftBtnPress()
	wid.Show(10, 10, 570, wid.Edit(&text, "Test", nil, nil))
	mouse.SimLeftBtnRelease()
	wid.Show(10, 10, 570, wid.Edit(&text, "Test", nil, nil))
	// Verify resulting image
	VerifyScreen(t, "TestEdit", 400, 200, saveScreen)
	// Place breakpoint here in order to look at the screen output.
	time.Sleep(time.Millisecond)

}
