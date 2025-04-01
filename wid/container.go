package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/theme"
)

type xContainerStyle struct {
	BorderRole     theme.UIRole
	BorderWidth    float32
	Role           theme.UIRole
	CornerRadius   float32
	InsidePadding  f32.Padding
	OutsidePadding f32.Padding
}

var DefaultContainerStyle = xContainerStyle{
	BorderRole:     theme.Outline,
	BorderWidth:    1.5,
	Role:           theme.PrimaryContainer,
	CornerRadius:   5.0,
	InsidePadding:  f32.Padding{4, 4, 4, 4},
	OutsidePadding: f32.Padding{4, 4, 4, 4},
}

func Container(style *xContainerStyle, widget Wid) Wid {
	return func(ctx Ctx) Dim {
		if style == nil {
			style = &DefaultContainerStyle
		}
		if ctx.Rect.H == 0 {
			dim := widget(ctx)
			ctx.Rect.H = dim.H
			ctx.Rect.W = dim.W
			return Dim{W: ctx.Rect.W, H: ctx.Rect.H, Baseline: 0}
		}
		ctx.Rect = ctx.Rect.Inset(style.OutsidePadding, style.BorderWidth)
		gpu.RoundedRect(ctx.Rect, style.CornerRadius, style.BorderWidth, style.Role.Bg(), theme.Outline.Fg())
		return widget(ctx)
	}
}
