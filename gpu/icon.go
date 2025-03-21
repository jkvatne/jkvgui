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
	// temp, err := png.Decode(f)
	// img.img = &temp
	// if err != nil {
	//	slog.Error("Failed to open ", filename)
	// }

	m, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	ii, ok := m.(*image.RGBA)
	if ok {
		img.img = ii
		bounds := m.Bounds()
		slog.Info("Image ", "name", filename, "w", bounds.Max.X, "h", bounds.Max.Y)
		iconProgram, err = shader.NewProgram(shader.VertexQuadShader, shader.FragmentImgShader)
	} else {
		slog.Error("Wrong image format")
	}
	if err != nil {
		slog.Error("Failed to link icon program: %v", err)
		return nil, err
	}
	// Generate texture
	ConfigureVaoVbo(&img.vao, &img.vbo, iconProgram)
	img.textureID = GenerateTexture(img.img)
	GetErrors()
	return &img, nil
}

func DrawImage(x, y, w float32, im *Img) {
	x *= ScaleX
	y *= ScaleY
	w *= ScaleX
	SetResolution(iconProgram)
	SetupDrawing(f32.Black, im.vao, iconProgram)
	RenderTexture(x, y, w, w, im.textureID, im.vbo)
}

func NewIcon(sz int, src []byte) *Icon {
	var err error
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
	iconProgram, err = shader.NewProgram(shader.VertexQuadShader, shader.FragmentQuadShader)
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

func DrawIcon(x, y, w float32, ic *Icon, color f32.Color) {
	x *= ScaleX
	y *= ScaleY
	w *= ScaleX
	SetResolution(iconProgram)
	SetupDrawing(color, ic.vao, iconProgram)
	RenderTexture(x, y, w, w, ic.textureID, ic.vbo)
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
