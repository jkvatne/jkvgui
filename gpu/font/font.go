package font

import (
	_ "embed"
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/gpu/font/freetype"
	"github.com/jkvatne/jkvgui/gpu/font/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var DebugFonts = flag.Bool("debugfonts", false, "Set to write font info to file")

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

// A Font allows rendering of a text.
type Font struct {
	FontChar     map[rune]*charInfo
	ttf          *truetype.Font
	Texture      uint32 // Holds the glyph texture id.
	color        f32.Color
	ascent       float32
	descent      float32
	Name         string
	Size         int
	dpi          float32
	Weight       float32
	No           int
	maxCharWidth int
	Height       float32
	Baseline     float32
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
	t := time.Now()
	LoadFontBytes(gpu.Normal14, "RobotoNormal", Roboto400, 14, 400)
	LoadFontBytes(gpu.Bold14, "RobotoBold", Roboto600, 14, 600)
	LoadFontBytes(gpu.Italic14, "RobotoItalic", RobotoItalic500, 14, 500)
	LoadFontBytes(gpu.Mono14, "RobotoMono", RobotoMono400, 14, 400)
	LoadFontBytes(gpu.Normal12, "RobotoNormal", Roboto400, 12, 400)
	LoadFontBytes(gpu.Bold12, "RobotoBold", Roboto600, 12, 600)
	LoadFontBytes(gpu.Italic12, "RobotoItalic", RobotoItalic500, 12, 500)
	LoadFontBytes(gpu.Mono12, "RobotoMono", RobotoMono400, 12, 400)
	LoadFontBytes(gpu.Normal10, "RobotoNormal", Roboto400, 10, 400)
	LoadFontBytes(gpu.Bold10, "RobotoBold", Roboto600, 10, 600)
	LoadFontBytes(gpu.Italic10, "RobotoItalic", RobotoItalic500, 10, 500)
	LoadFontBytes(gpu.Mono10, "RobotoMono", RobotoMono400, 10, 400)
	LoadFontBytes(gpu.Normal16, "RobotoNormal", Roboto400, 16, 400)
	LoadFontBytes(gpu.Bold16, "RobotoBold", Roboto600, 16, 600)
	LoadFontBytes(gpu.Normal20, "RobotoNormal", Roboto400, 20, 400)
	LoadFontBytes(gpu.Bold20, "RobotoBold", Roboto600, 20, 600)
	slog.Debug("LoadDefaultFonts()", "time", time.Since(t))
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
		slog.Error("Rune not found", "font", f.Name, "index", r)
	}
	return ch
}

// DrawText draws a string to the screen, takes a list of arguments like printf
// max is the maximum width. If longer, ellipsis is appended
// scale is the size relative to the default text size, typically between 0.7 and 2.5.
func (f *Font) DrawText(x, y float32, color f32.Color, maxW float32, dir gpu.Direction, str string) {
	runes := []rune(str)
	if len(runes) == 0 {
		return
	}
	f32.ExitIf(f == nil, "Font is nil")
	x *= gpu.ScaleX
	y *= gpu.ScaleY
	maxW *= gpu.ScaleX
	size := gpu.ScaleX * DefaultDpi / f.dpi
	gpu.SetupTexture(color, gpu.FontVao, gpu.FontVbo, gpu.FontProgram)
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
			gpu.RenderTexture(xPos, yPos, w, h, ch.TextureID, gpu.FontVbo, dir)
		} else if dir == gpu.TTB {
			xPos := x - bearingV
			yPos := y + offset + bearingH
			gpu.RenderTexture(xPos, yPos, h, w, ch.TextureID, gpu.FontVbo, dir)
		} else if dir == gpu.BTT {
			xPos := x - h + bearingV
			yPos := y - offset - w
			gpu.RenderTexture(xPos, yPos, h, w, ch.TextureID, gpu.FontVbo, dir)
		}
		offset += float32(ch.advance>>6) * size
		if ch == ellipsis {
			break
		}
	}
}

