package main

import (
	"fmt"
	"log/slog"
	"math"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
)

type Person struct {
	name    string
	age     int
	gender  string
	address string
}

var Persons [16]Person

var (
	lightMode     = true
	genders       = []string{"Male", "Female", "Both", "qyjpy", "Value5", "Value6", "Value7", "Value8", "Value9", "Value10", "Value11"}
	hint1         = "This is a hint word5 word6 word7 word8 qYyM9 qYyM10"
	hint2         = "This is a hint"
	hint3         = "This is a hint word5 word6 word7 word8 qYyM9 qYyM10 Word11 word12 jyword13"
	CurrentDialog [99]*wid.Wid
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
	// TODO CurrentDialog[sys.CurrentWno] = nil
	// sys.CurrentInfo.DialogVisible = false
	// sys.CurrentInfo.SuppressEvents = false

}

func DlgBtnClick() {
	// TODO w := dialog.YesNoDialog("Heading", "Some text", "Yes", "No", do, do)
	// TODO CurrentDialog[sys.CurrentWno] = &w
	// TODO sys.CurrentInfo.DialogVisible = true
	slog.Info("Created dialog")
}

func Monitor1BtnClick() {
	// TODO ms := sys.GetMonitors()
	// TODO x, y, _, _ := ms[0].GetWorkarea()
	// TODO sys.WindowList[0].SetPos(x+30, y+40)
}

func Monitor2BtnClick() {
	ms := sys.GetMonitors()
	if len(ms) > 1 {
		// TODO x, y, _, _ := ms[1].GetWorkarea()
		// TODO sys.WindowList[0].SetPos(x+30, y+40)
	}
}

func Maximize() {
	// TODO sys.MaximizeWindow(sys.WindowList[0])
}

func Minimize() {
	// TODO ys.MinimizeWindow(sys.WindowList[0])
}

func FullScreen1() {
	// TODO ms := sys.GetMonitors()
	// TODO sys.WindowList[0].SetMonitor(ms[0], 0, 0, 1024, 768, 0)
}

func FullScreen2() {
	// TODO ms := sys.GetMonitors()
	// TODO sys.WindowList[0].SetMonitor(ms[1], 0, 0, 1024, 768, 0)
}

func Restore() {
	// TODO sys.WindowList[0].SetMonitor(nil, 100, 100, 1024, 768, 0)
}

func ExitBtnClick() {
	os.Exit(0)
}

var mode string
var disabled bool

func DoPrimary() {
	slog.Info("Primary clicked")
}

func DoSecondary() {
	slog.Info("Secondary clicked")
}

func DoTextBtn() {
	slog.Info("Textbtn clicked")
}

func DoOutlineBtn() {
	slog.Info("OutlineBtn clicked")
}

func DoHomeBtn() {
	slog.Info("HomeBtn clicked")
}

var text = "abcdefg hijklmn opqrst"
var ss []wid.ScrollState

func Form(no int) wid.Wid {
	return wid.Scroller(&ss[no],
		wid.Label("Edit user information", wid.H1C),
		wid.Label("Use TAB to move focus, and Enter to save data", wid.I),

		// TODO wid.Label(fmt.Sprintf("Mouse pos = %0.0f, %0.0f", sys.Pos().X, sys.Pos().Y), wid.I),
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
			wid.Btn("Show dialogue", nil, DlgBtnClick, nil, hint1),
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
			wid.Btn("Primary", gpu.Home, DoPrimary, wid.Filled, ""),
			wid.Btn("Secondary", gpu.ContentOpen, DoSecondary, wid.Filled.Role(theme.Secondary), ""),
			wid.Btn("TextBtn", gpu.ContentSave, DoTextBtn, wid.Text, ""),
			wid.Btn("Outline", nil, DoOutlineBtn, wid.Outline, ""),
			wid.Btn("", gpu.Home, DoHomeBtn, wid.Round, ""),
		),
		wid.Label("Buttons with different fonts", nil),
		wid.Row(nil,
			wid.Btn("Primary", gpu.Home, DoPrimary, wid.Filled.Font(gpu.Normal10), ""),
			wid.Btn("Secondary", gpu.ContentOpen, DoSecondary, wid.Filled.Role(theme.Secondary).Font(gpu.Normal12), ""),
			wid.Btn("TextBtn", gpu.ContentSave, DoTextBtn, wid.Text.Font(gpu.Normal12), ""),
			wid.Btn("Outline", nil, DoOutlineBtn, wid.Outline, ""),
			wid.Btn("", gpu.Home, DoHomeBtn, wid.Round, ""),
		),
		wid.Label("Buttons with Elastic() between each", nil),
		wid.Row(nil,
			wid.Elastic(),
			wid.Btn("Primary", gpu.Home, DoPrimary, wid.Filled, "Primary"),
			wid.Elastic(),
			wid.Btn("Secondary", gpu.ContentOpen, DoSecondary, wid.Filled.Role(theme.Secondary), "Secondary"),
			wid.Elastic(),
			wid.Btn("TextBtn", gpu.ContentSave, DoTextBtn, wid.Text, "Text"),
			wid.Elastic(),
			wid.Btn("Outline", nil, DoOutlineBtn, wid.Outline, "Outline"),
			wid.Elastic(),
			wid.Btn("", gpu.Home, DoHomeBtn, wid.Round, ""),
			wid.Elastic(),
		),
	)
}

func Thread(self *sys.Window) {
	runtime.LockOSThread()
	self.Window.MakeContextCurrent()
	for !self.Window.ShouldClose() {
		self.StartFrame(theme.Surface.Bg())
		// Paint a frame around the whole window
		gpu.RoundedRect(gpu.ClientRectDp.Reduce(1), 7, 1, f32.Transparent, f32.Red)
		// TODO if self.CurrentDialog != nil {
		// TODO self.SuppressEvents = true
		// TODO }
		// Draw form
		Form(self.Wno)(wid.NewCtx(self))
		// if CurrentDialog[wno] != nil && sys.CurrentInfo.DialogVisible {
		//	dialog.Show(CurrentDialog[wno])
		// }
		if self.SuppressEvents {
			fmt.Printf("sys.Info[0].SuppressEvents=true\n")
		}
		self.EndFrame()
		self.PollEvents()
	}
}

// Demo using threads
func main() {
	var winCount = 1
	fmt.Printf("\nTesting drawing windows in different goroutines\n")
	fmt.Printf("Window count %d\n", winCount)
	fmt.Printf("CPU count=%d\n", runtime.NumCPU())
	fmt.Printf("ProcCount=%d\n", runtime.GOMAXPROCS(0))

	sys.Init()
	defer sys.Shutdown()

	for wno := range winCount {
		Persons[wno].gender = "Male"
		Persons[wno].name = "Ola Olsen" + strconv.Itoa(wno)
		Persons[wno].address = "Tulleveien " + strconv.Itoa(wno)
		Persons[wno].gender = "Male"
		Persons[wno].age = 10 + wno*5
		// We need a separate state for the scroller in each window.
		ss = append(ss, wid.ScrollState{})
	}

	for wno := range winCount {
		userScale := float32(math.Pow(1.5, float64(wno)))
		_ = sys.CreateWindow(wno*100, wno*100, int(750*userScale), int(400*userScale),
			"Rounded rectangle demo "+strconv.Itoa(wno+1), wno+1, userScale)
	}
	for wno := range winCount {
		go Thread(sys.WindowList[wno])
	}
	for sys.WindowCount.Load() > 0 {
		time.Sleep(time.Millisecond * 100)
	}
}
