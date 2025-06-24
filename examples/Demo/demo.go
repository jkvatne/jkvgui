package main

import (
	"fmt"
	"github.com/jkvatne/jkvgui/dialog"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gl"
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
	genders   = []string{"Male", "Female", "Both", "qyjpy", "Value5", "Value6", "Value7", "Value8", "Value9", "Value10", "Value11"}
	name      = "Olger Olsen"
	address   = "Stavanger"
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
	dialog.CurrentDialogue = dialog.YesNoDialog("Heading", "Some text", "Yes", "No", do, do)
	slog.Info("Cancel Btn clicked")
}

func Monitor1BtnClick() {
	ms := sys.GetMonitors()
	x, y, w, h := ms[0].GetWorkarea()
	sys.WindowList[0].SetSize(w, h)
	sys.WindowList[0].SetPos(x, y)
}

func Monitor2BtnClick() {
	ms := sys.GetMonitors()
	if len(ms) > 1 {
		x, y, w, h := ms[1].GetWorkarea()
		sys.WindowList[0].SetSize(w, h)
		sys.WindowList[0].SetPos(x, y)
	}
}

func Maximize() {
	sys.MaximizeWindow(sys.WindowList[0])
}
func Minimize() {
	sys.MinimizeWindow(sys.WindowList[0])
}

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

var text = "abcdefg hijklmn opqrst"
var ss1 = &wid.ScrollState{}
var ss2 = &wid.ScrollState{}

func Form(ss *wid.ScrollState) wid.Wid {
	return wid.Scroller(ss,
		wid.Label("Edit user information", wid.H1C),
		wid.Label("Use TAB to move focus, and Enter to save data", wid.I),

		wid.Label(fmt.Sprintf("Mouse pos = %0.0f, %0.0f", sys.Pos().X, sys.Pos().Y), wid.I),
		wid.Label(fmt.Sprintf("Switch rect = %0.0f, %0.0f, %0.0f, %0.0f",
			wid.SwitchRect.X, wid.SwitchRect.Y, wid.SwitchRect.W, wid.SwitchRect.H), wid.I),
		wid.Label("Extra text", wid.I),
		wid.DisableIf(&disabled,
			wid.Row(nil,
				wid.Btn("Maximize", nil, Maximize, nil, ""),
				wid.Btn("Minimize", nil, Minimize, nil, ""),
				wid.Btn("Monitor 1", nil, Monitor1BtnClick, nil, hint1),
				wid.Btn("Monitor 2", nil, Monitor2BtnClick, nil, hint1),
				wid.Elastic(),
				wid.Btn("ShowDialogue dialogue", nil, DlgBtnClick, nil, hint1),
				wid.Btn("DarkMode", nil, DarkModeBtnClick, nil, hint2),
				wid.Btn("LightMode", nil, LightModeBtnClick, nil, hint3),
			),
		),
		wid.Edit(&name, "Name", nil, wid.DefaultEdit.Size(100, 200)),
		wid.Edit(&address, "Address", nil, wid.DefaultEdit.Size(100, 200)),
		wid.Combo(&gender, genders, "Gender", wid.DefaultCombo.Size(100, 200)),
		wid.Edit(&text, "Test", nil, nil),
		wid.Label("FPS="+strconv.Itoa(sys.RedrawsPrSec()), nil),
		wid.Checkbox("Darkmode (g)", &lightMode, nil, ""),
		wid.Checkbox("Disabled", &disabled, nil, ""),
		wid.Row(nil,
			wid.RadioButton("Dark", &mode, "Dark", nil),
			wid.RadioButton("Light", &mode, "Light", nil),
			wid.Switch("Dark mode", &lightMode, nil, nil, ""),
		),
		wid.Label("14pt Buttons left adjusted (default row)", nil),
		wid.Row(nil,
			wid.Btn("Primary", gpu.Home, set0, wid.Filled, ""),
			wid.Btn("Secondary", gpu.ContentOpen, set1, wid.Filled.Role(theme.Secondary), ""),
			wid.Btn("TextBtn", gpu.ContentSave, set2, wid.Text, ""),
			wid.Btn("Outline", nil, set3, wid.Outline, ""),
			wid.Btn("", gpu.Home, set4, wid.Round, ""),
		),
		wid.Label("Buttons with different fonts", nil),
		wid.Row(nil,
			wid.Btn("Primary", gpu.Home, nil, wid.Filled.Font(gpu.Normal10), ""),
			wid.Btn("Secondary", gpu.ContentOpen, nil, wid.Filled.Role(theme.Secondary).Font(gpu.Normal12), ""),
			wid.Btn("TextBtn", gpu.ContentSave, set2, wid.Text.Font(gpu.Normal12), ""),
			wid.Btn("Outline", nil, nil, wid.Outline, ""),
			wid.Btn("", gpu.Home, nil, wid.Round, ""),
		),
		wid.Label("Buttons with Elastic() between each", nil),
		wid.Row(nil,
			wid.Elastic(),
			wid.Btn("Primary", gpu.Home, nil, wid.Filled, "Primary"),
			wid.Elastic(),
			wid.Btn("Secondary", gpu.ContentOpen, nil, wid.Filled.Role(theme.Secondary), "Secondary"),
			wid.Elastic(),
			wid.Btn("TextBtn", gpu.ContentSave, nil, wid.Text, "Text"),
			wid.Elastic(),
			wid.Btn("Outline", nil, nil, wid.Outline, "Outline"),
			wid.Elastic(),
			wid.Btn("", gpu.Home, set5, wid.Round, ""),
			wid.Elastic(),
		),

	)
}

func main() {
	sys.CreateWindow(666, 400, "Rounded rectangle demo 1", 1, 1.0)
	// Initialize gl
	if err := gl.Init(); err != nil {
		panic("Initialization error for OpenGL: " + err.Error())
	}
	s := gl.GetString(gl.VERSION)
	if s == nil {
		panic("Could get Open-GL version")
	}
	version := gl.GoStr(s)
	slog.Info("OpenGL", "version", version)

	gpu.InitGpu()
	font.LoadDefaultFonts()
	gpu.LoadIcons()

	sys.CreateWindow(1000, 600, "Rounded rectangle demo 2", 2, 1.0)
	defer sys.Shutdown()
	gpu.InitGpu()
	font.LoadDefaultFonts()
	gpu.LoadIcons()

	gpu.UpdateResolution(0)
	gpu.UpdateResolution(1)

	for sys.Running(0) {
		for wno, _ := range sys.WindowList {
			sys.MakeContextCurrent(wno)
			sys.StartFrame(theme.Surface.Bg())
			// Paint a frame around the whole window
			gpu.RoundedRect(gpu.CurrentInfo.WindowRect.Reduce(1), 10, 1, f32.Transparent, f32.Red)
			// Draw form
			if wno == 1 {
				Form(ss1)(wid.NewCtx(wno))
			} else {
				Form(ss2)(wid.NewCtx(wno))
			}
			dialog.ShowDialogue()
			sys.EndFrame(wno)
		}
	}
}
