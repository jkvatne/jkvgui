package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/focus"
	"github.com/jkvatne/jkvgui/font"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/icon"
	"github.com/jkvatne/jkvgui/mouse"
)

type CheckboxStyle struct {
	FontSize float32
	FontNo   int
	Color    f32.Color
	Padding  f32.Padding
}

var DefaultCheckbox = CheckboxStyle{
	FontSize: 1,
	FontNo:   0,
	Color:    f32.Color{R: 0, G: 0, B: 0, A: 1},
	Padding:  f32.Padding{L: 3, T: 3, R: 3, B: 3},
}

func Checkbox(text string, state *bool, style *CheckboxStyle, hint string) Wid {
	return func(ctx Ctx) Dim {
		if style == nil {
			style = &DefaultCheckbox
		}
		f := font.Fonts[style.FontNo]
		height := f.Height(style.FontSize) + style.Padding.T + style.Padding.B
		width := f.Width(style.FontSize, text)/2 + style.Padding.L + style.Padding.R
		baseline := f.Baseline(style.FontSize) + style.Padding.T
		iconRect := f32.Rect{X: ctx.Rect.X, Y: ctx.Rect.Y, W: height, H: height}

		if ctx.Rect.H == 0 {
			return Dim{W: height*6/5 + width, H: height, Baseline: baseline}
		}

		focused := focus.At(ctx.Rect, state)

		if mouse.LeftBtnReleased(ctx.Rect) {
			focus.Set(state)
			*state = !*state
		}
		if focused {
			gpu.Shade(iconRect.Reduce(-1), 5, f32.Shade, 5)
		}
		if mouse.Hovered(ctx.Rect) {
			Hint(hint, state)
		}
		if *state {
			icon.Draw(iconRect.X, iconRect.Y, height, icon.BoxChecked, style.Color)
		} else {
			icon.Draw(iconRect.X, iconRect.Y, height, icon.BoxUnchecked, style.Color)
		}
		f.Printf(ctx.Rect.X+style.Padding.L+height, ctx.Rect.Y+baseline, style.FontSize, 0, text)

		return Dim{W: ctx.Rect.W, H: ctx.Rect.H, Baseline: ctx.Baseline}
	}
}
