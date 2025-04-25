package gpu

import (
	"github.com/go-gl/gl/all-core/gl"
	"github.com/jkvatne/jkvgui/f32"
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

func SetTextColor(color f32.Color, program uint32) {
	gl.Uniform4f(gl.GetUniformLocation(program, gl.Str("textColor\x00")), color.R, color.G, color.B, color.A)
}

// SetupTexture
func SetupTexture(vao uint32, vbo uint32, program uint32) {
	// Activate corresponding render state
	gl.UseProgram(program)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	GetErrors("SetupTexture")
}

// RenderTexture will draw the texture given onto the frame buffer at given location and rotation.
func RenderTexture(x, y, w, h float32, texture uint32, vbo uint32, dir Direction) {
	// Render texture over quad
	gl.BindTexture(gl.TEXTURE_2D, texture)
	GetErrors("RenderTexture")
	// Update content of VBO memory
	if dir == TTB {
		vertices := []float32{
			x + w, y + h, 1.0, 0.0,
			x + w, y, 0.0, 0.0,
			x, y, 0.0, 1.0,

			x, y, 0.0, 1.0,
			x, y + h, 1.0, 1.0,
			x + w, y + h, 1.0, 0.0,
		}
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(vertices)*4, gl.Ptr(vertices)) // Be sure to use glBufferSubData and not glBufferData
	} else if dir == BTT {
		vertices := []float32{
			x, y, 1.0, 0.0,
			x, y + h, 0.0, 0.0,
			x + w, y + h, 0.0, 1.0,

			x + w, y + h, 0.0, 1.0,
			x + w, y, 1.0, 1.0,
			x, y, 1.0, 0.0,
		}
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(vertices)*4, gl.Ptr(vertices)) // Be sure to use glBufferSubData and not glBufferData
	} else if dir == LTR {
		vertices := []float32{
			x + w, y, 1.0, 0.0,
			x, y, 0.0, 0.0,
			x, y + h, 0.0, 1.0,

			x, y + h, 0.0, 1.0,
			x + w, y + h, 1.0, 1.0,
			x + w, y, 1.0, 0.0,
		}
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(vertices)*4, gl.Ptr(vertices)) // Be sure to use glBufferSubData and not glBufferData

	}
	// Render quad
	gl.DrawArrays(gl.TRIANGLES, 0, 16)
	GetErrors("RenderTexture")
	// Release buffer
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	GetErrors("RenderTexture")
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
	GetErrors("GenTexture")
	return texture
}
