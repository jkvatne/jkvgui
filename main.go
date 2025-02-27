package main

import (
	"github.com/jkvatne/jkvgui/gpu"
	"golang.org/x/image/colornames"
	"image/color"
)

func main() {
	window := gpu.InitWindow(91200, 99800, "Rounded rectangle demo", 1)
	defer gpu.Shutdown()
	gpu.InitOpenGL(colornames.White)
	for !window.ShouldClose() {
		gpu.StartFrame()
		gpu.Fonts[0].SetColor(0.0, 0.0, 3.0, 1.0)
		_ = gpu.Fonts[0].Printf(0, 100, 1.0, "Before frames")
		gpu.RoundedRect(50, 50, 550, 350, 20, 15, color.Transparent, colornames.Red)
		gpu.RoundedRect(650, 50, 350, 50, 10, 3, colornames.Lightgrey, colornames.Black)
		_ = gpu.Fonts[0].Printf(160, 150, 2.0, "After frames")
		gpu.HorLine(30, 850, 480, 1, colornames.Blue)
		gpu.HorLine(30, 850, 500, 10, colornames.Black)
		gpu.Rect(10, 10, float32(gpu.WindowWidth)-20, float32(gpu.WindowHeight)-20, 2, color.Transparent, colornames.Black)
		gpu.Text(660, 290, 30, 0, color.Black, "Roboto-Light")
		gpu.Text(660, 170, 22, 1, color.Black, "Roboto-Medium")
		gpu.Text(660, 210, 15, 2, color.Black, "Roboto-Regular")
		gpu.Text(660, 250, 30, 3, color.Black, "RobotoMono")
		gpu.EndFrame(500, window)
	}
}
