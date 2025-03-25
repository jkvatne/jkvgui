package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
)

type Dim struct {
	W        float32
	H        float32
	Baseline float32
}

type Ctx struct {
	Rect     f32.Rect
	Baseline float32
	Disabled bool
}

func (ctx Ctx) Alpha() float32 {
	if ctx.Disabled {
		return 0.3
	}
	return 1.0
}

func (ctx Ctx) Disable() Ctx {
	ctx.Disabled = true
	return ctx
}

func (ctx Ctx) Enable(enabled bool) Ctx {
	ctx.Disabled = !enabled
	return ctx
}

func DisableIf(disabler *bool, w Wid) Wid {
	return func(ctx Ctx) Dim {
		ctx.Disabled = *disabler
		return w(ctx)
	}
}

type Wid func(ctx Ctx) Dim

func Maximized() Ctx {
	return Ctx{Rect: f32.Rect{X: 0, Y: 0, W: gpu.WindowWidthDp, H: gpu.WindowHeightDp}, Baseline: 0}
}

type RowSetup int

const (
	Distribute RowSetup = iota
	Left
	Right
)

type ColSetup struct {
	Widths []float32
}

func Row(setup RowSetup, widgets ...Wid) Wid {
	return func(ctx Ctx) Dim {
		maxH := float32(0)
		maxB := float32(0)
		sumW := float32(0)
		ctx0 := Ctx{}
		emptyCount := 0
		dims := make([]Dim, len(widgets))
		for i, w := range widgets {
			dims[i] = w(ctx0)
			maxH = max(maxH, dims[i].H)
			maxB = max(maxB, dims[i].Baseline)
			sumW += dims[i].W
			if dims[i].W == 0 {
				emptyCount++
			}
		}
		if ctx.Rect.H == 0 {
			if setup == Distribute {
				return Dim{W: ctx.Rect.W / float32(len(widgets)), H: maxH, Baseline: maxB}
			} else {
				return Dim{W: sumW, H: maxH, Baseline: maxB}
			}
		}
		ctx1 := ctx
		ctx1.Rect.H = maxH
		ctx1.Baseline = maxB
		if setup == Left {
			// If empty elements found, the remaining space is distributed into the empty slots.
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
