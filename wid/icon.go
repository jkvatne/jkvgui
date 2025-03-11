package wid

import (
	"github.com/go-gl/gl/all-core/gl"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/shader"
	"golang.org/x/exp/shiny/iconvg"
	"image"
	"image/color"
	"image/draw"
)

type Icon struct {
	img       *image.RGBA
	imgSize   int
	imgColor  f32.Color
	vao       uint32
	vbo       uint32
	textureID uint32
	program   uint32
}

func NewIcon(sz int, c f32.Color, src []byte) *Icon {
	icon := new(Icon)
	icon.imgSize = sz
	m, _ := iconvg.DecodeMetadata(src)
	dx, dy := m.ViewBox.AspectRatio()
	icon.imgColor = c
	icon.img = image.NewRGBA(image.Rectangle{Max: image.Point{X: sz, Y: int(float32(sz) * dy / dx)}})
	var ico iconvg.Rasterizer
	ico.SetDstImage(icon.img, icon.img.Bounds(), draw.Src)
	m.Palette[0] = color.RGBA{R: uint8(c.R * 255), G: uint8(c.G * 255), B: uint8(c.B * 255), A: uint8(c.A * 255)}
	// m.Palette[1] = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	_ = iconvg.Decode(&ico, src, &iconvg.DecodeOptions{Palette: &m.Palette})
	icon.program, _ = shader.NewProgram(shader.VertexFontShader, shader.FragmentFontShader)

	/*f, err := os.Create("img.jpg")
	if err != nil {
		panic(err)
	}
	if err = jpeg.Encode(f, icon.img, nil); err != nil {
		log.Printf("failed to encode: %v", err)
	}
	f.Close()
	*/

	// Generate texture
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(icon.img.Rect.Dx()), int32(icon.img.Rect.Dy()), 0,
		gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(icon.img.Pix))

	icon.textureID = texture
	gl.BindTexture(gl.TEXTURE_2D, 0)

	// Configure VAO/VBO for texture quads
	gl.GenVertexArrays(1, &icon.vao)
	gl.GenBuffers(1, &icon.vbo)
	gl.BindVertexArray(icon.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, icon.vbo)

	gl.BufferData(gl.ARRAY_BUFFER, 6*4*4, nil, gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(icon.program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointerWithOffset(vertAttrib, 2, gl.FLOAT, false, 4*4, 0)
	defer gl.DisableVertexAttribArray(vertAttrib)

	texCoordAttrib := uint32(gl.GetAttribLocation(icon.program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointerWithOffset(texCoordAttrib, 2, gl.FLOAT, false, 4*4, 2*4)
	defer gl.DisableVertexAttribArray(texCoordAttrib)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
	gpu.GetErrors()
	return icon
}

func DrawIcon(xpos, ypos float32, ic *Icon) {
	// setup blending mode
	gpu.GetErrors()
	gpu.SetResolution(ic.program)

	gl.UseProgram(ic.program)
	gl.BindVertexArray(ic.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, ic.vbo)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	// set icon color
	r := gl.GetUniformLocation(ic.program, gl.Str("textColor\x00"))
	gl.Uniform4f(r, ic.imgColor.R, ic.imgColor.G, ic.imgColor.B, ic.imgColor.A)

	gl.ActiveTexture(gl.TEXTURE0)
	w := float32(ic.img.Rect.Max.X)
	h := float32(ic.img.Rect.Max.Y)
	vertices := []float32{
		xpos + w, ypos, 1.0, 0.0,
		xpos, ypos, 0.0, 0.0,
		xpos, ypos + h, 0.0, 1.0,

		xpos, ypos + h, 0.0, 1.0,
		xpos + w, ypos + h, 1.0, 1.0,
		xpos + w, ypos, 1.0, 0.0,
	}
	// Render glyph texture over quad
	gl.BindTexture(gl.TEXTURE_2D, ic.textureID)
	// BufferSubData(target Enum, offset int, data []byte)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW) // Be sure to use glBufferSubData and not glBufferData
	// Render quad
	gl.DrawArrays(gl.TRIANGLES, 0, 16)

	// clear opengl textures and programs
	// gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	// gl.BindVertexArray(0)
	// gl.BindTexture(gl.TEXTURE_2D, 0)
	// gl.UseProgram(0)
	gl.Disable(gl.BLEND)
	gpu.GetErrors()

}
