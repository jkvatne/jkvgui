package gpu

import (
	"encoding/binary"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gl"
	"image"
	"unsafe"
)

// Direction represents the direction in which strings should be rendered.
type Direction uint8

const (
	LTR Direction = iota
	RTL
	TTB
	BTT
)

// SetupAttributes
func SetupAttributes(color f32.Color, vao uint32, program uint32) {
	// Activate corresponding render state
	// fmt.Printf("Program handle: %d\n", program)
	ctx.UseProgram(program)
	GetErrors("SetupAttributes1")
	// setup blending mode
	ctx.Enable(gl.BLEND)
	GetErrors("SetupAttributes2")
	// ctx.BlendFuncSeparate(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA, gl.)
	// set text color
	c := []float32{color.R, color.G, color.B, color.A}
	l := ctx.GetUniformLocation(program, "textColor\x00")
	GetErrors("SetupAttributes3")
	ctx.Uniform4fv(l, c)
	ctx.ActiveTexture(gl.TEXTURE0)
	GetErrors("SetupAttributes4")
	ctx.BindVertexArray(vao)
	GetErrors("SetupAttributes5")
	// set screen resolution
	resUniform := ctx.GetUniformLocation(program, "resolution\x00")
	w := []float32{float32(WindowWidthPx), float32(WindowHeightPx)}
	ctx.Uniform2fv(resUniform, w)
	GetErrors("SetupAttributes6")
}

/*
// RenderTexture will draw the texture given onto the frame buffer at given location and rotation.
func RenderTexture(x, y, w, h float32, texture uint32, vbo uint32, dir Direction) {
	// Render texture over quad
	ctx.BindTexture(gl.TEXTURE_2D, texture)
	GetErrors("RenderTexture1")

	// Update content of VBO memory
	ctx.BindBuffer(gl.ARRAY_BUFFER, vbo)
	GetErrors("RenderTexture2")

	if dir == TTB {
		vertices := []float32{
			x + w, y + h, 1.0, 0.0,
			x + w, y, 0.0, 0.0,
			x, y, 0.0, 1.0,

			x, y, 0.0, 1.0,
			x, y + h, 1.0, 1.0,
			x + w, y + h, 1.0, 0.0,
		}
		b := unsafe.Slice((*byte)(unsafe.Pointer(&vertices[0])), len(vertices)*int(unsafe.Sizeof(vertices[0])))
		ctx.BufferSubData(gl.ARRAY_BUFFER, 0, b)

	} else if dir == BTT {
		vertices := []float32{
			x, y, 1.0, 0.0,
			x, y + h, 0.0, 0.0,
			x + w, y + h, 0.0, 1.0,

			x + w, y + h, 0.0, 1.0,
			x + w, y, 1.0, 1.0,
			x, y, 1.0, 0.0,
		}
		b := unsafe.Slice((*byte)(unsafe.Pointer(&vertices[0])), len(vertices)*int(unsafe.Sizeof(vertices[0])))
		ctx.BufferSubData(gl.ARRAY_BUFFER, 0, b)

	} else if dir == LTR {
		vertices := []float32{
			x + w, y, 1.0, 0.0,
			x, y, 0.0, 0.0,
			x, y + h, 0.0, 1.0,

			x, y + h, 0.0, 1.0,
			x + w, y + h, 1.0, 1.0,
			x + w, y, 1.0, 0.0,
		}
		b := unsafe.Slice((*byte)(unsafe.Pointer(&vertices[0])), len(vertices)*int(unsafe.Sizeof(vertices[0])))
		ctx.BufferSubData(gl.ARRAY_BUFFER, 0, b)
		GetErrors("RenderTexture3")

	}
	// Render quad
	ctx.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, 0)
	GetErrors("RenderTexture4")
	// Release buffer
	ctx.BindBuffer(gl.ARRAY_BUFFER, 0)
	GetErrors("RenderTexture5")
} */

