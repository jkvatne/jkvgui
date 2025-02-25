package main

import (
	"fmt"
	"github.com/go-gl/glfw/v3.3/glfw"
	"golang.org/x/image/colornames"
	"image/color"
	"jkvgui/gpu"
	"log"
)

// https://www.glfw.org/docs/latest/window_guide.html
func KeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	log.Printf("Key %v %v %v %v\n", key, scancode, action, mods)
}

var col [8]float32

var N = 10000

func MouseBtnCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		x, y := w.GetCursorPos()
		fmt.Printf("Mouse btn %d clicked at %0.1f,%0.1f\n", button, x, y)
	}
}

func ScrollCallback(w *glfw.Window, xoff float64, yoff float64) {
	fmt.Printf("Scroll dx=%v dy=%v\n", xoff, yoff)
}

func main() {
	window := gpu.InitWindow(1200, 800, "Rounded rectangle demo")
	defer gpu.Shutdown()
	gpu.InitOpenGL(colornames.White)
	gpu.LoadFont("Roboto-Medium.ttf")
	window.SetKeyCallback(KeyCallback)
	window.SetMouseButtonCallback(MouseBtnCallback)
	window.SetSizeCallback(gpu.SizeCallback)
	window.SetScrollCallback(ScrollCallback)

	for !window.ShouldClose() {
		gpu.StartFrame()
		gpu.Fonts[0].SetColor(0.0, 0.0, 1.0, 1.0)
		_ = gpu.Fonts[0].Printf(0, 100, 1.0, "Before frames"+"\x00")
		gpu.RoundedRect(50, 50, 550, 350, 20, 15, color.Transparent, colornames.Red)
		gpu.RoundedRect(650, 50, 350, 50, 10, 3, colornames.Lightgrey, colornames.Black)
		_ = gpu.Fonts[0].Printf(660, 90, 1.0, "After frames"+"\x00")
		gpu.HorLine(30, 850, 480, 1, colornames.Blue)
		gpu.HorLine(30, 850, 500, 10, colornames.Black)
		gpu.VertLine(20, 10, 600, 3, colornames.Black)
		gpu.Rect(10, 10, float32(gpu.WindowWidth)-20, float32(gpu.WindowHeight)-20, 2, color.Transparent, colornames.Black)
		gpu.EndFrame(50, window)
	}
}
