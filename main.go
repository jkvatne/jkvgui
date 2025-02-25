package main

import (
	"fmt"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"jkvgui/glfont"
	"jkvgui/gpu"
	"log"
	"runtime"
	"time"
)

const (
	windowWidth  = 1200
	windowHeight = 600
)

var vao uint32
var vbo uint32

var colors = []float32{
	1.0, 0.5, 0.5, 0.5,
	0.5, 1.0, 0.5, 0.5,
	0.1, 0.1, 0.1, 0.8,
	1.0, 0.5, 0.5, 0.2,
	0.5, 0.5, 0.5, 0.2,
	0.5, 0.5, 1.0, 0.2,
	0.5, 0.5, 1.0, 0.2,
	0.5, 0.5, 1.0, 0.2,
}

var triangles = []float32{
	50, 50, 1, 2, 20, 5, 50, 50, 550, 550,
	550, 50, 1, 2, 20, 5, 50, 50, 550, 550,
	50, 550, 1, 2, 20, 5, 50, 50, 550, 550,
	550, 550, 1, 2, 20, 5, 50, 50, 550, 550,
	550, 50, 1, 2, 20, 5, 50, 50, 550, 550,
	50, 550, 1, 2, 20, 5, 50, 50, 550, 550,

	650, 50, 0, 2, 40, 10, 650, 50, 1150, 450,
	1150, 50, 0, 2, 40, 10, 650, 50, 1150, 450,
	650, 450, 0, 2, 40, 10, 650, 50, 1150, 450,
	1150, 450, 0, 2, 40, 10, 650, 50, 1150, 450,
	1150, 50, 0, 2, 40, 10, 650, 50, 1150, 450,
	650, 450, 0, 2, 40, 10, 650, 50, 1150, 450,
}

var pos int

func add(a, b, c, d, e, f, g, h, i, j float32) {
	triangles[pos] = a
	triangles[pos+1] = b
	triangles[pos+2] = c
	triangles[pos+3] = d
	triangles[pos+4] = e
	triangles[pos+5] = f
	triangles[pos+6] = g
	triangles[pos+7] = h
	triangles[pos+8] = i
	triangles[pos+9] = j
	pos += 10
}

func setupFrame(x, y, w, h, rr, b, fillColor, frameColor float32) {
	add(x, y, fillColor, frameColor, rr, b, x, y, x+w, y+h)
	add(x+w, y, fillColor, frameColor, rr, b, x, y, x+w, y+h)
	add(x, y+h, fillColor, frameColor, rr, b, x, y, x+w, y+h)
	add(x+w, y+h, fillColor, frameColor, rr, b, x, y, x+w, y+h)
	add(x+w, y, fillColor, frameColor, rr, b, x, y, x+w, y+h)
	add(x, y+h, fillColor, frameColor, rr, b, x, y, x+w, y+h)
}

// https://www.glfw.org/docs/latest/window_guide.html
func KeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	log.Printf("Key %v %v %v %v\n", key, scancode, action, mods)
}

func InitKeys(window *glfw.Window) {
	window.SetKeyCallback(KeyCallback)
}

// initGlfw initializes glfw and returns a Window to use.
func initGlfw(width, height int, name string) *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.False)
	glfw.WindowHint(glfw.Floating, glfw.False) // Will keep window on top if true

	window, err := glfw.CreateWindow(width, height, name, nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	glfw.SwapInterval(1)
	scaleX, scaleY := window.GetContentScale()
	log.Printf("Window scaleX=%v, scaleY=%v\n", scaleX, scaleY)

	return window
}

func DrawTriangles(prog uint32) {
	gl.UseProgram(prog)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(triangles), gl.Ptr(triangles), gl.STATIC_DRAW)
	// gl.BufferSubData(gl.ARRAY_BUFFER, 0, 4*len(triangles), gl.Ptr(triangles))

	// position attribute
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 10*4, nil)
	gl.EnableVertexAttribArray(1)
	// color attribute gl.VertexAttribPointerWithOffset()
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 10*4, gl.PtrOffset(2*4))
	gl.EnableVertexAttribArray(2)
	// radius-width attribute
	gl.VertexAttribPointer(3, 2, gl.FLOAT, false, 10*4, gl.PtrOffset(4*4))
	gl.EnableVertexAttribArray(3)
	// rectangel
	gl.VertexAttribPointer(4, 4, gl.FLOAT, false, 10*4, gl.PtrOffset(6*4))
	gl.EnableVertexAttribArray(4)

	// set screen resolution
	resUniform := gl.GetUniformLocation(prog, gl.Str("resolution\x00"))
	gl.Uniform2f(resUniform, float32(windowWidth), float32(windowHeight))

	r2 := gl.GetUniformLocation(prog, gl.Str("colors\x00"))
	gl.Uniform4fv(r2, 12, &colors[0])

	// Do actual drawing
	gl.DrawArrays(gl.TRIANGLES, 0, 12)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
	gl.UseProgram(0)
}

var font *glfont.Font
var N = 10000

func LoadFonts() {
	var err error
	font, err = glfont.LoadFont("Roboto-Medium.ttf", 35, windowWidth, windowHeight)
	if err != nil {
		log.Panicf("LoadFont: %v", err)
	}
}

func main() {
	runtime.LockOSThread()
	window := initGlfw(windowWidth, windowHeight, "Demo")
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

	gpu.InitOpenGL()
	gl.ClearColor(0.95, 0.95, 0.86, 0.10)

	LoadFonts()
	InitKeys(window)
	rectProg := gpu.CreateProgram(gpu.RectangleVertShaderSource, gpu.RectangleFragShaderSource)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		// FPS=3 for 100*22*16=35200 labels!
		t := time.Now()
		font.SetColor(0.0, 0.0, 1.0, 1.0)
		_ = font.Printf(0, 100, 1.0, "Before frames"+"\x00")

		gl.Enable(gl.BLEND)
		gl.GenVertexArrays(1, &vao)
		gl.BindVertexArray(vao)
		gl.GenBuffers(1, &vbo)
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)

		for range N {
			DrawTriangles(rectProg)
		}
		gl.BindVertexArray(0)

		_ = font.Printf(0, 70, 1.0, "After frames"+"\x00")
		window.SwapBuffers()
		fmt.Printf("Frames pr second: %0.1f\r", float64(N)/time.Since(t).Seconds())

		glfw.PollEvents()
		runtime.GC()
		time.Sleep(100 * time.Millisecond)
	}

	glfw.Terminate()

}
