package gpu

import (
	"encoding/binary"
	"fmt"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gl"
	"github.com/jkvatne/jkvgui/gl/glutil"
	"github.com/jkvatne/jkvgui/glfw"
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
)

var ( // Public global variables
	WindowWidthPx  int
	WindowHeightPx int
	WindowWidthDp  float32
	WindowHeightDp float32
	LastRune       rune
	LastKey        glfw.Key
	WindowRect     f32.Rect
	WindowHasFocus         = true
	ScaleX         float32 = 1.0
	ScaleY         float32 = 1.0
	UserScale      float32 = 1.0
	Window         *glfw.Window
	DebugWidgets   bool
	Monitors       []Monitor
)

var ( // Private global variables
	rrprog     gl.Program
	shaderProg gl.Program
	// vao        uint32
	vbo gl.Buffer
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
	gl.ReadPixels(img.Pix, x, y, w, h,
		gl.RGBA, gl.UNSIGNED_BYTE)
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
	GetErrors()

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
	if err != nil {
		return nil, err
	}
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
	// TODO for _, p := range shader.Programs {
	//	SetResolution(p)
	// }
}

func SetResolution(program gl.Program) {
	if program.Value == 0 {
		panic("Program number must be greater than 0")
	}
	// Activate corresponding render state
	gl.UseProgram(program)
	// set screen resolution
	gl.Viewport(0, 0, WindowWidthPx, WindowHeightPx)
	resUniform := gl.GetUniformLocation(program, "resolution\x00")
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
	WindowRect = f32.Rect{W: WindowWidthDp, H: WindowHeightDp}
	slog.Info("UpdateSize", "w", width, "h", height, "scaleX", ScaleX, "ScaleY", ScaleY)
}

type Monitor struct {
	SizeMm image.Point
	SizePx image.Point
	ScaleX float32
	ScaleY float32
	Pos    image.Point
}

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
	var err error
	runtime.LockOSThread()
	theme.SetDefaultPallete(true)
	err = glfw.Init(gl.ContextWatcher)
	if err != nil {
		panic(err)
	}
	// Check all monitors and print size data
	ms := glfw.GetMonitors()
	for i, monitor := range ms {
		m := Monitor{}
		m.SizeMm.X, m.SizeMm.Y = monitor.Monitor.GetPhysicalSize()
		_, _, m.SizePx.X, m.SizePx.Y = monitor.Monitor.GetWorkarea()
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
	glfw.WindowHint(glfw.Resizable, 1)
	glfw.WindowHint(glfw.Samples, 4)
	glfw.WindowHint(glfw.Visible, 0)
	/*
		glfw.WindowHint(glfw.ContextVersionMajor, 3)
		glfw.WindowHint(glfw.ContextVersionMinor, 3)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.False)
		glfw.WindowHint(glfw.Floating, glfw.False) // True will keep window on top
		glfw.WindowHint(glfw.Maximized, glfw.False)
	*/
	// Create invisible windows so we can get scaling.
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

	// Initialize gl
	Window.MakeContextCurrent()
	glfw.SwapInterval(0)
	Window.Focus()

	slog.Info("OpenGl", "Renderer", gl.GetString(gl.RENDERER))
	slog.Info("OpenGl", "Version", gl.GetString(gl.VERSION))
	slog.Info("OpenGl", "ShadingLanguageVersion", gl.GetString(gl.SHADING_LANGUAGE_VERSION))

	gl.Enable(gl.BLEND)
	// gl.Enable(gl.MULTISAMPLE)
	gl.BlendEquation(gl.FUNC_ADD)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.ClearColor(1, 1, 1, 1)
	rrprog, _ = glutil.CreateProgram(shader.VertRectSource, shader.FragRectSource)
	shaderProg, _ = glutil.CreateProgram(shader.VertRectSource, shader.FragShadowSource)
	// gl.GenVertexArrays(1, &vao)
	// vbo = gl.CreateBuffer()
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	return Window
}

