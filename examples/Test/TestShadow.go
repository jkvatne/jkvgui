package main

import (
	"github.com/jkvatne/jkvgui/callback"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/theme"
)

func main() {
	theme.SetDefaultPallete(true)
	window := gpu.InitWindow(0, 0, "Test", 1)
	defer gpu.Shutdown()
	callback.Initialize(window)
	gpu.ScaleX = 5
	gpu.ScaleY = 5
	for !window.ShouldClose() {
		gpu.BackgroundColor(f32.White)
		gpu.StartFrame()
		r := f32.Rect{10, 10, 30, 20}
		gpu.RoundedRect(r, 0, 0.5, f32.Transparent, f32.Black)
		r.X += 50
		gpu.RoundedRect(r, 5, 0.5, f32.Transparent, f32.Black)
		r.X += 50
		gpu.RoundedRect(r, 5, 2, f32.Transparent, f32.Black)
		r.X += 50
		gpu.RoundedRect(r, 9999, 1, f32.Transparent, f32.Black)
		r.X += 50
		gpu.RoundedRect(r, 0, 0.5, f32.LightBlue, f32.Black)
		r.X += 50
		gpu.RoundedRect(r, 5, 0.51, f32.LightBlue, f32.Black)
		r.X += 50
		gpu.RoundedRect(r, 5, 2, f32.LightBlue, f32.Black)
		r.X += 50
		gpu.RoundedRect(r, 9999, 1, f32.LightBlue, f32.Black)
		r.X = 10
		r.Y += 50
		gpu.RoundedRect(r, 6, 0.5, f32.Transparent, f32.Black)
		gpu.Shade(r, 6, f32.Shade, 3)
		r.X += 50
		gpu.RoundedRect(r, 6, 0.5, f32.Transparent, f32.Black)
		gpu.Shade(r, 6, f32.Shade, 6)
		r.X += 50
		gpu.RoundedRect(r, 6, 0.5, f32.Transparent, f32.Black)
		gpu.Shade(r, 6, f32.Shade, 10)
		r.X += 50
		gpu.RoundedRect(r, 999, 0.5, f32.Transparent, f32.Black)
		gpu.Shade(r, 999, f32.Shade, 3)
		r.X += 50
		gpu.RoundedRect(r, 999, 0.5, f32.Transparent, f32.Black)
		gpu.Shade(r, 999, f32.Shade, 6)
		r.X += 50
		gpu.RoundedRect(r, 999, 0.5, f32.Transparent, f32.Black)
		gpu.Shade(r, 999, f32.Shade, 10)
		r.X += 50

		gpu.EndFrame(5)
	}
}
