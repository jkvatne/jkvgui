package font

import (
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
	dpi      float32
	weight   float32
	no       int
}

type charInfo struct {
	TextureID uint32 // ID handle of the glyph texture
	width     int    // glyph width
	height    int    // glyph height
	advance   int    // glyph advance
	bearingH  int    // glyph bearing horizontal
	bearingV  int    // glyph bearing vertical
}

// LoadDefaultFonts will load the default fonts from embedded data
// The user program can override these values by loading another font.
func LoadDefaultFonts() {
	LoadFontBytes(gpu.Normal14, "RobotoNormal", Roboto400, 14, 400)
	LoadFontBytes(gpu.Bold14, "RobotoBold", Roboto600, 14, 600)
	LoadFontBytes(gpu.Bold16, "RobotoBold", Roboto600, 16, 600)
	LoadFontBytes(gpu.Bold20, "RobotoBold", Roboto600, 20, 600)
	LoadFontBytes(gpu.Italic14, "RobotoItalic", RobotoItalic500, 14, 500)
	LoadFontBytes(gpu.Mono14, "RobotoMono", RobotoMono400, 14, 400)
	LoadFontBytes(gpu.Normal16, "RobotoNormal", Roboto400, 16, 400)
	LoadFontBytes(gpu.Normal12, "RobotoNormal", Roboto400, 12, 400)
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
		err := f.GenerateGlyphs(r, r)
		if err == nil {
			ch, ok = f.FontChar[r]
		}
	}
	// skip runes that are not in a font character range
	if !ok {
		slog.Error("Rune not found", "font", f.name, "index", r)
	}
	return ch
}

var done bool