func BackgroundRole(role theme.UIRole) {
	col := theme.Colors[role]
	gl.ClearColor(col.R, col.G, col.B, col.A)
}

func BackgroundColor(col f32.Color) {
	gl.ClearColor(col.R, col.G, col.B, col.A)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	UpdateResolution()
}

var InvalidateAt time.Time

func Invalidate(dt time.Duration) {
	if time.Since(InvalidateAt) <= 0 {
		// We passed the deadline. Set new
		InvalidateAt = time.Now().Add(dt)
	} else if time.Since(InvalidateAt) > dt {
		// There is a future deadline. Update only if the new one is earlier.
		InvalidateAt = time.Now().Add(dt)
	}
}

func Scale(fact float32, values ...*float32) {
	for _, x := range values {
		*x = *x * fact
	}
}

// Bytes returns the byte representation of float32 values in the given byte
// order. byteOrder must be either binary.BigEndian or binary.LittleEndian.
func Bytes(byteOrder binary.ByteOrder, values ...float32) []byte {
	le := false
	switch byteOrder {
	case binary.BigEndian:
	case binary.LittleEndian:
		le = true
	default:
		panic(fmt.Sprintf("invalid byte order %v", byteOrder))
	}

	b := make([]byte, 4*len(values))
	for i, v := range values {
		u := math.Float32bits(v)
		if le {
			b[4*i+0] = byte(u >> 0)
			b[4*i+1] = byte(u >> 8)
			b[4*i+2] = byte(u >> 16)
			b[4*i+3] = byte(u >> 24)
		} else {
			b[4*i+0] = byte(u >> 24)
			b[4*i+1] = byte(u >> 16)
			b[4*i+2] = byte(u >> 8)
			b[4*i+3] = byte(u >> 0)
		}
	}
	return b
}

func Shade(r f32.Rect, cornerRadius float32, fillColor f32.Color, shadowSize float32) {
	// Make the quad larger by the shadow width ss  and Correct for device independent pixels
	r.X = (r.X - shadowSize*0.75) * ScaleX
	r.Y = (r.Y - shadowSize*0.75) * ScaleX
	r.W = (r.W + shadowSize*1.5) * ScaleX
	r.H = (r.H + shadowSize*1.5) * ScaleX
	shadowSize *= ScaleX
	cornerRadius *= ScaleX
	if cornerRadius < 0 || cornerRadius > r.H/2 {
		cornerRadius = r.H / 2
	}

	gl.UseProgram(shaderProg)
	// gl.BindVertexArray(vao)
	vbo := gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.Enable(gl.BLEND)
	vertices := []float32{r.X + r.W, r.Y, r.X, r.Y, r.X, r.Y + r.H, r.X, r.Y + r.H,
		r.X + r.W, r.Y + r.H, r.X + r.W, r.Y}
	var col = make([]float32, 8)
	col[0] = fillColor.R
	col[1] = fillColor.G
	col[2] = fillColor.B
	col[3] = fillColor.A

	v := Bytes(binary.LittleEndian, vertices...)
	gl.BufferData(gl.ARRAY_BUFFER, v, gl.STATIC_DRAW)
	// position attribute
	gl.VertexAttribPointer(gl.Attrib{1}, 2, gl.FLOAT, false, 2*4, 0)
	gl.EnableVertexAttribArray(gl.Attrib{1})
	// Colors
	r2 := gl.GetUniformLocation(shaderProg, "colors")
	gl.Uniform4fv(r2, col)
	// Set pos data
	r3 := gl.GetUniformLocation(shaderProg, "pos")
	gl.Uniform2f(r3, r.X+r.W/2, r.Y+r.H/2)
	// Set halfbox
	r4 := gl.GetUniformLocation(shaderProg, "halfbox")
	gl.Uniform2f(r4, r.W/2, r.H/2)
	// Set radius/border width
	r5 := gl.GetUniformLocation(shaderProg, "rws")
	gl.Uniform4f(r5, cornerRadius, 0, shadowSize, 0)
	// Do actual drawing
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
	// Free memory
	gl.BindBuffer(gl.ARRAY_BUFFER, gl.Buffer{0})
	// gl.BindVertexArray(0)
	// gl.UseProgram(0)
	GetErrors()

}

