package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/theme"
)

func (style *ContainerStyle) W(w float32) *ContainerStyle {
	rr := *style
	rr.Width = w
	return &rr
}

func (style *ContainerStyle) H(h float32) *ContainerStyle {
	rr := *style
	rr.Height = h
	return &rr
}

func (style *ContainerStyle) R(c theme.UIRole) *ContainerStyle {
	rr := *style
	rr.Role = c
	return &rr
}

func (style *ContainerStyle) C(c f32.Color) *ContainerStyle {
	rr := *style
	rr.Color = c
	return &rr
}

func Row(style *ContainerStyle, widgets ...Wid) Wid {
	if style == nil {
		style = ContStyle
	}
	dims := make([]Dim, len(widgets))

	return func(ctx Ctx) Dim {
		if style.Height > 0 && ctx.Mode == CollectHeights {
			return Dim{W: ctx.W, H: style.Height}
		}

		ctx0 := ctx
		ctx0.Rect.W -= style.OutsidePadding.T + style.OutsidePadding.B + style.BorderWidth*2
		ctx0.Rect.H -= style.InsidePadding.L + style.InsidePadding.R + style.BorderWidth*2

		// Collect width for all children
		fracSumW := float32(0)
		sumW := float32(0)
		emptyCount := 0
		ctx0.Mode = CollectWidths
		for i, w := range widgets {
			ctx0.Rect.W = ctx.W * (1 - fracSumW)
			dims[i] = w(ctx0)
			if dims[i].W > 1.0 {
				sumW += dims[i].W
			} else if dims[i].W > 0.0 {
				fracSumW += dims[i].W
			} else {
				emptyCount++
			}
		}

		// Distribute Width
		freeW := max(ctx.Rect.W-sumW, 0)
		if fracSumW > 0.0 && freeW > 0.0 {
			// Distribute the free width according to fractions for each child
			for i := range widgets {
				if dims[i].W <= 1.0 {
					dims[i].W = freeW * dims[i].W / fracSumW
				}
			}
		} else if fracSumW == 0.0 && emptyCount > 0 && freeW > 0.0 {
			// Children with w=0 will share the free width equally
			for i := range widgets {
				if dims[i].W == 0.0 {
					dims[i].W = freeW / float32(emptyCount)
				}
			}
		}

		// Collect maxH for all children, given width
		ctx0.Mode = CollectHeights
		maxH := float32(0)
		maxB := float32(0)
		for i, w := range widgets {
			ctx0.Rect.W = dims[i].W
			temp := w(ctx0)
			if temp.H == 0.0 {
				temp.H = ctx.H
			}
			maxH = max(maxH, temp.H)
			maxB = max(maxB, dims[i].Baseline)
		}

		if ctx.Mode != RenderChildren {
			return Dim{W: style.Width, H: maxH}
		}

		ctx0.Mode = RenderChildren
		ctx0.Baseline = maxB
		ctx0.Rect.H = min(maxH, ctx0.Rect.H)
		sumW = 0.0
		for i, w := range widgets {
			ctx0.Rect.W = dims[i].W
			dim := w(ctx0)
			sumW += dim.W
			ctx0.Rect.X += dims[i].W
		}
		return Dim{W: sumW, H: maxH, Baseline: maxB}

	}
}
