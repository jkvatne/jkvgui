package gpu

import (
	_ "embed"
	"fmt"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jkvatne/jkvgui/lib"
	"github.com/jkvatne/jkvgui/shader"
	"golang.org/x/image/font/gofont/gomono"
	"image"
	"log"
	"runtime"
	"time"
)

//go:embed fonts/Roboto-Thin.ttf
var Roboto100 []byte // 100

//go:embed fonts/Roboto-ExtraLight.ttf
var Roboto200 []byte // 200

//go:embed fonts/Roboto-Light.ttf
var Roboto300 []byte // 300

//go:embed fonts/Roboto-Regular.ttf
var Roboto400 []byte // 400

//go:embed fonts/Roboto-Medium.ttf
var Roboto500 []byte // 500

//go:embed fonts/Roboto-SemiBold.ttf
var Roboto600 []byte // 600

//go:embed fonts/Roboto-Bold.ttf
var Roboto700 []byte // 700

//go:embed fonts/Roboto-Bold.ttf
var Roboto800 []byte // 800

//go:embed fonts/Roboto-Bold.ttf
var Roboto900 []byte // 900

//go:embed fonts/RobotoMono-Regular.ttf
var RobotoMono []byte

type Color struct {
	R float32
	G float32
	B float32
	A float32
}

var (
	Transparent  = Color{}
	Black        = Color{0, 0, 0, 1}
	Lightgrey    = Color{0.1, 0.1, 0.1, 0.1}
	Blue         = Color{0, 0, 1, 1}
	Red          = Color{1, 0, 0, 1}
	Green        = Color{0, 1, 0, 1}
	White        = Color{1, 1, 1, 1}
	startTime    time.Time
	vao          uint32
	vbo          uint32
	WindowWidth  int
	WindowHeight int
	InitialSize  float32 = 24

	Clickables       []Clickable
	MousePos         lib.Pos
	MouseBtnDown     bool
	MouseBtnReleased bool
)

type Clickable struct {
	Rect   lib.Rect
	Action func()
}

func SetResolution(program uint32) {
	if program == 0 {
		panic("Program number must be greater than 0")
	}
	// Activate corresponding render state
	gl.UseProgram(program)
	// set screen resolution
	gl.Viewport(0, 0, int32(WindowWidth), int32(WindowHeight))
	resUniform := gl.GetUniformLocation(program, gl.Str("resolution\x00"))
	gl.Uniform2f(resUniform, float32(WindowWidth), float32(WindowHeight))
}

func SizeCallback(w *glfw.Window, width int, height int) {
	WindowHeight = height
	WindowWidth = width
	if w != nil {
		Scale, _ = w.GetContentScale()
	}
	fmt.Printf("Size Callback w=%d, h=%d\n", WindowWidth, WindowHeight)
	// Must set viewport before changing resolution
	for _, f := range Fonts {
		SetResolution(f.program)
	}
	SetResolution(rrprog)
}

type Monitor struct {
	SizeMm image.Point
	SizePx image.Point
	Pos    image.Point
}

var Monitors = []Monitor{}

// initOpenGL initializes OpenGL and returns an intiialized program.
func InitOpenGL(bgColor Color) {

}

