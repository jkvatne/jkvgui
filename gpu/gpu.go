package gpu

import (
	"fmt"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gl"
	"image"
	"image/draw"
	"image/png"
	"log/slog"
	"math"
	"os"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

var ( // Public global variables
	WindowRect     f32.Rect
	WindowHeightPx int
	WindowWidthPx  int
	ScaleX         float32 = 1.0
	ScaleY         float32 = 1.0
	UserScale      float32 = 1.0
	Mutex          sync.Mutex
	InvalidateChan = make(chan time.Duration, 1)
)

var ( // Private global variables
	RRprog      uint32
	ShaderProg  uint32
	ImgProgram  uint32
	Vao         uint32
	Vbo         uint32
	FontProgram uint32
	FontVao     uint32
	FontVbo     uint32
)

const (
	Normal14 int = iota
	Bold14
	Bold16
	Bold20
	Italic14
	Mono14
	Normal12
	Normal16
	Normal20
	Bold12
	Italic12
	Mono12
	Normal10
	Bold10
	Italic10
	Mono10
)

var DeferredFunctions []func()
var HintActive bool

func Defer(f func()) {
	for _, g := range DeferredFunctions {
		if &f == &g {
			return
		}
	}
	DeferredFunctions = append(DeferredFunctions, f)
}

func RunDeferred() {
	for _, f := range DeferredFunctions {
		f()
	}
	DeferredFunctions = DeferredFunctions[0:0]
	HintActive = false
}

func NoClip() {
	gl.Disable(gl.SCISSOR_TEST)
}

func Clip(r f32.Rect) {
	ww := int32(float32(r.W) * ScaleX)
	hh := int32(float32(r.H) * ScaleY)
	xx := int32(float32(r.X) * ScaleX)
	yy := int32(WindowHeightPx) - hh - int32(float32(r.Y)*ScaleY)
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
	GetErrors("Capture")
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
	for _, p := range Programs {
		SetResolution(p)
	}
}

func SetResolution(program uint32) {
	if program == 0 {
		panic("Program number must be greater than 0")
	}
	// Activate the corresponding render state
	gl.UseProgram(program)
	// set screen resolution
	gl.Viewport(0, 0, int32(WindowWidthPx), int32(WindowHeightPx))
	resUniform := gl.GetUniformLocation(program, gl.Str("resolution\x00"))
	gl.Uniform2f(resUniform, float32(WindowWidthPx), float32(WindowHeightPx))
}

func InitGpu() {
	// Initialize gl
	if err := gl.Init(); err != nil {
		panic("Initialization error for OpenGL: " + err.Error())
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	slog.Info("OpenGL", "version", version)
	gl.Enable(gl.BLEND)
	gl.Enable(gl.MULTISAMPLE)
	gl.BlendEquation(gl.FUNC_ADD)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.ClearColor(1, 1, 1, 1)
	GetErrors("InitWindow start")

	// Set up rounded rectangle drawing and shader drawing
	RRprog, _ = NewProgram(VertRectSource, FragRectSource)
	ShaderProg, _ = NewProgram(VertRectSource, FragShadowSource)

	// Setup image drawing
	gl.GenVertexArrays(1, &Vao)
	gl.BindVertexArray(Vao)
	gl.GenBuffers(1, &Vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, Vbo)
	GetErrors("InitWindow Vbo Vao")
	ImgProgram, _ = NewProgram(VertQuadSource, FragImgSource)
	vertAttrib := uint32(gl.GetAttribLocation(ImgProgram, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointerWithOffset(vertAttrib, 2, gl.FLOAT, false, 4*4, 0)
	defer gl.DisableVertexAttribArray(vertAttrib)
	GetErrors("InitWindow vertexAttrib")
	texCoordAttrib := uint32(gl.GetAttribLocation(ImgProgram, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointerWithOffset(texCoordAttrib, 2, gl.FLOAT, false, 4*4, 2*4)
	defer gl.DisableVertexAttribArray(texCoordAttrib)
	GetErrors("InitWindow texCoord")
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
	GetErrors("InitWindow release buffers")

	// Setup font drawing
	gl.GenVertexArrays(1, &FontVao)
	gl.BindVertexArray(FontVao)
	gl.GenBuffers(1, &FontVbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, FontVbo)
	GetErrors("InitWindow setup FontVaoVbo")
	gl.BufferData(gl.ARRAY_BUFFER, 6*4*4, nil, gl.STATIC_DRAW)
	GetErrors("InitWindow font buffer data")
	FontProgram, _ = NewProgram(VertQuadSource, FragQuadSource)
	GetErrors("InitWindow FontProgram")
	vertAttrib = uint32(gl.GetAttribLocation(FontProgram, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointerWithOffset(vertAttrib, 2, gl.FLOAT, false, 4*4, 0)
	GetErrors("InitWindow font vertAttrib")
	defer gl.DisableVertexAttribArray(vertAttrib)
	texCoordAttrib = uint32(gl.GetAttribLocation(FontProgram, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointerWithOffset(texCoordAttrib, 2, gl.FLOAT, false, 4*4, 2*4)
	defer gl.DisableVertexAttribArray(texCoordAttrib)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	GetErrors("InitWindow texCoordAttrib")
	gl.BindVertexArray(0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	GetErrors("InitWindow exiting")
}

func SetBackgroundColor(col f32.Color) {
	gl.ClearColor(col.R, col.G, col.B, col.A)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func Shade(r f32.Rect, cornerRadius float32, fillColor f32.Color, shadowSize float32) {
	// Make the quad larger by the shadow width ss and Correct for device independent pixels
	r.X = (r.X - shadowSize*0.75) * ScaleX
	r.Y = (r.Y - shadowSize*0.75) * ScaleX
	r.W = (r.W + shadowSize*1.5) * ScaleX
	r.H = (r.H + shadowSize*1.5) * ScaleX
	shadowSize *= ScaleX
	cornerRadius *= ScaleX
	if cornerRadius < 0 {
		cornerRadius = r.H / 2
	}
	cornerRadius = max(0, min(min(r.H/2, r.W/2), cornerRadius+shadowSize))

	gl.UseProgram(ShaderProg)
	gl.BindVertexArray(Vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, Vbo)
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
	r2 := gl.GetUniformLocation(ShaderProg, gl.Str("colors\x00"))
	gl.Uniform4fv(r2, 16, &col[0])
	// Set pos data
	r3 := gl.GetUniformLocation(ShaderProg, gl.Str("pos\x00"))
	gl.Uniform2f(r3, r.X+r.W/2, r.Y+r.H/2)
	// Set halfbox
	r4 := gl.GetUniformLocation(ShaderProg, gl.Str("halfbox\x00"))
	gl.Uniform2f(r4, r.W/2, r.H/2)
	// Set radius/border width
	r5 := gl.GetUniformLocation(ShaderProg, gl.Str("rws\x00"))
	gl.Uniform4f(r5, cornerRadius, 0, shadowSize, 0)
	// Do actual drawing
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
	// Free memory
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
	gl.UseProgram(0)
	GetErrors("Shade")

}

var col [12]float32

func RoundedRect(r f32.Rect, cornerRadius float32, borderThickness float32, fillColor f32.Color, frameColor f32.Color) {
	RR(r, cornerRadius, borderThickness, fillColor, frameColor, f32.Transparent)
}

func i(x float32) float32 {
	return float32(int(x + 0.5))
}

func RR(r f32.Rect, cornerRadius, borderThickness float32, fillColor, frameColor f32.Color, surfaceColor f32.Color) {
	// Make the quad larger by the shadow width ss and Correct for device independent pixels
	r.X = i(r.X * ScaleX)
	r.Y = i(r.Y * ScaleX)
	r.W = i(r.W * ScaleX)
	r.H = i(r.H * ScaleX)
	cornerRadius *= ScaleX
	if cornerRadius < 0 || cornerRadius > r.H/2 {
		cornerRadius = r.H / 2
	}
	borderThickness = i(borderThickness * ScaleX)

	gl.UseProgram(RRprog)
	gl.BindVertexArray(Vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, Vbo)
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

	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)
	// position attribute
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 2*4, nil)
	gl.EnableVertexAttribArray(1)
	// Colors
	r2 := gl.GetUniformLocation(RRprog, gl.Str("colors\x00"))
	gl.Uniform4fv(r2, 16, &col[0])
	// Set pos data
	r3 := gl.GetUniformLocation(RRprog, gl.Str("pos\x00"))
	gl.Uniform2f(r3, r.X+r.W/2, r.Y+r.H/2)
	// Set halfbox
	r4 := gl.GetUniformLocation(RRprog, gl.Str("halfbox\x00"))
	gl.Uniform2f(r4, r.W/2, r.H/2)
	// Set radius/border width
	r5 := gl.GetUniformLocation(RRprog, gl.Str("rw\x00"))
	gl.Uniform2f(r5, cornerRadius, borderThickness)
	// Do actual drawing
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
	// Free memory
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
	gl.UseProgram(0)
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

func GetErrors(s string) {
	e := gl.GetError()
	if e != gl.NO_ERROR {
		slog.Error("OpenGl ", "error", e, "from", s)
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

func Invalidate(dt time.Duration) {
	select {
	case InvalidateChan <- dt:
		return
	default:
		return
	}

}

var BlinkFrequency = 2
var BlinkState atomic.Bool
var Blinking atomic.Bool

func blinker() {
	for {
		time.Sleep(time.Second / time.Duration(BlinkFrequency*2))
		b := BlinkState.Load()
		BlinkState.Store(!b)
		if Blinking.Load() {
			InvalidateChan <- 0
		}
	}
}

func init() {
	go blinker()
}
