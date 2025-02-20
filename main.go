package main

import (
	"fmt"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"log"
	"runtime"
	"strings"
	"testglfont/glfont"
)

const (
	fragmentShaderSource = `
		#version 400

		out vec4 frag_colour;
		uniform vec4 drawColor;

		void main() {
  			frag_colour = drawColor;
		}
	` + "\x00"

	vertexShaderSource = `
		#version 400

		in vec2 vp;
		//window res
		uniform vec2 resolution;
		
		void main() {
		   // convert the rectangle from pixels to 0.0 to 1.0
		   vec2 zeroToOne = vp / resolution;
		
		   // convert from 0->1 to 0->2
		   vec2 zeroToTwo = zeroToOne * 2.0;
		
		   // convert from 0->2 to -1->+1 (clipspace)
		   vec2 clipSpace = zeroToTwo - 1.0;
		
		   gl_Position = vec4(clipSpace * vec2(1, -1), 0, 1);
		}
	` + "\x00"

	windowWidth  = 2300
	windowHeight = 1200
)

var vao uint32

// https://github.com/go-gl/examples/blob/master/gl41core-cube/cube.go
func compileShader(source string, shaderType uint32) (uint32, error) {
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
		return 0, fmt.Errorf("failed to compile %v: %v", source, infoLog)
	}
	return shader, nil
}

// makeVao initializes and returns a vertex array from the points provided.
func makeVao(points []float32) {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 0, nil)
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
	glfw.WindowHint(glfw.Resizable, glfw.False)
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

// initOpenGL initializes OpenGL and returns an intiialized program.
func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.95, 0.95, 0.86, 1.0)

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog
}

// Because Gil specified "screen coordinates" (presumably with an upper-left origin), this short bit of
// code sets up the coordinate system to correspond to actual window coodrinates.  This code
// wouldn't be required if you chose a (more typical in 3D) abstract coordinate system.
func SetupTransform(w float64, h float64) {
	// Establish viewing area to cover entire window.
	gl.Viewport(0, 0, int32(w), int32(h))
	// Start modifying the projection matrix.
	gl.MatrixMode(gl.PROJECTION)
	// Reset project matrix.
	gl.LoadIdentity()
	// Map abstract coords directly to window coords
	gl.Ortho(0, w, 0, h, -1, 1)
	// Invert Y axis so increasing Y goes down.
	gl.Scalef(1, -1, 1)
	// Shift origin up to upper-left corner
	gl.Translatef(0, float32(-h), 0)
}

func draw(prog uint32) {
	// Set Drawing color
	vertexColorLocation := gl.GetUniformLocation(prog, gl.Str("drawColor\x00"))
	gl.Uniform4f(vertexColorLocation, 1.0, 1.0, 0.0, 1.0)
	// set screen resolution
	resUniform := gl.GetUniformLocation(prog, gl.Str("resolution\x00"))
	gl.Uniform2f(resUniform, float32(windowWidth), float32(windowHeight))
	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(triangle)/2))
}

var font *glfont.Font

func LoadFonts() {
	var err error
	font, err = glfont.LoadFont("Roboto-Medium.ttf", 35, windowWidth, windowHeight)
	if err != nil {
		log.Panicf("LoadFont: %v", err)
	}
}

var triangle = []float32{
	250, 250,
	50, 550,
	450, 550,
	450, 250,
	250, 550,
	650, 550,
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

	prog := initOpenGL()
	makeVao(triangle)

	LoadFonts()
	InitKeys(window)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(prog)
		draw(prog)
		// FPS=3 for 100*22*16=35200 labels! Dvs 10000 tests pr sec
		// set color and draw text
		font.SetColor(0.0, 0.0, 0.0, 1.0)
		_ = font.Printf(50, 50, 1.0, "Aøæ©")
		font.Printf(10, 50, 1.0, "Hello World!")
		window.SwapBuffers()
		glfw.PollEvents()
	}
	glfw.Terminate()

}
