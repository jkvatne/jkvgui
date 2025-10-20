package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/jkvatne/jkvgui/dialog"
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
	lightMode = true
	genders   = []string{"Male", "Female", "Both", "qyjpy", "Value5", "Value6", "Value7", "Value8", "Value9", "Value10", "Value11"}
	hint1     = "This is a hint word5 word6 word7 word8 qYyM9 qYyM10"
	hint2     = "This is a hint"
	hint3     = "This is a hint that is quite long, just to test word wrapping and hint location on screen. Should always be visible"
)

func createData() {
	for wno := range 16 {
		Persons[wno].gender = "Male"
		Persons[wno].name = "Ola Olsen" + strconv.Itoa(wno)
		Persons[wno].address = "Tulleveien " + strconv.Itoa(wno)
		Persons[wno].gender = "Male"
		Persons[wno].age = 10 + wno*5
		// We need a separate state for the scroller in each window.
		ss = append(ss, wid.ScrollState{Id: wno})
	}
}

func LightModeBtnClick() {
	lightMode = true
	theme.SetDefaultPallete(lightMode)
	slog.Info("LightModeBtnClick()")
	sys.Invalidate()
}

func DarkModeBtnClick() {
	lightMode = false
	theme.SetDefaultPallete(lightMode)
	slog.Info("DarkModeBtnClick()")
	sys.Invalidate()
}

func do() {
	dialog.Hide()
}

func DlgBtnClick() {
	w := dialog.YesNoDialog("Heading", "Some text", "Yes", "No", do, do)
	dialog.Show(&w)
	slog.Info("DlgBtnClick()")
}

func Monitor1BtnClick() {
	slog.Info("Monitor1BtnClick()")
	w := sys.GetCurrentContext()
	ms := sys.GetMonitors()
	if len(ms) > 1 {
		x, y, _, _ := ms[0].GetWorkarea()
		w.SetPos(x+30, y+40)
	}
}

func Monitor2BtnClick() {
	slog.Info("Monitor2BtnClick()")
	w := sys.GetCurrentContext()
	ms := sys.GetMonitors()
	if len(ms) > 1 {
		x, y, _, _ := ms[1].GetWorkarea()
		w.SetPos(x+30, y+40)
	}
}

func Maximize() {
	w := sys.GetCurrentContext()
	slog.Info("Maximize button handler start")
	sys.MaximizeWindow(w)
	slog.Info("Maximize button handler exit")
}

func Minimize() {
	slog.Info("Minimize()")
	w := sys.GetCurrentContext()
	sys.MinimizeWindow(w)
}

func FullScreen1() {
	slog.Info("FullScreen1()")
	w := sys.GetCurrentContext()
	ms := sys.GetMonitors()
	w.SetMonitor(ms[0], 0, 0, 1024, 768, 0)
}

func FullScreen2() {
	slog.Info("FullScreen2()")

	w := sys.GetCurrentContext()
	ms := sys.GetMonitors()
	w.SetMonitor(ms[1], 0, 0, 1024, 768, 0)
}

func Restore() {
	slog.Info("Restore()")
	w := sys.GetCurrentContext()
	w.SetMonitor(nil, 100, 100, int(750*1.5), int(400*1.5), 0)
}

