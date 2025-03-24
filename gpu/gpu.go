package gpu

import (
	"fmt"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/shader"
	"github.com/jkvatne/jkvgui/theme"
	"image"
	"image/draw"
	"image/png"
	"log/slog"
	"math"
	"os"
	"runtime"
	"time"
	"unsafe"
)

var ( // Public global variables
	Debugging      bool
	startTime      time.Time
	WindowWidthPx  int
	WindowHeightPx int
	WindowWidthDp  float32
	WindowHeightDp float32
	LastRune       rune
	LastKey        glfw.Key
	WindowRect     f32.Rect

	ScaleX         float32 = 1.0
	ScaleY         float32 = 1.0
	UserScale      float32 = 1.0
	Window         *glfw.Window
	WindowHasFocus = true
)

var ( // Private global variables
	rrprog     uint32
	shaderProg uint32
	vao        uint32
	vbo        uint32
)

const (
	Normal = 0
	Bold   = 1
	Italic = 2
	Mono   = 3
)

var DeferredFunctions []func()

func Defer(f func()) {
	DeferredFunctions = append(DeferredFunctions, f)
}

func RunDefered() {
	for _, f := range DeferredFunctions {
		f()
	}
	DeferredFunctions = DeferredFunctions[0:0]
}

func Return() bool {
	return LastKey == glfw.KeyEnter || LastKey == glfw.KeyKPEnter
}

func Clip(x, y, w, h float32) {
	if w == 0 {
		gl.Disable(gl.SCISSOR_TEST)
		return
	}
	ww := int32(float32(w) * ScaleX)
	hh := int32(float32(h) * ScaleY)
	xx := int32(float32(x) * ScaleX)
	yy := int32(WindowHeightPx) - hh - int32(float32(y)*ScaleY)
	gl.Scissor(xx, yy, ww, hh)
	gl.Enable(gl.SCISSOR_TEST)
}

func Capture(x, y, w, h int) *image.RGBA {
	x = int(float32(x) * ScaleX)
	y = int(float32(y) * ScaleY)
	w = int(float32(w) * ScaleX)
	h = int(float32(h) * ScaleY)
	y = WindowHeightPx - h - y

	img := image.NewRGBA(image.Rect(0, 0, w, h))
	gl.PixelStorei(gl.PACK_ALIGNMENT, 1)
	gl.ReadPixels(int32(x), int32(y), int32(w), int32(h),
		gl.RGBA, gl.UNSIGNED_BYTE, unsafe.Pointer(&img.Pix[0]))
	GetErrors()
	//  Upside down
	for y := 0; y < h/2-1; y++ {
		for x := 0; x < (4*w - 1); x++ {
			tmp := img.Pix[x+y*img.Stride]
			img.Pix[x+y*img.Stride] = img.Pix[x+(h-y-1)*img.Stride]
			img.Pix[x+(h-y-1)*img.Stride] = tmp
		}
	}
	// Scale by alfa
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			ofs := x*4 + y*img.Stride
			// Set alpha to 1.0 (dvs 255). It is not used in files
			img.Pix[ofs+3] = 255
		}
	}

	return img
}

func SaveImage(filename string, img *image.RGBA) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			slog.Error("Could not close", "file", filename)
		}
	}(f)
	return png.Encode(f, img)
}

func LoadImage(filename string) (*image.RGBA, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			slog.Error("Could not close", "file", filename)
		}
	}(f)
	img, _, err := image.Decode(f)
	m, ok := img.(*image.RGBA)
	if ok {
		return m, nil
	}
	b := img.Bounds()
	rgba := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(rgba, b, img, b.Min, draw.Src)
	return rgba, nil
}

func CaptureToFile(filename string, x, y, w, h int) error {
	img := Capture(x, y, w, h) // 1057-300
	return SaveImage(filename, img)
}

func ImgDiff(img1, img2 *image.RGBA) int {
	if img1.Bounds().Size() != img2.Bounds().Size() {
		return 256
	}
	maxDelta := 0
	for i := 0; i < len(img1.Pix); i++ {
		d := int(img1.Pix[i]) - int(img2.Pix[i])
		if d < 0 {
			d = -d
		}
		if d > maxDelta {
			maxDelta = d
		}
	}
	return maxDelta
}

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

