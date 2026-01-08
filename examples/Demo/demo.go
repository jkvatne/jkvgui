package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"runtime"
	"strconv"
	"sync"
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

var (
	Persons    [16]Person
	lightMode1 = true
	lightMode2 = true
	dummy      = false
	genders    = []string{"Male", "Female", "Both", "Any", "Value5", "Value6", "Value7", "Value8", "Value9", "Value10", "Value11", "Value12", "Value13", "Value14", "Value15"}
	hint1      = "This is a hint word5 word6 word7 word8 qYyM9 qYyM10"
	hint2      = "This is a hint"
	hint3      = "This is a hint that is quite long, just to test word wrapping and hint location on screen. Should always be visible"
	mode       string
	disabled   bool
	text       = "just some text"
	ss         []wid.ScrollState
	threaded   = flag.Bool("threaded", false, "Set to test with one go-routine pr window")
	n          = flag.Int("n", 2, "The number of windows used")
	Mutex      sync.Mutex
	progress   float32
)

func createData() {
	for wno := range 16 {
		Persons[wno].gender = "Male"
		Persons[wno].name = "Ola Olsen" + strconv.Itoa(wno)
		Persons[wno].address = "Skogveien " + strconv.Itoa(wno)
		Persons[wno].gender = "Male"
		Persons[wno].age = 10 + wno*5
		// We need a separate state for the scroller in each window.
		ss = append(ss, wid.ScrollState{Id: wno})
	}
}

func LightModeBtnClick() {
	lightMode1 = true
	theme.SetDefaultPalette(lightMode1)
	slog.Info("LightModeBtnClick()")
	sys.Invalidate()
}

func DarkModeBtnClick() {
	lightMode1 = false
	theme.SetDefaultPalette(lightMode1)
	slog.Info("DarkModeBtnClick()")
	sys.Invalidate()
}

func doYes() {
	slog.Info("Used clicked yes")
	dialog.Hide()
}

func doNo() {
	slog.Info("Used clicked no")
	dialog.Hide()
}

