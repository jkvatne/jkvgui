package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/shader"
	"golang.org/x/exp/shiny/iconvg"
	"image"
	"image/color"
	"image/draw"
	"log"
)

type Icon struct {
	Img       *image.RGBA
	ImgSize   int
	Color     f32.Color
	Vao       uint32
	Vbo       uint32
	TextureID uint32
}

func NewIcon(sz int, c f32.Color, src []byte) *Icon {
	var err error
	icon := new(Icon)
	icon.ImgSize = sz
	m, _ := iconvg.DecodeMetadata(src)
	dx, dy := m.ViewBox.AspectRatio()
	icon.Color = c
	icon.Img = image.NewRGBA(image.Rectangle{Max: image.Point{X: sz, Y: int(float32(sz) * dy / dx)}})
	var ico iconvg.Rasterizer
	ico.SetDstImage(icon.Img, icon.Img.Bounds(), draw.Src)
	m.Palette[0] = color.RGBA{R: uint8(c.R * 255), G: uint8(c.G * 255), B: uint8(c.B * 255), A: uint8(c.A * 255)}
	_ = iconvg.Decode(&ico, src, &iconvg.DecodeOptions{Palette: &m.Palette})
	// Make program for icon
	gpu.IconProgram, err = shader.NewProgram(shader.VertexQuadShader, shader.FragmentQuadShader)
	if err != nil {
		log.Panicf("Failed to link icon program: %v", err)
	}
	// Generate texture
	gpu.ConfigureVaoVbo(&icon.Vao, &icon.Vbo, gpu.IconProgram)
	icon.TextureID = gpu.GenerateTexture(icon.Img)
	gpu.GetErrors()
	return icon
}

func DrawIcon(xpos, ypos float32, ic *Icon, color f32.Color) {
	gpu.SetResolution(gpu.IconProgram)
	gpu.SetupDrawing(color, ic.Vao, gpu.IconProgram)
	gpu.RenderTexture(xpos, ypos, 100, 100, ic.TextureID, ic.Vbo)
}
