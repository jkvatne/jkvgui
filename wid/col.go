package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/theme"
)

type ContainerStyle struct {
	Width          float32
	Height         float32
	BorderRole     theme.UIRole
	BorderWidth    float32
	Role           theme.UIRole
	Color          f32.Color
	CornerRadius   float32
	InsidePadding  f32.Padding
	OutsidePadding f32.Padding
}

var ContStyle = &ContainerStyle{
	BorderRole:     theme.Transparent,
	BorderWidth:    0.0,
	Role:           theme.Transparent,
	CornerRadius:   0.0,
	InsidePadding:  f32.Padding{0, 0, 0, 0},
	OutsidePadding: f32.Padding{5, 5, 5, 5},
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
	if style == nil {
		style = ContStyle
	}
	sumH := style.OutsidePadding.T + style.OutsidePadding.B + 2*style.BorderWidth + style.InsidePadding.T + style.OutsidePadding.B
	dims := make([]Dim, len(widgets))
	fracSumH := float32(0.0)
	emptyCount := 0
	return func(ctx Ctx) Dim {
		if ctx.Mode == CollectWidths {
			return Dim{W: style.Width, H: style.Height}
		}
		// Correct for padding and border
		ctx0 := ctx
		ctx0.Rect.W -= style.OutsidePadding.T + style.OutsidePadding.B + style.BorderWidth*2
		ctx0.Rect.H -= style.OutsidePadding.L + style.OutsidePadding.R + style.BorderWidth*2
		ctx0.Rect.X += style.OutsidePadding.L + style.BorderWidth
		ctx0.Rect.Y += style.OutsidePadding.T + style.BorderWidth
		// Collect Height for all children
		ctx0.Mode = CollectHeights
		for i, w := range widgets {
			dims[i] = w(ctx0)
			if dims[i].H >= 1.0 {
				sumH += dims[i].H
				ctx0.H -= dims[i].H
			} else if dims[i].H > 0.0 {
				fracSumH += dims[i].H
			} else {
				emptyCount++
			}
		}

		// Distribute Height
		freeH := max(ctx.Rect.H-sumH, 0)
		if fracSumH > 0.0 && freeH > 0.0 {
			// Distribute the free height according to fractions for each child
			for i, _ := range widgets {
				if dims[i].H < 1.0 {
					dims[i].H = freeH * dims[i].H / fracSumH
				}
			}
		} else if fracSumH == 0.0 && emptyCount > 0 && freeH > 0.0 {
			// Children with H<1.0 will share the free width equally
			for i, _ := range widgets {
				if dims[i].H < 1.0 {
					dims[i].H = freeH / float32(emptyCount)
				}
			}
		}

		sumH = style.OutsidePadding.H(style.BorderWidth)
		sumH += style.InsidePadding.H(0)
		for i := range dims {
			sumH += dims[i].H
		}

		if ctx.Mode == CollectHeights {
			if style.Width < 1.0 {
				return Dim{W: ctx.W, H: sumH}
			}
			return Dim{W: style.Width, H: sumH}
		}

		// Render children with fixed Scroller/H
		ctx0 = ctx
		ctx0.Rect = ctx0.Rect.Inset(style.OutsidePadding, style.BorderWidth)
		// Draw frame
		if style.Role != theme.Transparent {
			gpu.RoundedRect(ctx0.Rect, style.CornerRadius, style.BorderWidth, style.Role.Bg(), theme.Outline.Fg())
		}
		ctx0.Rect = ctx0.Rect.Inset(style.InsidePadding, 0)
		ctx0.Mode = RenderChildren
		ctx0.Baseline = 0
		for i, w := range widgets {
			// ctx0.Rect.H = min(dims[i].H, ctx.Y+ctx.H-ctx0.Rect.Y)
			ctx0.Rect.H = dims[i].H
			dims[i] = w(ctx0)
			ctx0.Rect.Y += dims[i].H
		}
		return Dim{W: ctx.W, H: sumH, Baseline: 0}

	}
}