func RenderTexture(x, y, w, h float32, texture uint32, vao uint32, vbo uint32, program uint32, dir Direction) {
	// Clear color to see if we're rendering at all
	// ctx.UseProgram(program)
	vertices := []float32{
		x, y, 0.0, 0.0, // bottom-left
		x + w, y, 1.0, 0.0, // bottom-right
		x, y + h, 0.0, 1.0, // top-left
		x + w, y + h, 1.0, 1.0, // top-right
	}

	b := unsafe.Slice((*byte)(unsafe.Pointer(&vertices[0])), len(vertices)*4)

	ctx.BindVertexArray(vao) // Bind VAO first
	ctx.BindBuffer(gl.ARRAY_BUFFER, vbo)
	ctx.BufferSubData(gl.ARRAY_BUFFER, 0, b)

	ctx.BindTexture(gl.TEXTURE_2D, texture)
	GetErrors("RenderTexture1")

	ctx.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, 0)
	GetErrors("RenderTexture2")

	ctx.BindTexture(gl.TEXTURE_2D, 0)
	ctx.BindVertexArray(0) // Unbind VAO
}

// ConfigureVaoVbo for texture quads
func ConfigureVaoVbo(vao *uint32, vbo *uint32, program uint32, from string) {
	*vao = ctx.CreateVertexArray()
	ctx.BindVertexArray(*vao)

	// Create and initialize vertex buffer
	*vbo = ctx.CreateBuffer()
	ctx.BindBuffer(gl.ARRAY_BUFFER, *vbo)
	vertexSize := 16 * 4 // 4 vertices * 4 components * 4 bytes per float
	ctx.BufferInit(gl.ARRAY_BUFFER, vertexSize, gl.DYNAMIC_DRAW)
	GetErrors("CfgVabVbo1")

	// Create and initialize index buffer
	ebo := ctx.CreateBuffer()
	ctx.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	indexSize := 6 * 2 // 6 indices * 2 bytes per uint16
	ctx.BufferInit(gl.ELEMENT_ARRAY_BUFFER, indexSize, gl.DYNAMIC_DRAW)

	// Upload index data
	indices := []uint16{0, 1, 2, 1, 3, 2}
	indexBytes := make([]byte, len(indices)*2)
	for i, idx := range indices {
		binary.LittleEndian.PutUint16(indexBytes[i*2:], idx)
	}
	ctx.BufferSubData(gl.ELEMENT_ARRAY_BUFFER, 0, indexBytes)
	GetErrors("CfgVabVbo2")

	// Set up vertex attributes
	ctx.BindAttribLocation(program, 1, "vert\x00")
	ctx.EnableVertexAttribArray(1)
	ctx.VertexAttribPointer(1, 2, gl.FLOAT, false, 4*4, 0)
	GetErrors("CfgVabVbo4")

	ctx.BindAttribLocation(program, 2, "vertTexCoord\x00")
	ctx.EnableVertexAttribArray(2)
	ctx.VertexAttribPointer(2, 2, gl.FLOAT, false, 4*4, 2*4)
	GetErrors("CfgVabVbo5")

	ctx.BindBuffer(gl.ARRAY_BUFFER, 0)
	ctx.BindVertexArray(0)
	GetErrors("CfgVabVbo1 " + from)
}

// GenerateTexture will bind a rgba image to a texture and return its "name"
func GenerateTexture(rgba *image.RGBA) uint32 {
	var texture uint32
	texture = ctx.CreateTexture()
	ctx.BindTexture(gl.TEXTURE_2D, texture)
	ctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	ctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	// OBS ctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	// OBS ctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)

	// TexImage2D(target uint32, level int32, internalformat int32, width int32, height int32, format uint32, xtype uint32, pixels []byte)

	ctx.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA,
		int32(rgba.Rect.Dx()), int32(rgba.Rect.Dy()),
		gl.RGBA, gl.UNSIGNED_BYTE, rgba.Pix)
	ctx.BindTexture(gl.TEXTURE_2D, 0)
	GetErrors("GenTexture")
	return texture
}
