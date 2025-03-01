package main

import (
	"github.com/jkvatne/jkvgui/gpu"
)

type Dim struct {
	x float32
	y float32
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
		sumX := float32(0)
		ctx0 := Ctx{}
		for _, w := range widgets {
			dim := w(ctx0)
			maxY = max(maxY, dim.y)
			sumX += dim.x
		}
		ctx1 := ctx
		ctx1.height = maxY
		ctx1.width = sumX
		for _, w := range widgets {
			dim := w(ctx1)
			ctx1.x += dim.x
		}
		gpu.RoundedRect(ctx.x, ctx.y, ctx.width, ctx.height, 0, 1, gpu.Transparent, gpu.Black)
		return Dim{x: sumX, y: maxY}
	}
}

func Col(setup ColSetup, widgets ...Wid) Wid {
	return func(ctx Ctx) Dim {
		TotHeight := float32(0.0)
		maxY := float32(0.0)
		if ctx.height == 0 {
			for _, w := range widgets {
				h := w(ctx).y
				maxY = max(maxY, h)
				TotHeight += h
			}
			return Dim{ctx.width, maxY * float32(len(widgets))}
		} else {
			for _, w := range widgets {
				ctx.y += w(ctx).y
			}
			return Dim{100, TotHeight}
		}
	}
}

func Label(text string, size float32) Wid {
	return func(ctx Ctx) Dim {
		if ctx.height == 0 {
			height := size + 8.0
			return Dim{x: float32(len(text)) * size, y: height}
		} else {
			gpu.Fonts[0].SetColor(0.0, 0.0, 0.0, 1.0)
			gpu.Text(ctx.x, ctx.y+size, size, 1, gpu.Black, text)
			return Dim{x: float32(len(text)) * size, y: ctx.height}
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
		Label("Hello", 16),
		Label("World", 22),
		Elastic(),
		Label("Welcome!", 22),
	)
	return w
}

func Draw() {
	// Calculate sizes
	form := Form()
	ctx := Ctx{x: 50, y: 300, height: 260, width: 500}
	dim := form(ctx)
	gpu.Rect(ctx.x, ctx.y, dim.x, dim.y, 2, gpu.Transparent, gpu.Black)
	// gpu.Rect(ctx.x, ctx.y, ctx.width, ctx.height, 2, gpu.Transparent, gpu.Black)
}

func main() {
	window := gpu.InitWindow(91200, 99800, "Rounded rectangle demo", 1)
	defer gpu.Shutdown()
	gpu.InitOpenGL(gpu.White)

	for !window.ShouldClose() {
		gpu.StartFrame()

		gpu.Fonts[0].SetColor(0.0, 0.0, 3.0, 1.0)

		gpu.RoundedRect(650, 50, 350, 50, 10, 2, gpu.Lightgrey, gpu.Black)

		gpu.Text(50, 100, 12, 3, gpu.Black, "12 RobotoMono")
		gpu.Text(50, 124, 16, 1, gpu.Black, "16 Roboto-Medium")
		gpu.Text(50, 156, 24, 0, gpu.Black, "24 Roboto-Light")
		gpu.Text(50, 204, 32, 2, gpu.Black, "32 Roboto-Regular")
		// Black frame around the whole window
		gpu.Rect(10, 10, float32(gpu.WindowWidth)-20, float32(gpu.WindowHeight)-20, 2, gpu.Transparent, gpu.Black)
		Draw()

		gpu.EndFrame(500, window)
	}
}
