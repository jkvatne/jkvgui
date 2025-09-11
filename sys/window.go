package sys

import (
	"flag"
	"log/slog"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jkvatne/jkvgui/buildinfo"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/gpu/font"
	"github.com/jkvatne/jkvgui/theme"
)

// Pr window global variables.
type WinInfoStruct = struct {
	Name               string
	Window             *Window
	Wno                int
	ScaleX             float32
	ScaleY             float32
	UserScale          float32
	Mutex              sync.Mutex
	InvalidateCount    atomic.Int32
	HintActive         bool
	Focused            bool
	BlinkState         atomic.Bool
	Blinking           atomic.Bool
	Cursor             int
	CurrentTag         interface{}
	MoveToNext         bool
	MoveToPrevious     bool
	ToNext             bool
	LastTag            interface{}
	SuppressEvents     bool
	MousePos           f32.Pos
	LeftBtnDown        bool
	LeftBtnReleased    bool
	Dragging           bool
	LeftBtnDownTime    time.Time
	LeftBtnUpTime      time.Time
	LeftBtnDoubleClick bool
	ScrolledY          float32
	DialogVisible      bool
}

var (
	WinInfo     []*WinInfoStruct
	WindowCount atomic.Int32
	CurrentInfo *WinInfoStruct
)

// CreateWindow initializes glfw and returns a Window to use.
// MonitorNo is 1 or 0 for the primary monitor, 2 for secondary monitor etc.
// Size is given in dp (device independent pixels)
// Windows typically fills the screen in one of the following ways:
// - Constant aspect ratio, use as much of screen as possible (h=10000, w=10000)
// - Full screen. (Maximized window) (w=0, h=0)
// - Small window of a given size, shrunk if the screen is not big enough (h=200, w=200)
// - Use full screen height, but limit width (h=0, w=800)
// - Use full screen width, but limit height (h=800, w=0)
func CreateWindow(x, y, w, h int, name string, monitorNo int, userScale float32) {
	slog.Info("CreateWindow()", "Name", name, "Width", w, "Height", h)
	m := Monitors[max(0, min(monitorNo-1, len(Monitors)-1))]
	ScaleX, ScaleY := m.GetContentScale()
	if NoScaling {
		ScaleX, ScaleY = 1.0, 1.0
	}
	PosX, PosY, SizePxX, SizePxY := m.GetWorkarea()
	if w <= 0 {
		w = SizePxX
	} else {
		w = min(int(float32(w)*ScaleX), SizePxX)
	}
	if h <= 0 {
		h = SizePxY
	} else {
		h = min(int(float32(h)*ScaleY), SizePxY)
	}
	if x < 0 {
		PosX = PosX + (SizePxX-w)/2
	}
	if y < 0 {
		PosY = PosY + (SizePxY-h)/2
	}
	win := createInvisibleWindow(w, h, name, nil)
	WindowCount.Add(1)
	ScaleX, ScaleY = win.GetContentScale()
	WindowList = append(WindowList, win)
	info := WinInfoStruct{}
	info.LeftBtnUpTime = time.Now()
	wno := len(WindowList) - 1
	CurrentWindow = WindowList[wno]
	lb, tb, rb, bb := CurrentWindow.GetFrameSize()
	slog.Info("Borders", "lb", lb, "tb", tb, "rb", rb, "bb", bb)
	// Move the window to the selected monitor
	win.SetPos(PosX+x+lb, PosY+y+tb)
	win.SetSize(w+lb+rb, h+tb+bb)
	// Now we can update size and scaling
	info.UserScale = userScale
	WinListMutex.Lock()
	WinInfo = append(WinInfo, &info)
	WinListMutex.Unlock()
	info.Name = name
	info.Wno = wno
	info.Window = (*Window)(win)
	CurrentInfo = WinInfo[wno]
	SetupCursors()
	win.MakeContextCurrent()
	win.Show()
	slog.Info("CreateWindow()",
		"ScaleX", f32.F2S(WinInfo[wno].ScaleX, 2), ""+
			"ScaleY", f32.F2S(WinInfo[wno].ScaleY, 2),
		"Monitor", monitorNo, "UserScale",
		f32.F2S(userScale, 2), "W", w, "H", h,
		"WDp", int(gpu.ClientRectDp.W),
		"HDp", int(gpu.ClientRectDp.H))

	setCallbacks(win)
	win.Focus()
	CurrentInfo.Focused = true
	gpu.InitGpu()
	gpu.ScaleX = CurrentInfo.ScaleX
	gpu.ScaleY = CurrentInfo.ScaleY
	UpdateSize(win)
	font.LoadDefaultFonts()
	gpu.LoadIcons()
}

func Running() bool {
	for wno, win := range WindowList {
		if win.ShouldClose() {
			win.Destroy()
			WinListMutex.Lock()
			if CurrentInfo == WinInfo[wno] {
				CurrentInfo = nil
			}
			WindowList = append(WindowList[:wno], WindowList[wno+1:]...)
			WinInfo = append(WinInfo[:wno], WinInfo[wno+1:]...)
			WindowCount.Add(-1)
			if CurrentInfo == nil && len(WinInfo) > 0 {
				CurrentInfo = WinInfo[0]
				CurrentWindow = WindowList[0]
			}
			WinListMutex.Unlock()
		}
	}
	return len(WindowList) > 0
}

var BlinkFrequency = 2

func Blinker() {
	time.Sleep(time.Second * 2)
	for {
		time.Sleep(time.Second / time.Duration(BlinkFrequency*2))
		for wno := range WindowCount.Load() {
			WinListMutex.Lock()
			b := WinInfo[wno].BlinkState.Load()
			WinInfo[wno].BlinkState.Store(!b)
			if WinInfo[wno].Blinking.Load() {
				WinInfo[wno].InvalidateCount.Add(1)
				PostEmptyEvent()
			}
			WinListMutex.Unlock()
		}
	}
}

// Init will initialize the system.
// The pallete is set to the default values
// The GLFW hints are set to the default values
// The connected monitors are put into the Monitors slice.
// Monitor info is printed to slog.
func Init() {
	runtime.LockOSThread()
	flag.Parse()
	slog.SetLogLoggerLevel(slog.Level(*logLevel))
	InitializeProfiling()
	buildinfo.Get()
	if *maxFps {
		MaxDelay = 0
	} else {
		MaxDelay = time.Second
	}
	if err := glfwInit(); err != nil {
		panic(err)
	}
	theme.SetDefaultPallete(true)
	SetDefaultHints()
	// Check all monitors and print size data
	Monitors = GetMonitors()
	// Select monitor as given, or use primary monitor.
	for i, m := range Monitors {
		SizeMmX, SizeMmY := m.GetPhysicalSize()
		mScaleX, mScaleY := m.GetContentScale()
		PosX, PosY, SizePxX, SizePxY := m.GetWorkarea()
		slog.Info("GetMonitors() for ", "Monitor", i+1,
			"WidthMm", SizeMmX, "HeightMm", SizeMmY,
			"WidthPx", SizePxX, "HeightPx", SizePxY, "PosX", PosX, "PosY", PosY,
			"ScaleX", f32.F2S(mScaleX, 3), "ScaleY", f32.F2S(mScaleY, 3))
	}
	go Blinker()
}
