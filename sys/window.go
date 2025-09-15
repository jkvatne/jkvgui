package sys

import (
	"flag"
	"log/slog"
	"runtime"
	"time"

	"github.com/jkvatne/jkvgui/buildinfo"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gl"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/gpu/font"
	"github.com/jkvatne/jkvgui/theme"
	glfw "github.com/jkvatne/purego-glfw"
)

var OpenGlStarted bool

// CreateWindow initializes glfw and returns a Window to use.
// MonitorNo is 1 or 0 for the primary monitor, 2 for secondary monitor etc.
// Size is given in dp (device independent pixels)
// Windows typically fills the screen in one of the following ways:
// - Constant aspect ratio, use as much of screen as possible (h=10000, w=10000)
// - Full screen. (Maximized window) (w=0, h=0)
// - Small window of a given size, shrunk if the screen is not big enough (h=200, w=200)
// - Use full screen height, but limit width (h=0, w=800)
// - Use full screen width, but limit height (h=800, w=0)
func CreateWindow(x, y, w, h int, name string, monitorNo int, userScale float32) *Window {
	slog.Info("CreateWindow()", "Name", name, "Width", w, "Height", h)
	m := Monitors[max(0, min(monitorNo-1, len(Monitors)-1))]
	gpu.Gd.ScaleX, gpu.Gd.ScaleY = m.GetContentScale()
	if NoScaling {
		gpu.Gd.ScaleX, gpu.Gd.ScaleY = 1.0, 1.0
	}
	PosX, PosY, SizePxX, SizePxY := m.GetWorkarea()
	if w <= 0 {
		w = SizePxX
	} else {
		w = min(int(float32(w)*gpu.Gd.ScaleX), SizePxX)
	}
	if h <= 0 {
		h = SizePxY
	} else {
		h = min(int(float32(h)*gpu.Gd.ScaleY), SizePxY)
	}
	if x < 0 {
		PosX = PosX + (SizePxX-w)/2
	}
	if y < 0 {
		PosY = PosY + (SizePxY-h)/2
	}
	win := createInvisibleWindow(w, h, name, nil)
	gpu.Gd.ScaleX, gpu.Gd.ScaleY = win.GetContentScale()
	info := &Window{}
	info.Window = win
	info.LeftBtnUpTime = time.Now()
	lb, tb, rb, bb := info.Window.GetFrameSize()
	slog.Info("Borders", "lb", lb, "tb", tb, "rb", rb, "bb", bb)
	// Move the window to the selected monitor
	win.SetPos(PosX+x+lb, PosY+y+tb)
	win.SetSize(w+lb+rb, h+tb+bb)
	// Now we can update size and scaling
	info.UserScale = userScale
	// WinListMutex.Lock()
	WindowList = append(WindowList, info)
	wno := len(WindowList) - 1
	info.Wno = wno
	WindowCount.Add(1)
	// WinListMutex.Unlock()
	info.Name = name
	info.Window = win
	info.Trigger = make(chan bool, 1)
	SetupCursors()
	setCallbacks(win)
	win.Show()
	slog.Info("CreateWindow()",
		"ScaleX", f32.F2S(gpu.Gd.ScaleX, 2), ""+
			"ScaleY", f32.F2S(gpu.Gd.ScaleY, 2),
		"Monitor", monitorNo, "UserScale",
		f32.F2S(userScale, 2), "W", w, "H", h,
		"WDp", int(gpu.Gd.WidthDp),
		"HDp", int(gpu.Gd.HeightDp))

	win.Focus()
	return info
}

var BlinkFrequency = 2

func Blinker() {
	time.Sleep(time.Second * 2)
	for {
		time.Sleep(time.Second / time.Duration(BlinkFrequency*2))
		for wno := range WindowCount.Load() {
			// WinListMutex.Lock()
			b := WindowList[wno].BlinkState.Load()
			WindowList[wno].BlinkState.Store(!b)
			if WindowList[wno].Blinking.Load() {
				WindowList[wno].InvalidateCount.Add(1)
				PostEmptyEvent()
			}
			// WinListMutex.Unlock()
		}
	}
}

// Init will initialize the system.
func Init() {
	runtime.LockOSThread()
	flag.Parse()
	slog.SetLogLoggerLevel(slog.Level(*logLevel))
	InitializeProfiling()
	buildinfo.Get()
	if *maxFps == 0 {
		MinFrameDelay = 0
	} else {
		MinFrameDelay = time.Second / time.Duration(*maxFps)
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
	// go Blinker()
}

func (w *Window) UpdateSize() {
	width, height := w.Window.GetSize()
	if NoScaling {
		w.ScaleX = 1.0
		w.ScaleY = 1.0
	} else {
		w.ScaleX, w.ScaleY = w.Window.GetContentScale()
		w.ScaleX *= w.UserScale
		w.ScaleY *= w.UserScale
	}
	w.WidthPx = width
	w.HeightPx = height
	w.WidthDp = float32(width) / w.ScaleX
	w.HeightDp = float32(height) / w.ScaleY
}

func LoadOpenGl(w *Window) {
	slog.Info("gpu.Mutex.Lock in LoadOpenGl()")
	// gpu.Mutex.Lock()
	// defer gpu.Mutex.Unlock()
	if WindowCount.Load() == 0 {
		panic("LoadOpengl() must be called after at least one window is created")
	}
	w.MakeContextCurrent()
	if !OpenGlStarted {
		if err := gl.Init(); err != nil {
			panic("Initialization error for OpenGL: " + err.Error())
		}
		s := gl.GetString(gl.VERSION)
		if s == nil {
			panic("Could not get Open-GL version")
		}
		version := gl.GoStr(s)
		slog.Info("OpenGL", "version", version)
		font.LoadDefaultFonts(120)
		OpenGlStarted = true
	}
	gpu.InitGpu()
	slog.Info("gpu.Mutex.Unlock in LoadOpenGl()")
}

func GetCurrentContext() *glfw.Window {
	return glfw.GetCurrentContext()
}

func GetCurrentWindow() *Window {
	return GetWindow(GetCurrentContext())
}

func (w *Window) Running() bool {
	return !w.Window.ShouldClose()
}
