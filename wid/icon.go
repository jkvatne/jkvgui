package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/shader"
	"golang.org/x/exp/shiny/iconvg"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image"
	"image/color"
	"image/draw"
	"log"
)

type Icon struct {
	Img       *image.RGBA
	ImgSize   int
	Color     f32.Color
	Vao       uint32
	Vbo       uint32
	TextureID uint32
}

func NewIcon(sz int, src []byte) *Icon {
	var err error
	icon := new(Icon)
	icon.ImgSize = sz
	m, _ := iconvg.DecodeMetadata(src)
	dx, dy := m.ViewBox.AspectRatio()
	icon.Color = f32.White
	icon.Img = image.NewRGBA(image.Rectangle{Max: image.Point{X: sz, Y: int(float32(sz) * dy / dx)}})
	var ico iconvg.Rasterizer
	ico.SetDstImage(icon.Img, icon.Img.Bounds(), draw.Src)
	m.Palette[0] = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	_ = iconvg.Decode(&ico, src, &iconvg.DecodeOptions{Palette: &m.Palette})
	// Make program for icon
	gpu.IconProgram, err = shader.NewProgram(shader.VertexQuadShader, shader.FragmentQuadShader)
	if err != nil {
		log.Panicf("Failed to link icon program: %v", err)
	}
	// Generate texture
	gpu.ConfigureVaoVbo(&icon.Vao, &icon.Vbo, gpu.IconProgram)
	icon.TextureID = gpu.GenerateTexture(icon.Img)
	gpu.GetErrors()
	return icon
}

func DrawIcon(x, y, w float32, ic *Icon, color f32.Color) {
	x *= gpu.Scale
	y *= gpu.Scale
	w *= gpu.Scale
	gpu.SetResolution(gpu.IconProgram)
	gpu.SetupDrawing(color, ic.Vao, gpu.IconProgram)
	gpu.RenderTexture(x, y, w, w, ic.TextureID, ic.Vbo)
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
)

func LoadIcons() {
	Home = NewIcon(48, icons.ActionHome)
	BoxChecked = NewIcon(48, icons.ToggleCheckBox)
	BoxUnchecked = NewIcon(48, icons.ToggleCheckBoxOutlineBlank)
	RadioChecked = NewIcon(48, icons.ToggleRadioButtonChecked)
	RadioUnchecked = NewIcon(48, icons.ToggleRadioButtonUnchecked)
	ContentSave = NewIcon(48, icons.ContentSave)
	NavigationArrowDownward = NewIcon(48, icons.NavigationArrowDownward)
	NavigationArrowUpward = NewIcon(48, icons.NavigationArrowUpward)
	NavigationUnfoldMore = NewIcon(48, icons.NavigationUnfoldMore)
	NavigationArrowDropDown = NewIcon(48, icons.NavigationArrowDropDown)
	NavigationArrowDropUp = NewIcon(48, icons.NavigationArrowDropUp)
	ContentOpen = NewIcon(48, icons.FileFolderOpen)
}
