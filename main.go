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
	"time"
)

const (
	windowWidth  = 1200
	windowHeight = 600
)

var vao uint32
var vbo uint32

var colors = [32]float32{
	0.1, 0.1, 0.1, 0.5,
	1.0, 1.0, 0.0, 0.5,
	0.1, 0.1, 0.1, 0.8,
	1.0, 0.5, 0.5, 0.2,
	0.5, 0.5, 0.5, 0.2,
	0.5, 0.5, 1.0, 0.2,
	0.5, 0.5, 1.0, 0.2,
	0.5, 0.5, 1.0, 0.2,
}

var triangles = []float32{
	650, 50, 0, 2, 40, 10, 650, 50, 1150, 450,
	1150, 50, 0, 2, 40, 10, 650, 50, 1150, 450,
	650, 450, 0, 2, 40, 10, 650, 50, 1150, 450,
	1150, 450, 0, 2, 40, 10, 650, 50, 1150, 450,
	1150, 50, 0, 2, 40, 10, 650, 50, 1150, 450,
	650, 450, 0, 2, 40, 10, 650, 50, 1150, 450,
}

var rpos = [2]float32{300, 300}
var rw = [2]float32{20, 5}
var halfbox = [2]float32{250, 250}

var triangle = []float32{
	50, 50,
	550, 50,
	50, 550,
	550, 550,
	550, 50,
	50, 550,
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

func DrawTriangle(prog uint32) {
	gl.UseProgram(prog)
	gl.BufferData(gl.ARRAY_BUFFER, len(triangle)*4, gl.Ptr(triangle), gl.STATIC_DRAW)
	// position attribute
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 2*4, nil)
	gl.EnableVertexAttribArray(1)

	// set screen resolution
	r1 := gl.GetUniformLocation(prog, gl.Str("resolution\x00"))
	gl.Uniform2f(r1, float32(windowWidth), float32(windowHeight))
	// Colors
	r2 := gl.GetUniformLocation(prog, gl.Str("colors\x00"))
	gl.Uniform4fv(r2, 8, &colors[0])
	// Set pos data
	r3 := gl.GetUniformLocation(prog, gl.Str("pos\x00"))
	gl.Uniform2f(r3, rpos[0], rpos[1])
	// Set halfbox
	r4 := gl.GetUniformLocation(prog, gl.Str("halfbox\x00"))
	gl.Uniform2f(r4, halfbox[0], halfbox[1])
	// Set radius/border width
	r5 := gl.GetUniformLocation(prog, gl.Str("rw\x00"))
	gl.Uniform2f(r5, rw[0], rw[1])

	// Do actual drawing
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
}

func DrawTriangles(prog uint32) {
	gl.UseProgram(prog)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(triangles), gl.Ptr(triangles), gl.STATIC_DRAW)
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
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
}

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
	gpu.InitOpenGL(colornames.Skyblue)
	gpu.GetMonitors()

	gpu.BackgroundColor(colornames.Skyblue)
	font, err = glfont.LoadFont("Roboto-Medium.ttf", 35, windowWidth, windowHeight)
	panicOn(err, "Loading Rboto-Medium.ttf")
	InitKeys(window)
	rectProg := gpu.CreateProgram(gpu.RectangleVertShaderSource, gpu.RectangleFragShaderSource)
	rrProg := gpu.CreateProgram(gpu.RectVertShaderSource, gpu.RectFragShaderSource)

	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)

	for !window.ShouldClose() {

		gpu.StartFrame()
		font.SetColor(0.0, 0.0, 1.0, 1.0)
		_ = font.Printf(0, 100, 1.0, "Before frames"+"\x00")

		gl.BindVertexArray(vao)
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		gl.Enable(gl.BLEND)
		gl.BlendEquation(gl.FUNC_ADD)
		gl.BlendFunc(gl.SRC_ALPHA, gl.SRC_ALPHA)
		N = 1000

		t1 := time.Now()
		for range N {
			DrawTriangle(rrProg) // 34% with n=1000
		}
		dt1 := time.Since(t1)
		t2 := time.Now()
		for range 1 {
			DrawTriangles(rectProg) // 33% with N=1000
		}
		dt2 := time.Since(t2)
		fmt.Printf("With uniforms %vuS, Via coordds %vuS\r", dt1.Microseconds(), dt2.Microseconds())
		// Free memory
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		gl.BindVertexArray(0)
		gl.UseProgram(0)

		_ = font.Printf(0, 70, 1.0, "After frames"+"\x00")

		gpu.EndFrame(20, window)
	}
}
