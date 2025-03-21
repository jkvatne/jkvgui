package img

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/shader"
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

func Draw(x, y, w float32, img *Img) {
	x *= gpu.ScaleX
	y *= gpu.ScaleY
	w *= gpu.ScaleX
	gpu.SetupDrawing(f32.Black, img.vao, imgProgram)
	gpu.RenderTexture(x, y, w, w, img.textureID, img.vbo)
}
