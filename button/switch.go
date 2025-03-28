package button

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/focus"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/mouse"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
)

type SwitchStyle struct {
	Height float32
	Pad    f32.Padding
}

var DefaultSwitchStyle = &SwitchStyle{
	Height: 14,
	Pad:    f32.Padding{2, 2, 2, 2},
}

func Switch(state *bool, label string, action func(), style *SwitchStyle, hint string) wid.Wid {
	return func(ctx wid.Ctx) wid.Dim {
		if style == nil {
			style = DefaultSwitchStyle
		}
		width := style.Height*13/8 + style.Pad.R + style.Pad.L
		height := style.Height + style.Pad.T + style.Pad.B
		if ctx.Rect.H == 0 {
			return wid.Dim{W: width, H: height, Baseline: 0}
		}
		ctx.Rect.W = width
		ctx.Rect.H = height
		r1 := f32.Rect{ctx.Rect.X + style.Pad.R, ctx.Rect.Y + style.Pad.T,
			style.Height * 13 / 8, style.Height}
		r2 := f32.Rect{
			X: r1.X + style.Height/4,
			Y: r1.Y + style.Height/4,
			W: style.Height / 2,
			H: style.Height / 2,
		}
		if *state {
			r2.X = r1.X + style.Height*7/8
		}
		if mouse.Hovered(r2) || focus.At(ctx.Rect, state) {
			gpu.Shade(r2.Outset(f32.Pad(4)), 999, f32.Shade, 5)
		}
		if mouse.LeftBtnClick(ctx.Rect) {
			focus.Set(state)
			*state = !*state
		}
		if *state == false {
			gpu.RoundedRect(r1, 999, style.Height/32.0, theme.SurfaceContainer.Bg(), theme.Outline.Fg())
			gpu.RoundedRect(r2, 999, 0.0, theme.Outline.Fg(), theme.Outline.Fg())
		} else {
			gpu.RoundedRect(r1, 999, style.Height/32.0, theme.Primary.Bg(), theme.Primary.Bg())
			gpu.RoundedRect(r2, 999, 0.0, theme.Primary.Fg(), theme.Primary.Fg())
		}
		return wid.Dim{}
	}
}
