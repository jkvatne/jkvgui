package wid

import (
	"time"

	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/gpu/font"
	"github.com/jkvatne/jkvgui/theme"
)

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
	Delay:        time.Millisecond * 750,
}

// Hint is called if the mouse is inside a clickable widget
// i.e. when it is hovered.
func Hint(ctx Ctx, text string, tag any) {
	if text == "" {
		return
	}
	if !gpu.TagsEqual(ctx.Win.CurrentHint.Tag, tag) {
		ctx.Win.CurrentHint.T = time.Now()
	}
	ctx.Win.CurrentHint.Text = text
	ctx.Win.CurrentHint.Tag = tag
	ctx.Win.CurrentHint.WidgetRect = ctx.Rect
	ctx.Win.Defer(func() { showHint(ctx) })
}

// ShowHint is called at the end of the display loop.
// It will show the hint on top of everything else.
func showHint(ctx Ctx) {
	style := &DefaultHintStyle
	hint := ctx.Win.CurrentHint
	if time.Since(hint.T) > style.Delay {
		f := font.Get(style.FontNo)
		textHeight := f.Height
		w := textHeight * 8
		x := min(hint.WidgetRect.X+w+style.Padding.L+style.Padding.R, ctx.Win.WidthDp)
		x = max(0, x-w)
		lines := font.Split(hint.Text, w-style.Padding.L-style.Padding.R, f)
		h := textHeight*float32(len(lines)) + style.Padding.T + style.Padding.B + 2*style.BorderWidth
		// Nominal y location is below the original widget
		y := hint.WidgetRect.Y + hint.WidgetRect.H + textHeight/5
		// But if this is below the bottom of the window, put it above the original widget
		if y+h > ctx.Win.HeightDp {
			y = hint.WidgetRect.Y - h - textHeight/5
		}

		// Draw hint. Location given by x,y,w,h
		hintOutline := f32.Rect{X: x, Y: y, W: w, H: h}
		ctx.Win.Gd.RoundedRect(hintOutline, style.CornerRadius, style.BorderWidth, style.Color.Bg(), style.BorderColor.Fg())
		yb := y + style.Padding.T + f.Baseline
		for _, line := range lines {
			f.DrawText(ctx.Win.Gd,
				x+style.Padding.L+style.Padding.L+style.BorderWidth,
				yb,
				style.Color.Fg(),
				0, gpu.LTR, line)
			yb = yb + textHeight
		}
	}
}
