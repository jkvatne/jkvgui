package button

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/focus"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/gpu/font"
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

var DefaultRadioButton = RadioButtonStyle{
	FontSize: 1.0,
	FontNo:   0,
	Role:     theme.OnSurface,
	Padding:  f32.Padding{L: 5, T: 3, R: 8, B: 3},
}

func RadioButton(label string, value *string, key string, style *RadioButtonStyle) wid.Wid {
	return func(ctx wid.Ctx) wid.Dim {
		if style == nil {
			style = &DefaultRadioButton
		}
		f := font.Fonts[style.FontNo]
		fontHeight := f.Height(style.FontSize)
		height := fontHeight + style.Padding.T + style.Padding.B
		width := f.Width(style.FontSize, label) + style.Padding.L + style.Padding.R + height
		baseline := f.Baseline(style.FontSize) + style.Padding.T
		extRect := f32.Rect{X: ctx.Rect.X, Y: ctx.Rect.Y, W: width, H: height}
		iconRect := extRect.Inset(style.Padding, 0)
		iconRect.W = iconRect.H
		if ctx.Rect.H == 0 {
			return wid.Dim{W: height*6/5 + width + style.Padding.L, H: height, Baseline: baseline}
		}
		if gpu.DebugWidgets {
			gpu.RoundedRect(extRect, 0, 0.5, f32.Transparent, f32.Blue)
		}
		if mouse.LeftBtnClick(ctx.Rect) {
			focus.Set(value)
			if !ctx.Disabled {
				*value = key
			}
		}
		if focus.At(ctx.Rect, value) {
			gpu.Shade(iconRect.Move(0, -1), -1, f32.Shade, 5)
		} else if mouse.Hovered(ctx.Rect) {
			gpu.Shade(iconRect.Move(0, -1), -1, f32.Shade, 3)
		}
		if *value == key {
			icon.Draw(iconRect.X, iconRect.Y-1, iconRect.H, icon.RadioChecked, style.Role.Fg())
		} else {
			icon.Draw(iconRect.X, iconRect.Y-1, iconRect.H, icon.RadioUnchecked, style.Role.Fg())
		}
		f.DrawText(iconRect.X+fontHeight*6/5, ctx.Rect.Y+baseline, style.Role.Fg(), style.FontSize, 0, gpu.LeftToRight, label)

		return wid.Dim{W: ctx.Rect.W, H: ctx.Rect.H, Baseline: ctx.Baseline}
	}
}
