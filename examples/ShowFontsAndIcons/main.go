package main

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/gpu/font"
	"github.com/jkvatne/jkvgui/icon"
	"github.com/jkvatne/jkvgui/img"
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
	"log/slog"
)

var (
	lightMode = true
	gender    = "Male"
	genders   = []string{"Male", "Female", "Both", "qyjpy"}
	name      = "Ole Petter Olsen"
	address   = "Mo i Rana"
	hint1     = "This is a hint word5 word6 word7 word8 qYyM9 qYyM10"
	hint2     = "This is a hint"
	hint3     = "This is a hint word5 word6 word7 word8 qYyM9 qYyM10 Word11 word12 jyword13"
)

func YesBtnClick() {
	lightMode = true
	theme.SetDefaultPallete(lightMode)
	slog.Info("Yes Btn Clicked")
}

func CancelBtnClick() {
	slog.Info("Cancel Btn clicked")
}

func NoBtnClick() {
	lightMode = false
	theme.SetDefaultPallete(lightMode)
	slog.Info("No Btn Click\n")
}

var img1 *img.Img
var img2 *img.Img

func ShowIcons(x float32, y float32) {
	icon.Draw(x+25, y, 24, icon.ArrowDropDown, f32.Black)
	icon.Draw(x+50, y, 24, icon.Home, f32.Black)
	icon.Draw(x+75, y, 24, icon.BoxChecked, f32.Black)
	icon.Draw(x+100, y, 24, icon.BoxUnchecked, f32.Black)
	icon.Draw(x+125, y, 24, icon.RadioChecked, f32.Black)
	icon.Draw(x+150, y, 24, icon.RadioUnchecked, f32.Black)
	icon.Draw(x+175, y, 24, icon.ContentOpen, f32.Black)
	icon.Draw(x+200, y, 24, icon.ContentSave, f32.Black)
	icon.Draw(x+225, y, 24, icon.NavigationArrowDownward, f32.Black)
	icon.Draw(x+250, y, 24, icon.NavigationArrowUpward, f32.Black)
	icon.Draw(x+275, y, 24, icon.NavigationUnfoldMore, f32.Black)
	icon.Draw(x+300, y, 24, icon.NavigationArrowDropDown, f32.Black)
	icon.Draw(x+325, y, 50, icon.NavigationArrowDropUp, f32.Black)
	img.Draw(x+375, y, 100, 100, img1)
	img.Draw(x+500, y, 100, 100, img2)

}

func ShowFonts(x float32, y float32) {
	font.Fonts[gpu.Normal].SetColor(f32.Black)
	font.Fonts[gpu.Normal].Printf(x, y, 2, 0, gpu.LeftToRight, "24 Normal")
	gpu.HorLine(x, x+200, y, 2, f32.Blue)
	font.Fonts[gpu.Bold].SetColor(f32.Black)
	font.Fonts[gpu.Bold].Printf(x, y+30, 2, 0, gpu.LeftToRight, "24 Bold")
	font.Fonts[gpu.Mono].SetColor(f32.Black)
	font.Fonts[gpu.Mono].Printf(x, y+60, 2, 0, gpu.LeftToRight, "24 Mono")
	font.Fonts[gpu.Italic].SetColor(f32.Black)
	font.Fonts[gpu.Italic].Printf(x, y+90, 2, 0, gpu.LeftToRight, "24 Italic")
}

var window *glfw.Window

func main() {
	theme.SetDefaultPallete(lightMode)
	window = gpu.InitWindow(0, 0, "Fonts and images", 1)
	defer gpu.Shutdown()
	sys.Initialize(window, 16)
	img1, _ = img.New("mook-logo.png")
	img2, _ = img.New("music.jpg")

	for !window.ShouldClose() {
		sys.StartFrame(theme.Surface.Bg())
		// Paint a red frame around the whole window
		gpu.Rect(gpu.WindowRect.Reduce(2), 1, f32.Transparent, theme.PrimaryColor)
		ShowIcons(0, 10)
		ShowFonts(10, 100)

		font.Fonts[gpu.Normal].Printf(500, 200, 2, 250, gpu.TopToBottom, "TopToBottomTopToBottom")
		gpu.VertLine(500, 200, 200+180, 1, f32.Blue)

		font.Fonts[gpu.Normal].Printf(600, 400, 2, 250, gpu.BottomToTop, "BottomToTopBottomToTop")
		gpu.VertLine(600, 400-180, 400, 1, f32.Blue)

		font.Fonts[gpu.Normal].Printf(650, 50, 2, 360, gpu.LeftToRight, "TopToBottomTopToBottom")
		gpu.VertLine(650+360, 0, 50, 1, f32.Blue)

		for i := range 14 {
			w := float32(i)*5.0 + 120
			x := float32(20)
			y := 250 + float32(i)*35
			font.Fonts[gpu.Normal].Printf(x, y, 2, w, gpu.LeftToRight, "TruncatedTruncatedTruncatedTruncated")
			gpu.HorLine(x, x+w, y, 2, f32.Blue)
			gpu.VertLine(x+w, y-25, y, 1, f32.Blue)
		}
		sys.EndFrame(10)
	}
}
