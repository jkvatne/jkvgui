package main

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jkvatne/jkvgui/button"
	"github.com/jkvatne/jkvgui/callback"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/font"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
	"log/slog"
	"strconv"
)

var (
	lightMode        = true
	gender    string = "Male"
	genders          = []string{"Male", "Female", "Both", "qyjpy"}
	name             = "Ole Petter Olsen"
	address          = "Mo i Rana"
	hint1            = "This is a hint word5 word6 word7 word8 qYyM9 qYyM10"
	hint2            = "This is a hint"
	hint3            = "This is a hint word5 word6 word7 word8 qYyM9 qYyM10 Word11 word12 jyword13"
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

var img1 *gpu.Img

func ShowIcons(x float32, y float32) {
	/*gpu.DrawIcon(x+25, y, 24, gpu.ArrowDropDown, f32.Black)
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
	gpu.DrawIcon(x+300, y, 24, gpu.NavigationArrowDropDown, f32.Black)*/
	gpu.DrawImage(x+375, y, 100, img1)
	gpu.DrawIcon(x+325, y, 50, gpu.NavigationArrowDropUp, f32.Black)

}

func ShowFonts(x float32, y float32) {
	font.Fonts[gpu.Normal].SetColor(f32.Black)
	/*
		font.Fonts[gpu.Normal].Printf(x, y, 2, 0, "24 Normal")
		font.Fonts[gpu.Bold].SetColor(f32.Black)
		font.Fonts[gpu.Bold].Printf(x, y+30, 2, 0, "24 Bold")
		font.Fonts[gpu.Mono].SetColor(f32.Black)
		font.Fonts[gpu.Mono].Printf(x, y+60, 2, 0, "24 Mono")
		font.Fonts[gpu.Italic].SetColor(f32.Black)
		font.Fonts[gpu.Italic].Printf(x, y+90, 2, 0, "24 Italic")

	*/
}

func Form() wid.Wid {
	return wid.Col(nil,
		wid.Label("Edit user information", wid.H1C),
		wid.Label("Use TAB to move focus, and Enter to save data", wid.I),
		wid.Edit("Name", &name, nil, &wid.DefaultEdit),
		wid.Edit("Address", &address, nil, nil),
		wid.Combo(&gender, genders, nil),
		wid.Label("MpqyM2", nil),
		wid.Label(strconv.Itoa(gpu.RedrawsPrSec), nil),
		wid.Checkbox("Darkmode", &lightMode, nil, ""),
		wid.Row(wid.Left,
			wid.Label("Buttons", nil),
			wid.Elastic(),
			button.Filled("Cancel", nil, CancelBtnClick, nil, hint1),
			button.Filled("No", nil, NoBtnClick, &button.Btn, hint2),
			button.Filled("Yes", nil, YesBtnClick, &button.Btn, hint3),

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
	callback.Initialize(window)
	img1, _ = gpu.NewImg("mook-logo.png")

	for !window.ShouldClose() {
		gpu.BackgroundColor(theme.Surface)
		gpu.StartFrame()
		// Paint a red frame around the whole window
		gpu.Rect(gpu.WindowRect.Reduce(2), 1, f32.Transparent, theme.PrimaryColor)
		// Draw the screen widgets
		// Draw()
		ShowIcons(0, 250)
		ShowFonts(20, 300)
		// dialog.Show(nil)
		wid.ShowHint(nil)
		// focus.Update()
		gpu.EndFrame(30)
	}
}
