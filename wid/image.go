package wid

import (
	"bytes"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/theme"
)

type Img struct {
	img       *image.RGBA
	w, h      float32
	textureID uint32
}

type Fit int

const (
	FitAll Fit = iota
	FitHeight
	FitWidth
	Original
)

type ImgStyle struct {
	OutsidePadding f32.Padding
	BorderRole     theme.UIRole
	SurfaceRole    theme.UIRole
	BorderWidth    float32
	CornerRadius   float32
	Width          float32
	Height         float32
	Scaling        Fit
	Align          Alignment
}

var DefImg = &ImgStyle{
	Width:          0,
	Height:         0,
	OutsidePadding: f32.Padding{L: 5, T: 5, R: 5, B: 5},
	BorderRole:     theme.Outline,
	SurfaceRole:    theme.Surface,
	BorderWidth:    1.0,
	CornerRadius:   4.0,
	Scaling:        FitAll,
	Align:          AlignCenter,
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

func NewImageFrom(buffer []byte) (*Img, error) {
	var img = Img{}
	reader := bytes.NewReader(buffer)
	m, _, err := image.Decode(reader)
	f32.ExitOn(err, "Failed to decode image")
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
	return &img, nil
}

// NewImage generates a new image struct with the rgba image data
// It can later be displayed by using Draw()
func NewImage(filename string) (*Img, error) {
	f, err := os.Open(filename)
	f32.ExitOn(err, "Failed to open image file "+filename)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)
	var img = Img{}
	m, _, err := image.Decode(f)
	f32.ExitOn(err, "Failed to decode image "+filename)
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

	return &img, nil
}

// Draw will paint the image to the screen, and scale it
func Draw(Gd *gpu.GlData, x, y, w float32, h float32, img *Img) {
	if img.textureID == 0 {
		img.textureID = gpu.GenerateTexture(img.img)
	}
	f32.Scale(Gd.ScaleX, &x, &y, &w, &h)
	gpu.SetupTexture(f32.Red, Gd.FontVao, Gd.FontVbo, Gd.ImgProgram)
	gpu.RenderTexture(x, y, w, h, img.textureID, 0)
}

// Image is the widget for drawing images
func Image(img *Img, action func(), style *ImgStyle, altText string) Wid {
	if img.h == 0 {
		return nil
	}
	aspectRatio := img.w / img.h
	if style == nil {
		style = DefImg
	}
	return func(ctx Ctx) Dim {
		var w, h float32

		if ctx.Mode != RenderChildren {
			if style.Height > 1.0 && style.Width > 1.0 {
				return Dim{W: style.Width, H: style.Height}
			} else if style.Height > 1.0 {
				return Dim{W: h * aspectRatio, H: style.Height}
			} else if style.Width > 1.0 {
				return Dim{W: style.Width, H: style.Width / aspectRatio}
			} else if style.Height > 0.0 {
				return Dim{W: 0, H: 0}
			} else if ctx.Mode == CollectHeights {
				return Dim{W: ctx.W, H: ctx.W / aspectRatio}
			} else {
				return Dim{W: ctx.H * aspectRatio, H: ctx.H / aspectRatio}
			}
		}

		if action != nil && ctx.Win.LeftBtnClick(ctx.Rect) {
			ctx.Win.SetFocusedTag(action)
			if !ctx.Disabled {
				action()
				ctx.Win.Invalidate()
			}
		}

		ctx.Rect = ctx.Rect.Inset(style.OutsidePadding, style.BorderWidth)
		if style.Scaling == FitAll {
			if aspectRatio > ctx.Rect.W/ctx.Rect.H {
				// Too wide, scale down height
				w = ctx.Rect.W
				h = w / aspectRatio
			} else {
				// Too high, scale down width
				h = ctx.Rect.H
				w = h * aspectRatio
			}
		} else if style.Scaling == FitHeight {
			h = ctx.Rect.H
			w = h * aspectRatio
		} else if style.Scaling == FitWidth {
			w = ctx.Rect.W
			h = w / aspectRatio
		} else if style.Scaling == Original {
			w = img.w
			h = img.h
		} else {
			w = img.w
			h = img.h
		}

		x := ctx.Rect.X
		if style.Align == AlignCenter {
			x += (ctx.Rect.W - w) / 2
		} else if style.Align == AlignRight {
			x += ctx.Rect.W - w + style.OutsidePadding.L
		} else if style.Align == AlignLeft {
			x += style.OutsidePadding.L
		}
		if ctx.Mode == CollectWidths {
			if style.Width < 1.0 {
				return Dim{W: style.Width, H: style.Height}
			}
		} else if ctx.Mode == CollectHeights {
			if style.Height < 1.0 {
				return Dim{W: style.Width, H: style.Height}
			}
		} else if ctx.Mode == RenderChildren {
			Draw(&ctx.Win.Gd, x, ctx.Rect.Y, w, h, img)
			// Cover rounded corners with the background surface
			if style.BorderWidth > 0 {
				r := f32.Rect{X: x, Y: ctx.Y, W: w, H: h}
				ctx.Win.Gd.RR(r, style.CornerRadius, style.BorderWidth, f32.Transparent, style.SurfaceRole.Fg(), style.SurfaceRole.Bg())
			}
		}
		ctx0 := ctx
		ctx0.W = w
		ctx0.H = h
		if ctx.Win.Hovered(ctx0.Rect) {
			Hint(ctx, altText, img)
		}
		return Dim{
			w + style.OutsidePadding.L + style.OutsidePadding.R + style.BorderWidth,
			h + style.OutsidePadding.T + style.OutsidePadding.B + style.BorderWidth, 0}
	}
}