// Width returns the width of a piece of text in pixels
func (f *Font) Width(str string) float32 {
	var width int
	indices := []rune(str)
	if len(indices) == 0 {
		return 0
	}
	// Iterate through all characters in a string
	for i := range indices {
		ch := assertRune(f, indices[i])
		width += ch.advance
	}
	return float32(width) * DefaultDpi / f.dpi / 64
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

// Split will split a long string into an array of shorter strings that will fit within maxWidth
func Split(str string, maxWidth float32, font *Font) []string {
	maxW := int(maxWidth / DefaultDpi * font.dpi)
	if maxWidth == 0 || maxW > len(str)*font.maxCharWidth {
		return []string{str}
	}
	runes := []rune(str)
	// w is the running total for the width of the current line.
	w := 0
	lastSpace := 0
	start := 0
	lastW := 0
	maxW = maxW * 64
	var lines []string
	for i, r := range runes {
		ch := assertRune(font, r)
		adv := ch.advance
		if r == 32 {
			if w == 0 {
				// Skip leading whitespace
				continue
			}
			// Save position of last whitespace
			lastSpace = i
			lastW = w
			if w+adv >= maxW {
				// We have a space and will break on it
				lines = append(lines, string(runes[start:i]))
				start = i + 1
				w = 0
			} else {
				w += adv
			}

		} else {
			// Accumulate current line length
			w += adv
			if w > maxW {
				if lastSpace <= start {
					// No spaces, split within the current word
					lines = append(lines, string(runes[start:i]))
					start = i
					w = 0
				} else {
					lines = append(lines, string(runes[start:lastSpace]))
					start = lastSpace + 1
					w -= lastW
				}
			}
		}
	}
	lines = append(lines, string(runes[start:]))
	return lines
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
	c.SetFontSize(float64(f.Size))
	c.SetHinting(font.HintingFull)

	// create a new face to measure glyph dimensions
	ttfFace := truetype.NewFace(f.ttf, &truetype.Options{
		Size:    float64(f.Size),
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
			gBnd = f.ttf.Bounds(fixed.Int26_6(f.Size))
			// Make sure sizes are at least 1
			gw = max(1, int32((gBnd.Max.X-gBnd.Min.X)>>6))
			gh = max(1, int32((gBnd.Max.Y-gBnd.Min.Y)>>6))
		}
		// The glyph's ascent and descent equal -bounds.Min.Y and +bounds.Max.Y.
		gAscent := int(-gBnd.Min.Y) >> 6
		gDescent := int(gBnd.Max.Y) >> 6
		f.ascent = max(f.ascent, float32(gAscent))
		f.descent = max(f.descent, float32(gDescent))
		if f.dpi == 0 {
			panic("Font's dpi is zero")
		}
		f.Height = (f.ascent + f.descent) * DefaultDpi / f.dpi
		f.Baseline = f.ascent * DefaultDpi / f.dpi
		// set w,h and adv, bearing V and bearing H in char
		char.width = int(gw)
		char.height = int(gh)
		char.advance = int(gAdv)
		if char.advance > f.maxCharWidth {
			f.maxCharWidth = char.advance >> 6
		}
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
		if *DebugFonts {
			if ch == 'E' {
				slog.Info("Writing debug info to ./test-outputs")
				slog.Info("Letter E", "w", char.width, "h", char.height, "dpi", f.dpi, "default dpi", DefaultDpi, "scaleX", gpu.ScaleX, "f.size", f.Size)
				f32.AssertDir("test-outputs")
				file, err := os.Create("test-outputs/E-" + f.Name + "-" + strconv.Itoa(int(f.dpi)) + ".png")
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
	f.Size = size
	f.Name = name
	f.No = no
	f.Weight = weight
	_ = f.GenerateGlyphs(0x20, 0x7E)
	_ = f.GenerateGlyphs(197, 198)
	_ = f.GenerateGlyphs(216, 216)
	_ = f.GenerateGlyphs(229, 230)
	_ = f.GenerateGlyphs(248, 248)
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
