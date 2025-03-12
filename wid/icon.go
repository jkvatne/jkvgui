package wid

import (
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

	// Dummy image
	// icon.Img = image.NewRGBA(image.Rect(0, 0, 150, 150))
	// blue := color.RGBA{0, 0, 255, 255}
	// draw.Draw(icon.Img, icon.Img.Bounds(), &image.Uniform{blue}, image.ZP, draw.Src)

	// Write icon image to file for testing
	f, err := os.Create("img.jpg")
	if err != nil {
		panic(err)
	}
	if err = jpeg.Encode(f, icon.Img, nil); err != nil {
		log.Printf("failed to encode: %v", err)
	}
	f.Close()

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

func DrawIcon(xpos, ypos float32, ic *Icon) {
	gpu.GetErrors()
	gpu.SetResolution(gpu.IconProgram)
	gpu.SetupDrawing(ic.Color, ic.Vao, gpu.IconProgram)
	gpu.RenderTexture(xpos, ypos, 100, 100, ic.TextureID, ic.Vbo)
	gpu.GetErrors()

}