func DlgBtnClick() {
	w := dialog.YesNoDialog("Heading", "Some text", "Yes", "No", doYes, doNo)
	dialog.Show(&w, doYes, DlgBtnClick)
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
		x, y, _, _ := ms[min(1, len(ms))].GetWorkarea()
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
	w.SetMonitor(ms[min(1, len(ms))], 0, 0, 1024, 768, 0)
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

func DoPrimary1() {
	slog.Info("Primary 1 clicked")
}

func DoSecondary1() {
	slog.Info("Secondary 1 clicked")
}

func DoTextBtn1() {
	slog.Info("Text button 1 clicked")
}

func DoOutlineBtn1() {
	slog.Info("Outline button 1 clicked")
}

func DoHomeBtn1() {
	slog.Info("Home button 1 clicked")
}

func DoPrimary2() {
	slog.Info("Primary 2 clicked")
}

func DoSecondary2() {
	slog.Info("Secondary 2 clicked")
}

func DoTextBtn2() {
	slog.Info("Text button 2 clicked")
}

func DoOutlineBtn2() {
	slog.Info("Outline button 2 clicked")
}
func DoHomeBtn2() {
	slog.Info("Home button 2 clicked")
}

func Form(no int32) wid.Wid {
	sys.WinListMutex.RLock()
	defer sys.WinListMutex.RUnlock()
	return wid.Scroller(&ss[no],
		wid.Label(sys.WindowList[no].Name, wid.H1C),
		wid.Label("Use TAB to move focus, and Enter or space to click button", wid.L.Font(gpu.Normal10)),
		wid.Label(fmt.Sprintf("MousePos = %5.0f, %5.0f      FPS=%0.3f", sys.WindowList[no].MousePos().X, sys.WindowList[no].MousePos().Y, sys.WindowList[no].Fps()), nil),
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
		wid.ProgressBar(progress, nil),
		wid.Label("Fixed size edits with label size=100 and edit size=200", wid.L.Font(gpu.Normal10).Top(12)),
		wid.Edit(&Persons[no].name, "Name", nil, wid.DefaultEdit.Size(100, 200)),
		wid.Edit(&Persons[no].address, "Address", nil, wid.DefaultEdit.Size(100, 200)),
		wid.Combo(&Persons[no].gender, genders, "Gender", wid.DefaultCombo.Size(100, 200)),
		wid.Row(nil,
			wid.Checkbox("Darkmode", &lightMode1, nil, nil, hint3),
			wid.Checkbox("Disabled", &disabled, nil, nil, hint3),
			wid.RadioButton("Dark", &mode, "Dark", nil),
			wid.RadioButton("Light", &mode, "Light", nil),
			wid.Flex(),
			wid.Switch("Light mode", &lightMode2, nil, nil, hint3),
			wid.Switch("Dummy", &dummy, nil, nil, hint2),
		),
		wid.Label("Buttons with different fonts, left aligned", wid.L.Font(gpu.Normal10).Top(12)),
		wid.Row(nil,
			wid.Btn("Primary10", gpu.Home, DoPrimary1, wid.Filled.Font(gpu.Normal10), hint3),
			wid.Btn("Secondary12", gpu.ContentOpen, DoSecondary1, wid.Filled.Role(theme.Secondary).Font(gpu.Normal12), hint3),
			wid.Btn("TextBtn12", gpu.ContentSave, DoTextBtn1, wid.Text.Font(gpu.Normal12), hint3),
			wid.Btn("Outline14", nil, DoOutlineBtn1, wid.Outline, hint3),
			wid.Btn("", gpu.Home, DoHomeBtn1, wid.Round, hint3),
		),
		wid.Label("Buttons with Flex() between each", wid.L.Font(gpu.Normal10).Top(12)),
		wid.Row(nil,
			wid.Flex(),
			wid.Btn("Primary", gpu.Home, DoPrimary2, wid.Filled, "Primary"),
			wid.Flex(),
			wid.Btn("Secondary", gpu.ContentOpen, DoSecondary2, wid.Filled.Role(theme.Secondary), "Secondary"),
			wid.Flex(),
			wid.Btn("TextBtn", gpu.ContentSave, DoTextBtn2, wid.Text, "Text"),
			wid.Flex(),
			wid.Btn("Outline", nil, DoOutlineBtn2, wid.Outline, "Outline"),
			wid.Flex(),
			wid.Btn("", gpu.Home, DoHomeBtn2, wid.Round, hint3),
			wid.Flex(),
		),
	)
}

func show(wno int32) {
	if wno < sys.WindowCount.Load() {
		sys.WindowList[wno].StartFrame()
		wid.Show(Form(wno))
		dialog.Display(sys.WindowList[wno])
		sys.WindowList[wno].EndFrame()
	}
}

func Thread(wno int32) {
	runtime.LockOSThread()
	for !sys.WindowList[wno].Window.ShouldClose() {
		// We have to make sure only one thread at a time is using glfw.
		Mutex.Lock()
		show(wno)
		Mutex.Unlock()
	}
	slog.Info("Exit", "Thread", wno)
}

// Demo using threads
func main() {
	// Format slog output with time including microseconds, but no date.
	log.SetFlags(log.Lmicroseconds)
	slog.Info("Demo")
	sys.Init()
	sys.MinFrameDelay = time.Second / 20
	sys.MaxFrameDelay = time.Second / 2
	slog.SetLogLoggerLevel(slog.LevelInfo)
	defer sys.Shutdown()
	createData()
	sys.CreateWindow(100, 100, 1400, 1200, "Demo 1", 2, 2.0)
	if *n > 1 {
		sys.CreateWindow(200, 200, 750, 400, "Demo 2", 1, 1.0)
	}
	started := time.Now()
	if *threaded {
		go Thread(0)
		if *n > 1 {
			go Thread(1)
		}
		for sys.Running() {
			// We have to make sure only one thread at a time is using glfw.
			Mutex.Lock()
			sys.PollEvents()
			Mutex.Unlock()
		}
		slog.Info("Exit threaded demo")
	} else {
		for sys.Running() {
			show(0)
			if *n > 1 {
				show(1)
			}
			sys.PollEvents()
			progress = float32(time.Since(started).Seconds() / 10)
		}
		slog.Info("Exit non-threaded demo ")
	}
}
