package gpu

import (
	"fmt"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/freetype"
	"github.com/jkvatne/jkvgui/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/draw"
	"io"
)

// A Font allows rendering of text to an OpenGL context.
type Font struct {
	FontChar map[rune]*character
	ttf      *truetype.Font
	scale    int32
	Vao      uint32
	Vbo      uint32
	Program  uint32
	Texture  uint32 // Holds the glyph texture id.
	color    f32.Color
	Ascent   float32
	Descent  float32
	name     string
	weight   float32
}

type character struct {
	TextureID uint32 // ID handle of the glyph texture
	width     int    // glyph width
	height    int    // glyph height
	advance   int    // glyph advance
	bearingH  int    // glyph bearing horizontal
	bearingV  int    // glyph bearing vertical
}

func (f *Font) Name() string {
	return f.name
}

func (f *Font) Weight() float32 {
	return f.weight
}

// GenerateGlyphs builds a set of textures based on a ttf files gylphs
func (f *Font) GenerateGlyphs(low, high rune) error {
	// create a freetype context for drawing
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(f.ttf)
	c.SetFontSize(float64(f.scale))
	c.SetHinting(font.HintingFull)

	// create new face to measure glyph dimensions
	ttfFace := truetype.NewFace(f.ttf, &truetype.Options{
		Size:    float64(f.scale),
		DPI:     72,
		Hinting: font.HintingFull,
	})

	// make each glyph
	for ch := low; ch <= high; ch++ {
		char := new(character)

		gBnd, gAdv, ok := ttfFace.GlyphBounds(ch)
		if ok != true {
			fmt.Printf("ttf face glyphBounds error for ch=%d\n", int(ch))
			continue
		}

		gh := int32((gBnd.Max.Y - gBnd.Min.Y) >> 6)
		gw := int32((gBnd.Max.X - gBnd.Min.X) >> 6)

		// if gylph has no dimensions set to a max value
		if gw == 0 || gh == 0 {
			gBnd = f.ttf.Bounds(fixed.Int26_6(f.scale))
			gw = int32((gBnd.Max.X - gBnd.Min.X) >> 6)
			gh = int32((gBnd.Max.Y - gBnd.Min.Y) >> 6)

			// above can sometimes yield 0 for font smaller than 48pt, 1 is minimum
			if gw == 0 || gh == 0 {
				gw = 1
				gh = 1
			}
		}

		// The glyph's ascent and descent equal -bounds.Min.Y and +bounds.Max.Y.
		gAscent := int(-gBnd.Min.Y) >> 6
		gDescent := int(gBnd.Max.Y) >> 6
		f.Ascent = max(f.Ascent, float32(gAscent))
		f.Descent = max(f.Descent, float32(gDescent))

		// set w,h and adv, bearing V and bearing H in char
		char.width = int(gw)
		char.height = int(gh)
		char.advance = int(gAdv)
		char.bearingV = gDescent
		char.bearingH = (int(gBnd.Min.X) >> 6)

		// create image to draw glyph
		fg, bg := image.White, image.Black
		rect := image.Rect(0, 0, int(gw), int(gh))
		rgba := image.NewRGBA(rect)
		draw.Draw(rgba, rgba.Bounds(), bg, image.Point{}, draw.Src)

		// set the glyph dot
		px := 0 - (int(gBnd.Min.X) >> 6)
		py := (gAscent)
		pt := freetype.Pt(px, py)
		// Draw the text from mask to image
		c.SetClip(rgba.Bounds())
		c.SetDst(rgba)
		c.SetSrc(fg)
		_, err := c.DrawString(string(ch), pt)
		if err != nil {
			return err
		}
		// Generate texture
		char.TextureID = GenerateTexture(rgba)
		// add char to fontChar list
		f.FontChar[ch] = char
	}

	return nil
}

// LoadTrueTypeFont builds OpenGL buffers and glyph textures based on a ttf file
func LoadTrueTypeFont(program uint32, r io.Reader, scale int32, low, high rune, dir Direction) (*Font, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	// Read the truetype font.
	ttf, err := truetype.Parse(data)
	if err != nil {
		return nil, err
	}

	// make Font stuct type
	f := new(Font)
	f.FontChar = make(map[rune]*character)
	f.ttf = ttf
	f.scale = scale
	f.Program = program   // set shader program
	f.SetColor(f32.Black) // set default black

	err = f.GenerateGlyphs(low, high)
	if err != nil {
		return nil, err
	}
	ConfigureVaoVbo(&f.Vao, &f.Vbo, f.Program)
	return f, nil
}

// ConfigureVaoVbo for texture quads
func ConfigureVaoVbo(vao *uint32, vbo *uint32, program uint32) {
	gl.GenVertexArrays(1, vao)
	gl.BindVertexArray(*vao)
	gl.GenBuffers(1, vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, *vbo)

	gl.BufferData(gl.ARRAY_BUFFER, 6*4*4, nil, gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointerWithOffset(vertAttrib, 2, gl.FLOAT, false, 4*4, 0)
	defer gl.DisableVertexAttribArray(vertAttrib)

	texCoordAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointerWithOffset(texCoordAttrib, 2, gl.FLOAT, false, 4*4, 2*4)
	defer gl.DisableVertexAttribArray(texCoordAttrib)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
}

func GenerateTexture(rgba *image.RGBA) uint32 {
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA,
		int32(rgba.Rect.Dx()), int32(rgba.Rect.Dy()), 0,
		gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))
	GetErrors()
	gl.BindTexture(gl.TEXTURE_2D, 0)
	return texture
}
