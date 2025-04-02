package img

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/shader"
	"github.com/jkvatne/jkvgui/wid"
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

type Mode int

const (
	// FIT will stretch the image to fit the box given in ctx
	FIT Mode = iota
	// NATIVE will make the image size equal to the size in the file
	NATIVE
	SCALE
)

// W is the widget for drawing images
func W(img *Img, mode Mode, altText string) wid.Wid {
	return func(ctx wid.Ctx) wid.Dim {
		w := ctx.Rect.W
		h := ctx.Rect.W / img.w * img.h
		if !ctx.Draw {
			if mode == NATIVE {
				return wid.Dim{W: img.w, H: img.h}
			} else if mode == SCALE {
				return wid.Dim{W: w, H: h}
			}
		}
		if !ctx.Draw {
			return wid.Dim{W: img.w, H: img.h, Baseline: 0}
		}
		if mode == FIT {
			Draw(ctx.Rect.X, ctx.Rect.Y, ctx.Rect.W, ctx.Rect.H, img)
			return wid.Dim{W: ctx.Rect.W, H: ctx.Rect.H}
		} else if mode == NATIVE {
			Draw(ctx.Rect.X, ctx.Rect.Y, img.w, img.h, img)
			return wid.Dim{W: img.w, H: img.h}
		} else if mode == SCALE {
			Draw(ctx.Rect.X, ctx.Rect.Y, w, h, img)
			return wid.Dim{W: w, H: h}
		}
		return wid.Dim{}
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
