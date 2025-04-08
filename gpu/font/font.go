package font

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
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

var Dpi float32 = 108
var DefaultFontSize = 14

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

// DrawText draws a string to the screen, takes a list of arguments like printf
// max is the maximum width. If longer, ellipsis is appended
// scale is the size relative to the default text size.
func (f *Font) DrawText(x, y float32, color f32.Color, scale float32, maxW float32, dir gpu.Direction, fs string, argv ...interface{}) {
	indices := []rune(fmt.Sprintf(fs, argv...))
	if len(indices) == 0 {
		return
	}
	x *= gpu.ScaleX
	y *= gpu.ScaleY
	maxW *= gpu.ScaleX
	size := gpu.ScaleX * scale * 72 / Dpi
	gpu.SetupDrawing(color, f.Vao, f.Program)

	ellipsis, ok := f.FontChar[Ellipsis]
	if !ok {
		_ = f.GenerateGlyphs(Ellipsis, Ellipsis, Dpi)
		ellipsis, ok = f.FontChar[Ellipsis]
	}
	ellipsisWidth := float32(ellipsis.width+1) * size
	var offset float32
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
		bearingH := float32(ch.bearingH) * size
		bearingV := float32(ch.bearingV) * size
		w := float32(ch.width) * size
		h := float32(ch.height) * size
		if maxW > 0 && offset+ellipsisWidth+w+bearingH >= maxW && i < len(indices)-1 {
			ch = ellipsis
			w = float32(ch.width) * size
			h = float32(ch.height) * size
		}

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
	f, err := LoadTrueTypeFont(name, program, fd, size, 32, 127)
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
	f, err := LoadTrueTypeFont(name, program, fd, size, 32, 127)
	f32.ExitOn(err, "Could not load font bytes")
	f.name = name
	f.weight = weight
	f.size = size
	Fonts[no] = f
}

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