// DrawText draws a string to the screen, takes a list of arguments like printf
// max is the maximum width. If longer, ellipsis is appended
// scale is the size relative to the default text size, typically between 0.7 and 2.5.
func (f *Font) DrawText(x, y float32, color f32.Color, maxW float32, dir gpu.Direction, fs string, argv ...interface{}) {
	runes := []rune(fmt.Sprintf(fs, argv...))
	if len(runes) == 0 {
		return
	}
	x *= gpu.ScaleX
	y *= gpu.ScaleY
	maxW *= gpu.ScaleX
	size := gpu.ScaleX * DefaultDpi / f.dpi
	if !done {
		done = true
		slog.Info("DrawText", "no", f.no, "name", f.name, "f.size", f.size, "size", size, "f.dpi", f.dpi, "ScaleX", gpu.ScaleX)
	}
	gpu.SetupAttributes(color, f.Vao, f.Program)
	ellipsis := assertRune(f, Ellipsis)
	ellipsisWidth := float32(ellipsis.width+1) * size
	var offset float32
	// Iterate through all characters in a string
	for i := range runes {
		ch := assertRune(f, runes[i])
		bearingH := float32(ch.bearingH) * size
		bearingV := float32(ch.bearingV) * size
		if maxW > 0 && offset+ellipsisWidth+float32(ch.width)*size+bearingH >= maxW && i < len(runes)-1 {
			ch = ellipsis
		}
		w := float32(ch.width) * size
		h := float32(ch.height) * size
		// calculate the position and size for the current rune
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
func (f *Font) Width(fs string, argv ...interface{}) float32 {
	var width float32
	indices := []rune(fmt.Sprintf(fs, argv...))
	if len(indices) == 0 {
		return 0
	}
	// Iterate through all characters in a string
	for i := range indices {
		ch := assertRune(f, indices[i])
		width += float32(ch.advance >> 6)
	}
	return width * DefaultDpi / f.dpi
}

// RuneNo will give the rune number at pixel position x from the start
func (f *Font) RuneNo(x float32, s string) int {
	runes := []rune(s)
	width := float32(0)
	x = x * f.dpi / DefaultDpi
	// Iterate through all characters in a string
	for i, r := range runes {
		// find rune in fontChar list
		ch, ok := f.FontChar[r]
		// skip runes that are not in a font character range
		if !ok {
			fmt.Printf("%c %d\n", r, i)
			continue
		}
		width += float32(ch.advance >> 6)
		if width >= x {
			return i
		}
	}
	return len(runes)
}

// Height returns the font height at the given size
func (f *Font) Height() float32 {
	return (f.Ascent + f.Descent) * DefaultDpi / f.dpi
}

// Baseline returns the offset from the top to the font's baseline (the dot)
func (f *Font) Baseline() float32 {
	return f.Ascent * DefaultDpi / f.dpi
}

// Split will split a long string into an array of shorter strings that will fit within maxWidth
func Split(s string, maxWidth float32, font *Font) []string {
	var width float32
	lines := make([]string, 0)
	words := strings.Split(s, " ")
	line := ""
	for _, word := range words {
		if word == "" {
			continue
		}
		width = font.Width(line + " " + word)
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
					if font.Width(word[0:j]) > maxWidth {
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

// GenerateGlyphs builds a set of textures based on the ttf file glyphs
// The font has a fixed size in points, found in f.size.
// For normal text, that is a value between 10 and 16. The size is in points based on 72 points pr inch.
// The actual dpi is found in the last parameter and is typically much higher on modern screens.
// The size of the glyphs in physical pixels will be ca size*dpi/72
// (see truetype/face.go:206
func (f *Font) GenerateGlyphs(low, high rune) error {
	// create a freetype context for drawing
	c := freetype.NewContext()
	c.SetDPI(float64(f.dpi))
	c.SetFont(f.ttf)
	c.SetFontSize(float64(f.size))
	c.SetHinting(font.HintingFull)

	// create a new face to measure glyph dimensions
	ttfFace := truetype.NewFace(f.ttf, &truetype.Options{
		Size:    float64(f.size),
		DPI:     float64(f.dpi),
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
				fmt.Printf("Font %d, %s, letter E w=%d h=%d dpi=%0.2f default dpi=%0.2f  scaleX=%0.3f f.size=%d\n", f.no, f.name, char.width, char.height, f.dpi, DefaultDpi, gpu.ScaleX, f.size)
				file, err := os.Create("./test-outputs/E-" + f.name + "-" + strconv.Itoa(int(f.dpi)) + ".png")
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

// LoadFontBytes builds OpenGL buffers and glyph textures based on a ttf data array
func LoadFontBytes(no int, name string, data []byte, size int, weight float32) {
	ttf, err := truetype.Parse(data)
	f32.ExitOn(err, "Parsing font data failed")
	f := new(Font)
	f.FontChar = make(map[rune]*charInfo)
	f.ttf = ttf
	f.dpi = 72 * gpu.ScaleX
	f.size = size
	program, _ := gpu.NewProgram(gpu.VertQuadSource, gpu.FragQuadSource)
	f.Program = program
	f.name = name
	f.no = no
	f.weight = weight
	_ = f.GenerateGlyphs(0x20, 0x7E)
	gpu.ConfigureVaoVbo(&f.Vao, &f.Vbo, f.Program, "LoadFontBytes")
	Fonts[no] = f
}

// LoadFontFile loads the specified font at the given size (in pixels).
// The integer returned is the index to Fonts[]
// Will panic if font is not found
func LoadFontFile(no int, file string, size int, name string, weight float32) {
	f32.ExitIf(no < 0 || no > len(Fonts), "LoadFontFile: invalid index "+strconv.Itoa(no)+", must be between 0 and 31 ")
	fd, err := os.Open(file)
	f32.ExitOn(err, "Failed to open font file "+file)
	defer func(fd *os.File) {
		f32.ExitOn(fd.Close(), "Could not close file: "+file)
	}(fd)
	data := make([]byte, size)
	_, err = fd.Read(data)
	f32.ExitOn(err, "Failed to read font file "+file)
	LoadFontBytes(no, name, data, size, weight)
}
