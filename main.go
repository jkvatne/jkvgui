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
        in  vec2 aRadWidth;
        in  vec4 aRect;
		in  vec4 BorderColor;
		in  vec4 FillColor;
		layout(origin_upper_left) in vec4 gl_FragCoord;

		out vec4 fragColor;
		
		// b.x = half width, b.y = half height
		float sdRoundedBox( in vec2 p, in vec2 b, in float r ) {
			vec2 q = abs(p)-b+r;
			return min(max(q.x,q.y),0.0) + length(max(q,0.0)) - r;
		}

		void main() {
  			fragColor = FillColor;
            vec2 halfbox = vec2((aRect[2]-aRect[0])/2, (aRect[3]-aRect[1])/2);
            vec2 p = gl_FragCoord.xy;
            p = p-vec2((aRect[2]+aRect[0])/2, (aRect[3]+aRect[1])/2);
            float d1 = sdRoundedBox(p, halfbox, aRadWidth[0]);
            float w = aRadWidth[1];
            vec2 halfbox2 = vec2(halfbox.x-w*2, halfbox.y-2*w);
            float d2 = sdRoundedBox(p, halfbox2, aRadWidth[0]-w);
            if (d1>0.0) {
				discard;
            }
			if (d2<=0) {
				fragColor = BorderColor;
            }
		}
		` + "\x00"

	vertexShaderSource = `
		#version 400
		layout(location = 1) in vec2 inPos;
		layout(location = 2) in vec2 inColorIndex;
		layout(location = 3) in vec2 inRadWidth;
		layout(location = 4) in vec4 inRect;
        out  vec2 aRadWidth;
        out  vec4 aRect;
		out  vec4 BorderColor;
		out  vec4 FillColor;
		uniform vec2 resolution;
		uniform vec4 colors[8];

		void main() {
		    vec2 zeroToOne = inPos / resolution;
		    vec2 zeroToTwo = zeroToOne * 2.0;
	 	    vec2 clipSpace = zeroToTwo - 1.0;
		    gl_Position = vec4(clipSpace * vec2(1, -1), 0, 1);
			BorderColor =  colors[int(inColorIndex[0])];
			FillColor =  colors[int(inColorIndex[1])];
            aRadWidth = inRadWidth;
            aRect = inRect;
		}
		` + "\x00"

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

// https://github.com/go-gl/examples/blob/master/gl41core-cube/cube.go
func compileShader(source string, shaderType uint32) uint32 {
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
func initOpenGL() {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	gl.Enable(gl.BLEND)
	gl.BlendEquation(gl.FUNC_ADD)
	gl.BlendFunc(gl.SRC_ALPHA, gl.SRC_ALPHA)
	gl.ClearColor(0.95, 0.95, 0.86, 0.10)
}

func CreateProgram() uint32 {
	vertexShader := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	fragmentShader := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog
}

func DrawTriangles(prog uint32) {
	gl.UseProgram(prog)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(triangles), gl.Ptr(triangles), gl.STATIC_DRAW)
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
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

	gl.BindVertexArray(vao)

	gl.DrawArrays(gl.TRIANGLES, 0, 12)

	gl.BindVertexArray(0)
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

	initOpenGL()
	LoadFonts()
	InitKeys(window)
	prog := CreateProgram()
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		// FPS=3 for 100*22*16=35200 labels!
		// font.SetColor(0.0, 0.0, 1.0, 1.0)
		// _ = font.Printf(0, 100, 1.0, "Before frames"+"\x00")
		for range 100 {
			DrawTriangles(prog)
		}
		// _ = font.Printf(0, 70, 1.0, "After frames"+"\x00")
		// _ = font.Printf(0, 170, 1.0, "After frames"+"\x00")
		window.SwapBuffers()
		glfw.PollEvents()
	}
	glfw.Terminate()

}
