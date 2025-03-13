package gpu

import (
	"fmt"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/shader"
	"image"
	"log"
	"runtime"
	"time"
)

var (
	startTime           time.Time
	vao                 uint32
	vbo                 uint32
	WindowWidthPx       int
	WindowHeightPx      int
	WindowWidthDp       float32
	WindowHeightDp      float32
	InitialSize         float32 = 24 // * 1.75
	Clickables          []Clickable
	MousePos            f32.Pos
	MouseBtnDown        bool
	MouseBtnReleased    bool
	DefaultFont         = 1
	MoveFocusToNext     bool
	MoveFocusToPrevious bool
	FocusToNext         bool
	LastFocusable       interface{}
	LastRune            rune
	Backspace           bool
	WindowRect          f32.Rect
	rrprog              uint32
	IconProgram         uint32
	col                 [8]float32
	ScaleX              float32 = 1.75
	ScaleY              float32 = 1.75
)

type Clickable struct {
	Rect   f32.Rect
	Action func()
}

func SetResolution(program uint32) {
	if program == 0 {
		panic("Program number must be greater than 0")
	}
	// Activate corresponding render state
	gl.UseProgram(program)
	// set screen resolution
	gl.Viewport(0, 0, int32(WindowWidthPx), int32(WindowHeightPx))
	resUniform := gl.GetUniformLocation(program, gl.Str("resolution\x00"))
	gl.Uniform2f(resUniform, float32(WindowWidthPx), float32(WindowHeightPx))
}

func UpdateSize(w *glfw.Window, width int, height int) {
	WindowHeightPx = height
	WindowWidthPx = width
	ScaleX, ScaleY = w.GetContentScale()
	WindowWidthDp = float32(width) / ScaleX
	WindowHeightDp = float32(height) / ScaleY
	WindowRect = f32.Rect{0, 0, WindowWidthDp, WindowHeightDp}
	log.Printf("Size Callback w=%d, h=%d, scaleX=%0.2f, scaleY=%0.2f\n", width, height, ScaleX, ScaleY)

}

func SizeCallback(w *glfw.Window, width int, height int) {
	UpdateSize(w, width, height)
	// Must set viewport before changing resolution
	for _, f := range Fonts {
		SetResolution(f.Program)
	}
	SetResolution(rrprog)
	SetResolution(IconProgram)
}

func ScaleCallback(w *glfw.Window, x float32, y float32) {
	width, height := w.GetSize()
	SizeCallback(w, width, height)
}

type Monitor struct {
	SizeMm image.Point
	SizePx image.Point
	ScaleX float32
	ScaleY float32
	Pos    image.Point
}

var Monitors = []Monitor{}

// InitWindow initializes glfw and returns a Window to use.
// MonitorNo is 1 or 0 for the primary monitor, 2 for secondary monitor etc.
// Size is given in dp (device independent pixels)
// Windows typically fills the screen in one of the following ways:
// - Constant aspect ratio, use as much of screen as possible (h=10000, w=10000)
// - Full screen. (Maximized window) (w=0, h=0)
// - Small window of a given size, shrinked if screen is not big enough (h=200, w=200)
// - Use full screen height, but limit width (h=0, w=800)
// - Use full screen width, but limit height (h=800, w=0)
//
func InitWindow(width, height float32, name string, monitorNo int, bgColor f32.Color) *glfw.Window {
	runtime.LockOSThread()
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	// Check all monitors and print size data
	ms := glfw.GetMonitors()
	for i, monitor := range ms {
		m := Monitor{}
		m.SizeMm.X, m.SizeMm.Y = monitor.GetPhysicalSize()
		_, _, m.SizePx.X, m.SizePx.Y = monitor.GetWorkarea()
		m.ScaleX, m.ScaleY = monitor.GetContentScale()
		m.Pos.X, m.Pos.Y = monitor.GetPos()
		Monitors = append(Monitors, m)
		log.Printf("Monitor %d, %vmmx%vmm, %vx%vpx,  pos: %v, %v, scale: %0.2f, %0.2f\n",
			i+1, m.SizeMm.X, m.SizeMm.Y, m.SizePx.X, m.SizePx.Y, m.Pos.X, m.Pos.Y, m.ScaleX, m.ScaleY)
	}

	// Select monitor as given, or use primary monitor.
	monitorNo = max(0, min(monitorNo-1, len(Monitors)-1))
	m := Monitors[monitorNo]

	width = min(width*m.ScaleX, float32(m.SizePx.X))
	height = min(height*m.ScaleY, float32(m.SizePx.Y))

	// Configure glfw. First the window is NOT shown because we need to find window data.
	glfw.WindowHint(glfw.Visible, glfw.False)
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.False)
	glfw.WindowHint(glfw.Samples, 4)
	glfw.WindowHint(glfw.Floating, glfw.False) // True will keep window on top

	if width == 0 && height == 0 {
		glfw.WindowHint(glfw.Maximized, glfw.True)
	} else {
		glfw.WindowHint(glfw.Maximized, glfw.False)
	}

	// Create invisible windows so we can get scaling.
	window, err := glfw.CreateWindow(10, 10, name, nil, nil)
	if err != nil {
		panic(err)
	}
	// Move window to selected monitor
	window.SetPos(m.Pos.X, m.Pos.Y)
	_, top, _, _ := window.GetFrameSize()
	window.SetPos(m.Pos.X, m.Pos.Y+top)
	ww := m.SizePx.X
	hh := m.SizePx.Y - top
	if width > 0 {
		ww = min(int(width), ww)
	}
	if height > 0 {
		hh = min(int(height), hh)
	}
	window.SetSize(ww, hh)

	// Now we can update size and scaling
	w, h := window.GetSize()
	UpdateSize(window, w, h)
	window.Show()
	log.Printf("Window scaleX=%v, scaleY=%v\n", ScaleX, ScaleY)

	window.MakeContextCurrent()
	glfw.SwapInterval(1)
	window.SetKeyCallback(KeyCallback)
	window.SetCharCallback(CharCallback)
	window.SetSizeCallback(SizeCallback)
	window.SetScrollCallback(ScrollCallback)
	window.SetContentScaleCallback(ScaleCallback)

	window.SetMouseButtonCallback(MouseBtnCallback)
	window.SetCursorPosCallback(MousePosCallback)
	if err := gl.Init(); err != nil {
		panic("Initialization error for OpenGL: " + err.Error())
	}
	// Initialize gl
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)
	gl.Enable(gl.BLEND)
	gl.Enable(gl.MULTISAMPLE)
	gl.BlendEquation(gl.FUNC_ADD)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	BackgroundColor(bgColor)
	rrprog = shader.CreateProgram(shader.RectVertShaderSource, shader.RectFragShaderSource)
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	return window
}

