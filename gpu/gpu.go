package gpu

import (
	"fmt"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"image"
	"image/color"
	"log"
	"strings"
	"time"
)

var startTime time.Time
var vao uint32
var vbo uint32
var windowWidth int
var windowHeight int

// https://github.com/go-gl/examples/blob/master/gl41core-cube/cube.go
func CompileShader(source string, shaderType uint32) uint32 {
	shader := gl.CreateShader(shaderType)
	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		infoLog := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(infoLog))
		s := fmt.Sprintf("Failed to compile %v: %v", source, infoLog)
		panic(s)
	}
	return shader
}

func CreateProgram(vert, frag string) uint32 {
	vertexShader := CompileShader(vert, gl.VERTEX_SHADER)
	fragmentShader := CompileShader(frag, gl.FRAGMENT_SHADER)
	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog
}

type Monitor struct {
	SizeMm image.Point
	SizePx image.Point
	Pos    image.Point
}

var Monitors = []Monitor{}

func GetMonitors() {
	ms := glfw.GetMonitors()
	for i, monitor := range ms {
		m := Monitor{}
		m.SizeMm.X, m.SizeMm.Y = monitor.GetPhysicalSize()
		m.Pos.X, m.Pos.Y = monitor.GetPos()
		log.Printf("Monitor %d, %vmmx%vmm, %vx%vpx,  pos: %v, %v\n",
			i+1, m.SizeMm.X, m.SizeMm.Y,
			m.SizePx.X, m.SizePx.Y,
			m.Pos.X, m.Pos.Y)
	}

}

// initOpenGL initializes OpenGL and returns an intiialized program.
func InitOpenGL(bgColor color.Color) {
	if err := gl.Init(); err != nil {
		panic("Initialization error for OpenGL: " + err.Error())
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)
	gl.Enable(gl.BLEND)
	gl.BlendEquation(gl.FUNC_ADD)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	BackgroundColor(bgColor)
	rrprog = CreateProgram(RectVertShaderSource, RectFragShaderSource)
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)

}

// InitWindow initializes glfw and returns a Window to use.
func InitWindow(width, height int, name string) *glfw.Window {
	windowWidth, windowHeight = width, height
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
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

func BackgroundColor(col color.Color) {
	r, g, b, _ := col.RGBA()
	gl.ClearColor(float32(r)/65535.0, float32(g)/65535.0, float32(b)/65535.0, 1.0)
}

func StartFrame() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	startTime = time.Now()
}

func EndFrame(maxFrameRate int, window *glfw.Window) {
	window.SwapBuffers()
	glfw.PollEvents()
	t := time.Since(startTime)
	dt := time.Second/time.Duration(maxFrameRate) - t
	if dt < 0 {
		dt = 0
	}
	time.Sleep(dt)
}

var col [8]float32
var rrprog uint32

func RoundedRect(x, y, w, h, rr, t float32, fillColor, frameColor color.Color) {
	gl.UseProgram(rrprog)
	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.Enable(gl.BLEND)

	vertices := []float32{x + w, y, x, y, x, y + h, x, y + h, x + w, y + h, x + w, y}
	r, g, b, a := fillColor.RGBA()
	col[0] = float32(r) / 65535.0
	col[1] = float32(g) / 65535.0
	col[2] = float32(b) / 65535.0
	col[3] = float32(a) / 65535.0
	r, g, b, a = frameColor.RGBA()
	col[4] = float32(r) / 65535.0
	col[5] = float32(g) / 65535.0
	col[6] = float32(b) / 65535.0
	col[7] = float32(a) / 65535.0

	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)
	// position attribute
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 2*4, nil)
	gl.EnableVertexAttribArray(1)
	// set screen resolution
	r1 := gl.GetUniformLocation(rrprog, gl.Str("resolution\x00"))
	gl.Uniform2f(r1, float32(windowWidth), float32(windowHeight))
	// Colors
	r2 := gl.GetUniformLocation(rrprog, gl.Str("colors\x00"))
	gl.Uniform4fv(r2, 16, &col[0])
	// Set pos data
	r3 := gl.GetUniformLocation(rrprog, gl.Str("pos\x00"))
	gl.Uniform2f(r3, x+w/2, y+h/2)
	// Set halfbox
	r4 := gl.GetUniformLocation(rrprog, gl.Str("halfbox\x00"))
	gl.Uniform2f(r4, w/2, h/2)
	// Set radius/border width
	r5 := gl.GetUniformLocation(rrprog, gl.Str("rw\x00"))
	gl.Uniform2f(r5, rr, t)
	// Do actual drawing
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
	gl.UseProgram(0)
}

func HorLine(x1, x2, y, w float32, col color.Color) {
	RoundedRect(x1, y, x2-x1, w, 0, w, col, col)
}

func VertLine(x, y1, y2, w float32, col color.Color) {
	RoundedRect(x, y1, w, y2-y1, 0, w, col, col)
}

func Rect(x, y, w, h, t float32, fillColor, frameColor color.Color) {
	RoundedRect(x, y, w, h, 0, t, fillColor, frameColor)
}
