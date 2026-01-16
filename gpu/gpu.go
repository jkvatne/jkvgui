package gpu

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log/slog"
	"math"
	"os"

	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/purego-glfw/gl"
)

type IntRect struct{ X, Y, W, H int }

// GlData is Open-GL global variables
type GlData struct {
	RRprogram     uint32
	ShaderProgram uint32
	ImgProgram    uint32
	FontProgram   uint32
	Vao           uint32
	Vbo           uint32
	FontVao       uint32
	FontVbo       uint32
	ScaleX        float32
	ScaleY        float32
	HeightPx      int
	WidthPx       int
}

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

func (gd *GlData) Clip(r f32.Rect) {
	ww := r.W * gd.ScaleX
	hh := r.H * gd.ScaleY
	xx := r.X * gd.ScaleX
	yy := float32(gd.HeightPx) - hh - r.Y*gd.ScaleY
	gl.Scissor(int32(xx), int32(yy), int32(ww), int32(hh))
	gl.Enable(gl.SCISSOR_TEST)
}

func NoClip() {
	gl.Disable(gl.SCISSOR_TEST)
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

func SetResolution(program uint32, w, h int32) {
	if program == 0 {
		panic("Program number must be greater than 0")
	}
	// Activate the corresponding render state
	gl.UseProgram(program)
	// set screen resolution
	gl.Viewport(0, 0, w, h)
	resUniform := gl.GetUniformLocation(program, gl.Str("resolution\x00"))
	gl.Uniform2f(resUniform, float32(w), float32(h))
}

func (gd *GlData) InitGpu() {
	gl.Enable(gl.BLEND)
	gl.Enable(gl.MULTISAMPLE)
	gl.BlendEquation(gl.FUNC_ADD)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.ClearColor(1, 1, 1, 1)
	GetErrors("InitGpu() startup")
	// Set up the programs needed
	gd.RRprogram, _ = NewProgram(VertRectSource, FragRectSource)
	gd.ShaderProgram, _ = NewProgram(VertRectSource, FragShadowSource)
	gd.ImgProgram, _ = NewProgram(VertQuadSource, FragImgSource)
	gd.FontProgram, _ = NewProgram(VertQuadSource, FragQuadSource)
	// Setup image drawing
	gl.GenVertexArrays(1, &gd.Vao)
	gl.BindVertexArray(gd.Vao)
	gl.GenBuffers(1, &gd.Vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, gd.Vbo)
	GetErrors("InitGpu() Vbo Vao")
	vertAttrib := uint32(gl.GetAttribLocation(gd.ImgProgram, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointerWithOffset(vertAttrib, 2, gl.FLOAT, false, 4*4, 0)
	GetErrors("InitGpu() vertexAttrib")
	texCoordAttrib := uint32(gl.GetAttribLocation(gd.ImgProgram, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointerWithOffset(texCoordAttrib, 2, gl.FLOAT, false, 4*4, 2*4)
	GetErrors("InitGpu() texCoord")
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
	GetErrors("InitGpu() release buffers")
	// Setup font drawing
	gl.GenVertexArrays(1, &gd.FontVao)
	gl.BindVertexArray(gd.FontVao)
	gl.GenBuffers(1, &gd.FontVbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, gd.FontVbo)
	GetErrors("InitWindow setup FontVaoVbo")
	gl.BufferData(gl.ARRAY_BUFFER, 6*4*4, nil, gl.STATIC_DRAW)
	GetErrors("InitGpu() font buffer data")
	vertAttrib = uint32(gl.GetAttribLocation(gd.FontProgram, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointerWithOffset(vertAttrib, 2, gl.FLOAT, false, 4*4, 0)
	GetErrors("InitGpu() font vertAttrib")
	texCoordAttrib = uint32(gl.GetAttribLocation(gd.FontProgram, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointerWithOffset(texCoordAttrib, 2, gl.FLOAT, false, 4*4, 2*4)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	GetErrors("InitGpu() texCoordAttrib")
	gl.BindVertexArray(0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	GetErrors("InitGpu() exiting")
}

func SetBackgroundColor(col f32.Color) {
	gl.ClearColor(col.R, col.G, col.B, col.A)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func (gd *GlData) Shade(r f32.Rect, cornerRadius float32, fillColor f32.Color, shadowSize float32) {
	// Make the quad larger by the shadow width ss and Correct for device independent pixels
	r.X = (r.X - shadowSize*0.75) * gd.ScaleX
	r.Y = (r.Y - shadowSize*0.75) * gd.ScaleX
	r.W = (r.W + shadowSize*1.5) * gd.ScaleX
	r.H = (r.H + shadowSize*1.5) * gd.ScaleX
	shadowSize *= gd.ScaleX
	cornerRadius *= gd.ScaleX
	if cornerRadius < 0 {
		cornerRadius = r.H / 2
	}
	cornerRadius = max(0, min(min(r.H/2, r.W/2), cornerRadius+shadowSize))
	gl.UseProgram(gd.ShaderProgram)
	gl.BindVertexArray(gd.Vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, gd.Vbo)
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
	r2 := gl.GetUniformLocation(gd.ShaderProgram, gl.Str("colors\x00"))
	gl.Uniform4fv(r2, 16, &col[0])
	// Set pos data
	r3 := gl.GetUniformLocation(gd.ShaderProgram, gl.Str("pos\x00"))
	gl.Uniform2f(r3, r.X+r.W/2, r.Y+r.H/2)
	// Set halfbox
	r4 := gl.GetUniformLocation(gd.ShaderProgram, gl.Str("halfbox\x00"))
	gl.Uniform2f(r4, r.W/2, r.H/2)
	// Set radius/border width
	r5 := gl.GetUniformLocation(gd.ShaderProgram, gl.Str("rws\x00"))
	gl.Uniform4f(r5, cornerRadius, 0, shadowSize, 0)
	// Do actual drawing
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
	// Free memory
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
	gl.UseProgram(0)
	GetErrors("Shade")
}

func (gd *GlData) RoundedRect(r f32.Rect, cornerRadius float32, borderThickness float32, fillColor f32.Color, frameColor f32.Color) {
	gd.RR(r, cornerRadius, borderThickness, fillColor, frameColor, f32.Transparent)
}

func i(x float32) float32 {
	return float32(int(x + 0.5))
}

func (gd *GlData) RR(r f32.Rect, cornerRadius, borderThickness float32, fillColor, frameColor f32.Color, surfaceColor f32.Color) {
	// Make the quad larger by the shadow width ss and Correct for device independent pixels
	r.X = i(r.X * gd.ScaleX)
	r.Y = i(r.Y * gd.ScaleX)
	r.W = i(r.W * gd.ScaleX)
	r.H = i(r.H * gd.ScaleX)
	cornerRadius *= gd.ScaleX
	if cornerRadius < 0 || cornerRadius > r.H/2 {
		cornerRadius = r.H / 2
	}
	borderThickness = i(borderThickness * gd.ScaleX)

	gl.UseProgram(gd.RRprogram)
	gl.BindVertexArray(gd.Vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, gd.Vbo)
	gl.Enable(gl.BLEND)

	vertices := []float32{r.X + r.W, r.Y, r.X, r.Y, r.X, r.Y + r.H, r.X, r.Y + r.H,
		r.X + r.W, r.Y + r.H, r.X + r.W, r.Y}
	if borderThickness == 0.0 {
		frameColor = fillColor
	}
	var col [12]float32
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
	r2 := gl.GetUniformLocation(gd.RRprogram, gl.Str("colors\x00"))
	gl.Uniform4fv(r2, 16, &col[0])
	// Set pos data
	r3 := gl.GetUniformLocation(gd.RRprogram, gl.Str("pos\x00"))
	gl.Uniform2f(r3, r.X+r.W/2, r.Y+r.H/2)
	// Set halfbox
	r4 := gl.GetUniformLocation(gd.RRprogram, gl.Str("halfbox\x00"))
	gl.Uniform2f(r4, r.W/2, r.H/2)
	// Set radius/border width
	r5 := gl.GetUniformLocation(gd.RRprogram, gl.Str("rw\x00"))
	gl.Uniform2f(r5, cornerRadius, borderThickness)
	// Do actual drawing
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
	// Free memory
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
	gl.UseProgram(0)
	GetErrors("RoundedRect")
}

func (gd *GlData) HorLine(x1, x2, y, w float32, col f32.Color) {
	r := f32.Rect{X: x1, Y: y, W: x2 - x1, H: w}
	gd.RoundedRect(r, 0, w, col, col)
}

func (gd *GlData) VertLine(x, y1, y2, w float32, col f32.Color) {
	r := f32.Rect{X: x, Y: y1, W: w, H: y2 - y1}
	gd.RoundedRect(r, 0, w, col, col)
}

func (gd *GlData) SolidRect(r f32.Rect, fillColor f32.Color) {
	gd.RoundedRect(r, 0, 0, fillColor, fillColor)
}

func (gd *GlData) OutlinedRect(r f32.Rect, frameThickness float32, frameColor f32.Color) {
	gd.RoundedRect(r, 0, frameThickness, f32.Transparent, frameColor)
}

func GetErrors(s string) {
	for {
		e := gl.GetError()
		if e == gl.NO_ERROR {
			return
			slog.Error("OpenGl", "error", e, "from", s)
		}
	}
}

func sqDiff(x, y uint8) uint64 {
	d := uint64(x) - uint64(y)
	return d * d
}

func Compare(img1, img2 *image.RGBA) (int64, error) {
	if img1 == nil || img2 == nil {
		return 0, fmt.Errorf("images can not be nil")
	}
	if img1.Bounds() != img2.Bounds() {
		return 0, fmt.Errorf("image bounds not equal: %+v, %+v", img1.Bounds(), img2.Bounds())
	}
	accumError := int64(0)
	for i := 0; i < len(img1.Pix); i++ {
		accumError += int64(sqDiff(img1.Pix[i], img2.Pix[i]))
	}
	return int64(math.Sqrt(float64(accumError))), nil
}
