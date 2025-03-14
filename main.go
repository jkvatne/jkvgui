package main

import (
	_ "embed"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jkvatne/jkvgui/dialog"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/wid"
	"log/slog"
	"strconv"
)

//go:embed gpu/fonts/Roboto-Thin.ttf
var Roboto100 []byte // 100

//go:embed gpu/fonts/Roboto-ExtraLight.ttf
var Roboto200 []byte // 200

//go:embed gpu/fonts/Roboto-Light.ttf
var Roboto300 []byte // 300

//go:embed gpu/fonts/Roboto-Regular.ttf
var Roboto400 []byte // 400

//go:embed gpu/fonts/Roboto-Medium.ttf
var Roboto500 []byte // 500

//go:embed gpu/fonts/Roboto-SemiBold.ttf
var Roboto600 []byte // 600

//go:embed gpu/fonts/Roboto-Bold.ttf
var Roboto700 []byte // 700

//go:embed gpu/fonts/Roboto-Bold.ttf
var Roboto800 []byte // 800

//go:embed gpu/fonts/Roboto-Bold.ttf
var Roboto900 []byte // 900

//go:embed gpu/fonts/RobotoMono-Regular.ttf
var RobotoMono400 []byte

//go:embed gpu/fonts/RobotoMono-Bold.ttf
var RobotoMono600 []byte

//go:embed gpu/fonts/RobotoMono-Light.ttf
var RobotoMono200 []byte

var P = f32.Padding{2, 2, 2, 2}

func YesBtnClick() {
	slog.Info("Yes Btn Clicked")
	gpu.UserScale *= 1.25
	gpu.UpdateSize(window, gpu.WindowWidthPx, gpu.WindowHeightPx)
}

func CancelBtnClick() {
	slog.Info("Cancel Btn clicked")
}

func NoBtnClick() {
	slog.Info("No Btn Click\n")
	gpu.UserScale /= 1.25
	gpu.UpdateSize(window, gpu.WindowWidthPx, gpu.WindowHeightPx)
}

var name = "Ole Petter Olsen"
var address = "Mo i Rana"
var hint1 = "This is a hint word5 word6 word7 word8 qYyM9 qYyM10"
var hint2 = "This is a hint"
var hint3 = "This is a hint word5 word6 word7 word8 qYyM9 qYyM10 Word11 word12 jyword13"

func ShowIcons() {
	wid.DrawIcon(50, 20, 24, wid.Home, f32.Blue)
	wid.DrawIcon(75, 20, 24, wid.BoxChecked, f32.Black)
	wid.DrawIcon(100, 20, 24, wid.BoxUnchecked, f32.Black)
	wid.DrawIcon(125, 20, 24, wid.RadioChecked, f32.Black)
	wid.DrawIcon(150, 20, 24, wid.RadioUnchecked, f32.Black)
	wid.DrawIcon(175, 20, 24, wid.ContentSave, f32.Black)
	wid.DrawIcon(200, 20, 24, wid.NavigationArrowDownward, f32.Black)
	wid.DrawIcon(225, 20, 24, wid.NavigationArrowUpward, f32.Black)
	wid.DrawIcon(250, 20, 24, wid.NavigationUnfoldMore, f32.Black)
	wid.DrawIcon(275, 20, 24, wid.NavigationArrowDropDown, f32.Black)
	wid.DrawIcon(300, 20, 24, wid.NavigationArrowDropUp, f32.Black)
}

// From freetype.go, line 263, Her c.dpi is allways 72.
// c.scale = fixed.Int26_6(0.5 + (c.fontSize * c.dpi * 64 / 72))
// size = fontsize  in pixels.

func LoadFonts() {
	_ = gpu.LoadFontBytes(Roboto200, 24, "Roboto", 200)
	_ = gpu.LoadFontBytes(Roboto400, 24, "Roboto", 400)
	_ = gpu.LoadFontBytes(Roboto600, 24, "Roboto", 600)
	_ = gpu.LoadFontBytes(RobotoMono200, 24, "RobotoMono", 200)
	_ = gpu.LoadFontBytes(RobotoMono400, 24, "RobotoMono", 400)
	_ = gpu.LoadFontBytes(RobotoMono600, 24, "RobotoMono", 600)
}

func ShowFonts() {
	gpu.Fonts[0].Printf(50, 100, 2, 0, "24 Roboto200")      // Thin
	gpu.Fonts[1].Printf(50, 130, 2, 0, "24 Roboto400")      // Regular
	gpu.Fonts[2].Printf(50, 160, 2, 0, "24 Roboto600")      // Bold
	gpu.Fonts[3].Printf(50, 190, 2, 0, "24 RobotoMono200")  // Thin
	gpu.Fonts[4].Printf(50, 220, 2, 0, "24 RobotoMono400")  // Regular
	gpu.Fonts[5].Printf(50, 250, 2, 0, "24 RobotoMono600")  // Bold
	gpu.Fonts[0].Printf(350, 100, 1, 0, "12 Roboto200")     // Thin
	gpu.Fonts[1].Printf(350, 130, 1, 0, "12 Roboto400")     // Regular
	gpu.Fonts[2].Printf(350, 160, 1, 0, "12 Roboto600")     // Bold
	gpu.Fonts[3].Printf(350, 190, 1, 0, "12 RobotoMono200") // Thin
	gpu.Fonts[4].Printf(350, 220, 1, 0, "12 RobotoMono400") // Regular
	gpu.Fonts[5].Printf(350, 250, 1, 0, "12 RobotoMono600") // Bold
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
	ctx := wid.Ctx{Rect: f32.Rect{X: 50, Y: 300, W: 400, H: 200}, Baseline: 0}
	gpu.Rect(ctx.Rect, 1, f32.Transparent, f32.LightBlue)
	_ = form(ctx)
}

var window *glfw.Window

func main() {

	window = gpu.InitWindow(0, 0, "Rounded rectangle demo", 1, f32.LightGrey)
	defer gpu.Shutdown()
	LoadFonts()
	slog.Info("hello, world")
	wid.LoadIcons()
	gpu.UpdateResolution()
	for !window.ShouldClose() {
		gpu.StartFrame()
		// Paint a red frame around the whole window
		gpu.Rect(gpu.WindowRect.Reduce(10), 2, f32.Transparent, f32.Red)
		// Draw the screen widgets
		Draw()
		ShowFonts()
		ShowIcons()
		dialog.Show(nil)
		wid.ShowHint(nil)
		gpu.EndFrame(30)
	}
}
