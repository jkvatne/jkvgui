package wid

import (
	"github.com/jkvatne/jkvgui/f32"
)

type Dim struct {
	W        float32
	H        float32
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
			maxH = max(maxH, dims[i].H)
			maxB = max(maxB, dims[i].baseline)
			sumW += dims[i].W
			if dims[i].W == 0 {
				ne++
			}
		}
		if ctx.Rect.H == 0 {
			return Dim{W: sumW, H: maxH, baseline: maxB}
		}
		if ne > 0 {
			remaining := ctx.Rect.W - sumW
			for i, d := range dims {
				if d.W == 0 {
					dims[i].W = remaining / float32(ne)
				}
			}
		}
		ctx1 := ctx
		ctx1.Rect.H = maxH
		ctx1.Baseline = maxB
		for i, w := range widgets {
			ctx1.Rect.W = dims[i].W
			_ = w(ctx1)
			ctx1.Rect.X += dims[i].W
		}
		return Dim{W: sumW, H: maxH, baseline: maxB}
	}
}

func Col(setup *ColSetup, widgets ...Wid) Wid {
	return func(ctx Ctx) Dim {
		sumH := float32(0.0)
		ctx0 := ctx
		ctx0.Rect.H = 0
		ne := 0
		maxW := float32(0)
		dims := make([]Dim, len(widgets))
		for i, w := range widgets {
			dims[i] = w(ctx0)
			maxW = max(maxW, dims[i].W)
			sumH += dims[i].H
			if dims[i].W == 0 {
				ne++
			}
		}
		if ctx.Rect.H == 0 {
			return Dim{W: maxW, H: sumH, baseline: 0}
		}
		if ne > 0 {
			remaining := ctx.Rect.H - sumH
			for i, d := range dims {
				if d.H == 0 {
					dims[i].H = remaining / float32(ne)
				}
			}
		}
		for i, w := range widgets {
			ctx.Rect.H = dims[i].H
			w(ctx)
			ctx.Rect.Y += dims[i].H
		}
		return Dim{100, sumH, 0}
	}
}

func Elastic() Wid {
	return func(ctx Ctx) Dim {
		return Dim{}
	}
}
