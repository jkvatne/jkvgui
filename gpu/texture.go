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

func SetupDrawing(color f32.Color, vao uint32, program uint32) {
	// Activate corresponding render state
	gl.UseProgram(program)
	// setup blending mode
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	// set text color
	gl.Uniform4f(gl.GetUniformLocation(program, gl.Str("textColor\x00")), color.R, color.G, color.B, color.A)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindVertexArray(vao)
	// set screen resolution
	gl.Viewport(0, 0, int32(WindowWidthPx), int32(WindowHeightPx))
	resUniform := gl.GetUniformLocation(program, gl.Str("resolution\x00"))
	gl.Uniform2f(resUniform, float32(WindowWidthPx), float32(WindowHeightPx))
	GetErrors("SetupDrawing")
}

func RenderTexture(x, y, w, h float32, texture uint32, vbo uint32, dir Direction) {
	// Render texture over quad
	gl.BindTexture(gl.TEXTURE_2D, texture)
	// Update content of VBO memory
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
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
	// Release buffer
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	GetErrors("RenderTexture")
}

// ConfigureVaoVbo for texture quads
func ConfigureVaoVbo(vao *uint32, vbo *uint32, program uint32, from string) {
	gl.GenVertexArrays(1, vao)
	gl.BindVertexArray(*vao)
	GetErrors("CfgVabVbo1 " + from)
	gl.GenBuffers(1, vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, *vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 6*4*4, nil, gl.STATIC_DRAW)
	GetErrors("CfgVabVbo2" + from)
	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	GetErrors("CfgVabVbo3 " + from)
	gl.EnableVertexAttribArray(vertAttrib)
	GetErrors("CfgVabVbo4 " + from)
	gl.VertexAttribPointerWithOffset(vertAttrib, 2, gl.FLOAT, false, 4*4, 0)
	GetErrors("CfgVabVbo5 " + from)
	defer gl.DisableVertexAttribArray(vertAttrib)
	texCoordAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertTexCoord\x00")))
	GetErrors("CfgVabVbo6 " + from)
	gl.EnableVertexAttribArray(texCoordAttrib)
	GetErrors("CfgVabVbo7 " + from)
	gl.VertexAttribPointerWithOffset(texCoordAttrib, 2, gl.FLOAT, false, 4*4, 2*4)
	GetErrors("CfgVabVbo8 " + from)
	defer gl.DisableVertexAttribArray(texCoordAttrib)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	GetErrors("CfgVabVbo9 " + from)
	gl.BindVertexArray(0)
	GetErrors("CfgVabVbo10 " + from)
}

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
