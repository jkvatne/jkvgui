package test

import (
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
	"log/slog"
	"testing"
	"time"
)

var text = "abcdefg hijklmn opqrst"

func TestEdit(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelError)
	if sys.CurrentWindow == nil {
		sys.CreateWindow(0, 0, 600, 70, "Test", 2, 1.0)
	}
	sys.StartFrame(theme.Canvas.Bg())
	// Simulate click between j and k
	sys.SimPos(420, 30)
	sys.SimLeftBtnPress()
	// Draw buttons'
	wid.Show(10, 10, 570, wid.Edit(&text, "Test", nil, nil))
	sys.SimLeftBtnRelease()
	wid.Show(10, 10, 570, wid.Edit(&text, "Test", nil, nil))
	// Verify resulting image
	VerifyScreen(t, "TestEditCursor", 600, 70, saveScreen)
	// Place breakpoint here in order to look at the screen output.
	time.Sleep(time.Microsecond)

	sys.StartFrame(theme.Canvas.Bg())
	// Simulate doubleclick between j and k
	sys.SimPos(420, 30)
	sys.SimLeftBtnPress()
	// Draw buttons'
	wid.Show(10, 10, 570, wid.Edit(&text, "Test", nil, nil))
	sys.SimLeftBtnRelease()
	wid.Show(10, 10, 570, wid.Edit(&text, "Test", nil, nil))
	sys.SimLeftBtnPress()
	wid.Show(10, 10, 570, wid.Edit(&text, "Test", nil, nil))
	sys.SimLeftBtnRelease()
	wid.Show(10, 10, 570, wid.Edit(&text, "Test", nil, nil))
	// Verify resulting image
	VerifyScreen(t, "TestEdit", 600, 70, saveScreen)
	// Place breakpoint here in order to look at the screen output.
	time.Sleep(time.Millisecond)

}
