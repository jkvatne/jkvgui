package gpu

import (
	"image"

	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/purego-glfw/gl"
)

// Direction represents the direction in which strings should be rendered.
type Direction uint8

const (
	// LTR for normal display or left to right text,
	LTR Direction = iota
	// TTB Top-To-Bottom, rotates the image 90 degrees for top-to-bottom text
	TTB
	// INV Invert, rotates by 180 degrees
	INV
	// BTT Bottom-to-Top text, rotates by 270 degrees
	BTT
	// RTL Right-To-Left will mirror image the image without rotation
	RTL
	// MTB Mirror image and rotate by 90 degrees
	MTB
	// MIV   Mirror image and rotate by 180 degrees
	MIV
	// MBT   Mirror image and rotate by 270 degrees
	MBT
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
	GetErrors("SetupTexture")
}

func set(vertices *[24]float32, x0, y0, x1, y1, x2, y2, x3, y3 float32) {
	vertices[4] = x0
	vertices[5] = y0
	vertices[0] = x1
	vertices[1] = y1
	vertices[20] = x1
	vertices[21] = y1
	vertices[16] = x2
	vertices[17] = y2
	vertices[8] = x3
	vertices[9] = y3
	vertices[12] = x3
	vertices[13] = y3
}

// RenderTexture will draw the texture given onto the frame buffer at given location and rotation.
// The direction parameter gives the rotation and mirroring.
func RenderTexture(x, y, w, h float32, texture uint32, direction Direction) {
	// vertices has the texture coordinates identical for all quads, in 2,3, 6,7 etc
	var vertices = [24]float32{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0, 1, 0}
	// Render texture over quad
	gl.BindTexture(gl.TEXTURE_2D, texture)
	if direction == LTR {
		set(&vertices, x, y, x+w, y, x+w, y+h, x, y+h)
	} else if direction == TTB {
		set(&vertices, x+w, y, x+w, y+h, x, y+h, x, y)
	} else if direction == INV {
		set(&vertices, x+w, y+h, x, y+h, x, y, x+w, y)
	} else if direction == BTT {
		set(&vertices, x, y+h, x, y, x+w, y, x+w, y+h)
	} else if direction == RTL {
		set(&vertices, x+w, y, x, y, x, y+h, x+w, y+h)
	} else if direction == MTB {
		set(&vertices, x, y, x, y+h, x+w, y+h, x+w, y)
	} else if direction == MIV {
		set(&vertices, x, y+h, x+w, y+h, x+w, y, x, y)
	} else if direction == MBT {
		set(&vertices, x+w, y+h, x+w, y, x, y, x, y+h)
	}
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(vertices)*4, gl.Ptr(&vertices[0])) // Be sure to use glBufferSubData and not glBufferData
	// Render quad
	gl.DrawArrays(gl.TRIANGLES, 0, 16)
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
	GetErrors("GenerateTexture")
	return texture
}
