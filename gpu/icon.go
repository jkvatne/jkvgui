package gpu

import (
	"github.com/jkvatne/jkvgui/f32"
	"golang.org/x/exp/shiny/iconvg"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image"
	"image/color"
	"image/draw"
)

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
var iconProgram uint32

// Icon is the data structure containing the rgba image and
// other persistent data. The color is specified while draing it.
type Icon struct {
	img       *image.RGBA
	vao       uint32
	vbo       uint32
	textureID uint32
}

// New creates a new Icon structure containing the rgba image of it at a given size
// The size should be big enough so later scaling while drawing will not make it fuzzy.
func New(sz int, src []byte) *Icon {
	icon := new(Icon)
	m, err := iconvg.DecodeMetadata(src)
	f32.ExitOn(err, "Failed to decode icon metadata: %v", err)
	dx, dy := m.ViewBox.AspectRatio()
	// icon.color = f32.White
	icon.img = image.NewRGBA(image.Rectangle{Max: image.Point{X: sz, Y: int(float32(sz) * dy / dx)}})
	var ico iconvg.Rasterizer
	ico.SetDstImage(icon.img, icon.img.Bounds(), draw.Src)
	m.Palette[0] = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	err = iconvg.Decode(&ico, src, &iconvg.DecodeOptions{Palette: &m.Palette})
	f32.ExitOn(err, "Failed to decode icon metadata: %v", err)
	if iconProgram == 0 {
		iconProgram, err = NewProgram(VertQuadSource, FragQuadSource)
		f32.ExitOn(err, "Failed to link icon program: %v", err)
	}
	ConfigureVaoVbo(&icon.vao, &icon.vbo, iconProgram)
	icon.textureID = GenerateTexture(icon.img)
	return icon
}

// Draw will paint the icon to the screen, and scale it
func Draw(x, y, w float32, icon *Icon, color f32.Color) {
	Scale(ScaleX, &x, &y, &w)
	SetupDrawing(color, icon.vao, iconProgram)
	RenderTexture(x, y, w, w, icon.textureID, icon.vbo, 0)
}

// LoadIcons will pre-load some often used icons
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
