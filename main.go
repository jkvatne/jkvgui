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

func draw20(prog uint32) {
	gl.UseProgram(prog)
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
	gl.DrawArrays(gl.TRIANGLES, 0, 12)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
	gl.UseProgram(0)
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

		for range N {
			DrawTriangles(rectProg)
		}

		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		gl.BindVertexArray(0)

		_ = font.Printf(0, 70, 1.0, "After frames"+"\x00")

		gpu.EndFrame(20, window)
	}
}
