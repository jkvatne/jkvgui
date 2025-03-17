package main

import (
	_ "embed"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jkvatne/jkvgui/callback"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/focus"
	"github.com/jkvatne/jkvgui/font"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/wid"
	"log/slog"
	"strconv"
)

//go:embed font/fonts/Roboto-Thin.ttf
var Roboto100 []byte // 100

//go:embed font/fonts/Roboto-ExtraLight.ttf
var Roboto200 []byte // 200

//go:embed font/fonts/Roboto-Light.ttf
var Roboto300 []byte // 300

//go:embed font/fonts/Roboto-Regular.ttf
var Roboto400 []byte // 400

//go:embed font/fonts/Roboto-Medium.ttf
var Roboto500 []byte // 500

//go:embed font/fonts/Roboto-SemiBold.ttf
var Roboto600 []byte // 600

//go:embed font/fonts/Roboto-Bold.ttf
var Roboto700 []byte // 700

//go:embed font/fonts/Roboto-Bold.ttf
var Roboto800 []byte // 800

//go:embed font/fonts/Roboto-Bold.ttf
var Roboto900 []byte // 900

//go:embed font/fonts/RobotoMono-Regular.ttf
var RobotoMono400 []byte

//go:embed font/fonts/RobotoMono-Bold.ttf
var RobotoMono600 []byte

//go:embed font/fonts/RobotoMono-Light.ttf
var RobotoMono200 []byte

var P = f32.Padding{2, 2, 2, 2}

func YesBtnClick() {
	slog.Info("Yes Btn Clicked")
}

func CancelBtnClick() {
	slog.Info("Cancel Btn clicked")
}

func NoBtnClick() {
	slog.Info("No Btn Click\n")
}

var name = "Ole Petter Olsen"
var address = "Mo i Rana"
var hint1 = "This is a hint word5 word6 word7 word8 qYyM9 qYyM10"
var hint2 = "This is a hint"
var hint3 = "This is a hint word5 word6 word7 word8 qYyM9 qYyM10 Word11 word12 jyword13"

func ShowIcons(x float32, y float32) {
	wid.DrawIcon(x+50, y, 24, wid.Home, f32.Blue)
	wid.DrawIcon(x+75, y, 24, wid.BoxChecked, f32.Black)
	wid.DrawIcon(x+100, y, 24, wid.BoxUnchecked, f32.Black)
	wid.DrawIcon(x+125, y, 24, wid.RadioChecked, f32.Black)
	wid.DrawIcon(x+150, y, 24, wid.RadioUnchecked, f32.Black)
	wid.DrawIcon(x+175, y, 24, wid.ContentSave, f32.Black)
	wid.DrawIcon(x+200, y, 24, wid.NavigationArrowDownward, f32.Black)
	wid.DrawIcon(x+225, y, 24, wid.NavigationArrowUpward, f32.Black)
	wid.DrawIcon(x+250, y, 24, wid.NavigationUnfoldMore, f32.Black)
	wid.DrawIcon(x+275, y, 24, wid.NavigationArrowDropDown, f32.Black)
	wid.DrawIcon(x+300, y, 24, wid.NavigationArrowDropUp, f32.Black)
}

// From freetype.go, line 263, Her c.dpi is allways 72.
// c.scale = fixed.Int26_6(0.5 + (c.fontSize * c.dpi * 64 / 72))
// size = fontsize  in pixels.

func LoadFonts() {
	_ = font.LoadFontBytes(Roboto200, 24, "Roboto", 200)
	_ = font.LoadFontBytes(Roboto400, 24, "Roboto", 400)
	_ = font.LoadFontBytes(Roboto600, 24, "Roboto", 600)
	_ = font.LoadFontBytes(RobotoMono200, 24, "RobotoMono", 200)
	_ = font.LoadFontBytes(RobotoMono400, 24, "RobotoMono", 400)
	_ = font.LoadFontBytes(RobotoMono600, 24, "RobotoMono", 600)
}

