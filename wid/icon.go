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
	"image/jpeg"
	"log"
	"os"
)

type Icon struct {
	Img       *image.RGBA
	imgSize   int
	imgColor  f32.Color
	Vao       uint32
	Vbo       uint32
	TextureID uint32
	Program   uint32
}

func NewIcon(sz int, c f32.Color, src []byte) *Icon {
	icon := new(Icon)
	icon.imgSize = sz
	m, _ := iconvg.DecodeMetadata(src)
	dx, dy := m.ViewBox.AspectRatio()
	icon.imgColor = c
	icon.Img = image.NewRGBA(image.Rectangle{Max: image.Point{X: sz, Y: int(float32(sz) * dy / dx)}})
	var ico iconvg.Rasterizer
	ico.SetDstImage(icon.Img, icon.Img.Bounds(), draw.Src)
	m.Palette[0] = color.RGBA{R: uint8(c.R * 255), G: uint8(c.G * 255), B: uint8(c.B * 255), A: uint8(c.A * 255)}
	_ = iconvg.Decode(&ico, src, &iconvg.DecodeOptions{Palette: &m.Palette})

	icon.Img = image.NewRGBA(image.Rect(0, 0, 150, 150))
	blue := color.RGBA{0, 0, 255, 255}
	draw.Draw(icon.Img, icon.Img.Bounds(), &image.Uniform{blue}, image.ZP, draw.Src)

	f, err := os.Create("img.jpg")
	if err != nil {
		panic(err)
	}
	if err = jpeg.Encode(f, icon.Img, nil); err != nil {
		log.Printf("failed to encode: %v", err)
	}
	f.Close()

	icon.Program, _ = shader.NewProgram(shader.VertexFontShader, shader.FragmentFontShader)
	icon.TextureID = gpu.GenerateTexture(icon.Img)
	gpu.ConfigureVaoVbo(&icon.Vao, &icon.Vbo, icon.Program)
	return icon
}

func DrawIcon(xpos, ypos float32, ic *Icon) {
	// setup blending mode
	gpu.GetErrors()
	gpu.SetResolution(ic.Program)

	gl.UseProgram(ic.Program)
	gl.BindVertexArray(ic.Vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, ic.Vbo)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	// set icon color
	r := gl.GetUniformLocation(ic.Program, gl.Str("textColor\x00"))
	gl.Uniform4f(r, ic.imgColor.R, ic.imgColor.G, ic.imgColor.B, ic.imgColor.A)

	gl.ActiveTexture(gl.TEXTURE0)
	w := float32(ic.Img.Rect.Max.X)
	h := float32(ic.Img.Rect.Max.Y)
	vertices := []float32{
		xpos + w, ypos, 1.0, 0.0,
		xpos, ypos, 0.0, 0.0,
		xpos, ypos + h, 0.0, 1.0,

		xpos, ypos + h, 0.0, 1.0,
		xpos + w, ypos + h, 1.0, 1.0,
		xpos + w, ypos, 1.0, 0.0,
	}
	// Render glyph texture over quad
	gl.BindTexture(gl.TEXTURE_2D, ic.TextureID)
	// BufferSubData(target Enum, offset int, data []byte)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW) // Be sure to use glBufferSubData and not glBufferData
	// Render quad
	gl.DrawArrays(gl.TRIANGLES, 0, 16)

	// clear opengl textures and programs
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	gl.UseProgram(0)
	gl.Disable(gl.BLEND)
	gpu.GetErrors()

}
