package gpu

import (
	"fmt"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"image/color"
	"log"
	"strings"
	"time"
)

var startTime time.Time

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

// initOpenGL initializes OpenGL and returns an intiialized program.
func InitOpenGL() {
	if err := gl.Init(); err != nil {
		panic("Initialization error for OpenGL: " + err.Error())
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	gl.Enable(gl.BLEND)
	gl.BlendEquation(gl.FUNC_ADD)
	gl.BlendFunc(gl.SRC_ALPHA, gl.SRC_ALPHA)
}

// InitWindow initializes glfw and returns a Window to use.
func InitWindow(width, height int, name string) *glfw.Window {
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

func BackgroundColor(col color.RGBA) {
	gl.ClearColor(float32(col.R)/255.0, float32(col.G)/255.0, float32(col.B)/255.0, float32(col.A)/255.0)
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
