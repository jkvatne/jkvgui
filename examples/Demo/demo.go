package main

import (
	"fmt"
	"github.com/jkvatne/jkvgui/dialog"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
	"log/slog"
	"os"
	"strconv"
)

type Person struct {
	name    string
	age     int
	gender  string
	address string
}

var Persons [16]Person

var (
	lightMode = true
	genders   = []string{"Male", "Female", "Both", "qyjpy", "Value5", "Value6", "Value7", "Value8", "Value9", "Value10", "Value11"}
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
	x, y, _, _ := ms[0].GetWorkarea()
	sys.WindowList[0].SetPos(x+30, y+40)
}

func Monitor2BtnClick() {
	ms := sys.GetMonitors()
	if len(ms) > 1 {
		x, y, _, _ := ms[1].GetWorkarea()
		sys.WindowList[0].SetPos(x+30, y+40)
	}
}

func Maximize() {
	sys.MaximizeWindow(sys.WindowList[0])
}

func Minimize() {
	sys.MinimizeWindow(sys.WindowList[0])
}

func FullScreen1() {
	ms := sys.GetMonitors()
	sys.WindowList[0].SetMonitor(ms[0], 0, 0, 1024, 768, 0)
}

func FullScreen2() {
	ms := sys.GetMonitors()
	sys.WindowList[0].SetMonitor(ms[1], 0, 0, 1024, 768, 0)
}

func Restore() {
	sys.WindowList[0].SetMonitor(nil, 100, 100, 1024, 768, 0)
}

func ExitBtnClick() {
	os.Exit(0)
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
var ss []wid.ScrollState

func Form(no int) wid.Wid {
	return wid.Scroller(&ss[no],
		wid.Label("Edit user information", wid.H1C),
		wid.Label("Use TAB to move focus, and Enter to save data", wid.I),

		wid.Label(fmt.Sprintf("Mouse pos = %0.0f, %0.0f", sys.Pos().X, sys.Pos().Y), wid.I),
		wid.Label(fmt.Sprintf("Switch rect = %0.0f, %0.0f, %0.0f, %0.0f",
			wid.SwitchRect.X, wid.SwitchRect.Y, wid.SwitchRect.W, wid.SwitchRect.H), wid.I),
		wid.Label("Extra text", wid.I),
		wid.Row(nil,
			wid.Btn("Maximize", nil, Maximize, nil, ""),
			wid.Btn("Minimize", nil, Minimize, nil, ""),
			wid.Btn("Full screen 1", nil, FullScreen1, nil, ""),
			wid.Btn("Full screen 2", nil, FullScreen2, nil, ""),
			wid.Btn("Windowed", nil, Restore, nil, ""),
			wid.Btn("Monitor 1", nil, Monitor1BtnClick, nil, hint1),
			wid.Btn("Monitor 2", nil, Monitor2BtnClick, nil, hint1)),
		wid.Row(nil,
			wid.Btn("ShowDialogue dialogue", nil, DlgBtnClick, nil, hint1),
			wid.Btn("DarkMode", nil, DarkModeBtnClick, nil, hint2),
			wid.Btn("LightMode", nil, LightModeBtnClick, nil, hint3),
			wid.Btn("Exit", nil, ExitBtnClick, nil, ""),
		),
		wid.Edit(&Persons[no].name, "Name", nil, wid.DefaultEdit.Size(100, 200)),
		wid.Edit(&Persons[no].address, "Address", nil, wid.DefaultEdit.Size(100, 200)),
		wid.Combo(&Persons[no].gender, genders, "Gender", wid.DefaultCombo.Size(100, 200)),
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
	sys.Init()
	defer sys.Shutdown()

	var winCount = 1
	for wno := range winCount {
		sys.CreateWindow(wno*100, wno*100, 600, 400,
			"Rounded rectangle demo "+strconv.Itoa(wno+1), 2, 1.0)
		Persons[wno].gender = "Male"
		Persons[wno].name = "Ola Olsen" + strconv.Itoa(wno)
		Persons[wno].address = "Tulleveien " + strconv.Itoa(wno)
		Persons[wno].gender = "Male"
		Persons[wno].age = 10 + wno*5
	}

	for range winCount {
		ss = append(ss, wid.ScrollState{})
	}
	contentRect := gpu.CurrentInfo.WindowContentRectDp

	for sys.Running() {
		for sys.CurrentWno, _ = range sys.WindowList {
			sys.StartFrame(theme.Surface.Bg())
			// Paint a frame around the whole window
			contentRect = gpu.CurrentInfo.WindowContentRectDp
			gpu.RoundedRect(contentRect.Reduce(1), 7, 1, f32.Transparent, f32.Red)
			// Draw form
			Form(sys.CurrentWno)(wid.NewCtx())
			dialog.ShowDialogue()
			sys.EndFrame()
		}
		sys.PollEvents()
	}
}
