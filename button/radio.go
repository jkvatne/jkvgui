package button

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/focus"
	"github.com/jkvatne/jkvgui/font"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/icon"
	"github.com/jkvatne/jkvgui/mouse"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
)

type RadioButtonStyle struct {
	FontSize float32
	FontNo   int
	Role     theme.UIRole
	Padding  f32.Padding
}

var DefaultRadioButton RadioButtonStyle = RadioButtonStyle{
	FontSize: 1,
	FontNo:   0,
	Role:     theme.OnSurface,
	Padding:  f32.Padding{5, 7, 15, 7},
}

func RadioButton(text string, value *string, key string, style *RadioButtonStyle) wid.Wid {
	return func(ctx wid.Ctx) wid.Dim {
		if style == nil {
			style = &DefaultRadioButton
		}
		f := font.Fonts[style.FontNo]
		height := f.Height(style.FontSize) + style.Padding.T + style.Padding.B
		width := f.Width(style.FontSize, text)/2 + style.Padding.L + style.Padding.R
		baseline := f.Baseline(style.FontSize) + style.Padding.T
		iconRect := f32.Rect{X: ctx.Rect.X, Y: ctx.Rect.Y, W: height, H: height}

		if ctx.Rect.H == 0 {
			return wid.Dim{W: height*6/5 + width + style.Padding.L, H: height, Baseline: baseline}
		}

		if mouse.LeftBtnReleased(ctx.Rect) {
			focus.Set(value)
			if !ctx.Disabled {
				*value = key
			}
		}
		if focus.At(ctx.Rect, value) {
			gpu.Shade(iconRect.Reduce(-1), 999, f32.Shade, 5)
		}
		if *value == key {
			icon.Draw(iconRect.X, iconRect.Y, height, icon.RadioChecked, style.Role.Fg())
		} else {
			icon.Draw(iconRect.X, iconRect.Y, height, icon.RadioUnchecked, style.Role.Fg())
		}
		f.Printf(ctx.Rect.X+style.Padding.L+height, ctx.Rect.Y+baseline, style.FontSize, 0, text)

		return wid.Dim{W: ctx.Rect.W, H: ctx.Rect.H, Baseline: ctx.Baseline}
	}
}
