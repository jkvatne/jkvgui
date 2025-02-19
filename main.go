package main

import (
	"fmt"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"log"
	"runtime"
	"testglfont/glfont"
)

const windowWidth = 2300
const windowHeight = 1200

var oldState glfw.Action

func init() {
	runtime.LockOSThread()
}

func KeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	fmt.Printf("%v %v %v %v\n", key, scancode, action, mods)
}

func main() {

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Floating, glfw.True)

	monitors := glfw.GetMonitors()
	for i, monitor := range monitors {
		x, y := monitor.GetPhysicalSize()
		fmt.Printf("%d Monitor size: %v, %v\n", i+1, x, y)
		x, y = monitor.GetPos()
		fmt.Printf("%d Monitor pos: %v, %v\n", i+1, x, y)
		mode := monitor.GetVideoMode()
		h := mode.Height
		w := mode.Width
		fmt.Printf("%d w=%d, h=%d\n", i+1, w, h)
	}
	window, _ := glfw.CreateWindow(int(windowWidth), int(windowHeight), "Font demo", nil, nil)

	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	if err := gl.Init(); err != nil {
		panic(err)
	}

	font, err := glfont.LoadFont("Roboto-Medium.ttf", 35, windowWidth, windowHeight)
	if err != nil {
		log.Panicf("LoadFont: %v", err)
	}

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

	window.SetKeyCallback(KeyCallback)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		// t := time.Now()
		// FPS=3 for 100*22*16=35200 labels!
		for i := range 22 * 16 {
			// set color and draw text
			font.SetColor(0.0, 0.0, 0.0, 1.0)
			_ = font.Printf(float32((i&0xF)*120+10), float32(25+(i>>4)*50), 1.0, "Aøæ©")
		}
		window.SwapBuffers()
		glfw.PollEvents()
		// fmt.Printf("FPS: %v\r", 1.0/time.Since(t).Seconds())
		state := window.GetKey(glfw.KeyA)
		if state != oldState {
			fmt.Printf("Escape %v\n", state)
			oldState = state
		}
		oldState = state
	}
}
