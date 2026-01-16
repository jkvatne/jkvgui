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

func TestCombo(t *testing.T) {
	slog.Info("Test combo")
	sys.Init()
	defer sys.Shutdown()
	sys.NoScaling = false
	win := sys.CreateWindow(0, 0, 400, 400, "Test", 0, 2.0)
	win.StartFrame()
	styles := make([]wid.ComboStyle, len(cases))
	values := make([]string, len(cases))
	for i, c := range cases {
		values[i] = "<none>"
		styles[i] = wid.DefaultCombo
		styles[i].LabelRightAdjust = c.rightAdjust
		styles[i].LabelSize = c.labelSize
		styles[i].LabelSpacing = c.labelSpacing
		styles[i].EditSize = c.editSize
		styles[i].BorderWidth = c.borderWidth
		styles[i].BorderColor = c.borderColor
		wid.Display(win, 0, float32(i)*22, 200, wid.Combo(&values[i], list, "Label"+strconv.Itoa(i), &styles[i]))
	}
	// Verify resulting image
	VerifyScreen(t, win, "TestCombo", 600, 400, 500)
	win.EndFrame()
	// Place breakpoint here in order to look at the screen output.
	time.Sleep(time.Microsecond)
}
