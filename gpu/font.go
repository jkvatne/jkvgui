package gpu

import (
	"bytes"
	"fmt"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/shader"
	"log"
	"os"
)

// Direction represents the direction in which strings should be rendered.
type Direction uint8

// Known directions.
const (
	LeftToRight Direction = iota // E.g.: Latin
	RightToLeft                  // E.g.: Arabic
	TopToBottom                  // E.g.: Chinese
)

type color struct {
	r float32
	g float32
	b float32
	a float32
}

// SetColor allows you to set the text color to be used when you draw the text
func (f *Font) SetColor(c f32.Color) {
	f.color.R = c.R
	f.color.G = c.G
	f.color.B = c.B
	f.color.A = c.A
}

// Printf draws a string to the screen, takes a list of arguments like printf
func (f *Font) Printf(x, y float32, points float32, max float32, fs string, argv ...interface{}) {
	indices := []rune(fmt.Sprintf(fs, argv...))
	if len(indices) == 0 {
		return
	}
	x *= Scale
	y *= Scale
	if max > 0 {
		max = max*Scale + x
	}
	size := Scale * points / float32(InitialSize)
	SetupDrawing(f.color, f.Vao, f.Program)
	// Iterate through all characters in string
	for i := range indices {
		// get rune
		runeIndex := indices[i]
		if max > 0 && x > max-points*Scale {
			runeIndex = rune(0x2026)
		}

		// find rune in fontChar list
		ch, ok := f.FontChar[runeIndex]
		// load missing runes in batches of 32
		if !ok {
			low := runeIndex - (runeIndex % 32)
			_ = f.GenerateGlyphs(low, low+31)
			ch, ok = f.FontChar[runeIndex]
		}
		// skip runes that are not in font chacter range
		if !ok {
			log.Printf("%c %d\n", runeIndex, runeIndex)
			continue
		}

		// calculate position and size for current rune
		xpos := x + float32(ch.bearingH)*size
		ypos := y - float32(ch.height-ch.bearingV)*size
		w := float32(ch.width) * size
		h := float32(ch.height) * size
		RenderTexture(xpos, ypos, w, h, ch.TextureID, f.Vbo)
		// Now advance cursors for next glyph (note that advance is number of 1/64 pixels)
		x += float32((ch.advance >> 6)) * size // Bitshift by 6 to get value in pixels (2^6 = 64 (divide amount of 1/64th pixels by 64 to get amount of pixels))
		if runeIndex == rune(0x2026) {
			break
		}
	}

}

// Width returns the width of a piece of text in pixels
func (f *Font) Width(scale float32, fs string, argv ...interface{}) float32 {
	var width float32
	indices := []rune(fmt.Sprintf(fs, argv...))
	if len(indices) == 0 {
		return 0
	}
	// Iterate through all characters in string
	for i := range indices {
		// get rune
		runeIndex := indices[i]
		// find rune in fontChar list
		ch, ok := f.FontChar[runeIndex]
		// load missing runes in batches of 32
		if !ok {
			low := runeIndex & rune(32-1)
			_ = f.GenerateGlyphs(low, low+31)
			ch, ok = f.FontChar[runeIndex]
		}
		// skip runes that are not in font chacter range
		if !ok {
			fmt.Printf("%c %d\n", runeIndex, runeIndex)
			continue
		}
		// Now advance cursors for next glyph (note that advance is number of 1/64 pixels)
		width += float32((ch.advance >> 6)) * scale // Bitshift by 6 to get value in pixels (2^6 = 64 (divide amount of 1/64th pixels by 64 to get amount of pixels))
	}
	return width
}

var Fonts []*Font

// LoadFontBytes loads the specified font bytes at the given scale.
func LoadFontBytes(buf []byte, scale float32) (*Font, error) {
	program, _ := shader.NewProgram(shader.VertexQuadShader, shader.FragmentQuadShader)
	fd := bytes.NewReader(buf)
	return LoadTrueTypeFont(program, fd, int32(scale), 32, 127, LeftToRight)
}

// LoadFont loads the specified font at the given scale.
func LoadFontFile(file string, scale int32) (*Font, error) {
	fd, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	program, _ := shader.NewProgram(shader.VertexQuadShader, shader.FragmentQuadShader)
	return LoadTrueTypeFont(program, fd, scale, 32, 127, LeftToRight)
}

func LoadFont(buf []byte, size float32, name string, weight float32) {
	var f *Font
	var err error
	f, err = LoadFontBytes(buf, size)
	if err != nil {
		panic(err)
	}
	f.SetColor(f32.Black)
	f.name = name
	f.weight = weight
	Fonts = append(Fonts, f)
}
