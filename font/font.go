package font

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/shader"
	"log/slog"
	"os"
	"strconv"
)

const Ellipsis = rune(0x2026)

//go:embed fonts/Roboto-Regular.ttf
var Roboto400 []byte // 400

//go:embed fonts/Roboto-MediumItalic.ttf
var RobotoItalic500 []byte

//go:embed fonts/Roboto-SemiBold.ttf
var Roboto600 []byte // 600

//go:embed fonts/RobotoMono-Regular.ttf
var RobotoMono400 []byte

var Fonts [32]*Font

// Direction represents the direction in which strings should be rendered.
type Direction uint8

const (
	LeftToRight Direction = iota
	// RightToLeft
)

var Dpi float32 = 164
var DefaultFontSizePt = 12

// LoadFonts will load the default fonts from embedded data
func LoadFonts() {
	LoadFontBytes(gpu.Normal, Roboto400, DefaultFontSizePt, "RobotoNormal", 400)
	LoadFontBytes(gpu.Bold, Roboto600, DefaultFontSizePt, "RobotoBold", 600)
	LoadFontBytes(gpu.Italic, RobotoItalic500, DefaultFontSizePt, "RobotoItalic", 500)
	LoadFontBytes(gpu.Mono, RobotoMono400, DefaultFontSizePt, "RobotoMono", 400)
}

// Get returns the font with the given number and sets its color.
func Get(no int, color f32.Color) *Font {
	f := Fonts[no]
	f.SetColor(color)
	return f
}

// SetColor allows you to set the text color to be used when you draw the text
func (f *Font) SetColor(c f32.Color) {
	f.color.R = c.R
	f.color.G = c.G
	f.color.B = c.B
	f.color.A = c.A
}

// GetColor returns current font color
func (f *Font) GetColor() f32.Color {
	return f.color
}

// Printf draws a string to the screen, takes a list of arguments like printf
// max is the maximum width. If longer, ellipsis is appended
// scale is the size relative to the default text size.
func (f *Font) Printf(x, y float32, scale float32, maxX float32, fs string, argv ...interface{}) {
	indices := []rune(fmt.Sprintf(fs, argv...))
	if len(indices) == 0 {
		return
	}
	x *= gpu.ScaleX
	y *= gpu.ScaleY
	if maxX > 0 {
		maxX = maxX*gpu.ScaleX + x
	}
	size := gpu.ScaleX * scale * 72 / Dpi
	gpu.SetupDrawing(f.color, f.Vao, f.Program)

	ellipsis, ok := f.FontChar[Ellipsis]
	if !ok {
		_ = f.GenerateGlyphs(Ellipsis, Ellipsis, Dpi)
		ellipsis, ok = f.FontChar[Ellipsis]
	}
	ellipsisWidth := float32(ellipsis.width+1) * size

	// Iterate through all characters in string
	for i := range indices {
		// get rune
		runeIndex := indices[i]

		// find rune in fontChar list
		ch, ok := f.FontChar[runeIndex]
		// load missing runes in batches of 32
		if !ok {
			low := runeIndex - (runeIndex % 32)
			_ = f.GenerateGlyphs(low, low+31, Dpi)
			ch, ok = f.FontChar[runeIndex]
		}
		// skip runes that are not in font character range
		if !ok {
			slog.Error("Illegal rune in printf", "index", runeIndex)
			continue
		}
		//  if x+w+ellipsisw > maxx and not last character then print ellipsis
		//
		// calculate position and size for current rune
		xPos := x + float32(ch.bearingH)*size
		yPos := y - float32(ch.height-ch.bearingV)*size
		w := float32(ch.width) * size
		h := float32(ch.height) * size
		if xPos+w+ellipsisWidth >= maxX && i < len(indices)-1 && maxX > 0 {
			ch = ellipsis
			yPos = y - float32(ch.height-ch.bearingV)*size
			w = float32(ch.width) * size
			h = float32(ch.height) * size
		}
		gpu.RenderTexture(xPos, yPos, w, h, ch.TextureID, f.Vbo)
		// Now advance cursors for next glyph (note that advance is number of 1/64 pixels)
		x += float32(ch.advance>>6) * size // Bitshift by 6 to get value in pixels (2^6 = 64 (divide amount of 1/64th pixels by 64 to get amount of pixels))
		if ch == ellipsis {
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
			_ = f.GenerateGlyphs(low, low+31, Dpi)
			ch, ok = f.FontChar[runeIndex]
		}
		// skip runes that are not in font character range
		if !ok {
			fmt.Printf("%c %d\n", runeIndex, runeIndex)
			continue
		}
		// Now advance cursors for next glyph (note that advance is number of 1/64 pixels)
		width += float32(ch.advance >> 6) // Bitshift by 6 to get value in pixels (2^6 = 64 (divide amount of 1/64th pixels by 64 to get amount of pixels))
	}
	return width * scale * 72 / Dpi
}

// RuneNo will give the rune number at pixel position x from the start
func (f *Font) RuneNo(x float32, scale float32, s string) int {
	runes := []rune(s)
	width := float32(0)
	x = x * Dpi / 72
	// Iterate through all characters in string
	for i, r := range runes {
		// find rune in fontChar list
		ch, ok := f.FontChar[r]
		// skip runes that are not in font character range
		if !ok {
			fmt.Printf("%c %d\n", r, i)
			continue
		}
		// Now advance cursors for next glyph (note that advance is number of 1/64 pixels)
		width += float32(ch.advance>>6) * scale // Bitshift by 6 to get value in pixels (2^6 = 64 (divide amount of 1/64th pixels by 64 to get amount of pixels))
		if width >= x {
			return i
		}
	}
	return len(runes)
}

func (f *Font) Height(size float32) float32 {
	return (f.Ascent + f.Descent) * size * 72 / Dpi
}

func (f *Font) Baseline(size float32) float32 {
	return f.Ascent * size * 72 / Dpi
}

// LoadFontFile loads the specified font at the given size (in pixels).
// The integer returned is the index to Fonts[]
// Will panic if font is not found
func LoadFontFile(no int, file string, size int, name string, weight float32) {
	if no < 0 || no > len(Fonts) {
		panic("LoadFontFile: invalid index " + strconv.Itoa(no))
	}
	program, _ := shader.NewProgram(shader.VertQuadSource, shader.FragQuadSource)
	fd, err := os.Open(file)
	if err != nil {
		panic("Font file not found: " + file)
	}
	defer func(fd *os.File) {
		err := fd.Close()
		if err != nil {
			panic("Could not close file: " + file)
		}
	}(fd)
	f, err := LoadTrueTypeFont(name, program, fd, size, 32, 127, LeftToRight)
	if err != nil {
		panic("Could not load font bytes: " + err.Error())
	}
	f.name = name
	f.weight = weight
	f.size = size
	Fonts[no] = f
}

// LoadFontBytes loads the specified font at the given size (in pixels).
// The integer returned is the index to Fonts[]
// Will panic if font is not found
func LoadFontBytes(no int, buf []byte, size int, name string, weight float32) {
	if no < 0 || no > len(Fonts) {
		panic("LoadFontFile: invalid index " + strconv.Itoa(no))
	}
	program, _ := shader.NewProgram(shader.VertQuadSource, shader.FragQuadSource)
	fd := bytes.NewReader(buf)
	f, err := LoadTrueTypeFont(name, program, fd, size, 32, 127, LeftToRight)
	if err != nil {
		panic("Could not load font bytes: " + err.Error())
	}
	f.SetColor(f32.Black)
	f.name = name
	f.weight = weight
	f.size = size
	Fonts[no] = f
}