func ExitBtnClick() {
	slog.Info("Exit()")
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
	sys.WinListMutex.RLock()
	defer sys.WinListMutex.RUnlock()
	return wid.Scroller(&ss[no],
		wid.Label(sys.WindowList[no].Name, wid.H1C),
		wid.Label("Use TAB to move focus, and Enter to save data", wid.I),
		wid.Label(fmt.Sprintf("Mouse pos = %0.0f, %0.0f", sys.WindowList[no].MousePos().X, sys.WindowList[no].MousePos().Y), wid.I),
		wid.Label("Extra text", wid.I),
		wid.Row(nil,
			wid.Btn("Maximize", nil, Maximize, nil, hint3),
			wid.Btn("Minimize", nil, Minimize, nil, hint3),
			wid.Btn("Full screen 1", nil, FullScreen1, nil, hint3),
			wid.Btn("Full screen 2", nil, FullScreen2, nil, hint3),
			wid.Btn("Windowed", nil, Restore, nil, hint3),
			wid.Btn("Monitor 1", nil, Monitor1BtnClick, nil, hint1),
			wid.Btn("Monitor 2", nil, Monitor2BtnClick, nil, hint1)),
		wid.Row(nil,
			wid.Btn("Show dialogue", nil, DlgBtnClick, nil, hint1),
			wid.Btn("DarkMode", nil, DarkModeBtnClick, nil, hint2),
			wid.Btn("LightMode", nil, LightModeBtnClick, nil, hint3),
			wid.Btn("Exit", nil, ExitBtnClick, nil, hint3),
		),
		wid.Edit(&Persons[no].name, "Name", nil, wid.DefaultEdit.Size(100, 200)),
		wid.Edit(&Persons[no].address, "Address", nil, wid.DefaultEdit.Size(100, 200)),
		wid.Combo(&Persons[no].gender, genders, "Gender", wid.DefaultCombo.Size(100, 200)),
		wid.Edit(&text, "Test", nil, nil),
		wid.Label("FPS="+fmt.Sprintf("%0.3f", sys.WindowList[no].Fps()), nil),
		wid.Checkbox("Darkmode (g)", &lightMode, nil, hint3),
		wid.Checkbox("Disabled", &disabled, nil, hint3),
		wid.Row(nil,
			wid.RadioButton("Dark", &mode, "Dark", nil),
			wid.RadioButton("Light", &mode, "Light", nil),
			wid.Switch("Dark mode", &lightMode, nil, nil, hint3),
		),
		wid.Label("14pt Buttons left adjusted (default row)", nil),
		wid.Row(nil,
			wid.Btn("Primary", gpu.Home, DoPrimary, wid.Filled, hint3),
			wid.Btn("Secondary", gpu.ContentOpen, DoSecondary, wid.Filled.Role(theme.Secondary), hint3),
			wid.Btn("TextBtn", gpu.ContentSave, DoTextBtn, wid.Text, hint3),
			wid.Btn("Outline", nil, DoOutlineBtn, wid.Outline, hint3),
			wid.Btn("", gpu.Home, DoHomeBtn, wid.Round, hint3),
		),
		wid.Label("Buttons with different fonts", nil),
		wid.Row(nil,
			wid.Btn("Primary", gpu.Home, DoPrimary, wid.Filled.Font(gpu.Normal10), hint3),
			wid.Btn("Secondary", gpu.ContentOpen, DoSecondary, wid.Filled.Role(theme.Secondary).Font(gpu.Normal12), hint3),
			wid.Btn("TextBtn", gpu.ContentSave, DoTextBtn, wid.Text.Font(gpu.Normal12), hint3),
			wid.Btn("Outline", nil, DoOutlineBtn, wid.Outline, hint3),
			wid.Btn("", gpu.Home, DoHomeBtn, wid.Round, hint3),
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
			wid.Btn("", gpu.Home, DoHomeBtn, wid.Round, hint3),
			wid.Elastic(),
		),
	)
}

var threaded = flag.Bool("threaded", true, "Set to test with one go-routine pr window")

func show(win *sys.Window) {
	win.StartFrame(theme.OnCanvas.Bg())
	wid.Show(Form(win.Wno))
	dialog.Display(win)
	win.EndFrame()
}

func Thread(win *sys.Window) {
	runtime.LockOSThread()
	for !win.Window.ShouldClose() {
		gpu.Mutex.Lock()
		show(win)
		gpu.Mutex.Unlock()
	}
	win.Destroy()
}

// Demo using threads
func main() {
	log.SetFlags(log.Lmicroseconds)
	sys.Init()
	defer sys.Shutdown()
	createData()
	win1 := sys.CreateWindow(100, 100, int(750), int(400), "Demo 1", 1, 1.0)
	win2 := sys.CreateWindow(200, 200, int(750*2.0), int(400*2.0), "Demo 2", 1, 2.0)
	if *threaded {
		go Thread(win1)
		go Thread(win2)
		for sys.WindowCount.Load() > 0 {
			time.Sleep(20 * time.Millisecond)
			gpu.Mutex.Lock()
			sys.PollEvents()
			gpu.Mutex.Unlock()
		}
		slog.Info("Exit Threaded()")
	} else {
		for sys.Running() {
			show(win1)
			show(win2)
			sys.PollEvents()
		}
	}
	slog.Info("Exit from main()")
}
