package gpu

import (
	"encoding/binary"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gl"
	"image"
)

// Direction represents the direction in which strings should be rendered.
type Direction uint8

const (
	LTR Direction = iota
	RTL
	TTB
	BTT
)

func SetupDrawing(color f32.Color, vao *gl.Buffer, program gl.Program) {
	GetErrors()
	// Activate corresponding render state
	gl.UseProgram(program)
	GetErrors()
	// setup blending mode
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	GetErrors()
	// set text color
	gl.Uniform4f(gl.GetUniformLocation(program, "textColor"), color.R, color.G, color.B, color.A)
	GetErrors()
	gl.ActiveTexture(gl.TEXTURE0)
	GetErrors()
	// gl.BindVertexArray(vao)
	// set screen resolution
	gl.Viewport(0, 0, int(WindowWidthPx), int(WindowHeightPx))
	GetErrors()
	resUniform := gl.GetUniformLocation(program, "resolution")
	GetErrors()
	gl.Uniform2f(resUniform, float32(WindowWidthPx), float32(WindowHeightPx))
	GetErrors()
}

func RenderTexture(x, y, w, h float32, texture gl.Texture, vbo gl.Buffer, dir Direction) {
	// Render texture over quad
	gl.BindTexture(gl.TEXTURE_2D, texture)
	GetErrors()
	// Update content of VBO memory
	vbo2 := gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo2)
	GetErrors()
	if dir == TTB {
		vertices := []float32{
			x + w, y + h, 1.0, 0.0,
			x + w, y, 0.0, 0.0,
			x, y, 0.0, 1.0,

			x, y, 0.0, 1.0,
			x, y + h, 1.0, 1.0,
			x + w, y + h, 1.0, 0.0,
		}
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, Bytes(binary.LittleEndian, vertices...))
		GetErrors()
	} else if dir == BTT {
		vertices := []float32{
			x, y, 1.0, 0.0,
			x, y + h, 0.0, 0.0,
			x + w, y + h, 0.0, 1.0,

			x + w, y + h, 0.0, 1.0,
			x + w, y, 1.0, 1.0,
			x, y, 1.0, 0.0,
		}
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, Bytes(binary.LittleEndian, vertices...))
		GetErrors()
	} else if dir == LTR {
		vertices := []float32{
			x + w, y, 1.0, 0.0,
			x, y, 0.0, 0.0,
			x, y + h, 0.0, 1.0,

			x, y + h, 0.0, 1.0,
			x + w, y + h, 1.0, 1.0,
			x + w, y, 1.0, 0.0,
		}
		gl.BufferData(gl.ARRAY_BUFFER, Bytes(binary.LittleEndian, vertices...), gl.STATIC_DRAW)
		GetErrors()

	}
	// Render quad
	gl.DrawArrays(gl.TRIANGLES, 0, 16)
	GetErrors()
	// Release buffer
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	GetErrors()
}

// SetupBuffers for texture quads
func SetupBuffers(vao *gl.Buffer, vbo *gl.Buffer, program gl.Program) {
	*vbo = gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, *vbo)
	// gl.BufferData(gl.ARRAY_BUFFER, , gl.STATIC_DRAW)
	vertAttrib := gl.GetAttribLocation(program, "vert")
	gl.EnableVertexAttribArray(vertAttrib)
	// gl.VertexAttribPointerWithOffset(vertAttrib, 2, gl.FLOAT, false, 4*4, 0)
	defer gl.DisableVertexAttribArray(vertAttrib)

	texCoordAttrib := gl.GetAttribLocation(program, "vertTexCoord")
	gl.EnableVertexAttribArray(texCoordAttrib)
	// gl.VertexAttribPointerWithOffset(texCoordAttrib, 2, gl.FLOAT, false, 4*4, 2*4)
	defer gl.DisableVertexAttribArray(texCoordAttrib)
	GetErrors()
}

func GenerateTexture(rgba *image.RGBA) gl.Texture {
	texture := gl.CreateTexture()
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexImage2D(gl.TEXTURE_2D, 0, rgba.Rect.Dx(), rgba.Rect.Dy(),
		gl.RGBA, gl.UNSIGNED_BYTE, rgba.Pix)
	GetErrors()
	return texture
}
