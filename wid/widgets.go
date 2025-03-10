package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
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
type ColSetup struct {
	Widths []float32
}

func Row(setup RowSetup, widgets ...Wid) Wid {
	return func(ctx Ctx) Dim {
		maxY := float32(0)
		maxB := float32(0)
		sumW := float32(0)
		ctx0 := Ctx{}
		ne := 0
		dims := make([]Dim, len(widgets))
		for i, w := range widgets {
			dims[i] = w(ctx0)
			maxY = max(maxY, dims[i].h)
			maxB = max(maxB, dims[i].baseline)
			sumW += dims[i].w
			if dims[i].w == 0 {
				ne++
			}
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
		ctx1.Rect.H = maxY
		ctx1.Baseline = maxB
		for i, w := range widgets {
			_ = w(ctx1)
			ctx1.Rect.X += dims[i].w
		}
		gpu.RoundedRect(ctx.Rect.X, ctx.Rect.Y, ctx.Rect.W, maxY, 0, 1, f32.Transparent, f32.Color{0, 1, 0, 0.2})
		return Dim{w: sumW, h: maxY, baseline: maxB}
	}
}

func Col(setup ColSetup, widgets ...Wid) Wid {
	return func(ctx Ctx) Dim {
		TotHeight := float32(0.0)
		maxY := float32(0.0)
		if ctx.Rect.H == 0 {
			for _, w := range widgets {
				h := w(ctx).h
				maxY = max(maxY, h)
				TotHeight += h
			}
			return Dim{ctx.Rect.W, maxY * float32(len(widgets)), 0}
		} else {
			for _, w := range widgets {
				ctx.Rect.Y += w(ctx).h
			}
			return Dim{100, TotHeight, 0}
		}
	}
}

func Label(text string, size float32, p f32.Padding, fontNo int) Wid {
	return func(ctx Ctx) Dim {
		if ctx.Rect.H == 0 {
			height := (gpu.Fonts[fontNo].Ascent+gpu.Fonts[fontNo].Descent)*size/gpu.InitialSize + p.T + p.B
			width := gpu.Fonts[fontNo].Width(size, text)/gpu.InitialSize + p.L + p.R
			return Dim{w: width, h: height, baseline: gpu.Fonts[fontNo].Ascent*size/gpu.InitialSize + p.T}
		} else {
			gpu.Fonts[fontNo].SetColor(f32.Black)
			gpu.Fonts[fontNo].Printf(ctx.Rect.X+p.L, ctx.Rect.Y+p.T+ctx.Baseline, size, 0, text)
			return Dim{}
		}
	}
}

func Elastic() Wid {
	return func(ctx Ctx) Dim {
		return Dim{}
	}
}

func RR(r f32.RRect, t float32, fillColor f32.Color, frameColor f32.Color) {
	gpu.RoundedRect(r.X, r.Y, r.W, r.H, r.RR, t, fillColor, frameColor)
}
