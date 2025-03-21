package gpu

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/shader"
	"golang.org/x/exp/shiny/iconvg"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"log/slog"
	"os"
)

type Icon struct {
	img       *image.RGBA
	imgSize   int
	color     f32.Color
	vao       uint32
	vbo       uint32
	textureID uint32
}

type Img struct {
	img       *image.RGBA
	w, h      float32
	vao       uint32
	vbo       uint32
	textureID uint32
}

func NewImg(filename string) (*Img, error) {
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
	ConfigureVaoVbo(&img.vao, &img.vbo, imgProgram)
	img.textureID = GenerateTexture(img.img)
	// SetResolution(imgProgram)
	GetErrors()
	return &img, nil
}

func DrawImage(x, y, w float32, img *Img) {
	x *= ScaleX
	y *= ScaleY
	w *= ScaleX
	SetupDrawing(f32.Black, img.vao, imgProgram)
	RenderTexture(x, y, w, w, img.textureID, img.vbo)
}

func NewIcon(sz int, src []byte) *Icon {
	icon := new(Icon)
	icon.imgSize = sz
	m, _ := iconvg.DecodeMetadata(src)
	dx, dy := m.ViewBox.AspectRatio()
	icon.color = f32.White
	icon.img = image.NewRGBA(image.Rectangle{Max: image.Point{X: sz, Y: int(float32(sz) * dy / dx)}})
	var ico iconvg.Rasterizer
	ico.SetDstImage(icon.img, icon.img.Bounds(), draw.Src)
	m.Palette[0] = color.RGBA{R: 0, G: 0, B: 255, A: 255}
	m.Palette[0] = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	_ = iconvg.Decode(&ico, src, &iconvg.DecodeOptions{Palette: &m.Palette})
	// Make program for icon
	var err error
	iconProgram, err = shader.NewProgram(shader.VertQuadSource, shader.FragQuadSource)
	if err != nil {
		slog.Error("Failed to link icon program: %v", err)
		os.Exit(1)
	}
	// Generate texture
	ConfigureVaoVbo(&icon.vao, &icon.vbo, iconProgram)
	icon.textureID = GenerateTexture(icon.img)
	GetErrors()
	return icon
}

func DrawIcon(x, y, w float32, icon *Icon, color f32.Color) {
	x *= ScaleX
	y *= ScaleY
	w *= ScaleX
	if iconProgram == 0 {
		var err error
		// Make program for icon
		iconProgram, err = shader.NewProgram(shader.VertQuadSource, shader.FragQuadSource)
		if err != nil {
			slog.Error("Failed to link icon program: %v", err)
			os.Exit(1)
		}
		// Generate texture
		ConfigureVaoVbo(&icon.vao, &icon.vbo, iconProgram)
		icon.textureID = GenerateTexture(icon.img)
		GetErrors()
	}
	// SetResolution(iconProgram)
	SetupDrawing(color, icon.vao, iconProgram)
	RenderTexture(x, y, w, w, icon.textureID, icon.vbo)
}

var (
	Home                    *Icon
	BoxChecked              *Icon
	BoxUnchecked            *Icon
	RadioChecked            *Icon
	RadioUnchecked          *Icon
	ContentSave             *Icon
	ContentOpen             *Icon
	NavigationArrowDownward *Icon
	NavigationArrowUpward   *Icon
	NavigationUnfoldMore    *Icon
	NavigationArrowDropDown *Icon
	NavigationArrowDropUp   *Icon
	ArrowDropDown           *Icon
)

var arrowDropDownData = []byte{
	0x89, 0x49, 0x56, 0x47, 0x02, 0x0a, 0x00, 0x50, 0x50, 0xb0, 0xb0,
	0xc0, 0x62, 0x70, // Start point at -15, -8
	0x21, 0x9E, 0x9E, // Lineto 15,15
	0x9E, 0x62, // Lineto 15,-15
	0xe1,
}

func LoadIcons() {
	NavigationArrowDropDown = NewIcon(48, icons.NavigationArrowDropDown)
	Home = NewIcon(48, icons.ActionHome)
	BoxChecked = NewIcon(48, icons.ToggleCheckBox)
	BoxUnchecked = NewIcon(48, icons.ToggleCheckBoxOutlineBlank)
	RadioChecked = NewIcon(48, icons.ToggleRadioButtonChecked)
	RadioUnchecked = NewIcon(48, icons.ToggleRadioButtonUnchecked)
	ContentSave = NewIcon(48, icons.ContentSave)
	ContentOpen = NewIcon(48, icons.FileFolderOpen)
	ArrowDropDown = NewIcon(48, arrowDropDownData)
	NavigationArrowDownward = NewIcon(48, icons.NavigationArrowDownward)
	NavigationArrowUpward = NewIcon(48, icons.NavigationArrowUpward)
	NavigationUnfoldMore = NewIcon(48, icons.NavigationUnfoldMore)
	NavigationArrowDropUp = NewIcon(48, icons.NavigationArrowDropUp)
}
