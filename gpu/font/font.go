package font

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/gpu/font/freetype"
	"github.com/jkvatne/jkvgui/gpu/font/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/draw"
	"image/png"
	"io"
	"log/slog"
	"os"
	"strconv"
	"strings"
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

// DefaultDpi is the value used by the freetype library
var DefaultDpi float32 = 72

var Dpi float32 = 108

var DefaultFontSize = 14

// A Font allows rendering of text to an OpenGL context.
type Font struct {
	FontChar map[rune]*charInfo
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

type charInfo struct {
	TextureID uint32 // ID handle of the glyph texture
	width     int    // glyph width
	height    int    // glyph height
	advance   int    // glyph advance
	bearingH  int    // glyph bearing horizontal
	bearingV  int    // glyph bearing vertical
}

// LoadFonts will load the default fonts from embedded data
func LoadFonts(Fontsize int) {
	LoadFontBytes(gpu.Normal, Roboto400, Fontsize, "RobotoNormal", 400)
	LoadFontBytes(gpu.Bold, Roboto600, Fontsize, "RobotoBold", 600)
	LoadFontBytes(gpu.Italic, RobotoItalic500, Fontsize, "RobotoItalic", 500)
	LoadFontBytes(gpu.Mono, RobotoMono400, Fontsize, "RobotoMono", 400)
}

// Get returns the font with the given number
func Get(no int) *Font {
	f := Fonts[no]
	return f
}

// GetColor returns current font color
func (f *Font) GetColor() f32.Color {
	return f.color
}

func assertRune(f *Font, r rune) *charInfo {
	ch, ok := f.FontChar[r]
	if !ok {
		err := f.GenerateGlyphs(r, r, Dpi)
		if err == nil {
			ch, ok = f.FontChar[r]
		}
	}
	// skip runes that are not in font character range
	if !ok {
		slog.Error("Illegal rune", "index", r)
	}
	return ch
}

// DrawText draws a string to the screen, takes a list of arguments like printf
// max is the maximum width. If longer, ellipsis is appended
// scale is the size relative to the default text size, typically between 0.7 and 2.5.
func (f *Font) DrawText(x, y float32, color f32.Color, scale float32, maxW float32, dir gpu.Direction, fs string, argv ...interface{}) {
	runes := []rune(fmt.Sprintf(fs, argv...))
	if len(runes) == 0 {
		return
	}
	x *= gpu.ScaleX
	y *= gpu.ScaleY
	maxW *= gpu.ScaleX
	size := gpu.ScaleX * scale * DefaultDpi / Dpi
	gpu.SetupDrawing(color, f.Vao, f.Program)
	ellipsis := assertRune(f, Ellipsis)
	ellipsisWidth := float32(ellipsis.width+1) * size
	var offset float32
	// Iterate through all characters in string
	for i := range runes {
		ch := assertRune(f, runes[i])
		bearingH := float32(ch.bearingH) * size
		bearingV := float32(ch.bearingV) * size
		if maxW > 0 && offset+ellipsisWidth+float32(ch.width)*size+bearingH >= maxW && i < len(runes)-1 {
			ch = ellipsis
		}
		w := float32(ch.width) * size
		h := float32(ch.height) * size
		// calculate position and size for current rune
		if dir == gpu.LTR {
			xPos := x + offset + bearingH
			yPos := y - h + bearingV
			gpu.RenderTexture(xPos, yPos, w, h, ch.TextureID, f.Vbo, dir)
		} else if dir == gpu.TTB {
			xPos := x - bearingV
			yPos := y + offset + bearingH
			gpu.RenderTexture(xPos, yPos, h, w, ch.TextureID, f.Vbo, dir)
		} else if dir == gpu.BTT {
			xPos := x - h + bearingV
			yPos := y - offset - w
			gpu.RenderTexture(xPos, yPos, h, w, ch.TextureID, f.Vbo, dir)
		}
		offset += float32(ch.advance>>6) * size
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
		ch := assertRune(f, indices[i])
		width += float32(ch.advance >> 6)
	}
	return width * scale * DefaultDpi / Dpi
}

// RuneNo will give the rune number at pixel position x from the start
func (f *Font) RuneNo(x float32, scale float32, s string) int {
	runes := []rune(s)
	width := float32(0)
	x = x * Dpi / DefaultDpi
	// Iterate through all characters in string
	for i, r := range runes {
		// find rune in fontChar list
		ch, ok := f.FontChar[r]
		// skip runes that are not in font character range
		if !ok {
			fmt.Printf("%c %d\n", r, i)
			continue
		}
		width += float32(ch.advance>>6) * scale
		if width >= x {
			return i
		}
	}
	return len(runes)
}

// Height returns the font height at the given size
func (f *Font) Height(size float32) float32 {
	return (f.Ascent + f.Descent) * size * DefaultDpi / Dpi
}

// Baseline returns the offset from the top to the font's baseline (the dot)
func (f *Font) Baseline(size float32) float32 {
	return f.Ascent * size * DefaultDpi / Dpi
}

// LoadFontFile loads the specified font at the given size (in pixels).
// The integer returned is the index to Fonts[]
// Will panic if font is not found
func LoadFontFile(no int, file string, size int, name string, weight float32) {
	if no < 0 || no > len(Fonts) {
		panic("LoadFontFile: invalid index " + strconv.Itoa(no))
	}
	program, _ := gpu.NewProgram(gpu.VertQuadSource, gpu.FragQuadSource)
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
	f, err := LoadTrueTypeFont(name, program, fd, size)
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
	f32.ExitIf(no < 0 || no > len(Fonts), "LoadFontFile: invalid index "+strconv.Itoa(no))
	program, err := gpu.NewProgram(gpu.VertQuadSource, gpu.FragQuadSource)
	f32.ExitOn(err, "Could not generate font shader program")
	fd := bytes.NewReader(buf)
	f, err := LoadTrueTypeFont(name, program, fd, size)
	f32.ExitOn(err, "Could not load font bytes")
	f.name = name
	f.weight = weight
	f.size = size
	Fonts[no] = f
}

// Split will split a long string into an array of shorter strings that will fit within maxWidth
func Split(s string, maxWidth float32, font *Font, scale float32) []string {
	var width float32
	lines := make([]string, 0)
	words := strings.Split(s, " ")
	line := ""
	for _, word := range words {
		if word == "" {
			continue
		}
		width = font.Width(scale, line+" "+word)
		if width <= maxWidth {
			line = line + word + " "
		} else {
			if len(line) > 0 {
				// Use words up to the current word
				lines = append(lines, line)
				line = word + " "
			} else {
				// Hard break a very long word
				for j := len(word) - 1; j >= 1; j-- {
					if font.Width(scale, word[0:j]) > maxWidth {
						line = word[0:j]
						word = word[j:]
						break
					}
				}
				lines = append(lines, word)
			}
		}
	}
	lines = append(lines, line)
	return lines
}

// Name returns the font's name
func (f *Font) Name() string {
	return f.name
}

// Weight returns the font's weight where 400 is normal.
func (f *Font) Weight() float32 {
	return f.weight
}

// GenerateGlyphs builds a set of textures based on a ttf files gylphs
// The font has a fixed size in points, found in f.size.
// For normal text, that is a value between 10 and 16. The size is in points based on 72 points pr inch.
// The actual dpi is found in the last parameter, and is typically much higher on modern screens.
// The size of the glyphs in physical pixels will be ca size*dpi/72
// (see truetype/face.go:206
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
		char := new(charInfo)

		gBnd, gAdv, ok := ttfFace.GlyphBounds(ch)
		if ok != true {
			slog.Error("ttf face glyphBounds error", "rune", int(ch))
			continue
		}

		gh := int32((gBnd.Max.Y - gBnd.Min.Y) >> 6)
		gw := int32((gBnd.Max.X - gBnd.Min.X) >> 6)

		// if gylph has no dimensions, set it to the max value
		if gw == 0 || gh == 0 {
			gBnd = f.ttf.Bounds(fixed.Int26_6(f.size))
			// Make sure sizes are at least 1
			gw = max(1, int32((gBnd.Max.X-gBnd.Min.X)>>6))
			gh = max(1, int32((gBnd.Max.Y-gBnd.Min.Y)>>6))
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
		char.bearingH = int(gBnd.Min.X) >> 6

		// create image to draw glyph
		fg, bg := image.White, image.Black
		rect := image.Rect(0, 0, int(gw), int(gh))
		rgba := image.NewRGBA(rect)
		draw.Draw(rgba, rgba.Bounds(), bg, image.Point{}, draw.Src)
		// set the glyph dot
		px := 0 - (int(gBnd.Min.X) >> 6)
		py := gAscent
		pt := freetype.Pt(px, py)
		// Draw the text from mask to image
		c.SetClip(rgba.Bounds())
		c.SetDst(rgba)
		c.SetSrc(fg)
		_, err := c.DrawString(string(ch), pt)
		if err != nil {
			return err
		}
		if gpu.DebugWidgets {
			if ch == 'E' {
				file, err := os.Create("./test-outputs/E-" + f.name + "-" + strconv.Itoa(int(dpi)) + ".png")
				if err != nil {
					slog.Error(err.Error())
				} else {
					_ = png.Encode(file, rgba)
					_ = file.Close()
				}
			}
		}
		// Generate texture
		char.TextureID = gpu.GenerateTexture(rgba)
		// add char to fontChar list
		f.FontChar[ch] = char
	}
	return nil
}

// LoadTrueTypeFont builds OpenGL buffers and glyph textures based on a ttf file
func LoadTrueTypeFont(name string, program uint32, r io.Reader, size int) (*Font, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	// Read the truetype font data from the given io.Reader r
	ttf, err := truetype.Parse(data)
	if err != nil {
		return nil, err
	}

	f := new(Font)
	f.FontChar = make(map[rune]*charInfo)
	f.ttf = ttf
	f.size = size
	f.Program = program
	f.name = name
	gpu.ConfigureVaoVbo(&f.Vao, &f.Vbo, f.Program)
	return f, nil
}
