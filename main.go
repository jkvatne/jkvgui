package main

import (
	"fmt"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"log"
	"runtime"
	"testglfont/glfont"

	// "testglfont/glfont"
)

const windowWidth = 2300
const windowHeight = 1200

var oldState glfw.Action
var triangle = []float32{
	0, 0.5, 0, // top
	-0.5, -0.5, 0, // left
	0.5, -0.5, 0, // right
}

func init() {
	runtime.LockOSThread()
}

func CreateWindow() *glfw.Window {
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Floating, glfw.True)
	window, _ := glfw.CreateWindow(int(windowWidth), int(windowHeight), "Font demo", nil, nil)
	window.MakeContextCurrent()
	glfw.SwapInterval(1)
	scaleX, scaleY := window.GetContentScale()
	log.Printf("Window scaleX=%v, scaleY=%v\n", scaleX, scaleY)
	return window
}

// https://www.glfw.org/docs/latest/window_guide.html
func KeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	log.Println("%v %v %v %v\n", key, scancode, action, mods)
}

func InitKeys(window *glfw.Window) {
	window.SetKeyCallback(KeyCallback)
}

// initOpenGL initializes OpenGL and returns an intiialized program.
func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

	prog := gl.CreateProgram()
	gl.LinkProgram(prog)
	return prog
}

func draw(window *glfw.Window, prog uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(prog)

	glfw.PollEvents()
	window.SwapBuffers()
}

func main() {

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	monitors := glfw.GetMonitors()
	for i, monitor := range monitors {
		mw, mh := monitor.GetPhysicalSize()
		x, y := monitor.GetPos()
		mode := monitor.GetVideoMode()
		h := mode.Height
		w := mode.Width
		log.Printf("Monitor %d, %vmmx%vmm, %vx%vpx,  pos: %v, %v\n", i+1, mw, mh, w, h, x, y)
	}

	window := CreateWindow()
	_ = initOpenGL()

	font, err := glfont.LoadFont("Roboto-Medium.ttf", 35, windowWidth, windowHeight)
	if err != nil {
		log.Panicf("LoadFont: %v", err)
	}

	InitKeys(window)

	// prog := initOpenGL()
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0.0, windowWidth, 0.0, windowHeight, -1.0, 1.0)
	gl.MatrixMode(gl.MODELVIEW)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.Color3f(1.0, 0.0, 1.0)
		gl.Begin(gl.POLYGON)
		gl.Vertex2i(-50, -90)
		gl.Vertex2i(100, -90)
		gl.Vertex2i(100, 150)
		gl.Vertex2i(-50, 150)
		gl.End()
		gl.Flush()

		// FPS=3 for 100*22*16=35200 labels! Dvs 10000 tests pr sec
		// set color and draw text
		font.SetColor(0.0, 0.0, 0.0, 1.0)
		// _ = font.Printf(float32((i&0xF)*120+10), float32(25+(i>>4)*50), 1.0, "Aøæ©")
		font.Printf(10, 50, 1.0, "Hello World!")

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