var col = make([]float32, 12)

func SolidRR(r f32.Rect, cornerRadius float32, fillColor f32.Color) {
	RR(r, cornerRadius, 0, fillColor, f32.Transparent, f32.Transparent)
}

func RoundedRect(r f32.Rect, cornerRadius float32, borderThickness float32, fillColor f32.Color, frameColor f32.Color) {
	RR(r, cornerRadius, borderThickness, fillColor, frameColor, f32.Transparent)
}

func RR(r f32.Rect, cornerRadius, borderThickness float32, fillColor, frameColor f32.Color, surfaceColor f32.Color) {
	// Make the quad larger by the shadow width ss  and Correct for device independent pixels
	r.X = r.X * ScaleX
	r.Y = r.Y * ScaleX
	r.W = r.W * ScaleX
	r.H = r.H * ScaleX
	cornerRadius *= ScaleX
	if cornerRadius < 0 || cornerRadius > r.H/2 {
		cornerRadius = r.H / 2
	}
	borderThickness *= ScaleX

	gl.UseProgram(rrprog)
	vbo := gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.Enable(gl.BLEND)

	vertices := []float32{r.X + r.W, r.Y, r.X, r.Y, r.X, r.Y + r.H, r.X, r.Y + r.H,
		r.X + r.W, r.Y + r.H, r.X + r.W, r.Y}
	if borderThickness == 0.0 {
		frameColor = fillColor
	}
	col[0] = fillColor.R
	col[1] = fillColor.G
	col[2] = fillColor.B
	col[3] = fillColor.A
	col[4] = frameColor.R
	col[5] = frameColor.G
	col[6] = frameColor.B
	col[7] = frameColor.A
	col[8] = surfaceColor.R
	col[9] = surfaceColor.G
	col[10] = surfaceColor.B
	col[11] = surfaceColor.A

	gl.BufferData(gl.ARRAY_BUFFER, Bytes(binary.LittleEndian, vertices...), gl.STATIC_DRAW)
	// position attribute
	vertexPositionAttrib := gl.GetAttribLocation(rrprog, "aVertexPosition")
	gl.EnableVertexAttribArray(vertexPositionAttrib)
	gl.VertexAttribPointer(vertexPositionAttrib, 2, gl.FLOAT, false, 8, 0)
	// Colors
	r2 := gl.GetUniformLocation(rrprog, "colors")
	gl.Uniform4fv(r2, col)
	// Set pos data
	r3 := gl.GetUniformLocation(rrprog, "pos")
	gl.Uniform2f(r3, r.X+r.W/2, r.Y+r.H/2)
	// Set halfbox
	r4 := gl.GetUniformLocation(rrprog, "halfbox")
	gl.Uniform2f(r4, r.W/2, r.H/2)
	// Set radius/border width
	r5 := gl.GetUniformLocation(rrprog, "rw")
	gl.Uniform2f(r5, cornerRadius, borderThickness)
	// Do actual drawing
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
	// Free memory
	// gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	// gl.BindVertexArray(0)
	// gl.UseProgram(0)
}

func HorLine(x1, x2, y, w float32, col f32.Color) {
	r := f32.Rect{X: x1, Y: y, W: x2 - x1, H: w}
	RoundedRect(r, 0, w, col, col)
}

func VertLine(x, y1, y2, w float32, col f32.Color) {
	r := f32.Rect{X: x, Y: y1, W: w, H: y2 - y1}
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
	for {
		e := gl.GetError()
		if e == gl.NO_ERROR {
			break
		}
		slog.Error("OpenGl", "error", e)
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

func SetupLogging(defaultLevel slog.Level) {
	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelError)
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: lvl,
	}))
	slog.SetDefault(logger)
	slog.Info("Test output of info")
}
