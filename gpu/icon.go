package gpu

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/shader"
	"golang.org/x/exp/shiny/iconvg"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image"
	"image/color"
	"image/draw"
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
