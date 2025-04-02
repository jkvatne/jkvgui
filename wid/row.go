package wid

type RowStyle struct {
	Dist Distribute
	W    []float32
}

type Distribute uint8

const (
	Start Distribute = iota
	End
	Middle
	Even
)

var DefaultRowStyle RowStyle = RowStyle{
	Dist: 0,
	W:    []float32{},
}

func Row(style *RowStyle, widgets ...Wid) Wid {
	return func(ctx Ctx) Dim {
		if style == nil {
			style = &DefaultRowStyle
		}
		maxH := float32(0)
		maxB := float32(0)
		sumW := float32(0)
		ctx0 := Ctx{}
		emptyCount := 0
		dims := make([]Dim, len(widgets))
		for i, w := range widgets {
			if len(style.W) > 0 {
				ctx0.Rect.W = ctx.Rect.W * style.W[i]
			}
			dims[i] = w(ctx0)
			maxH = max(maxH, dims[i].H)
			maxB = max(maxB, dims[i].Baseline)
			sumW += dims[i].W
			if dims[i].W == 0 {
				emptyCount++
			}
		}
		if !ctx.Draw {
			if len(style.W) != 0 {
				return Dim{W: ctx.Rect.W, H: ctx.Rect.H}
			} else if style.Dist == Even {
				return Dim{W: ctx.Rect.W / float32(len(widgets)), H: maxH, Baseline: maxB}
			} else {
				return Dim{W: sumW, H: maxH, Baseline: maxB}
			}
		}

		ctx1 := ctx
		ctx1.Rect.H = maxH
		ctx1.Baseline = maxB
		if len(style.W) > 0 {
			for i, w := range widgets {
				width := ctx.Rect.W * style.W[i]
				ctx1.Rect.W = width
				dim := w(ctx1)
				ctx1.Rect.X += dim.W
			}
			return Dim{W: sumW, H: maxH, Baseline: maxB}
		} else if style.Dist == End {
			// If empty elements are found, the remaining space is distributed into the empty slots.
			ctx1.Rect.X += ctx.Rect.W - sumW
			for i, w := range widgets {
				ctx1.Rect.W = dims[i].W
				_ = w(ctx1)
				ctx1.Rect.X += dims[i].W
			}
			return Dim{W: sumW, H: maxH, Baseline: maxB}
		} else if style.Dist == Start {
			// If empty elements are found, the remaining space is distributed into the empty slots.
			if emptyCount > 0 {
				remaining := ctx.Rect.W - sumW
				for i, d := range dims {
					if d.W == 0 {
						dims[i].W = remaining / float32(emptyCount)
					}
				}
			}
			for i, w := range widgets {
				ctx1.Rect.W = dims[i].W
				_ = w(ctx1)
				ctx1.Rect.X += dims[i].W
			}
			return Dim{W: sumW, H: maxH, Baseline: maxB}
		} else {
			// Distribute evenly in equal-sized widgets
			ctx1.Rect.W = ctx.Rect.W / float32(len(widgets))
			for _, w := range widgets {
				_ = w(ctx1)
				ctx1.Rect.X += ctx1.Rect.W
			}
			return Dim{W: sumW, H: maxH, Baseline: maxB}

		}
	}
}
