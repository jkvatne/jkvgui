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
		out vec4 colour;
		//uniform vec4 drawColor;
		in vec4 drawColor;
		void main() {
  			colour = drawColor;
		}
	` + "\x00"

	vertexShaderSource = `
		#version 400
		layout(location = 1) in vec2 aPos;
		layout(location = 2) in vec4 aColor;
		out  vec4 drawColor;
		uniform vec2 resolution;

		void main() {
		    // convert the rectangle from pixels to 0.0 to 1.0
		    vec2 zeroToOne = aPos / resolution;
		
		    // convert from 0->1 to 0->2
		    vec2 zeroToTwo = zeroToOne * 2.0;
		
		    // convert from 0->2 to -1->+1 (clipspace)
	 	    vec2 clipSpace = zeroToTwo - 1.0;
		
		    gl_Position = vec4(clipSpace * vec2(1, -1), 0, 1);
			drawColor = aColor;
		}
	` + "\x00"

	windowWidth  = 2300
	windowHeight = 1200
)

var vao uint32

var triangle = []float32{
	250, 250, 0.1, 0.2, 0.0, 1.0,
	50, 550, 0.2, 0.2, 0.0, 1.0,
	450, 550, 0.3, 0.2, 0.0, 1.0,
	450, 250, 0.4, 0.2, 1.0, 1.0,
	250, 550, 0.5, 0.2, 1.0, 1.0,
	650, 550, 0.6, 0.2, 1.0, 1.0,
}

var colors = []float32{
	0.2, 0.0, 0.0, 1.0,
	0.0, 0.2, 0.0, 1.0,
	0.6, 0.0, 0.2, 1.0,
	0.5, 0.0, 0.0, 1.0,
	0.0, 0.5, 0.0, 1.0,
	0.0, 0.0, 0.5, 1.0,
}

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
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	// position attribute
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 6*4, nil)
	gl.EnableVertexAttribArray(1)
	// color attribute
	gl.VertexAttribPointer(2, 4, gl.FLOAT, false, 6*4, gl.PtrOffset(2*4))
	gl.EnableVertexAttribArray(2)
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

func draw(prog uint32) {
	makeVao(triangle)
	// Set Drawing color
	vertexColorLocation := gl.GetUniformLocation(prog, gl.Str("drawColor\x00"))
	gl.Uniform4f(vertexColorLocation, 1.0, 1.0, 0.0, 1.0)
	// set screen resolution
	resUniform := gl.GetUniformLocation(prog, gl.Str("resolution\x00"))
	gl.Uniform2f(resUniform, float32(windowWidth), float32(windowHeight))
	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
}

var font *glfont.Font

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

	prog := initOpenGL()
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
