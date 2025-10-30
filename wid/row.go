package wid

func Row(style *ContainerStyle, widgets ...Wid) Wid {
	Default(&style, ContStyle)
	w := make([]float32, len(widgets))

	return func(ctx Ctx) Dim {
		if style.Height > 0 && ctx.Mode == CollectHeights {
			return Dim{W: ctx.W, H: min(ctx.H, style.Height)}
		}

		ctx0 := ctx
		ctx0.Rect.W -= style.OutsidePadding.T + style.OutsidePadding.B + style.BorderWidth*2
		ctx0.Rect.H -= style.InsidePadding.L + style.InsidePadding.R + style.BorderWidth*2

		// Collect width for all children
		fracSumW := float32(0)
		sumW := float32(0)
		emptyCount := 0
		ctx0.Mode = CollectWidths
		for i, widget := range widgets {
			ctx0.Rect.W = ctx.W * (1 - fracSumW)
			dim := widget(ctx0)
			w[i] = dim.W
			if w[i] > 1.0 {
				sumW += w[i]
			} else if w[i] > 0.0 {
				fracSumW += w[i]
			} else {
				emptyCount++
			}
		}

		// Distribute Width
		freeW := max(ctx.Rect.W-sumW, 0)
		if fracSumW > 0.0 && freeW > 0.0 {
			// Distribute the free width according to fractions for each child
			for i := range widgets {
				if w[i] <= 1.0 {
					w[i] = freeW * w[i] / fracSumW
				}
			}
		} else if fracSumW == 0.0 && emptyCount > 0 && freeW > 0.0 {
			// Children with w=0 will share the free width equally
			for i := range widgets {
				if w[i] == 0.0 {
					w[i] = freeW / float32(emptyCount)
				}
			}
		}

		// Collect maxH for all children, given width
		ctx0.Mode = CollectHeights
		maxH := float32(0)
		maxB := float32(0)
		for i, widget := range widgets {
			ctx0.Rect.W = w[i]
			dim := widget(ctx0)
			if dim.H == 0.0 {
				dim.H = ctx.H
			}
			maxH = max(maxH, dim.H)
			maxB = max(maxB, dim.Baseline)
		}

		if ctx.Mode != RenderChildren {
			return Dim{W: style.Width, H: maxH}
		}

		ctx0.Mode = RenderChildren
		ctx0.Baseline = maxB
		ctx0.Rect.H = min(maxH, ctx0.Rect.H)
		sumW = 0.0
		for i, widget := range widgets {
			ctx0.Rect.W = w[i]
			dim := widget(ctx0)
			sumW += dim.W
			ctx0.Rect.X += w[i]
		}
		return Dim{W: sumW, H: maxH, Baseline: maxB}

	}
}
