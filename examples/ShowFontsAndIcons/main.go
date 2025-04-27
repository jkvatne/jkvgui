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
	"strconv"
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
	gpu.DrawIcon(x+25, y, 24, gpu.ArrowDropDown, f32.Black)
	gpu.DrawIcon(x+50, y, 24, gpu.Home, f32.Black)
	gpu.DrawIcon(x+75, y, 24, gpu.BoxChecked, f32.Black)
	gpu.DrawIcon(x+100, y, 24, gpu.BoxUnchecked, f32.Black)
	gpu.DrawIcon(x+125, y, 24, gpu.RadioChecked, f32.Black)
	gpu.DrawIcon(x+150, y, 24, gpu.RadioUnchecked, f32.Black)
	gpu.DrawIcon(x+175, y, 24, gpu.ContentOpen, f32.Black)
	gpu.DrawIcon(x+200, y, 24, gpu.ContentSave, f32.Black)
	gpu.DrawIcon(x+225, y, 24, gpu.NavigationArrowDownward, f32.Black)
	gpu.DrawIcon(x+250, y, 24, gpu.NavigationArrowUpward, f32.Black)
	gpu.DrawIcon(x+275, y, 24, gpu.NavigationUnfoldMore, f32.Black)
	gpu.DrawIcon(x+300, y, 24, gpu.NavigationArrowDropDown, f32.Black)
	gpu.DrawIcon(x+325, y, 24, gpu.NavigationArrowDropUp, f32.Black)
}

func ShowFonts(x float32, y float32) {
	for _, f := range font.Fonts {
		if f != nil {
			f.DrawText(x, y, f32.Black, 0, gpu.LTR, strconv.Itoa(f.No)+" "+f.Name+" "+strconv.Itoa(f.Size))
			y += 25
		}
	}
}

var window *glfw.Window

func main() {
	// Setting this true will draw a light blue frame around widgets.
	gpu.DebugWidgets = false
	theme.SetDefaultPallete(lightMode)
	window = gpu.InitWindow(0, 0, "Fonts and images", 1)
	defer gpu.Shutdown()
	sys.Initialize(window)
	for !window.ShouldClose() {
		sys.StartFrame(theme.Surface.Bg())
		// Paint a red frame around the whole window
		gpu.Rect(gpu.WindowRect.Reduce(2), 1, f32.Transparent, theme.PrimaryColor)
		ShowIcons(0, 10)
		ShowFonts(10, 100)

		font.Fonts[gpu.Normal14].DrawText(400, 250, f32.Black, 250, gpu.BTT, "BottomToTopBottomToTop")
		font.Fonts[gpu.Normal14].DrawText(400, 100, f32.Black, 250, gpu.TTB, "TopToBottomTopToBottom")

		for i := range 14 {
			w := float32(i)*5.0 + 120
			x := float32(450)
			y := 100 + float32(i)*15
			font.Fonts[gpu.Normal14].DrawText(x, y, f32.Black, w, gpu.LTR, "TruncatedTruncatedTruncatedTruncated")
			gpu.VertLine(x+w, y-15, y, 1, f32.Blue)
		}
		sys.EndFrame(10)
	}
}
