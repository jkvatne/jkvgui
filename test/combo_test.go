package test

import (
	"log/slog"
	"strconv"
	"testing"
	"time"

	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
)

var (
	list = []string{"value1", "value2", "value3", "value4", "value5"}
)

type setup struct {
	labelSize    float32
	editSize     float32
	rightAdjust  bool
	maxDropdown  int
	labelSpacing float32
	borderWidth  float32
	borderColor  theme.UIRole
}

var cases = []setup{
	setup{0.3, 0.7, false, 3, 0, 1, theme.Outline},
	setup{70, 100, false, 3, 0, 0, theme.Outline},
	setup{50, 150, true, 3, 0, 1, theme.Outline},
	setup{0, 0, true, 3, 5, 1, theme.Primary},
	setup{0, 0, true, 3, 15, 1, theme.Primary},
}

func MakeStyle(i int) wid.ComboStyle {
	c := cases[i]
	s := wid.DefaultCombo
	s.LabelRightAdjust = c.rightAdjust
	s.LabelSize = c.labelSize
	s.LabelSpacing = c.labelSpacing
	s.EditSize = c.editSize
	s.BorderWidth = c.borderWidth
	s.BorderColor = c.borderColor
	return s
}

func TestCombo(t *testing.T) {
	slog.Info("Test combo")
	sys.Init()
	defer sys.Shutdown()
	sys.NoScaling = false
	win := sys.CreateWindow(0, 0, 200, 200, "Test", 0, 1.0)
	win.StartFrame()
	styles := make([]wid.ComboStyle, len(cases))
	values := make([]string, len(cases))
	for i, _ := range cases {
		values[i] = "<none>"
		styles[i] = MakeStyle(i)
		wid.Display(win, 0, float32(i)*22, 200, wid.Combo(&values[i], list, "Label"+strconv.Itoa(i), &styles[i]))
	}

	// Simulate double click
	sys.BlinkState.Store(true)
	win.SimLeftDoubleClick(180, 12)
	wid.Display(win, 0, 0, 200, wid.Combo(&values[0], list, "Label"+strconv.Itoa(0), &styles[0]))
	// Verify resulting image
	VerifyScreen(t, win, "Combo", 200, 200, 500)
	win.EndFrame()
	// Place breakpoint here in order to look at the screen output.
	time.Sleep(time.Microsecond)
}

func TestComboDoubleClick(t *testing.T) {
	sys.Init()
	defer sys.Shutdown()
	sys.NoScaling = false
	win := sys.CreateWindow(0, 0, 200, 75, "Test", 0, 1.0)
	win.StartFrame()
	style := MakeStyle(0)
	value := "<none>"
	wid.Display(win, 0, 0, 200, wid.Combo(&value, list, "Label0", &style))
	// Simulate double click
	sys.BlinkState.Store(true)
	win.SimLeftDoubleClick(180, 12)
	wid.Display(win, 0, 0, 200, wid.Combo(&value, list, "Label0", &style))
	// Verify resulting image
	VerifyScreen(t, win, "TestComboDoubleClick", 200, 75, 500)
	win.EndFrame()
	// Place breakpoint here in order to look at the screen output.
	time.Sleep(time.Microsecond)
}

func TestComboClick(t *testing.T) {
	sys.Init()
	defer sys.Shutdown()
	sys.NoScaling = false
	win := sys.CreateWindow(0, 0, 200, 100, "Test", 0, 1.0)
	win.StartFrame()
	style := MakeStyle(0)
	value := "<none>"
	wid.Display(win, 0, 65, 200, wid.Combo(&value, list, "Label0", &style))
	// Simulate double click
	win.SimLeftClick(184, 77)
	sys.BlinkState.Store(true)
	wid.Display(win, 0, 65, 200, wid.Combo(&value, list, "Label0", &style))
	// Verify resulting image
	VerifyScreen(t, win, "TestComboClick", 200, 100, 500)
	win.EndFrame()
	// Place breakpoint here in order to look at the screen output.
	time.Sleep(time.Microsecond)
}
