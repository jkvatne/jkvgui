package gpu

import (
	"image"

	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gl"
)

// Direction represents the direction in which strings should be rendered.
type Direction uint8

const (
	LTR Direction = iota
	RTL
	TTB
	BTT
)

// SetupTexture will set up vao for the program
func SetupTexture(color f32.Color, vao uint32, vbo uint32, program uint32) {
	// Activate corresponding render state
	gl.UseProgram(program)
	// setup blending mode
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	// set text color
	gl.Uniform4f(gl.GetUniformLocation(program, gl.Str("textColor\x00")), color.R, color.G, color.B, color.A)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
}

// var vertices = []float32{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0, 1, 0}

// RenderTexture will draw the texture given onto the frame buffer at given location and rotation.
func RenderTexture(x, y, w, h float32, texture uint32, vbo uint32, dir Direction) {
	var vertices []float32
	// Render texture over quad
	gl.BindTexture(gl.TEXTURE_2D, texture)
	if dir == TTB {
		vertices[0] = x + w
		vertices[1] = y + h
		vertices[4] = x + w
		vertices[5] = y
		vertices[8] = x
		vertices[9] = y
		vertices[12] = x
		vertices[13] = y
		vertices[16] = x
		vertices[17] = y + h
		vertices[20] = x + w
		vertices[21] = y + h
	} else if dir == BTT {
		vertices[0] = x
		vertices[1] = y
		vertices[4] = x
		vertices[5] = y + h
		vertices[8] = x + w
		vertices[9] = y + h
		vertices[12] = x + w
		vertices[13] = y + h
		vertices[16] = x + w
		vertices[17] = y
		vertices[20] = x
		vertices[21] = y
	} else if dir == LTR {
		vertices[0] = x + w
		vertices[1] = y
		vertices[4] = x
		vertices[5] = y
		vertices[8] = x
		vertices[9] = y + h
		vertices[12] = x
		vertices[13] = y + h
		vertices[16] = x + w
		vertices[17] = y + h
		vertices[20] = x + w
		vertices[21] = y
	}
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(vertices)*4, gl.Ptr(vertices)) // Be sure to use glBufferSubData and not glBufferData
	// Render quad
	gl.DrawArrays(gl.TRIANGLES, 0, 16)
}

// GenerateTexture will bind a rgba image to a texture and return its "name"
func GenerateTexture(rgba *image.RGBA) uint32 {
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA,
		int32(rgba.Rect.Dx()), int32(rgba.Rect.Dy()), 0,
		gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))
	gl.BindTexture(gl.TEXTURE_2D, 0)
	return texture
}