func focusCallback(w *glfw.Window, focused bool) {
	WindowHasFocus = focused
	Invalidate(0)
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
func InitWindow(width, height float32, name string, monitorNo int) *glfw.Window {
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
	glfw.WindowHint(glfw.Maximized, glfw.False)

	// Create invisible windows so we can get scaling.
	var err error
	Window, err = glfw.CreateWindow(m.SizePx.X, m.SizePx.Y, name, nil, nil)
	if err != nil {
		panic(err)
	}
	// Move window to selected monitor
	Window.SetPos(m.Pos.X, m.Pos.Y)
	_, top, _, _ := Window.GetFrameSize()

	if width == 0 && height == 0 {
		Window.SetPos(m.Pos.X, m.Pos.Y+top)
		Window.SetSize(m.SizePx.X, m.SizePx.Y-top)
	} else {
		Window.SetPos(
			m.Pos.X+(m.SizePx.X-int(width))/2,
			top+m.Pos.Y+(m.SizePx.Y-int(height))/2)
		ww := m.SizePx.X
		hh := m.SizePx.Y - top
		if width > 0 {
			ww = min(int(width), ww)
		}
		if height > 0 {
			hh = min(int(height), hh)
		}
		Window.SetSize(ww, hh)
	}

	// Now we can update size and scaling
	w, h := Window.GetSize()
	UpdateSize(Window, w, h)
	Window.Show()
	slog.Info("New window", "ScaleX", ScaleX, "ScaleY", ScaleY, "W", w, "H", h)

	Window.MakeContextCurrent()
	glfw.SwapInterval(1)
	Window.SetContentScaleCallback(scaleCallback)
	Window.SetFocusCallback(focusCallback)
	Window.SetSizeCallback(sizeCallback)
	if err := gl.Init(); err != nil {
		panic("Initialization error for OpenGL: " + err.Error())
	}
	Window.Focus()
	// Initialize gl
	version := gl.GoStr(gl.GetString(gl.VERSION))
	slog.Info("OpenGL", "version", version)
	gl.Enable(gl.BLEND)
	gl.Enable(gl.MULTISAMPLE)
	gl.BlendEquation(gl.FUNC_ADD)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.ClearColor(1, 1, 1, 1)
	rrprog, _ = shader.NewProgram(shader.VertRectSource, shader.FragRectSource)
	shaderProg, _ = shader.NewProgram(shader.VertRectSource, shader.FragShadowSource)
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	return Window
}

func BackgroundRole(role theme.UIRole) {
	col := theme.Colors[role]
	gl.ClearColor(col.R, col.G, col.B, col.A)
}

func BackgroundColor(col f32.Color) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(col.R, col.G, col.B, col.A)
	UpdateResolution()
}

var invalidateAt time.Time

func Invalidate(dt time.Duration) {
	if time.Since(invalidateAt) <= 0 {
		// We passed the deadline. Set new
		invalidateAt = time.Now().Add(dt)
	} else if time.Since(invalidateAt) > dt {
		// There is a future deadline. Update only if the new one is earlier.
		invalidateAt = time.Now().Add(dt)
	}
}

var Redraws int
var RedrawStart time.Time
var RedrawsPrSec int

type Clickable struct {
	Rect   f32.Rect
	Action any
}

var Clickables []Clickable

func StartFrame(color f32.Color) {
	startTime = time.Now()
	Redraws++
	if time.Since(RedrawStart).Seconds() >= 1 {
		RedrawsPrSec = Redraws
		RedrawStart = time.Now()
		Redraws = 0
	}
	Clickables = Clickables[0:0]
	BackgroundColor(color)
}

// EndFrame will do buffer swapping and focus updates
// Then it will loop and sleep until an event happens
// The event could be an invalidate call
func EndFrame(maxFrameRate int) {
	RunDefered()
	LastKey = 0
	Window.SwapBuffers()
	for {
		dt := max(0, time.Second/time.Duration(maxFrameRate)-time.Since(startTime))
		time.Sleep(dt)
		startTime = time.Now()
		glfw.PollEvents()
		// Could use glfw.WaitEventsTimeout(0.03)
		if time.Since(invalidateAt) >= 0 {
			invalidateAt = time.Now().Add(time.Second)
			break
		}
	}
}