// InitWindow initializes glfw and returns a Window to use.
// MonitorNo is 1 or 0 for the primary monitor, 2 for secondary monitor etc.
func InitWindow(width, height int, name string, monitorNo int, bgColor Color) *glfw.Window {
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
			i+1, m.SizeMm.X, m.SizeMm.Y, m.SizePx.X, m.SizePx.Y, m.Pos.X, m.Pos.Y)
	}
	monitorNo = max(0, min(monitorNo-1, len(Monitors)-1))
	width = min(width, Monitors[monitorNo].SizePx.X)
	height = min(height, Monitors[monitorNo].SizePx.Y)

	glfw.WindowHint(glfw.Visible, glfw.False)
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.False)
	if width == Monitors[monitorNo].SizePx.X && height == Monitors[monitorNo].SizePx.Y {
		glfw.WindowHint(glfw.Maximized, glfw.True)
	} else {
		glfw.WindowHint(glfw.Maximized, glfw.False)
	}

	glfw.WindowHint(glfw.Floating, glfw.False) // True will keep window on top

	window, err := glfw.CreateWindow(width, height, name, nil, nil)
	if err != nil {
		panic(err)
	}
	if width != Monitors[monitorNo].SizePx.X || height != Monitors[monitorNo].SizePx.Y {
		left, top, right, bottom := window.GetFrameSize()
		width = width - (left + right)
		height = height - (top + bottom)
		x := Monitors[monitorNo].Pos.X + left
		y := Monitors[monitorNo].Pos.Y + top
		window.SetPos(x, y)
		window.SetSize(width, height)
	}
	window.Show()
	scaleX, scaleY := window.GetContentScale()
	Scale = scaleY
	log.Printf("Window scaleX=%v, scaleY=%v\n", scaleX, scaleY)
	WindowWidth, WindowHeight = window.GetSize()
	window.MakeContextCurrent()
	glfw.SwapInterval(1)
	window.SetKeyCallback(KeyCallback)
	window.SetSizeCallback(SizeCallback)
	window.SetScrollCallback(ScrollCallback)

	window.SetMouseButtonCallback(MouseBtnCallback)
	window.SetCursorPosCallback(MousePosCallback)
	if err := gl.Init(); err != nil {
		panic("Initialization error for OpenGL: " + err.Error())
	}
	// Initialize gl
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)
	gl.Enable(gl.BLEND)
	gl.BlendEquation(gl.FUNC_ADD)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	BackgroundColor(bgColor)
	rrprog = shader.CreateProgram(shader.RectVertShaderSource, shader.RectFragShaderSource)
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	LoadFont(Roboto100, InitialSize)
	LoadFont(Roboto200, InitialSize)
	LoadFont(Roboto300, InitialSize)
	LoadFont(Roboto400, InitialSize)
	LoadFont(Roboto500, InitialSize)
	LoadFont(Roboto600, InitialSize)
	LoadFont(Roboto700, InitialSize)
	LoadFont(Roboto800, InitialSize)
	LoadFont(gomono.TTF, InitialSize)
	fmt.Printf("Initial size w=%d, h=%d\n", WindowWidth, WindowHeight)
	SizeCallback(window, WindowWidth, WindowHeight)

	return window
}

func BackgroundColor(col Color) {
	gl.ClearColor(col.R, col.G, col.B, col.A)
}

func StartFrame() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	startTime = time.Now()
	Clickables = Clickables[0:0]
}

func EndFrame(maxFrameRate int, window *glfw.Window) {
	window.SwapBuffers()
	if MoveFocusToNext {
		FocusToNext = true
		MoveFocusToNext = false
	}
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
var Scale float32 = 1.75

func RoundedRect(x, y, w, h, rr, t float32, fillColor, frameColor Color) {
	x *= Scale
	y *= Scale
	w *= Scale
	h *= Scale
	rr *= Scale
	t *= Scale
	gl.UseProgram(rrprog)
	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.Enable(gl.BLEND)

	vertices := []float32{x + w, y, x, y, x, y + h, x, y + h, x + w, y + h, x + w, y}
	col[0] = fillColor.R
	col[1] = fillColor.G
	col[2] = fillColor.B
	col[3] = fillColor.A
	col[4] = frameColor.R
	col[5] = frameColor.G
	col[6] = frameColor.B
	col[7] = frameColor.A

	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)
	// position attribute
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 2*4, nil)
	gl.EnableVertexAttribArray(1)
	SetResolution(rrprog)
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
	gl.UseProgram(0)
}

func HorLine(x1, x2, y, w float32, col Color) {
	RoundedRect(x1, y, x2-x1, w, 0, w, col, col)
}

func VertLine(x, y1, y2, w float32, col Color) {
	RoundedRect(x, y1, w, y2-y1, 0, w, col, col)
}

func Rect(x, y, w, h, t float32, fillColor, frameColor Color) {
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

var LastKey glfw.Key
var MoveFocusToNext bool
var MoveFocusToPrevious bool
var FocusToNext bool
var LastFocusable interface{}

// https://www.glfw.org/docs/latest/window_guide.html
func KeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	log.Printf("Key %v %v %v %v\n", key, scancode, action, mods)
	if key == glfw.KeyTab && action == glfw.Release {
		if mods != glfw.ModShift {
			MoveFocusToNext = true
		} else {
			MoveFocusToPrevious = true
		}
	}
}

var N = 10000

func ScrollCallback(w *glfw.Window, xoff float64, yoff float64) {
	fmt.Printf("Scroll dx=%v dy=%v\n", xoff, yoff)
}
