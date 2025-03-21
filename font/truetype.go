package font

import (
	"fmt"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/freetype"
	"github.com/jkvatne/jkvgui/freetype/truetype"
	"github.com/jkvatne/jkvgui/gpu"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/draw"
	"image/png"
	"io"
	"log/slog"
	"os"
	"strconv"
)

// A Font allows rendering of text to an OpenGL context.
type Font struct {
	FontChar map[rune]*character
	ttf      *truetype.Font
	Vao      uint32
	Vbo      uint32
	Program  uint32
	Texture  uint32 // Holds the glyph texture id.
	color    f32.Color
	Ascent   float32
	Descent  float32
	name     string
	size     int
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
func (f *Font) GenerateGlyphs(low, high rune, dpi float32) error {
	// create a freetype context for drawing
	c := freetype.NewContext()
	c.SetDPI(float64(dpi))
	c.SetFont(f.ttf)
	c.SetFontSize(float64(f.size))
	c.SetHinting(font.HintingFull)

	// create new face to measure glyph dimensions
	ttfFace := truetype.NewFace(f.ttf, &truetype.Options{
		Size:    float64(f.size),
		DPI:     float64(dpi),
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
			gBnd = f.ttf.Bounds(fixed.Int26_6(f.size))
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
		if ch == 'E' {
			file, err := os.Create("./test-outputs/E-" + f.name + "-" + strconv.Itoa(int(dpi)) + ".png")
			if err != nil {
				slog.Error(err.Error())
			}
			_ = png.Encode(file, rgba)
			_ = file.Close()
		}
		// Generate texture
		char.TextureID = gpu.GenerateTexture(rgba)
		// add char to fontChar list
		f.FontChar[ch] = char
	}
	return nil
}

// LoadTrueTypeFont builds OpenGL buffers and glyph textures based on a ttf file
func LoadTrueTypeFont(name string, program uint32, r io.Reader, size int, low, high rune, dir Direction) (*Font, error) {
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
	f.size = size
	f.Program = program   // set shader program
	f.SetColor(f32.Black) // set default black
	f.name = name
	err = f.GenerateGlyphs(low, high, Dpi)
	if err != nil {
		return nil, err
	}
	gpu.ConfigureVaoVbo(&f.Vao, &f.Vbo, f.Program)
	return f, nil
}
