package main

import (
	"github.com/jkvatne/jkvgui/button"
	"github.com/jkvatne/jkvgui/callback"
	"github.com/jkvatne/jkvgui/dialog"
	"github.com/jkvatne/jkvgui/gpu"
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
var disabled bool

func set0() {
	// n1.WriteObject(0x4000, 0, 1, 0, "Set schedule 0")
}

func set1() {
	// n1.WriteObject(0x4000, 0, 1, 1, "Set schedule 1")
}

func set2() {
	// n1.WriteObject(0x4000, 0, 1, 2, "Set schedule 2")
}

func set3() {
	// n1.WriteObject(0x4000, 0, 1, 3, "Set schedule 3")
}

func set4() {
	// n1.WriteObject(0x4000, 0, 1, 4, "Set schedule 4")
}

func set5() {
}

func Form() wid.Wid {
	return wid.Col(nil,
		wid.Row(1,
			button.RadioButton("Dark", &mode, "Dark", nil),
			button.RadioButton("Light", &mode, "Light", nil),
			button.Switch(&on, nil, nil, ""),
			wid.Label("Switch", nil),
		),
		/*
			wid.Label("Edit user information", wid.H1C),
			wid.Label("Use TAB to move focus, and Enter to save data", wid.I),
			wid.Edit(&name, "", nil, nil),
			wid.Edit(&address, "", nil, nil),
			wid.Combo(&gender, genders, nil),
			wid.Label("MpqyM2", nil),
			wid.Label("FPS="+strconv.Itoa(gpu.RedrawsPrSec), nil),
			button.Checkbox("Darkmode (g)", &lightMode, nil, ""),
			button.Checkbox("Disabled", &disabled, nil, ""),
			wid.DisableIf(&disabled,
				wid.Row(1,
					wid.Elastic(),
					wid.Label("Buttons", wid.H1R),
					button.Filled("Show dialogue", nil, DlgBtnClick, nil, hint1),
					button.Filled("DarkMode", nil, DarkModeBtnClick, nil, hint2),
					button.Filled("LightMode", nil, LightModeBtnClick, nil, hint3),
				),
			),
			wid.Row(wid.Distribute,
				button.Filled("Primary", icon.Home, set0, button.Role(theme.Primary), ""),
				button.Filled("Secondary", icon.ContentOpen, set1, button.Role(theme.Secondary), ""),
				button.Filled("Surface", icon.ContentSave, set2, button.Role(theme.Surface), ""),
				button.Filled("Container", icon.RadioChecked, set3, button.Role(theme.SurfaceContainer), ""),
				button.Round(icon.Home, set5, nil, ""),
			),
			wid.Row(wid.Left,
				wid.Elastic(),
				button.Filled("Primary", icon.Home, set0, button.Role(theme.Primary), ""),
				wid.Elastic(),
				button.Filled("Secondary", icon.ContentOpen, set1, button.Role(theme.Secondary), ""),
				wid.Elastic(),
				button.Filled("Surface", icon.ContentSave, set2, button.Role(theme.Surface), ""),
				wid.Elastic(),
				button.Filled("Container", icon.RadioChecked, set3, button.Role(theme.SurfaceContainer), ""),
				wid.Elastic(),
				button.Round(icon.Home, set5, nil, ""),
				wid.Elastic(),
			),
		*/
	)
}

func main() {
	theme.SetDefaultPallete(lightMode)
	gpu.UserScale = 1.5
	window := gpu.InitWindow(0, 0, "Rounded rectangle demo", 1)
	defer gpu.Shutdown()
	callback.Initialize(window)
	for !window.ShouldClose() {
		gpu.StartFrame(theme.Surface.Bg())
		Form()(wid.Maximized())
		wid.ShowHint(nil)
		dialog.Show(nil)
		gpu.EndFrame(50)
	}
}
