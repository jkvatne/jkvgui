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
)

//go:embed fonts/Roboto-Thin.ttf
var Roboto100 []byte // 100

//go:embed fonts/Roboto-ExtraLight.ttf
var Roboto200 []byte // 200

//go:embed fonts/Roboto-Light.ttf
var Roboto300 []byte // 300

//go:embed fonts/Roboto-Regular.ttf
var Roboto400 []byte // 400

//go:embed fonts/Roboto-Medium.ttf
var Roboto500 []byte // 500

//go:embed fonts/Roboto-MediumItalic.ttf
var RobotoItalic500 []byte

//go:embed fonts/Roboto-SemiBold.ttf
var Roboto600 []byte // 600

//go:embed fonts/Roboto-Bold.ttf
var Roboto700 []byte // 700

//go:embed fonts/Roboto-Bold.ttf
var Roboto800 []byte // 800

//go:embed fonts/Roboto-Bold.ttf
var Roboto900 []byte // 900

//go:embed fonts/RobotoMono-Regular.ttf
var RobotoMono400 []byte

//go:embed fonts/RobotoMono-Bold.ttf
var RobotoMono600 []byte

//go:embed fonts/RobotoMono-Light.ttf
var RobotoMono200 []byte

var Fonts []*Font

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

var OverSampling = float32(2.0)

// From freetype.go, line 263, Her c.dpi is allways 72.
// c.scale = fixed.Int26_6(0.5 + (c.fontSize * c.dpi * 64 / 72))
// size = fontsize  in pixels.
func LoadFonts() {
	LoadFontBytes(gpu.Normal, Roboto500, 24, "RobotoNormal", 400)
	LoadFontBytes(gpu.Bold, Roboto600, 24, "RobotoBold", 600)
	LoadFontBytes(gpu.Italic, RobotoItalic500, 24, "RobotoItalic", 500)
	LoadFontBytes(gpu.Mono, RobotoMono400, 24, "RobotoMono", 400)
}

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
// max is the maximum width. If longer, elipsis is appended
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
	size := gpu.ScaleX * scale / OverSampling
	gpu.SetupDrawing(f.color, f.Vao, f.Program)

	ch, ok := f.FontChar[rune(0x2026)]
	if !ok {
		_ = f.GenerateGlyphs(rune(0x2026), rune(0x2026))
		ch, ok = f.FontChar[rune(0x2026)]
	}
	elipsisWidth := float32(ch.width) + 1

	// Iterate through all characters in string
	for i := range indices {
		// get rune
		runeIndex := indices[i]
		if maxX > 0 && x > (maxX-elipsisWidth) {
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
			slog.Error("Illegal rune in printf", "index", runeIndex)
			continue
		}

		// calculate position and size for current rune
		xpos := x + float32(ch.bearingH)*size
		ypos := y - float32(ch.height-ch.bearingV)*size
		w := float32(ch.width) * size
		h := float32(ch.height) * size
		gpu.RenderTexture(xpos, ypos, w, h, ch.TextureID, f.Vbo)
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
		width += float32((ch.advance >> 6)) // Bitshift by 6 to get value in pixels (2^6 = 64 (divide amount of 1/64th pixels by 64 to get amount of pixels))
	}
	return width * scale / OverSampling
}

// RuneNo will will give the rune number at pixel posision x from the start
func (f *Font) RuneNo(x float32, scale float32, s string) int {
	runes := []rune(s)
	width := float32(0)
	// Iterate through all characters in string
	for i, rune := range runes {
		// find rune in fontChar list
		ch, ok := f.FontChar[rune]
		// skip runes that are not in font chacter range
		if !ok {
			fmt.Printf("%c %d\n", rune, i)
			continue
		}
		// Now advance cursors for next glyph (note that advance is number of 1/64 pixels)
		width += float32((ch.advance >> 6)) * scale / OverSampling // Bitshift by 6 to get value in pixels (2^6 = 64 (divide amount of 1/64th pixels by 64 to get amount of pixels))
		if width >= x {
			return i
		}
	}
	return len(runes)
}

func (f *Font) Height(size float32) float32 {
	return (f.Ascent + f.Descent) * size / OverSampling
}

func (f *Font) Baseline(size float32) float32 {
	return f.Ascent * size / OverSampling
}

// LoadFontFile loads the specified font at the given size (in pixels).
// The integer returened is the index to Fonts[]
func LoadFontFile(file string, size int, name string, weight float32) int {
	program, _ := shader.NewProgram(shader.VertexQuadShader, shader.FragmentQuadShader)
	fd, err := os.Open(file)
	if err != nil {
		panic("Font file not found: " + file)
	}
	defer fd.Close()
	f, err := LoadTrueTypeFont(program, fd, size, 32, 127, LeftToRight)
	if err != nil {
		panic("Could not load font bytes: " + err.Error())
	}
	f.name = name
	f.weight = weight
	f.size = size
	Fonts = append(Fonts, f)
	return len(Fonts) - 1
}

// LoadFontBytesloads the specified font at the given size (in pixels).
// The integer returened is the index to Fonts[]
func LoadFontBytes(no int, buf []byte, size int, name string, weight float32) {
	program, _ := shader.NewProgram(shader.VertexQuadShader, shader.FragmentQuadShader)
	fd := bytes.NewReader(buf)
	f, err := LoadTrueTypeFont(program, fd, size, 32, 127, LeftToRight)
	if err != nil {
		panic("Could not load font bytes: " + err.Error())
	}
	f.SetColor(f32.Black)
	f.name = name
	f.weight = weight
	f.size = size
	if len(Fonts) == 0 {
		Fonts = make([]*Font, 16)
	}
	// for no >= len(Fonts) {Fonts = append(Fonts, nil)}
	Fonts[no] = f
}
