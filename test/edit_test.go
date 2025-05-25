package main

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/input"
	"github.com/jkvatne/jkvgui/wid"
	"log/slog"
	"testing"
	"time"
)

var text = "abcdefg hijklmn opqrst"

func TestEditCursor(t *testing.T) {
	input.Initialize()
	slog.SetLogLoggerLevel(slog.LevelError)
	gpu.InitWindow(600, 70, "Test", 2, 1.0)
	defer input.Shutdown()
	input.InitializeWindow()
	// Simulate click between j and k
	input.SimPos(420, 30)
	input.SimLeftBtnPress()
	// Draw buttons'
	wid.Show(10, 10, 570, wid.Edit(&text, "Test", nil, nil))
	input.SimLeftBtnRelease()
	wid.Show(10, 10, 570, wid.Edit(&text, "Test", nil, nil))
	// Verify resulting image
	VerifyScreen(t, "TestEditCursor", 600, 70, saveScreen)
	// Place breakpoint here in order to look at the screen output.
	time.Sleep(time.Microsecond)
}

func TestEdit(t *testing.T) {
	input.Initialize()
	slog.SetLogLoggerLevel(slog.LevelError)
	gpu.InitWindow(600, 70, "Test", 2, 1.0)
	defer input.Shutdown()
	input.InitializeWindow()
	gpu.SetBackgroundColor(f32.White)
	// Simulate doubleclick between j and k
	input.SimPos(420, 30)
	input.SimLeftBtnPress()
	// Draw buttons'
	wid.Show(10, 10, 570, wid.Edit(&text, "Test", nil, nil))
	input.SimLeftBtnRelease()
	wid.Show(10, 10, 570, wid.Edit(&text, "Test", nil, nil))
	input.SimLeftBtnPress()
	wid.Show(10, 10, 570, wid.Edit(&text, "Test", nil, nil))
	input.SimLeftBtnRelease()
	wid.Show(10, 10, 570, wid.Edit(&text, "Test", nil, nil))
	// Verify resulting image
	VerifyScreen(t, "TestEdit", 400, 200, saveScreen)
	// Place breakpoint here in order to look at the screen output.
	time.Sleep(time.Millisecond)

}
