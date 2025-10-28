package wid

import (
	"github.com/jkvatne/jkvgui/f32"
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
	InsidePadding:  f32.Padding{},
	OutsidePadding: f32.Padding{L: 2, T: 2, R: 2, B: 2},
}

var Primary = ContainerStyle{
	BorderRole:     theme.Outline,
	BorderWidth:    1,
	Role:           theme.PrimaryContainer,
	CornerRadius:   0.0,
	InsidePadding:  f32.Padding{L: 2, T: 2, R: 2, B: 2},
	OutsidePadding: f32.Padding{L: 2, T: 2, R: 2, B: 2},
}

var Secondary = ContainerStyle{
	BorderRole:     theme.Outline,
	BorderWidth:    0,
	Role:           theme.SecondaryContainer,
	CornerRadius:   9.0,
	InsidePadding:  f32.Padding{L: 4, T: 4, R: 4, B: 4},
	OutsidePadding: f32.Padding{L: 4, T: 4, R: 4, B: 4},
}

func (style *ContainerStyle) Size(w, h, bw float32) *ContainerStyle {
	ss := *style
	ss.Width = w
	ss.Height = h
	ss.BorderWidth = bw
	return &ss
}

func Col(style *ContainerStyle, widgets ...Wid) Wid {
	if style == nil {
		style = ContStyle
	}
	hPad := style.OutsidePadding.T + style.OutsidePadding.B + 2*style.BorderWidth + style.InsidePadding.T + style.OutsidePadding.B
	dims := make([]Dim, len(widgets))
	h := make([]float32, len(widgets))
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
		fracSumH := float32(0.0)
		sumH := hPad
		for i, w := range widgets {
			dim := w(ctx0)
			h[i] = dim.H
			if h[i] >= 1.0 {
				sumH += h[i]
				ctx0.H -= h[i]
			} else if h[i] > 0.0 {
				fracSumH += h[i]
			} else {
				emptyCount++
			}
		}

		// Distribute Height

		freeH := max(ctx.Rect.H-sumH, 0)
		if fracSumH > 0.0 && freeH > 0.0 {
			// Distribute the free height according to fractions for each child
			for i := range widgets {
				if h[i] < 1.0 {
					h[i] = freeH * h[i] / fracSumH
				}
			}
		} else if fracSumH == 0.0 && emptyCount > 0 && freeH > 0.0 {
			// Children with H<1.0 will share the free width equally
			for i := range widgets {
				if h[i] < 1.0 {
					h[i] = freeH / float32(emptyCount)
				}
			}
		}

		sumH = style.OutsidePadding.H(style.BorderWidth)
		sumH += style.InsidePadding.H(0)
		for i := range dims {
			sumH += h[i]
		}
		sumH += style.OutsidePadding.T + style.BorderWidth*2

		if ctx.Mode == CollectHeights {
			if style.Width < 1.0 {
				return Dim{W: ctx.W, H: sumH}
			}
			return Dim{W: style.Width, H: sumH}
		}

		// Render children with fixed Scroller/H
		ctx0 = ctx
		ctx0.Rect = ctx0.Rect.Inset(style.OutsidePadding, style.BorderWidth)
		ctx0.Y += style.OutsidePadding.T + style.BorderWidth
		ctx0.H = sumH
		// Draw frame
		ctx.Win.Gd.RoundedRect(ctx0.Rect, style.CornerRadius, style.BorderWidth, style.Role.Bg(), theme.Outline.Fg())
		ctx0.Rect = ctx0.Rect.Inset(style.InsidePadding, 0)
		ctx0.Mode = RenderChildren
		ctx0.Baseline = 0
		for i, w := range widgets {
			ctx0.Rect.H = h[i]
			dims[i] = w(ctx0)
			ctx0.Rect.Y += dims[i].H // h[i]
		}
		sumH += style.OutsidePadding.B
		return Dim{W: ctx.W, H: sumH, Baseline: 0}

	}
}
