package gpu

import (
	"fmt"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/shader"
	"image"
	"log/slog"
	"runtime"
	"time"
)

var (
	startTime      time.Time
	vao            uint32
	vbo            uint32
	WindowWidthPx  int
	WindowHeightPx int
	WindowWidthDp  float32
	WindowHeightDp float32
	DefaultFont    = 1
	LastRune       rune
	LastKey        glfw.Key
	WindowRect     f32.Rect
	Rrprog         uint32
	IconProgram    uint32
	ScaleX         float32 = 1.75
	ScaleY         float32 = 1.75
	UserScale      float32 = 1.4
	Window         *glfw.Window
)

func UpdateResolution() {
	for _, p := range shader.Programs {
		SetResolution(p)
	}
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
	ScaleX *= UserScale
	ScaleY *= UserScale
	WindowWidthDp = float32(width) / ScaleX
	WindowHeightDp = float32(height) / ScaleY
	WindowRect = f32.Rect{0, 0, WindowWidthDp, WindowHeightDp}
	slog.Info("UpdateSize", "w", width, "h", height, "scaleX", ScaleX, "ScaleY", ScaleY)
}

func sizeCallback(w *glfw.Window, width int, height int) {
	UpdateSize(w, width, height)
	UpdateResolution()
	Invalidate(0)
}

func scaleCallback(w *glfw.Window, x float32, y float32) {
	width, height := w.GetSize()
	sizeCallback(w, width, height)
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
		slog.Info("InitWindow()", "Monitor", i+1,
			"WidthMm", m.SizeMm.X, "HeightMm", m.SizeMm.Y,
			"WidthPx", m.SizePx.X, "HeightMm", m.SizePx.Y, "PosX", m.Pos.X, "PosY", m.Pos.Y,
			"ScaleX", m.ScaleX, "ScaleY", m.ScaleY)
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
	var err error
	Window, err = glfw.CreateWindow(m.SizePx.X, m.SizePx.Y, name, nil, nil)
	if err != nil {
		panic(err)
	}
	// Move window to selected monitor
	Window.SetPos(m.Pos.X, m.Pos.Y)
	_, top, _, _ := Window.GetFrameSize()
	Window.SetPos(m.Pos.X, m.Pos.Y+top)
	ww := m.SizePx.X
	hh := m.SizePx.Y - top
	if width > 0 {
		ww = min(int(width), ww)
	}
	if height > 0 {
		hh = min(int(height), hh)
	}
	Window.SetSize(ww, hh)

	// Now we can update size and scaling
	w, h := Window.GetSize()
	UpdateSize(Window, w, h)
	Window.Show()
	slog.Info("New window", "ScaleX", ScaleX, "ScaleY", ScaleY)

	Window.MakeContextCurrent()
	glfw.SwapInterval(1)
	Window.SetContentScaleCallback(scaleCallback)
	Window.SetSizeCallback(sizeCallback)
	if err := gl.Init(); err != nil {
		panic("Initialization error for OpenGL: " + err.Error())
	}
	// Initialize gl
	version := gl.GoStr(gl.GetString(gl.VERSION))
	slog.Info("OpenGL", "version", version)
	gl.Enable(gl.BLEND)
	gl.Enable(gl.MULTISAMPLE)
	gl.BlendEquation(gl.FUNC_ADD)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	BackgroundColor(bgColor)
	Rrprog, _ = shader.NewProgram(shader.RectVertShaderSource, shader.RectFragShaderSource)
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	return Window
}

func BackgroundColor(col f32.Color) {
	gl.ClearColor(col.R, col.G, col.B, col.A)
}

var invalidate time.Duration

func Invalidate(time time.Duration) {
	invalidate = time
}

var Redraws int
var RedrawStart time.Time
var RedrawsPrSec int

func StartFrame() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	startTime = time.Now()
	Redraws++
	if time.Since(RedrawStart).Seconds() >= 1 {
		RedrawsPrSec = Redraws
		RedrawStart = time.Now()
		Redraws = 0
	}

}

// EndFrame will do buffer swapping and focus updates
// Then it will loop and sleep until an event happens
// The event could be an invalidate call
func EndFrame(maxFrameRate int) {
	LastKey = 0
	Window.SwapBuffers()
	for {
		dt := max(0, time.Second/time.Duration(maxFrameRate)-time.Since(startTime))
		time.Sleep(dt)
		startTime = time.Now()
		invalidate -= dt
		glfw.PollEvents()
		// Could use glfw.WaitEventsTimeout(0.03)
		if invalidate <= 0 {
			invalidate = 1 * time.Second
			break
		}
	}
}

func RoundedRect(r f32.Rect, cornerRadius, borderThickness float32, fillColor, frameColor f32.Color, shadowSize float32, shadowColor float32) {
	// Make the quad larger by the shadow width ss  and Correct for device independent pixels
	r.X = (r.X - shadowSize) * ScaleX
	r.Y = (r.Y - shadowSize) * ScaleX
	r.W = (r.W + shadowSize + shadowSize) * ScaleX
	r.H = (r.H + shadowSize + shadowSize) * ScaleX
	cornerRadius *= ScaleX
	borderThickness *= ScaleX

	gl.UseProgram(Rrprog)
	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.Enable(gl.BLEND)

	vertices := []float32{r.X + r.W, r.Y, r.X, r.Y, r.X, r.Y + r.H, r.X, r.Y + r.H,
		r.X + r.W, r.Y + r.H, r.X + r.W, r.Y}
	var col [8]float32
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
	SetResolution(Rrprog)
	// Colors
	r2 := gl.GetUniformLocation(Rrprog, gl.Str("colors\x00"))
	gl.Uniform4fv(r2, 16, &col[0])
	// Set pos data
	r3 := gl.GetUniformLocation(Rrprog, gl.Str("pos\x00"))
	gl.Uniform2f(r3, r.X+r.W/2, r.Y+r.H/2)
	// Set halfbox
	r4 := gl.GetUniformLocation(Rrprog, gl.Str("halfbox\x00"))
	gl.Uniform2f(r4, r.W/2, r.H/2)
	// Set radius/border width
	r5 := gl.GetUniformLocation(Rrprog, gl.Str("rws\x00"))
	gl.Uniform4f(r5, cornerRadius, borderThickness, shadowSize*ScaleX, shadowColor)
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

func GetErrors() {
	e := gl.GetError()
	if e != gl.NO_ERROR {
		slog.Error("OpenGl ", "error", e)
	}
}
