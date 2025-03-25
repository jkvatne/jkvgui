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
	"log"
	"log/slog"
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
)

func W(img *Img, mode Mode, altText string) wid.Wid {
	return func(ctx wid.Ctx) wid.Dim {
		if ctx.Rect.W == 0 && ctx.Rect.H == 0 {
			return wid.Dim{W: img.w, H: img.h, Baseline: 0}
		}
		if mode == FIT {
			Draw(ctx.Rect.X, ctx.Rect.Y, ctx.Rect.W, ctx.Rect.H, img)
		} else {
			Draw(ctx.Rect.X, ctx.Rect.Y, img.w, img.h, img)
		}
		return wid.Dim{W: ctx.Rect.W, H: ctx.Rect.H}
	}
}

func New(filename string) (*Img, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var img = Img{}
	m, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	var ok bool
	img.img, ok = m.(*image.RGBA)
	if !ok {
		b := m.Bounds()
		img.img = image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
		draw.Draw(img.img, b, m, b.Min, draw.Src)
	}
	bounds := m.Bounds()
	img.w = float32(bounds.Dx())
	img.h = float32(bounds.Dy())
	slog.Info("Image ", "name", filename, "w", bounds.Max.X, "h", bounds.Max.Y)
	imgProgram, err = shader.NewProgram(shader.VertQuadSource, shader.FragImgSource)
	if err != nil {
		slog.Error("Failed to link icon program: %v", err)
		return nil, err
	}
	// Generate texture
	gpu.ConfigureVaoVbo(&img.vao, &img.vbo, imgProgram)
	img.textureID = gpu.GenerateTexture(img.img)
	// SetResolution(imgProgram)
	gpu.GetErrors()
	return &img, nil
}

func Draw(x, y, w float32, h float32, img *Img) {
	x *= gpu.ScaleX
	y *= gpu.ScaleY
	w *= gpu.ScaleX
	h *= gpu.ScaleY
	gpu.SetupDrawing(f32.Black, img.vao, imgProgram)
	gpu.RenderTexture(x, y, w, h, img.textureID, img.vbo)
}