func BackgroundColor(col f32.Color) {
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

func RoundedRect(r f32.Rect, rr, t float32, fillColor, frameColor f32.Color, ss float32, sc float32) {
	// Make the quad larger by the shadow width ss  and Correct for device independent pixels
	r.X = (r.X - ss) * ScaleX
	r.Y = (r.Y - ss) * ScaleX
	r.W = (r.W + ss + ss) * ScaleX
	r.H = (r.H + ss + ss) * ScaleX
	rr *= ScaleX
	t *= ScaleX

	gl.UseProgram(rrprog)
	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.Enable(gl.BLEND)

	vertices := []float32{r.X + r.W, r.Y, r.X, r.Y, r.X, r.Y + r.H, r.X, r.Y + r.H,
		r.X + r.W, r.Y + r.H, r.X + r.W, r.Y}
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
	gl.Uniform2f(r3, r.X+r.W/2, r.Y+r.H/2)
	// Set halfbox
	r4 := gl.GetUniformLocation(rrprog, gl.Str("halfbox\x00"))
	gl.Uniform2f(r4, r.W/2, r.H/2)
	// Set radius/border width
	r5 := gl.GetUniformLocation(rrprog, gl.Str("rws\x00"))
	gl.Uniform4f(r5, rr, t, ss*ScaleX, sc)
	// Do actual drawing
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
	// Free memory
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
	gl.UseProgram(0)
}

func HorLine(x1, x2, y, w float32, col f32.Color) {
	r := f32.Rect{x1, y, x2 - x1, w}
	RoundedRect(r, 0, w, col, col, 0, 0)
}

func VertLine(x, y1, y2, w float32, col f32.Color) {
	r := f32.Rect{x, y1, w, y2 - y1}
	RoundedRect(r, 0, w, col, col, 0, 0)
}

func Rect(r f32.Rect, t float32, fillColor, frameColor f32.Color) {
	RoundedRect(r, 0, t, fillColor, frameColor, 0, 0)
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
	// log.Printf("Key %v %v %v %v\n", key, scancode, action, mods)
	if key == glfw.KeyTab && action == glfw.Release {
		if mods != glfw.ModShift {
			MoveFocusToNext = true
		} else {
			MoveFocusToPrevious = true
		}
	}
	if key == glfw.KeyBackspace && action == glfw.Release {
		Backspace = true
	}
}

func CharCallback(w *glfw.Window, char rune) {
	log.Printf("Rune=%d\n", int(char))
	LastRune = char
}

var N = 10000

func ScrollCallback(w *glfw.Window, xoff float64, yoff float64) {
	fmt.Printf("Scroll dx=%v dy=%v\n", xoff, yoff)
}

func GetErrors() {
	e := gl.GetError()
	if e != gl.NO_ERROR {
		log.Printf("OpenGl Error: %x\n", e)
	}
}
