package wid

import (
	"time"

	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/gpu/font"
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
)

var CurrentHint struct {
	Pos  f32.Pos // Mouse position at time of pop-up
	Text string
	T    time.Time
	Tag  any
}

type HintStyle struct {
	FontNo       int
	FontSize     float32
	Color        theme.UIRole
	CornerRadius float32
	BorderColor  theme.UIRole
	BorderWidth  float32
	Padding      f32.Padding
	Delay        time.Duration
}

var DefaultHintStyle = HintStyle{
	FontNo:       gpu.Normal10,
	FontSize:     0.9,
	Color:        theme.SecondaryContainer,
	CornerRadius: 5,
	BorderColor:  theme.Outline,
	BorderWidth:  1,
	Padding:      f32.Padding{L: 3, T: 3, R: 1, B: 2},
	Delay:        time.Millisecond * 1000,
}

// Hint is called if the mouse is inside a clickable widget
// i.e. when it is hovered.
func Hint(text string, tag any) {
	if text == "" {
		return
	}
	if !gpu.TagsEqual(CurrentHint.Tag, tag) {
		CurrentHint.T = time.Now()
	}
	CurrentHint.Text = text
	CurrentHint.Tag = tag
	CurrentHint.Pos = sys.Pos()
	sys.CurrentInfo.HintActive = true
	gpu.Defer(showHint)
}

// ShowHint is called at the end of the display loop.
// It will show the hint on top of everything else.
func showHint() {
	style := &DefaultHintStyle
	if time.Since(CurrentHint.T) > style.Delay && sys.CurrentInfo.HintActive {
		f := font.Get(style.FontNo)
		textHeight := f.Height
		w := textHeight * 8
		x := min(CurrentHint.Pos.X+w+style.Padding.L+style.Padding.R, sys.WindowWidthDp())
		x = max(float32(0), x-w)
		lines := font.Split(CurrentHint.Text, w-style.Padding.L-style.Padding.R, f)
		h := textHeight*float32(len(lines)) + style.Padding.T + style.Padding.B + 2*style.BorderWidth
		y := min(CurrentHint.Pos.Y+h, sys.WindowHeightDp())
		y = max(0, y-h)
		yb := y + style.Padding.T + f.Baseline
		r := f32.Rect{X: x, Y: y, W: w, H: h}
		gpu.RoundedRect(r, style.CornerRadius, style.BorderWidth, style.Color.Bg(), style.BorderColor.Fg())
		for _, line := range lines {
			f.DrawText(
				x+style.Padding.L+style.Padding.L+style.BorderWidth,
				yb,
				style.Color.Fg(),
				0, gpu.LTR, line)
			yb = yb + textHeight
		}
	}
}
