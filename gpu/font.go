package gpu

import (
	"bytes"
	"fmt"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/jkvatne/jkvgui/lib"
	"github.com/jkvatne/jkvgui/shader"
	"os"
	"strings"
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

// Use default preapration for exported functions like `LoadFont` and `LoadFontFromBytes`
func configureDefaults(windowWidth int, windowHeight int) uint32 {
	// Configure the default font vertex and fragment shaders
	program, err := shader.NewProgram(shader.VertexFontShader, shader.FragmentFontShader)
	if err != nil {
		panic(err)
	}
	// Activate corresponding render state
	gl.UseProgram(program)
	// set screen resolution
	resUniform := gl.GetUniformLocation(program, gl.Str("resolution\x00"))
	gl.Uniform2f(resUniform, float32(windowWidth), float32(windowHeight))
	return program
}

// LoadFontBytes loads the specified font bytes at the given scale.
func LoadFontBytes(buf []byte, scale int32, windowWidth int, windowHeight int) (*Font, error) {
	program := configureDefaults(windowWidth, windowHeight)
	fd := bytes.NewReader(buf)
	return LoadTrueTypeFont(program, fd, scale, 32, 127, LeftToRight)
}

// SetColor allows you to set the text color to be used when you draw the text
func (f *Font) SetColor(red float32, green float32, blue float32, alpha float32) {
	f.color.R = red
	f.color.G = green
	f.color.B = blue
	f.color.A = alpha
}

// UpdateResolution used to recalibrate fonts for new window size
func (f *Font) UpdateResolution(windowWidth int, windowHeight int) {
	gl.UseProgram(f.program)
	resUniform := gl.GetUniformLocation(f.program, gl.Str("resolution\x00"))
	gl.Uniform2f(resUniform, float32(windowWidth), float32(windowHeight))
	gl.UseProgram(0)
}

// Printf draws a string to the screen, takes a list of arguments like printf
func (f *Font) Printf(x, y float32, scale float32, fs string, argv ...interface{}) {
	indices := []rune(fmt.Sprintf(fs, argv...))
	if len(indices) == 0 {
		return
	}
	// setup blending mode
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	// Activate corresponding render state
	gl.UseProgram(f.program)
	// set text color
	gl.Uniform4f(gl.GetUniformLocation(f.program, gl.Str("textColor\x00")), f.color.R, f.color.G, f.color.B, f.color.A)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindVertexArray(f.vao)
	// Iterate through all characters in string
	for i := range indices {
		// get rune
		runeIndex := indices[i]
		// find rune in fontChar list
		ch, ok := f.fontChar[runeIndex]
		// load missing runes in batches of 32
		if !ok {
			low := runeIndex - (runeIndex % 32)
			_ = f.GenerateGlyphs(low, low+31)
			ch, ok = f.fontChar[runeIndex]
		}
		// skip runes that are not in font chacter range
		if !ok {
			fmt.Printf("%c %d\n", runeIndex, runeIndex)
			continue
		}

		// calculate position and size for current rune
		xpos := x + float32(ch.bearingH)*scale
		ypos := y - float32(ch.height-ch.bearingV)*scale
		w := float32(ch.width) * scale
		h := float32(ch.height) * scale
		vertices := []float32{
			xpos + w, ypos, 1.0, 0.0,
			xpos, ypos, 0.0, 0.0,
			xpos, ypos + h, 0.0, 1.0,

			xpos, ypos + h, 0.0, 1.0,
			xpos + w, ypos + h, 1.0, 1.0,
			xpos + w, ypos, 1.0, 0.0,
		}

		// Render glyph texture over quad
		gl.BindTexture(gl.TEXTURE_2D, ch.textureID)
		// Update content of VBO memory
		gl.BindBuffer(gl.ARRAY_BUFFER, f.vbo)

		// BufferSubData(target Enum, offset int, data []byte)
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(vertices)*4, gl.Ptr(vertices)) // Be sure to use glBufferSubData and not glBufferData
		// Render quad
		gl.DrawArrays(gl.TRIANGLES, 0, 16)

		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		// Now advance cursors for next glyph (note that advance is number of 1/64 pixels)
		x += float32((ch.advance >> 6)) * scale // Bitshift by 6 to get value in pixels (2^6 = 64 (divide amount of 1/64th pixels by 64 to get amount of pixels))

	}

	// clear opengl textures and programs
	gl.BindVertexArray(0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	gl.UseProgram(0)
	gl.Disable(gl.BLEND)
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
		ch, ok := f.fontChar[runeIndex]
		// load missing runes in batches of 32
		if !ok {
			low := runeIndex & rune(32-1)
			_ = f.GenerateGlyphs(low, low+31)
			ch, ok = f.fontChar[runeIndex]
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

// LoadFont loads the specified font at the given scale.
func LoadFontFile(file string, scale int32, windowWidth int, windowHeight int) (*Font, error) {
	fd, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	program := configureDefaults(windowWidth, windowHeight)
	return LoadTrueTypeFont(program, fd, scale, 32, 127, LeftToRight)
}

func LoadFont(name string, scale int32) {
	var f *Font
	var err error
	if strings.EqualFold(name, "Roboto-Medium") {
		f, err = LoadFontBytes(RobotoMedium, scale, WindowWidth, WindowHeight)
	} else if strings.EqualFold(name, "Roboto") {
		f, err = LoadFontBytes(RobotoMedium, scale, WindowWidth, WindowHeight)
	} else if strings.EqualFold(name, "Roboto-Light") {
		f, err = LoadFontBytes(RobotoLight, scale, WindowWidth, WindowHeight)
	} else if strings.EqualFold(name, "Roboto-Regular") {
		f, err = LoadFontBytes(RobotoRegular, scale, WindowWidth, WindowHeight)
	} else if strings.EqualFold(name, "RobotoMono") {
		f, err = LoadFontBytes(RobotoMono, scale, WindowWidth, WindowHeight)
	} else {
		f, err = LoadFontFile(name, scale, WindowWidth, WindowHeight)
	}
	f.SetColor(0.0, 0.0, 0.0, 1.0)
	lib.PanicOn(err, "Loading "+name)
	Fonts = append(Fonts, f)
}
