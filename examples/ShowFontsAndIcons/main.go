package main

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/gpu/font"
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
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

var img1 *wid.Img
var img2 *wid.Img

func ShowIcons(x float32, y float32) {
	gpu.Draw(x+25, y, 24, gpu.ArrowDropDown, f32.Black)
	gpu.Draw(x+50, y, 24, gpu.Home, f32.Black)
	gpu.Draw(x+75, y, 24, gpu.BoxChecked, f32.Black)
	gpu.Draw(x+100, y, 24, gpu.BoxUnchecked, f32.Black)
	gpu.Draw(x+125, y, 24, gpu.RadioChecked, f32.Black)
	gpu.Draw(x+150, y, 24, gpu.RadioUnchecked, f32.Black)
	gpu.Draw(x+175, y, 24, gpu.ContentOpen, f32.Black)
	gpu.Draw(x+200, y, 24, gpu.ContentSave, f32.Black)
	gpu.Draw(x+225, y, 24, gpu.NavigationArrowDownward, f32.Black)
	gpu.Draw(x+250, y, 24, gpu.NavigationArrowUpward, f32.Black)
	gpu.Draw(x+275, y, 24, gpu.NavigationUnfoldMore, f32.Black)
	gpu.Draw(x+300, y, 24, gpu.NavigationArrowDropDown, f32.Black)
	gpu.Draw(x+325, y, 50, gpu.NavigationArrowDropUp, f32.Black)
	gpu.Draw(x+375, y, 100, 100, img1)
	gpu.Draw(x+500, y, 100, 100, img2)

}

func ShowFonts(x float32, y float32) {
	font.Fonts[gpu.Normal].DrawText(x, y, f32.Black, 2, 0, gpu.LeftToRight, "24 Normal")
	gpu.HorLine(x, x+200, y, 2, f32.Blue)
	font.Fonts[gpu.Bold].DrawText(x, y+30, f32.Blue, 2, 0, gpu.LeftToRight, "24 Bold")
	font.Fonts[gpu.Mono].DrawText(x, y+60, f32.Black, 2, 0, gpu.LeftToRight, "24 Mono")
	font.Fonts[gpu.Italic].DrawText(x, y+90, f32.Black, 2, 0, gpu.LeftToRight, "24 Italic")
}

var window *glfw.Window

func main() {
	theme.SetDefaultPallete(lightMode)
	window = gpu.InitWindow(0, 0, "Fonts and images", 1)
	defer gpu.Shutdown()
	sys.Initialize(window, 16)
	img1, _ = gpu.New("mook-logo.png")
	img2, _ = gpu.New("music.jpg")

	for !window.ShouldClose() {
		sys.StartFrame(theme.Surface.Bg())
		// Paint a red frame around the whole window
		gpu.Rect(gpu.WindowRect.Reduce(2), 1, f32.Transparent, theme.PrimaryColor)
		ShowIcons(0, 10)
		ShowFonts(10, 100)

		font.Fonts[gpu.Normal].DrawText(500, 200, f32.Black, 2, 250, gpu.TopToBottom, "TopToBottomTopToBottom")
		gpu.VertLine(500, 200, 200+180, 1, f32.Blue)

		font.Fonts[gpu.Normal].DrawText(600, 400, f32.Black, 2, 250, gpu.BottomToTop, "BottomToTopBottomToTop")
		gpu.VertLine(600, 400-180, 400, 1, f32.Blue)

		font.Fonts[gpu.Normal].DrawText(650, 50, f32.Black, 2, 360, gpu.LeftToRight, "TopToBottomTopToBottom")
		gpu.VertLine(650+360, 0, 50, 1, f32.Blue)

		for i := range 14 {
			w := float32(i)*5.0 + 120
			x := float32(20)
			y := 250 + float32(i)*35
			font.Fonts[gpu.Normal].DrawText(x, y, f32.Black, 2, w, gpu.LeftToRight, "TruncatedTruncatedTruncatedTruncated")
			gpu.HorLine(x, x+w, y, 2, f32.Blue)
			gpu.VertLine(x+w, y-25, y, 1, f32.Blue)
		}
		sys.EndFrame(10)
	}
}