func Shade(r f32.Rect, cornerRadius float32, fillColor f32.Color, shadowSize float32) {
	// Make the quad larger by the shadow width ss  and Correct for device independent pixels
	r.X = (r.X - shadowSize*0.75) * ScaleX
	r.Y = (r.Y - shadowSize*0.75) * ScaleX
	r.W = (r.W + shadowSize*1.5) * ScaleX
	r.H = (r.H + shadowSize*1.5) * ScaleX
	maxRR := min(r.H/2, r.W/2)
	shadowSize *= ScaleX
	cornerRadius *= ScaleX
	cornerRadius = max(0, min(maxRR, cornerRadius+shadowSize))

	gl.UseProgram(shaderProg)
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

	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)
	// position attribute
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 2*4, nil)
	gl.EnableVertexAttribArray(1)
	// Colors
	r2 := gl.GetUniformLocation(shaderProg, gl.Str("colors\x00"))
	gl.Uniform4fv(r2, 16, &col[0])
	// Set pos data
	r3 := gl.GetUniformLocation(shaderProg, gl.Str("pos\x00"))
	gl.Uniform2f(r3, r.X+r.W/2, r.Y+r.H/2)
	// Set halfbox
	r4 := gl.GetUniformLocation(shaderProg, gl.Str("halfbox\x00"))
	gl.Uniform2f(r4, r.W/2, r.H/2)
	// Set radius/border width
	r5 := gl.GetUniformLocation(shaderProg, gl.Str("rws\x00"))
	gl.Uniform4f(r5, cornerRadius, 0, shadowSize, 0)
	// Do actual drawing
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
	// Free memory
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
	gl.UseProgram(0)
	GetErrors()

}

var col [8]float32

func SolidRR(r f32.Rect, cornerRadius float32, fillColor f32.Color) {
	RoundedRect(r, cornerRadius, 0, fillColor, f32.Transparent)
}

func RoundedRect(r f32.Rect, cornerRadius, borderThickness float32, fillColor, frameColor f32.Color) {
	// Make the quad larger by the shadow width ss  and Correct for device independent pixels
	r.X = r.X * ScaleX
	r.Y = r.Y * ScaleX
	r.W = r.W * ScaleX
	r.H = r.H * ScaleX
	cornerRadius *= ScaleX
	cornerRadius = max(0, min(r.H/2, cornerRadius))
	borderThickness *= ScaleX

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
	r5 := gl.GetUniformLocation(rrprog, gl.Str("rw\x00"))
	gl.Uniform2f(r5, cornerRadius, borderThickness)
	// Do actual drawing
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
	// Free memory
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
	gl.UseProgram(0)
}

func HorLine(x1, x2, y, w float32, col f32.Color) {
	r := f32.Rect{x1, y, x2 - x1, w}
	RoundedRect(r, 0, w, col, col)
}

func VertLine(x, y1, y2, w float32, col f32.Color) {
	r := f32.Rect{x, y1, w, y2 - y1}
	RoundedRect(r, 0, w, col, col)
}

func Rect(r f32.Rect, t float32, fillColor, frameColor f32.Color) {
	RoundedRect(r, 0, t, fillColor, frameColor)
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

func sqDiff(x, y uint8) uint64 {
	d := uint64(x) - uint64(y)
	return d * d
}

func Compare(img1, img2 *image.RGBA) (int64, error) {
	if img1.Bounds() != img2.Bounds() {
		return 0, fmt.Errorf("image bounds not equal: %+v, %+v", img1.Bounds(), img2.Bounds())
	}
	accumError := int64(0)
	for i := 0; i < len(img1.Pix); i++ {
		accumError += int64(sqDiff(img1.Pix[i], img2.Pix[i]))
	}
	return int64(math.Sqrt(float64(accumError))), nil
}
