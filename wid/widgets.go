package wid

import (
	"github.com/jkvatne/jkvgui/f32"
)

type Dim struct {
	w        float32
	h        float32
	baseline float32
}

type Ctx struct {
	Rect     f32.Rect
	Baseline float32
}

type Wid func(ctx Ctx) Dim

type RowSetup struct {
	Height float32
}

var DefaultRowSetup = RowSetup{
	Height: 0,
}

type ColSetup struct {
	Widths []float32
}

func Row(setup *RowSetup, widgets ...Wid) Wid {
	if setup == nil {
		setup = &DefaultRowSetup
	}
	return func(ctx Ctx) Dim {
		maxH := float32(0)
		maxB := float32(0)
		sumW := float32(0)
		ctx0 := Ctx{}
		ne := 0
		dims := make([]Dim, len(widgets))
		for i, w := range widgets {
			dims[i] = w(ctx0)
			maxH = max(maxH, dims[i].h)
			maxB = max(maxB, dims[i].baseline)
			sumW += dims[i].w
			if dims[i].w == 0 {
				ne++
			}
		}
		if ctx.Rect.H == 0 {
			return Dim{w: sumW, h: maxH, baseline: maxB}
		}
		if ne > 0 {
			remaining := ctx.Rect.W - sumW
			for i, d := range dims {
				if d.w == 0 {
					dims[i].w = remaining / float32(ne)
				}
			}
		}
		ctx1 := ctx
		ctx1.Rect.H = maxH
		ctx1.Baseline = maxB
		for i, w := range widgets {
			_ = w(ctx1)
			ctx1.Rect.X += dims[i].w
		}
		// gpu.RoundedRect(ctx.Rect.X, ctx.Rect.Y, ctx.Rect.W, maxY, 0, 1, f32.Transparent, f32.Color{0, 1, 0, 0.2}, 0)
		return Dim{w: sumW, h: maxH, baseline: maxB}
	}
}

func Col(setup *ColSetup, widgets ...Wid) Wid {
	return func(ctx Ctx) Dim {
		TotHeight := float32(0.0)
		ctx0 := ctx
		ctx0.Rect.H = 0
		h := make([]float32, len(widgets))
		for i, w := range widgets {
			h[i] = w(ctx0).h
			TotHeight += h[i]
		}
		for i, w := range widgets {
			ctx.Rect.H = h[i]
			w(ctx)
			ctx.Rect.Y += h[i]
		}
		return Dim{100, TotHeight, 0}
	}
}

func Elastic() Wid {
	return func(ctx Ctx) Dim {
		return Dim{}
	}
}
