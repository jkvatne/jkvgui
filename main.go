package main

import (
	"github.com/jkvatne/jkvgui/gpu"
)

type Dim struct {
	w float32
	h float32
}

type Wid func(ctx Ctx) Dim

type Ctx struct {
	x, y          float32
	width, height float32
}

type RowSetup struct {
	Height float32
}
type ColSetup struct {
	Widths []float32
}

func Row(setup RowSetup, widgets ...Wid) Wid {
	return func(ctx Ctx) Dim {
		maxY := float32(0)
		sumW := float32(0)
		ctx0 := Ctx{}
		ne := 0
		dims := make([]Dim, len(widgets))
		for i, w := range widgets {
			dims[i] = w(ctx0)
			maxY = max(maxY, dims[i].h)
			sumW += dims[i].w
			if dims[i].w == 0 {
				ne++
			}
		}

		if ne > 0 {
			remaining := ctx.width - sumW
			for i, d := range dims {
				if d.w == 0 {
					dims[i].w = remaining / float32(ne)
				}
			}
		}
		ctx1 := ctx
		ctx1.height = maxY
		ctx1.width = sumW
		for i, w := range widgets {
			_ = w(ctx1)
			ctx1.x += dims[i].w
		}
		gpu.RoundedRect(ctx.x, ctx.y, ctx.width, maxY, 0, 1, gpu.Transparent, gpu.Color{0, 1, 0, 0.2})
		return Dim{w: sumW, h: maxY}
	}
}

func Col(setup ColSetup, widgets ...Wid) Wid {
	return func(ctx Ctx) Dim {
		TotHeight := float32(0.0)
		maxY := float32(0.0)
		if ctx.height == 0 {
			for _, w := range widgets {
				h := w(ctx).h
				maxY = max(maxY, h)
				TotHeight += h
			}
			return Dim{ctx.width, maxY * float32(len(widgets))}
		} else {
			for _, w := range widgets {
				ctx.y += w(ctx).h
			}
			return Dim{100, TotHeight}
		}
	}
}

func Label(text string, size float32) Wid {
	return func(ctx Ctx) Dim {
		fontNo := 1
		if ctx.height == 0 {
			height := (gpu.Fonts[fontNo].Ascent + gpu.Fonts[fontNo].Descent) * size / gpu.InitialSize
			width := gpu.Fonts[fontNo].Width(size, text) / gpu.InitialSize
			return Dim{w: width, h: height}
		} else {
			gpu.Fonts[0].SetColor(0.0, 0.0, 0.0, 1.0)
			gpu.Fonts[1].Printf(ctx.x, ctx.y+gpu.Fonts[fontNo].Ascent*size/gpu.InitialSize, size, text)
			return Dim{}
		}
	}
}

func Elastic() Wid {
	return func(ctx Ctx) Dim {
		return Dim{}
	}
}

func Form() Wid {
	r := RowSetup{}
	w := Row(r,
		Label("Hellogyjpq", 16),
		Label("Worljpqy", 22),
		Label("Welcome!", 12),
		Label("Worljpqy", 32),
		Elastic(),
		Label("Welcome!", 12),
	)
	return w
}

func Draw() {
	// Calculate sizes
	form := Form()
	ctx := Ctx{x: 50, y: 300, height: 260, width: 500}
	_ = form(ctx)
	// gpu.Rect(ctx.x, ctx.y, dim.x, dim.y, 2, gpu.Transparent, gpu.Black)
	// gpu.Rect(ctx.x, ctx.y, ctx.width, ctx.height, 2, gpu.Transparent, gpu.Black)
}

func main() {
	window := gpu.InitWindow(91200, 99800, "Rounded rectangle demo", 1)
	defer gpu.Shutdown()
	gpu.InitOpenGL(gpu.White)

	for !window.ShouldClose() {
		gpu.StartFrame()
		gpu.RoundedRect(650, 50, 350, 50, 10, 2, gpu.Lightgrey, gpu.Blue)
		gpu.Fonts[3].Printf(50, 100, 12, "12 RobotoMono")
		gpu.Fonts[1].Printf(50, 124, 16, "16 Roboto-Medium")
		gpu.Fonts[0].Printf(50, 156, 24, "24 Roboto-Light")
		gpu.Fonts[2].Printf(50, 204, 32, "32 Roboto-Regular")
		// Black frame around the whole window
		gpu.Rect(10, 10, float32(gpu.WindowWidth)-20, float32(gpu.WindowHeight)-20, 2, gpu.Transparent, gpu.Red)
		Draw()

		gpu.EndFrame(500, window)
	}
}
