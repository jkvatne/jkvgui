package main

import (
	"fmt"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"golang.org/x/image/colornames"
	"jkvgui/glfont"
	"jkvgui/gpu"
	"log"
	"runtime"
)

const (
	windowWidth  = 1200
	windowHeight = 600
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

func main() {
	var err error
	runtime.LockOSThread()
	window := gpu.InitWindow(windowWidth, windowHeight, "Rounded rectangle demo")
	defer glfw.Terminate()
	gpu.InitOpenGL(colornames.White)
	gpu.GetMonitors()

	font, err = glfont.LoadFont("Roboto-Medium.ttf", 35, windowWidth, windowHeight)
	panicOn(err, "Loading Rboto-Medium.ttf")
	window.SetKeyCallback(KeyCallback)

	for !window.ShouldClose() {

		gpu.StartFrame()
		font.SetColor(0.0, 0.0, 1.0, 1.0)
		_ = font.Printf(0, 100, 1.0, "Before frames"+"\x00")

		gpu.StartRR()

		gpu.DrawRoundedRect(50, 50, 550, 350, 20, 5, colornames.Black, colornames.Red)
		gpu.DrawRoundedRect(650, 50, 350, 50, 10, 3, colornames.Lightgrey, colornames.Black)

		// Free memory
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		gl.BindVertexArray(0)
		gl.UseProgram(0)

		_ = font.Printf(660, 90, 1.0, "After frames"+"\x00")

		gpu.EndFrame(20, window)
	}
}
