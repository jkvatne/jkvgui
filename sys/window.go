package sys

import (
	"flag"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jkvatne/jkvgui/buildinfo"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/gpu/font"
)

// Pr window global variables.
type WinInfo = struct {
	Name                string
	Wno                 int
	ScaleX              float32
	ScaleY              float32
	UserScale           float32
	Mutex               sync.Mutex
	InvalidateCount     atomic.Int32
	HintActive          bool
	Focused             bool
	BlinkState          atomic.Bool
	Blinking            atomic.Bool
	Cursor              int
	CurrentTag          interface{}
	MoveToNext          bool
	MoveToPrevious      bool
	ToNext              bool
	LastTag             interface{}
	SuppressEvents      bool
	MousePos            f32.Pos
	LeftBtnDown         bool
	LeftBtnReleased     bool
	Dragging            bool
	LeftBtnDownTime     time.Time
	LeftBtnUpTime       time.Time
	LeftBtnDoubleClick  bool
	ScrolledY           float32
	WindowOuterRectPx   gpu.IntRect
	WindowContentRectDp f32.Rect
	DialogVisible       bool
}

func WindowHeightDp() float32 {
	return gpu.WindowContentRectDp.H
}

func WindowWidthDp() float32 {
	return gpu.WindowContentRectDp.W
}

var (
	Info        []*WinInfo
	WindowCount atomic.Int32
	CurrentInfo *WinInfo
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
	info := WinInfo{}
	info.LeftBtnUpTime = time.Now()
	wno := len(WindowList) - 1
	CurrentWindow = WindowList[wno]
	_, top, _, _ := CurrentWindow.GetFrameSize()
	// Move the window to the selected monitor
	win.SetPos(PosX+x, PosY+y+top)
	win.SetSize(w, h)
	// Now we can update size and scaling
	info.UserScale = userScale
	Info = append(Info, &info)
	info.Name = name
	info.Wno = wno
	CurrentInfo = Info[wno]
	UpdateSize(len(WindowList) - 1)
	SetupCursors()
	win.MakeContextCurrent()
	win.Show()
	slog.Info("CreateWindow()",
		"ScaleX", f32.F2S(Info[wno].ScaleX, 2), ""+
			"ScaleY", f32.F2S(Info[wno].ScaleY, 2),
		"Monitor", monitorNo, "UserScale",
		f32.F2S(userScale, 2), "W", w, "H", h,
		"WDp", int(gpu.WindowContentRectDp.W),
		"HDp", int(gpu.WindowContentRectDp.H))

	setCallbacks(win)
	win.Focus()
	CurrentInfo.Focused = true
	gpu.InitGpu()
	gpu.ScaleX = CurrentInfo.ScaleX
	gpu.ScaleY = CurrentInfo.ScaleY
	font.LoadDefaultFonts()
	gpu.LoadIcons()
}

func Running() bool {
	for wno, win := range WindowList {
		if win.ShouldClose() {
			win.Destroy()
			WinListMutex.Lock()
			WindowList = append(WindowList[:wno], WindowList[wno+1:]...)
			Info = append(Info[:wno], Info[wno+1:]...)
			WindowCount.Add(-1)
			WinListMutex.Unlock()
		}
	}
	return len(WindowList) > 0
}

func UpdateSize(wno int) {
	width, height := WindowList[wno].GetSize()
	if NoScaling {
		Info[wno].ScaleX = 1.0
		Info[wno].ScaleY = 1.0
	} else {
		Info[wno].ScaleX, Info[wno].ScaleY = WindowList[wno].GetContentScale()
		Info[wno].ScaleX *= Info[wno].UserScale
		Info[wno].ScaleY *= Info[wno].UserScale
	}
	gpu.WindowOuterRectPx = gpu.IntRect{0, 0, width, height}
	gpu.WindowContentRectDp = f32.Rect{
		W: float32(width) / Info[wno].ScaleX,
		H: float32(height) / Info[wno].ScaleY}
	gpu.ScaleX, gpu.ScaleY = Info[wno].ScaleX, Info[wno].ScaleY
}

func init() {
	flag.Parse()
	slog.SetLogLoggerLevel(slog.Level(*logLevel))
	InitializeProfiling()
	buildinfo.Get()
}
