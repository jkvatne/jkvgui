package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/gpu/font"
	"github.com/jkvatne/jkvgui/theme"
)

type CbStyle struct {
	FontNo          int
	Role            theme.UIRole
	Padding         f32.Padding
	BorderThickness float32
}

var DefaultCb = CbStyle{
	FontNo:  gpu.Normal12,
	Role:    theme.OnSurface,
	Padding: f32.Padding{L: 3, T: 1, R: 2, B: 1},
}

var DefaultCheckbox = CbStyle{
	FontNo:  gpu.Normal12,
	Role:    theme.OnSurface,
	Padding: f32.Padding{L: 3, T: 1, R: 2, B: 1},
}

var GridCheckBox = CbStyle{
	FontNo:          gpu.Normal12,
	Role:            theme.PrimaryContainer,
	Padding:         f32.Padding{L: 1, T: 1, R: 1, B: 1},
	BorderThickness: 0.0,
}

func Checkbox(label string, state *bool, action func(), style *CbStyle, hint string) Wid {
	if style == nil {
		style = &DefaultCheckbox
	}
	f := font.Fonts[style.FontNo]
	fontHeight := f.Height
	baseline := f.Baseline
	height := fontHeight + style.Padding.T + style.Padding.B
	width := f.Width(label) + style.Padding.L + style.Padding.R + fontHeight

	return func(ctx Ctx) Dim {
		extRect := f32.Rect{X: ctx.Rect.X, Y: ctx.Rect.Y, W: width, H: height}
		iconRect := extRect.Inset(style.Padding, 0)
		iconRect.W = iconRect.H
		if ctx.Mode != RenderChildren {
			return Dim{W: width, H: height, Baseline: baseline}
		}
		if ctx.Win.LeftBtnPressed(ctx.Rect) {
			ctx.Win.SetFocusedTag(state)
		}
		if ctx.Win.At(state) && IsKeyClick(ctx) {
			*state = !*state
			if action != nil {
				action()
			}
		}
		if ctx.Win.At(state) {
			ctx.Win.Gd.Shade(iconRect, 4, f32.Shade, 3)
		}
		if ctx.Win.Hovered(ctx.Rect) {
			ctx.Win.Gd.Shade(iconRect, 4, f32.Shade, 3)
			Hint(ctx, hint, state)
		}
		if *state {
			ctx.Win.Gd.DrawIcon(iconRect.X, iconRect.Y, iconRect.H, gpu.BoxChecked, style.Role.Fg())
		} else {
			ctx.Win.Gd.DrawIcon(iconRect.X, iconRect.Y, iconRect.H, gpu.BoxUnchecked, style.Role.Fg())
		}
		f.DrawText(ctx.Win.Gd, iconRect.X+fontHeight*6/5, ctx.Rect.Y+baseline, style.Role.Fg(), 0, gpu.LTR, label)
		// Draw frame around value
		ctx.Win.Gd.RoundedRect(ctx.Rect, 0.0, style.BorderThickness, f32.Transparent, style.Role.Fg())
		DrawDebuggingInfo(ctx, iconRect, iconRect, ctx.Rect)

		return Dim{W: ctx.Rect.W, H: ctx.Rect.H, Baseline: ctx.Baseline}
	}
}
