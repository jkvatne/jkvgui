package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/theme"
)

var ProgressStyle = ContainerStyle{
	BorderRole:     theme.Outline,
	BorderWidth:    1,
	Role:           theme.PrimaryContainer,
	CornerRadius:   2.0,
	InsidePadding:  f32.Padding{L: 2, T: 1, R: 2, B: 1},
	OutsidePadding: f32.Padding{L: 2, T: 2, R: 2, B: 2},
	Height:         16,
}

func ProgressBar(fraction float32, style *ContainerStyle) Wid {
	Default(&style, &ProgressStyle)
	fraction = min(1.0, max(fraction, 0))
	return func(ctx Ctx) Dim {
		h := style.Height
		if h < 1.0 {
			h = 16
		}
		if ctx.Mode != RenderChildren {
			return Dim{W: style.Width, H: h}
		}
		barRect := ctx.Rect
		barRect.H = h
		barRect = barRect.Inset(style.OutsidePadding, style.BorderWidth)
		// Draw track
		ctx.Win.Gd.RoundedRect(barRect, style.CornerRadius, 0, style.Role.Bg(), style.Role.Fg())
		// Draw progress bar
		barRect = barRect.Inset(style.InsidePadding, style.BorderWidth)
		barRect.W *= fraction
		ctx.Win.Gd.RoundedRect(barRect, style.CornerRadius, 0, style.Role.Fg(), style.Role.Fg())

		return Dim{W: ctx.W, H: h}
	}
}
