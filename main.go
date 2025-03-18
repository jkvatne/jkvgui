package main

import (
	_ "embed"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jkvatne/jkvgui/callback"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/font"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/mouse"
	"github.com/jkvatne/jkvgui/theme"
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

//go:embed font/fonts/Roboto-MediumItalic.ttf
var RobotoItalic500 []byte

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
var lightMode = true
var gender string = "Male"
var genders = []string{"Male", "Female", "Both", "qyjpy"}

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

var name = "Ole Petter Olsen"
var address = "Mo i Rana"
var hint1 = "This is a hint word5 word6 word7 word8 qYyM9 qYyM10"
var hint2 = "This is a hint"
var hint3 = "This is a hint word5 word6 word7 word8 qYyM9 qYyM10 Word11 word12 jyword13"

func ShowIcons(x float32, y float32) {
	gpu.DrawIcon(x+50, y, 24, gpu.Home, f32.Grey)
	gpu.DrawIcon(x+75, y, 24, gpu.BoxChecked, f32.Grey)
	gpu.DrawIcon(x+100, y, 24, gpu.BoxUnchecked, f32.Grey)
	gpu.DrawIcon(x+125, y, 24, gpu.RadioChecked, f32.Grey)
	gpu.DrawIcon(x+150, y, 24, gpu.RadioUnchecked, f32.Grey)
	gpu.DrawIcon(x+175, y, 24, gpu.ContentSave, f32.Grey)
	gpu.DrawIcon(x+200, y, 24, gpu.NavigationArrowDownward, f32.Grey)
	gpu.DrawIcon(x+225, y, 24, gpu.NavigationArrowUpward, f32.Grey)
	gpu.DrawIcon(x+250, y, 24, gpu.NavigationUnfoldMore, f32.Grey)
	gpu.DrawIcon(x+275, y, 24, gpu.NavigationArrowDropDown, f32.Grey)
	gpu.DrawIcon(x+300, y, 24, gpu.NavigationArrowDropUp, f32.Grey)
	gpu.DrawIcon(x+325, y, 24, gpu.ArrowDropDown, f32.Grey)
	gpu.DrawIcon(x+350, y, 24, gpu.ContentOpen, f32.Grey)
}

// From freetype.go, line 263, Her c.dpi is allways 72.
// c.scale = fixed.Int26_6(0.5 + (c.fontSize * c.dpi * 64 / 72))
// size = fontsize  in pixels.
func LoadFonts() {
	font.LoadFontBytes(gpu.Normal, Roboto500, 24, "RobotoNormal", 400)
	font.LoadFontBytes(gpu.Bold, Roboto600, 24, "RobotoBold", 600)
	font.LoadFontBytes(gpu.Italic, RobotoItalic500, 24, "RobotoItalic", 500)
	font.LoadFontBytes(gpu.Mono, RobotoMono400, 24, "RobotoMono", 400)
}

func ShowFonts(x float32, y float32) {
	font.Fonts[gpu.Normal].SetColor(f32.Grey)
	font.Fonts[gpu.Normal].Printf(x, y, 2, 0, "24 Normal")
	font.Fonts[gpu.Bold].SetColor(f32.Grey)
	font.Fonts[gpu.Bold].Printf(x, y+30, 2, 0, "24 Bold")
	font.Fonts[gpu.Mono].SetColor(f32.Grey)
	font.Fonts[gpu.Mono].Printf(x, y+60, 2, 0, "24 Mono")
	font.Fonts[gpu.Italic].SetColor(f32.Grey)
	font.Fonts[gpu.Italic].Printf(x, y+90, 2, 0, "24 Italic")
}

func Form() wid.Wid {
	return wid.Col(nil,
		wid.Label("Edit user information", wid.H1),
		wid.Label("Use TAB to move focus, and Enter to save data", wid.I),
		wid.Edit(&name, nil, &wid.DefaultEdit),
		wid.Edit(&address, nil, nil),
		wid.Combo(&gender, genders, nil),
		wid.Label("MpqyM2", nil),
		wid.Label(strconv.Itoa(gpu.RedrawsPrSec), nil),
		wid.Checkbox("Darkmode", &lightMode, nil, ""),
		wid.Row(nil,
			wid.Label("Buttons", nil),
			wid.Elastic(),
			wid.Button("Cancel", CancelBtnClick, nil, hint1),
			wid.Button("No", NoBtnClick, &wid.OkBtn, hint2),
			wid.Button("Yes", YesBtnClick, &wid.OkBtn, hint3),

		),
	)
}

func Draw() {
	// Calculate sizes
	form := Form()
	ctx := wid.Ctx{Rect: f32.Rect{X: 20, Y: 20, W: 400, H: 300}, Baseline: 0}
	_ = form(ctx)
}

var window *glfw.Window

func main() {

	theme.SetDefaultPallete(lightMode)
	window = gpu.InitWindow(0, 0, "Rounded rectangle demo", 1)
	defer gpu.Shutdown()
	window.SetMouseButtonCallback(mouse.BtnCallback)
	window.SetCursorPosCallback(mouse.PosCallback)
	window.SetKeyCallback(callback.KeyCallback)
	window.SetCharCallback(callback.CharCallback)
	window.SetScrollCallback(callback.ScrollCallback)

	LoadFonts()
	gpu.LoadIcons()
	gpu.UpdateResolution()
	for !window.ShouldClose() {
		gpu.BackgroundColor(theme.Surface)
		gpu.StartFrame()
		// Paint a red frame around the whole window
		gpu.Rect(gpu.WindowRect.Reduce(2), 1, f32.Transparent, theme.PrimaryColor)
		// Draw the screen widgets
		Draw()
		ShowFonts(50, 400)
		ShowIcons(50, 350)
		// dialog.Show(nil)
		wid.ShowHint(nil)
		// focus.Update()
		gpu.EndFrame(30)
	}
}
