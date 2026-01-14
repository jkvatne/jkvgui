package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/theme"
)

func Col(style *ContainerStyle, widgets ...Wid) Wid {
	Default(&style, ContStyle)
	hPad := style.TotalVerticalPadding()
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

		sumH = f32.Sum(h...) + style.TotalVerticalPadding()
		if ctx.Mode == CollectHeights {
			if style.Width < 1.0 {
				return Dim{W: ctx.W, H: sumH}
			}
			return Dim{W: style.Width, H: sumH}
		}

		// Render children with fixed Scroller/H
		ctx0 = ctx
		ctx0.Rect = ctx0.Rect.Inset(style.OutsidePadding, style.BorderWidth)
		ctx0.Y += style.InsidePadding.T
		ctx0.H = sumH
		// Draw frame
		ctx.Win.Gd.RoundedRect(ctx0.Rect, style.CornerRadius, style.BorderWidth, style.Role.Bg(), theme.Outline.Bg())
		ctx0.Rect = ctx0.Rect.Inset(style.InsidePadding, 0)
		ctx0.Mode = RenderChildren
		ctx0.Baseline = 0
		for i, w := range widgets {
			if h[i] < 0 {
				h[i] = 0
			}
			ctx0.Rect.H = h[i]
			dims[i] = w(ctx0)
			ctx0.Rect.Y += h[i]
		}
		return Dim{W: ctx.W, H: sumH, Baseline: 0}
	}
}
