package gpu

import (
	_ "embed"
	"fmt"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jkvatne/jkvgui/glfont"
	"github.com/jkvatne/jkvgui/shader"
	"image"
	"image/color"
	"log"
	"runtime"
	"strings"
	"time"
)

//go:embed Roboto-Medium.ttf
var RobotoMedium []byte

//go:embed Roboto-Light.ttf
var RobotoLight []byte

//go:embed Roboto-Regular.ttf
var RobotoRegular []byte

//go:embed RobotoMono-Regular.ttf
var RobotoMono []byte

var Fonts []*glfont.Font

var startTime time.Time
var vao uint32
var vbo uint32
var WindowWidth int
var WindowHeight int
var InitialSize int32 = 30

func LoadFont(name string, scale int32) {
	var f *glfont.Font
	var err error
	if strings.EqualFold(name, "Roboto-Medium") {
		f, err = glfont.LoadFontBytes(RobotoMedium, scale, WindowWidth, WindowHeight)
	} else if strings.EqualFold(name, "Roboto") {
		f, err = glfont.LoadFontBytes(RobotoMedium, scale, WindowWidth, WindowHeight)
	} else if strings.EqualFold(name, "Roboto-Light") {
		f, err = glfont.LoadFontBytes(RobotoLight, scale, WindowWidth, WindowHeight)
	} else if strings.EqualFold(name, "Roboto-Regular") {
		f, err = glfont.LoadFontBytes(RobotoRegular, scale, WindowWidth, WindowHeight)
	} else if strings.EqualFold(name, "RobotoMono") {
		f, err = glfont.LoadFontBytes(RobotoMono, scale, WindowWidth, WindowHeight)
	} else {
		f, err = glfont.LoadFont(name, scale, WindowWidth, WindowHeight)
	}
	panicOn(err, "Loading "+name)
	Fonts = append(Fonts, f)
}

func SizeCallback(w *glfw.Window, width int, height int) {
	WindowHeight = height
	WindowWidth = width
	gl.Viewport(0, 0, int32(width), int32(height))
	for _, f := range Fonts {
		f.UpdateResolution(WindowWidth, WindowHeight)
	}
}

type Monitor struct {
	SizeMm image.Point
	SizePx image.Point
	Pos    image.Point
}

var Monitors = []Monitor{}

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
	rrprog = shader.CreateProgram(shader.RectVertShaderSource, shader.RectFragShaderSource)
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)

	LoadFont("Roboto-Light", InitialSize)
	LoadFont("Roboto-Medium", InitialSize)
	LoadFont("Roboto-Regular", InitialSize)
	LoadFont("RobotoMono", InitialSize)
}

// InitWindow initializes glfw and returns a Window to use.
func InitWindow(width, height int, name string, monitorNo int) *glfw.Window {
	runtime.LockOSThread()
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	ms := glfw.GetMonitors()
	for i, monitor := range ms {
		m := Monitor{}
		m.SizeMm.X, m.SizeMm.Y = monitor.GetPhysicalSize()
		_, _, m.SizePx.X, m.SizePx.Y = monitor.GetWorkarea()
		m.Pos.X, m.Pos.Y = monitor.GetPos()
		Monitors = append(Monitors, m)
		log.Printf("Monitor %d, %vmmx%vmm, %vx%vpx,  pos: %v, %v\n",
			i+1, m.SizeMm.X, m.SizeMm.Y,
			m.SizePx.X, m.SizePx.Y,
			m.Pos.X, m.Pos.Y)
	}
	if monitorNo >= len(Monitors) {
		monitorNo = 0
	}
	if width > Monitors[monitorNo].SizePx.X {
		width = Monitors[monitorNo].SizePx.X
	}
	if height > Monitors[monitorNo].SizePx.Y {
		height = Monitors[monitorNo].SizePx.Y
	}
	glfw.WindowHint(glfw.Visible, glfw.False)
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
	left, top, right, bottom := window.GetFrameSize()
	width = width - (left+right)*4/7
	height = height - (top+bottom)*4/7
	window.MakeContextCurrent()
	glfw.SwapInterval(1)
	scaleX, scaleY := window.GetContentScale()
	log.Printf("Window scaleX=%v, scaleY=%v\n", scaleX, scaleY)
	x := Monitors[monitorNo].Pos.X + left
	y := Monitors[monitorNo].Pos.Y + top
	window.SetPos(x, y)
	window.SetSize(width, height)
	window.Show()
	window.SetKeyCallback(KeyCallback)
	window.SetMouseButtonCallback(MouseBtnCallback)
	window.SetSizeCallback(SizeCallback)
	window.SetScrollCallback(ScrollCallback)
	WindowWidth, WindowHeight = width, height

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

var rrprog uint32
var col [8]float32

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
	gl.Uniform2f(r1, float32(WindowWidth), float32(WindowHeight))
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
	// Free memory
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
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

func Shutdown() {
	glfw.Terminate()
}

func panicOn(err error, s string) {
	if err != nil {
		panic(fmt.Sprintf("%s: %v", s, err))
	}
}

// https://www.glfw.org/docs/latest/window_guide.html
func KeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	log.Printf("Key %v %v %v %v\n", key, scancode, action, mods)
}

var N = 10000

func MouseBtnCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		x, y := w.GetCursorPos()
		fmt.Printf("Mouse btn %d clicked at %0.1f,%0.1f\n", button, x, y)
	}
}

func ScrollCallback(w *glfw.Window, xoff float64, yoff float64) {
	fmt.Printf("Scroll dx=%v dy=%v\n", xoff, yoff)
}

func Text(x, y float32, Size float32, fontNr int, color color.Color, text string) {
	r, g, b, a := color.RGBA()
	Fonts[fontNr].SetColor(float32(r)/65535.0, float32(g)/65535.0, float32(b)/65535.0, float32(a)/65535.0)
	_ = Fonts[fontNr].Printf(x, y, Size/float32(InitialSize), text)
}
