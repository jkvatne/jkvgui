package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/theme"
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
	textureID uint32
}

type ImgStyle struct {
	OutsidePadding f32.Padding
	BorderRole     theme.UIRole
	SurfaceRole    theme.UIRole
	BorderWidth    float32
	CornerRadius   float32
	Width          float32
	Height         float32
}

var DefImg = &ImgStyle{
	Width:          0.5,
	OutsidePadding: f32.Padding{L: 5, T: 3, R: 4, B: 3},
	BorderRole:     theme.Outline,
	SurfaceRole:    theme.Surface,
	BorderWidth:    0.0,
	CornerRadius:   15.0,
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

func (b *ImgStyle) Bg(r theme.UIRole) *ImgStyle {
	bb := *b
	bb.SurfaceRole = r
	return &bb
}

// NewImage generates a new image struct with the rgba image data
// It can later be displayed by using Draw()
func NewImage(filename string) (*Img, error) {
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
	img.textureID = gpu.GenerateTexture(img.img)
	return &img, nil
}

// Draw will paint the image to the screen, and scale it
func Draw(x, y, w float32, h float32, img *Img) {
	gpu.Scale(gpu.ScaleX, &x, &y, &w, &h)
	gpu.SetupTexture(f32.Red, gpu.FontVao, gpu.ImgProgram)
	gpu.RenderTexture(x, y, w, h, img.textureID, gpu.FontVbo, 0)
}

// Image is the widget for drawing images
func Image(img *Img, style *ImgStyle, altText string) Wid {
	aspectRatio := float32(img.w) / float32(img.h)
	if style == nil {
		style = DefImg
	}

	return func(ctx Ctx) Dim {
		var w, h float32
		ctx.Rect = ctx.Rect.Inset(style.OutsidePadding, style.BorderWidth)

		if aspectRatio > ctx.Rect.W/ctx.Rect.H {
			// Too wide, scale down height
			w = ctx.Rect.W
			h = w / aspectRatio
		} else {
			// Too high, scale down width
			h = ctx.Rect.H
			w = h * aspectRatio
		}

		ctx.Rect.W = w
		ctx.Rect.H = h

		if ctx.Mode == CollectWidths {
			if style.Width < 1.0 {
				return Dim{W: style.Width, H: style.Height}
			}
			return Dim{W: w, H: h}
		} else if ctx.Mode == CollectHeights {
			return Dim{W: w, H: h}
		} else {
			Draw(ctx.Rect.X, ctx.Rect.Y, w, h, img)
			// Cover rounded corners with the background surface
			gpu.RR(ctx.Rect, style.CornerRadius, 2.0, f32.Transparent, style.SurfaceRole.Fg(), style.SurfaceRole.Bg())
			return Dim{W: w, H: h}
		}

	}
}
