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
	height         float32
	OutsidePadding f32.Padding
}

var DefaultSwitchStyle = &SwitchStyle{
	height:         24,
	OutsidePadding: f32.Padding{5, 5, 5, 5},
}

func Switch(state *bool, action func(), style *SwitchStyle, hint string) wid.Wid {
	return func(ctx wid.Ctx) wid.Dim {
		if style == nil {
			style = DefaultSwitchStyle
		}
		if ctx.Rect.H == 0 {
			return wid.Dim{W: style.height*52/32 + style.OutsidePadding.R + style.OutsidePadding.L,
				H: style.height + style.OutsidePadding.T + style.OutsidePadding.B, Baseline: 0}
		}
		r1 := f32.Rect{ctx.Rect.X + style.OutsidePadding.R, ctx.Rect.Y + style.OutsidePadding.T,
			style.height * 52 / 32, style.height}
		r2 := f32.Rect{
			X: r1.X + style.height/4,
			Y: r1.Y + style.height/4,
			W: style.height / 2,
			H: style.height / 2,
		}
		if *state {
			r2.X = r1.X + style.height*7/8
		}
		if mouse.LeftBtnPressed(ctx.Rect) {
			// gpu.Shade(r.Outset(f32.Padding{4, 4, 4, 4}).Move(0, 0), cr, f32.Shade, 4)
			// b += 1
		} else if mouse.Hovered(r1) {
			gpu.Shade(r2.Outset(f32.Padding{8, 8, 8, 8}), 5, f32.Shade, 4)
		}
		if mouse.LeftBtnReleased(ctx.Rect) {
			focus.Set(state)
			*state = !*state
		}
		if focus.At(ctx.Rect, state) {
			// b += 1
		}

		if *state == false {
			gpu.RoundedRect(r1, 999, style.height/32.0, theme.SurfaceContainer.Bg(), theme.Outline.Fg())
			gpu.RoundedRect(r2, 999, 0.0, theme.Outline.Fg(), theme.Outline.Fg())
		} else {
			gpu.RoundedRect(r1, 999, style.height/32.0, theme.Primary.Bg(), theme.Primary.Bg())
			gpu.RoundedRect(r2, 999, 0.0, theme.Primary.Fg(), theme.Primary.Fg())
		}

		r1 = r1.Out(12)
		gpu.Shade(r1, 999, f32.Shade, 3)

		return wid.Dim{}
	}
}
