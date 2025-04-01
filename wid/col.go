package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/theme"
)

type ContainerStyle struct {
	Widths         []float32
	BorderRole     theme.UIRole
	BorderWidth    float32
	Role           theme.UIRole
	CornerRadius   float32
	InsidePadding  f32.Padding
	OutsidePadding f32.Padding
	Dist           Distribute
}

var Default = ContainerStyle{
	BorderRole:     theme.Outline,
	BorderWidth:    0.0,
	Role:           theme.Surface,
	CornerRadius:   0.0,
	InsidePadding:  f32.Padding{0, 0, 0, 0},
	OutsidePadding: f32.Padding{0, 0, 0, 0},
	Dist:           0,
}

var Primary = ContainerStyle{
	BorderRole:     theme.Outline,
	BorderWidth:    0,
	Role:           theme.PrimaryContainer,
	CornerRadius:   9.0,
	InsidePadding:  f32.Padding{4, 4, 4, 4},
	OutsidePadding: f32.Padding{4, 4, 4, 4},
}

var Secondary = ContainerStyle{
	BorderRole:     theme.Outline,
	BorderWidth:    0,
	Role:           theme.SecondaryContainer,
	CornerRadius:   9.0,
	InsidePadding:  f32.Padding{4, 4, 4, 4},
	OutsidePadding: f32.Padding{4, 4, 4, 4},
}

func Col(style *ContainerStyle, widgets ...Wid) Wid {
	return func(ctx Ctx) Dim {
		if style == nil {
			style = &Default
		}
		sumH := float32(0.0)
		ctx0 := Ctx{}
		ctx0.Rect.W = ctx.Rect.W
		ne := 0
		maxW := float32(0)
		dims := make([]Dim, len(widgets))
		// Calculate sum of minimum heights for all children
		for i, w := range widgets {
			dims[i] = w(ctx0)
			maxW = max(maxW, dims[i].W)
			sumH += dims[i].H
			if dims[i].W == 0 {
				ne++
			}
		}
		sumH += style.OutsidePadding.T + style.OutsidePadding.B + style.BorderWidth*2
		sumH += style.InsidePadding.T + style.InsidePadding.B
		maxW += style.InsidePadding.L + style.InsidePadding.R + style.BorderWidth*2
		if ctx.Rect.H == 0 {
			return Dim{W: maxW, H: sumH, Baseline: 0}
		}

		if ne > 0 {
			remaining := ctx.Rect.H - sumH
			for i, d := range dims {
				if d.H == 0 {
					dims[i].H = remaining / float32(ne)
				}
			}
		}
		ctx1 := ctx
		ctx1.Rect = ctx.Rect.Inset(style.OutsidePadding, style.BorderWidth)
		gpu.RoundedRect(ctx1.Rect, style.CornerRadius, style.BorderWidth, style.Role.Bg(), theme.Outline.Fg())
		ctx1.Rect = ctx1.Rect.Inset(style.InsidePadding, 0)
		for i, w := range widgets {
			ctx1.Rect.H = dims[i].H
			w(ctx1)
			ctx1.Rect.Y += dims[i].H
		}
		return Dim{ctx.Rect.W, sumH, 0}
	}
}
