package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/gpu/font"
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
)

type RadioButtonStyle struct {
	FontNo  int
	Role    theme.UIRole
	Padding f32.Padding
}

var DefaultRadioButton = RadioButtonStyle{
	FontNo:  gpu.Normal12,
	Role:    theme.OnSurface,
	Padding: f32.Padding{L: 3, T: 1, R: 2, B: 1},
}

func RadioButton(label string, value *string, key string, style *RadioButtonStyle) Wid {
	return func(ctx Ctx) Dim {
		if style == nil {
			style = &DefaultRadioButton
		}
		f := font.Fonts[style.FontNo]
		fontHeight := f.Height
		height := fontHeight + style.Padding.T + style.Padding.B
		width := f.Width(label) + style.Padding.L + style.Padding.R + height
		baseline := f.Baseline + style.Padding.T
		extRect := f32.Rect{X: ctx.Rect.X, Y: ctx.Rect.Y, W: width, H: height}
		iconRect := extRect.Inset(style.Padding, 0)
		iconRect.W = iconRect.H
		if ctx.Mode != RenderChildren {
			return Dim{W: height*6/5 + width + style.Padding.L, H: height, Baseline: baseline}
		}
		if *gpu.DebugWidgets {
			gpu.RoundedRect(extRect, 0, 0.5, f32.Transparent, f32.Blue)
		}
		if sys.LeftBtnClick(ctx.Rect) {
			sys.SetFocusedTag(value)
			if !ctx.Disabled {
				*value = key
			}
		}
		if sys.At(ctx.Rect, value) {
			gpu.Shade(iconRect.Move(0, -1), -1, f32.Shade, 5)
		} else if sys.Hovered(ctx.Rect) {
			gpu.Shade(iconRect.Move(0, -1), -1, f32.Shade, 3)
		}
		if *value == key {
			gpu.DrawIcon(iconRect.X, iconRect.Y-1, iconRect.H, gpu.RadioChecked, style.Role.Fg())
		} else {
			gpu.DrawIcon(iconRect.X, iconRect.Y-1, iconRect.H, gpu.RadioUnchecked, style.Role.Fg())
		}
		f.DrawText(iconRect.X+fontHeight*6/5, ctx.Rect.Y+baseline, style.Role.Fg(), 0, gpu.LTR, label)

		return Dim{W: ctx.Rect.W, H: ctx.Rect.H, Baseline: ctx.Baseline}
	}
}
