package main

import (
	"github.com/jkvatne/jkvgui/btn"
	"github.com/jkvatne/jkvgui/dialog"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/icon"
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
		wid.Label("Edit user information", wid.H1C),
		wid.Label("Use TAB to move focus, and Enter to save data", wid.I),
		wid.Edit(&name, "", nil, nil),
		wid.Edit(&address, "", nil, nil),
		wid.Combo(&gender, genders, "Gender", wid.DefaultCombo.Size(100)),
		wid.Label("MpqyM2", nil),
		wid.Label("FPS="+strconv.Itoa(sys.RedrawsPrSec), nil),
		btn.Checkbox("Darkmode (g)", &lightMode, nil, ""),
		btn.Checkbox("Disabled", &disabled, nil, ""),
		wid.Row(nil,
			btn.RadioButton("Dark", &mode, "Dark", nil),
			btn.RadioButton("Light", &mode, "Light", nil),
			btn.Switch("Dark mode", &lightMode, nil, nil, ""),
		),
		wid.Row(nil,
			btn.Btn("Primary", icon.Home, set0, btn.Role(theme.Primary), ""),
			btn.Btn("Secondary", icon.ContentOpen, set1, btn.Role(theme.Secondary), ""),
			btn.Btn("Surface", icon.ContentSave, set2, btn.Role(theme.Surface), ""),
			btn.Btn("Container", icon.RadioChecked, set3, btn.Role(theme.SurfaceContainer), ""),
			btn.Btn("", icon.Home, set5, &btn.Round, ""),
		),
		wid.Row(nil,
			wid.Elastic(),
			btn.Btn("Primary", icon.Home, set0, btn.Role(theme.Primary), ""),
			wid.Elastic(),
			btn.Btn("Secondary", icon.ContentOpen, set1, btn.Role(theme.Secondary), ""),
			wid.Elastic(),
			btn.Btn("Surface", icon.ContentSave, set2, btn.Role(theme.Surface), ""),
			wid.Elastic(),
			btn.Btn("Container", icon.RadioChecked, set3, btn.Role(theme.SurfaceContainer), ""),
			wid.Elastic(),
			btn.Btn("", icon.Home, set5, &btn.Round, ""),
			wid.Elastic(),
		),
		wid.DisableIf(&disabled,
			wid.Row(nil,
				wid.Elastic(),
				wid.Label("Buttons", wid.H1R),
				btn.Btn("Show dialogue", nil, DlgBtnClick, nil, hint1),
				btn.Btn("DarkMode", nil, DarkModeBtnClick, nil, hint2),
				btn.Btn("LightMode", nil, LightModeBtnClick, nil, hint3),
			),
		),
	)
}

func main() {
	// Setting this true will draw a light blue frame around widgets.
	gpu.DebugWidgets = false
	theme.SetDefaultPallete(lightMode)
	// This is a user defined zoom level. Can be used to set higher
	// zoom factor than normal. Nice for people with reduced vision.
	// This value can be changed by using ctrl+scrollwheel
	gpu.UserScale = 1.5
	window := gpu.InitWindow(0, 0, "Rounded rectangle demo", 1)
	defer gpu.Shutdown()
	sys.Initialize(window, 14)
	for !window.ShouldClose() {
		sys.StartFrame(theme.Surface.Bg())
		Form()(wid.NewCtx())
		wid.ShowHint(nil)
		dialog.Show(nil)
		sys.EndFrame(50)
	}
}
