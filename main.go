package main

import (
	"fmt"
	"github.com/go-gl/glfw/v3.3/glfw"
	"golang.org/x/image/colornames"
	"image/color"
	"jkvgui/glfont"
	"jkvgui/gpu"
	"log"
	"runtime"
)

// https://www.glfw.org/docs/latest/window_guide.html
func KeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	log.Printf("Key %v %v %v %v\n", key, scancode, action, mods)
}

var col [8]float32

var font *glfont.Font
var N = 10000

func panicOn(err error, s string) {
	if err != nil {
		panic(fmt.Sprintf("%s: %v", s, err))
	}
}

func MouseBtnCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		fmt.Printf("Mouse btn %v clicked\n", button)
	}
}

func ScrollCallback(w *glfw.Window, xoff float64, yoff float64) {
	fmt.Printf("Scroll dx=%v dy=%v\n", xoff, yoff)
}

func main() {
	var err error
	runtime.LockOSThread()
	window := gpu.InitWindow(1200, 800, "Rounded rectangle demo")
	defer glfw.Terminate()
	gpu.InitOpenGL(colornames.White)

	font, err = glfont.LoadFont("Roboto-Medium.ttf", 35, gpu.WindowWidth, gpu.WindowHeight)
	panicOn(err, "Loading Rboto-Medium.ttf")
	window.SetKeyCallback(KeyCallback)
	window.SetMouseButtonCallback(MouseBtnCallback)
	window.SetSizeCallback(gpu.SizeCallback)
	window.SetScrollCallback(ScrollCallback)

	for !window.ShouldClose() {
		font.UpdateResolution(gpu.WindowWidth, gpu.WindowHeight)
		gpu.StartFrame()
		font.SetColor(0.0, 0.0, 1.0, 1.0)
		_ = font.Printf(0, 100, 1.0, "Before frames"+"\x00")
		gpu.RoundedRect(50, 50, 550, 350, 20, 15, color.Transparent, colornames.Red)
		gpu.RoundedRect(650, 50, 350, 50, 10, 3, colornames.Lightgrey, colornames.Black)
		_ = font.Printf(660, 90, 1.0, "After frames"+"\x00")
		gpu.HorLine(30, 850, 480, 1, colornames.Blue)
		gpu.HorLine(30, 850, 500, 10, colornames.Black)
		gpu.VertLine(20, 10, 600, 3, colornames.Black)
		gpu.Rect(10, 10, float32(gpu.WindowWidth)-20, float32(gpu.WindowHeight)-20, 2, color.Transparent, colornames.Black)
		gpu.EndFrame(50, window)
	}
}
