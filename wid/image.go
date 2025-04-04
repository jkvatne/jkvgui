package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/shader"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

var imgProgram uint32

type Img struct {
	img       *image.RGBA
	w, h      float32
	vao       uint32
	vbo       uint32
	textureID uint32
}

type ImgStyle struct {
	Width        float32
	Height       float32
	CornerRadius float32
}

var DefaultImgStyle = &ImgStyle{
	Width: 0.5,
}

func (b *ImgStyle) W(w float32) *ImgStyle {
	bb := *b
	bb.Width = w
	return &bb
}

func (b *ImgStyle) H(h float32) *ImgStyle {
	bb := *b
	bb.Height = h
	return &bb
}

// Image is the widget for drawing images
func Image(img *Img, style *ImgStyle, altText string) Wid {
	aspectRatio := float32(img.w) / float32(img.h)
	return func(ctx Ctx) Dim {
		var w, h float32
		if aspectRatio > ctx.Rect.W/ctx.Rect.H {
			// Too wide, scale down height
			h = ctx.Rect.H / aspectRatio
			w = ctx.Rect.W
		} else {
			h = ctx.Rect.H
			w = ctx.Rect.W * aspectRatio
		}
		if ctx.Mode == CollectWidths {
			if style.Width < 1.0 {
				return Dim{W: style.Width, H: style.Height}
			}
			return Dim{W: w, H: h}
		} else if ctx.Mode == CollectHeights {
			return Dim{W: w, H: h}
		} else {
			Draw(ctx.Rect.X, ctx.Rect.Y, ctx.Rect.W, ctx.Rect.H, img)
			return Dim{W: ctx.Rect.W, H: ctx.Rect.H}
		}
	}
}

// New generates a new image struct with the rgba image data
// It can later be displayed by using Draw()
func New(filename string) (*Img, error) {
	f, err := os.Open(filename)
	f32.ExitOn(err, "Failed to open image file %s", filename)
	defer f.Close()
	var img = Img{}
	m, _, err := image.Decode(f)
	f32.ExitOn(err, "Failed to decode image %s", filename)
	var ok bool
	img.img, ok = m.(*image.RGBA)
	if !ok {
		// The decoded image was not rgba. Draw it into a new rgba image
		b := m.Bounds()
		img.img = image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
		draw.Draw(img.img, b, m, b.Min, draw.Src)
	}
	bounds := m.Bounds()
	img.w = float32(bounds.Dx())
	img.h = float32(bounds.Dy())
	if imgProgram == 0 {
		imgProgram, err = shader.NewProgram(shader.VertQuadSource, shader.FragImgSource)
		f32.ExitOn(err, "Failed to link icon program: %v", err)
	}
	gpu.ConfigureVaoVbo(&img.vao, &img.vbo, imgProgram)
	img.textureID = gpu.GenerateTexture(img.img)
	return &img, nil
}

// Draw will paint the image to the screen, and scale it
func Draw(x, y, w float32, h float32, img *Img) {
	gpu.Scale(gpu.ScaleX, &x, &y, &w, &h)
	gpu.SetupDrawing(f32.Black, img.vao, imgProgram)
	gpu.RenderTexture(x, y, w, h, img.textureID, img.vbo, 0)
}
