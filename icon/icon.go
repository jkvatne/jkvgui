package icon

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/shader"
	"golang.org/x/exp/shiny/iconvg"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image"
	"image/color"
	"image/draw"
	"log/slog"
	"os"
)

var iconProgram uint32

type Icon struct {
	img       *image.RGBA
	imgSize   int
	color     f32.Color
	vao       uint32
	vbo       uint32
	textureID uint32
}

func New(sz int, src []byte) *Icon {
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
	gpu.ConfigureVaoVbo(&icon.vao, &icon.vbo, iconProgram)
	icon.textureID = gpu.GenerateTexture(icon.img)
	gpu.GetErrors()
	return icon
}

func Draw(x, y, w float32, icon *Icon, color f32.Color) {
	x *= gpu.ScaleX
	y *= gpu.ScaleY
	w *= gpu.ScaleX
	if iconProgram == 0 {
		var err error
		// Make program for icon
		iconProgram, err = shader.NewProgram(shader.VertQuadSource, shader.FragQuadSource)
		if err != nil {
			slog.Error("Failed to link icon program: %v", err)
			os.Exit(1)
		}
		// Generate texture
		gpu.ConfigureVaoVbo(&icon.vao, &icon.vbo, iconProgram)
		icon.textureID = gpu.GenerateTexture(icon.img)
		gpu.GetErrors()
	}
	// SetResolution(iconProgram)
	gpu.SetupDrawing(color, icon.vao, iconProgram)
	gpu.RenderTexture(x, y, w, w, icon.textureID, icon.vbo)
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
	NavigationArrowDropDown = New(48, icons.NavigationArrowDropDown)
	Home = New(48, icons.ActionHome)
	BoxChecked = New(48, icons.ToggleCheckBox)
	BoxUnchecked = New(48, icons.ToggleCheckBoxOutlineBlank)
	RadioChecked = New(48, icons.ToggleRadioButtonChecked)
	RadioUnchecked = New(48, icons.ToggleRadioButtonUnchecked)
	ContentSave = New(48, icons.ContentSave)
	ContentOpen = New(48, icons.FileFolderOpen)
	ArrowDropDown = New(48, arrowDropDownData)
	NavigationArrowDownward = New(48, icons.NavigationArrowDownward)
	NavigationArrowUpward = New(48, icons.NavigationArrowUpward)
	NavigationUnfoldMore = New(48, icons.NavigationUnfoldMore)
	NavigationArrowDropUp = New(48, icons.NavigationArrowDropUp)
}
