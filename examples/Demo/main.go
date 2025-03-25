package main

import (
	"github.com/jkvatne/jkvgui/button"
	"github.com/jkvatne/jkvgui/callback"
	"github.com/jkvatne/jkvgui/dialog"
	"github.com/jkvatne/jkvgui/gpu"
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

func LightModeBtnClick() {
	lightMode = true
	theme.SetDefaultPallete(lightMode)
	slog.Info("Yes Btn Clicked")
}

func DarkModeBtnClick() {
	lightMode = false
	theme.SetDefaultPallete(lightMode)
	slog.Info("No Btn Click\n")
}

func do() {
	dialog.Exit()
}

func DlgBtnClick() {
	dialog.Current = dialog.YesNoDialog("Heading", "Some text", "Yes", "No", do, do)
	slog.Info("Cancel Btn clicked")
}

var on bool
var mode string

func Form() wid.Wid {
	return wid.Col(nil,
		button.Switch(&on, nil, nil, ""),
		wid.Label("Edit user information", wid.H1C),
		wid.Label("Use TAB to move focus, and Enter to save data", wid.I),
		wid.Edit("", &name, nil, nil),
		wid.Edit("", &address, nil, nil),
		wid.Combo(&gender, genders, nil),
		wid.Label("MpqyM2", nil),
		wid.Label("FPS="+strconv.Itoa(gpu.RedrawsPrSec), nil),
		wid.Row(1,
			button.RadioButton("Dark", &mode, "Dark", nil),
			button.RadioButton("Light", &mode, "Light", nil),
		),
		button.Checkbox("Darkmode (g)", &lightMode, nil, ""),
		// func(ctx wid.Ctx) wid.Dim {
		//		return
		wid.Row(1,
			wid.Elastic(),
			wid.Label("Buttons", wid.H1R),
			button.Filled("Show dialogue", nil, DlgBtnClick, nil, hint1),
			button.Filled("DarkMode", nil, DarkModeBtnClick, &button.Btn, hint2),
			button.Filled("LightMode", nil, LightModeBtnClick, &button.Btn, hint3),
		),
		// (ctx.Enable(true))
		//	},
	)
}

func main() {
	theme.SetDefaultPallete(lightMode)
	window := gpu.InitWindow(0, 0, "Rounded rectangle demo", 1)
	defer gpu.Shutdown()
	callback.Initialize(window)
	for !window.ShouldClose() {
		gpu.StartFrame(theme.Surface.Bg())
		Form()(wid.Maximized())
		wid.ShowHint(nil)
		dialog.Show(nil)
		gpu.EndFrame(30)
	}
}