func ShowFonts(x float32, y float32) {
	font.Fonts[0].Printf(x, y, 2, 0, "24 Roboto200")             // Thin
	font.Fonts[1].Printf(x, y+30, 2, 0, "24 Roboto400")          // Regular
	font.Fonts[2].Printf(x, y+60, 2, 0, "24 Roboto600")          // Bold
	font.Fonts[3].Printf(x, y+90, 2, 0, "24 RobotoMono200")      // Thin
	font.Fonts[4].Printf(x, y+120, 2, 0, "24 RobotoMono400")     // Regular
	font.Fonts[5].Printf(x, y+150, 2, 0, "24 RobotoMono600")     // Bold
	font.Fonts[0].Printf(x+300, y, 1, 0, "12 Roboto200")         // Thin
	font.Fonts[1].Printf(x+300, y+30, 1, 0, "12 Roboto400")      // Regular
	font.Fonts[2].Printf(x+300, y+60, 1, 0, "12 Roboto600")      // Bold
	font.Fonts[3].Printf(x+300, y+90, 1, 0, "12 RobotoMono200")  // Thin
	font.Fonts[4].Printf(x+300, y+120, 1, 0, "12 RobotoMono400") // Regular
	font.Fonts[5].Printf(x+300, y+150, 1, 0, "12 RobotoMono600") // Bold
}

var darkmode bool

func Form() wid.Wid {
	return wid.Col(nil,
		wid.Edit(&name, nil, &wid.DefaultEdit),
		wid.Edit(&address, nil, &wid.DefaultEdit),
		wid.Label("MpqyM1", 2, &P, 1),
		wid.Label("MpqyM2", 2, &P, 1),
		wid.Label("Mpqy3", 1, &P, 1),
		wid.Label(strconv.Itoa(gpu.RedrawsPrSec), 1, &P, 1),
		wid.Checkbox("Darkmode", &darkmode, nil, ""),
		wid.Row(nil,
			wid.Label("Buttons", 2, &P, 4),
			wid.Elastic(),
			wid.Button("Cancel", CancelBtnClick, wid.PrimaryBtn, hint1),
			wid.Button("No", NoBtnClick, wid.PrimaryBtn, hint2),
			wid.Button("Yes", YesBtnClick, wid.PrimaryBtn, hint3),

		),
	)
}

func Draw() {
	// Calculate sizes
	form := Form()
	ctx := wid.Ctx{Rect: f32.Rect{X: 20, Y: 20, W: 400, H: 300}, Baseline: 0}
	gpu.Rect(ctx.Rect, 1, f32.Transparent, f32.LightBlue)
	_ = form(ctx)
}

var window *glfw.Window

func main() {

	window = gpu.InitWindow(0, 0, "Rounded rectangle demo", 1, f32.Blue)
	defer gpu.Shutdown()
	window.SetMouseButtonCallback(focus.MouseBtnCallback)
	window.SetCursorPosCallback(focus.MousePosCallback)
	window.SetKeyCallback(callback.KeyCallback)
	window.SetCharCallback(callback.CharCallback)
	window.SetScrollCallback(callback.ScrollCallback)

	LoadFonts()
	wid.LoadIcons()
	gpu.UpdateResolution()
	// for !window.ShouldClose() {
	// Test paint a shadow
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gpu.RoundedRect(f32.Rect{200, 200, 260, 260}, 12, 6, f32.Red, f32.Green)
	gpu.Shade(f32.Rect{220, 220, 50, 50}, 12, f32.Shadow, 8)
	_ = gpu.CaptureToFile("shadow.png", 180, 180, 300, 300)
	gpu.EndFrame(1)
	// }

	/*
		for !window.ShouldClose() {
			gpu.StartFrame()
			focus.Clickables = focus.Clickables[0:0]
			// Paint a red frame around the whole window
			gpu.Rect(gpu.WindowRect.Reduce(10), 2, f32.Transparent, f32.Red)
			// Draw the screen widgets
			Draw()
			ShowFonts(50, 400)
			ShowIcons(50, 350)
			// dialog.Show(nil)
			wid.ShowHint(nil)
			// focus.Update()
			gpu.EndFrame(30)
		} */
}
