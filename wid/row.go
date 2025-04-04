package wid

func (r *ContainerStyle) W(w float32) *ContainerStyle {
	rr := *r
	rr.Width = w
	return &rr
}

func (r *ContainerStyle) H(h float32) *ContainerStyle {
	rr := *r
	rr.Height = h
	return &rr
}

func Row(style *ContainerStyle, widgets ...Wid) Wid {
	if style == nil {
		style = ContStyle
	}
	maxH := float32(0)
	maxB := float32(0)
	sumW := float32(0)
	fracSumW := float32(0)
	emptyCount := 0
	dims := make([]Dim, len(widgets))

	return func(ctx Ctx) Dim {
		if ctx.Mode != RenderChildren {
			return Dim{W: style.Width, H: ctx.Rect.H}
		}
		ctx0 := ctx
		ctx0.Rect.W -= style.OutsidePadding.T + style.OutsidePadding.B + style.BorderWidth*2
		ctx0.Rect.H -= style.InsidePadding.L + style.InsidePadding.R + style.BorderWidth*2

		// Collect W for all children
		ctx0.Mode = CollectWidths
		sumW = 0.0
		fracSumW = 0.0
		emptyCount = 0
		for i, w := range widgets {
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
			for i, _ := range widgets {
				if dims[i].W < 1.0 {
					dims[i].W = freeW * dims[i].W / fracSumW
				}
			}
		} else if fracSumW == 0.0 && emptyCount > 0 && freeW > 0.0 {
			// Children with W=0 will share the free width equaly
			for i, _ := range widgets {
				if dims[i].W == 0.0 {
					dims[i].W = freeW / float32(emptyCount)
				}
			}
		}

		// Collect maxH for all children, given width
		ctx0.Mode = CollectHeights
		maxB, maxH = 0.0, 0.0
		for i, w := range widgets {
			ctx0.Rect.W = dims[i].W
			temp := w(ctx0)
			maxH = max(maxH, temp.H)
			maxB = max(maxB, dims[i].Baseline)
		}

		// Render children with fixed W/H
		ctx0.Mode = RenderChildren
		ctx0.Baseline = maxB
		ctx0.Rect.H = maxH
		sumW = 0.0
		for i, w := range widgets {
			ctx0.Rect.W = dims[i].W
			ctx0.Rect.H = maxH
			dims[i] = w(ctx0)
			sumW += dims[i].W
			ctx0.Rect.X += dims[i].W
		}
		return Dim{W: sumW, H: maxH, Baseline: maxB}

	}
}
