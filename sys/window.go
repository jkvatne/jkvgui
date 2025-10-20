package sys

import (
	"flag"
	"image"
	"log/slog"
	"runtime"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/jkvatne/jkvgui/buildinfo"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gl"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/gpu/font"
	"github.com/jkvatne/jkvgui/theme"
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
func CreateWindow(x, y, w, h int, name string, monitorNo int, userScale float32) *Window {
	slog.Debug("CreateWindow()", "Name", name, "Width", w, "Height", h)
	win := &Window{}
	m := Monitors[max(0, min(monitorNo-1, len(Monitors)-1))]
	win.Gd.ScaleX, win.Gd.ScaleY = m.GetContentScale()
	if NoScaling {
		win.Gd.ScaleX, win.Gd.ScaleY = 1.0, 1.0
	}
	PosX, PosY, SizePxX, SizePxY := m.GetWorkarea()
	if w <= 0 {
		w = SizePxX
	} else {
		w = min(int(float32(w)*win.Gd.ScaleX), SizePxX)
	}
	if h <= 0 {
		h = SizePxY
	} else {
		h = min(int(float32(h)*win.Gd.ScaleY), SizePxY)
	}
	if x < 0 {
		PosX = PosX + (SizePxX-w)/2
	}
	if y < 0 {
		PosY = PosY + (SizePxY-h)/2
	}
	win.Window = createInvisibleWindow(w, h, name, nil)
	win.Gd.ScaleX, win.Gd.ScaleY = win.Window.GetContentScale()
	win.LeftBtnUpTime = time.Now()
	lb, tb, rb, bb := win.Window.GetFrameSize()
	slog.Debug("Borders", "lb", lb, "tb", tb, "rb", rb, "bb", bb)
	// Move the window to the selected monitor
	win.Window.SetPos(PosX+x+lb, PosY+y+tb)
	// Now we can update size and scaling
	win.UserScale = userScale
	win.Window.SetSize(w+lb+rb, h+tb+bb)
	win.UpdateSize(w+lb+rb, h+tb+bb)
	WinListMutex.Lock()
	WindowList = append(WindowList, win)
	wno := len(WindowList) - 1
	WinListMutex.Unlock()
	win.Wno = wno
	WindowCount.Add(1)
	win.Name = name
	win.Trigger = make(chan bool, 1)
	SetupCursors()
	setCallbacks(win.Window)
	win.Window.Show()
	slog.Debug("CreateWindow()",
		"ScaleX", f32.F2S(win.Gd.ScaleX, 2, 4), ""+
			"ScaleY", f32.F2S(win.Gd.ScaleY, 2, 4),
		"Monitor", monitorNo, "UserScale",
		f32.F2S(userScale, 2, 4), "W", w, "H", h,
		"WDp", int(win.WidthDp),
		"HDp", int(win.HeightDp))

	win.Window.Focus()
	LoadOpenGl(win)
	slog.Debug("CreateWindow() done", "Name", name)
	return win
}

var BlinkFrequency = 2
var BlinkState atomic.Bool

func Blinker() {
	for {
		time.Sleep(time.Second / time.Duration(BlinkFrequency*2))
		BlinkState.Store(!BlinkState.Load())
		Invalidate()
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
		slog.Debug("GetMonitors() for ", "Monitor", i+1,
			"WidthMm", SizeMmX, "HeightMm", SizeMmY,
			"WidthPx", SizePxX, "HeightPx", SizePxY, "PosX", PosX, "PosY", PosY,
			"ScaleX", f32.F2S(mScaleX, 3, 4), "ScaleY", f32.F2S(mScaleY, 3, 4))
	}
	// go Blinker()
}

func (w *Window) UpdateSizeDp() {
	w.Gd.ScaleX, w.Gd.ScaleY = w.Window.GetContentScale()
	w.Gd.ScaleX *= w.UserScale
	w.Gd.ScaleY *= w.UserScale
	if NoScaling {
		w.Gd.ScaleX = 1.0
		w.Gd.ScaleY = 1.0
	}
	w.WidthDp = float32(w.WidthPx) / w.Gd.ScaleX
	w.HeightDp = float32(w.HeightPx) / w.Gd.ScaleY
}

func (w *Window) UpdateSize(width, height int) {
	w.WidthPx = width
	w.HeightPx = height
	w.UpdateSizeDp()
}

func LoadOpenGl(w *Window) {
	w.MakeContextCurrent()
	if !OpenGlStarted {
		OpenGlStarted = true
		if err := gl.Init(); err != nil {
			panic("Initialization error for OpenGL: " + err.Error())
		}
		s := gl.GetString(gl.VERSION)
		if s == nil {
			panic("Could not get Open-GL version")
		}
		version := gl.GoStr(s)
		slog.Debug("OpenGL", "version", version)
	}
	w.Gd.InitGpu()
	font.LoadDefaultFonts(font.DefaultDpi * w.Gd.ScaleX)
	gpu.LoadIcons()
	DetachCurrentContext()
}

func GetCurrentWindow() *Window {
	return GetWindow(GetCurrentContext())
}

func (w *Window) ClientRectDp() f32.Rect {
	return f32.Rect{0, 0, w.WidthDp, w.HeightDp}
}

func (w *Window) Defer(f func()) {
	for _, g := range w.DeferredFunctions {
		if &f == &g {
			return
		}
	}
	w.DeferredFunctions = append(w.DeferredFunctions, f)
}

func (w *Window) RunDeferred() {
	for _, f := range w.DeferredFunctions {
		f()
	}
	w.DeferredFunctions = w.DeferredFunctions[0:0]
}

func (w *Window) MakeContextCurrent() {
	w.Window.MakeContextCurrent()
}

func (w *Window) SetCursor(c int) {
	w.Cursor = c
}

// Invalidate will trigger all windows to paint their contenst
func Invalidate() {
	WinListMutex.RLock()
	defer WinListMutex.RUnlock()
	for _, w := range WindowList {
		w.Invalidate()
	}
}

func Running() bool {
	for wno, win := range WindowList {
		if win.Window.ShouldClose() {
			win.Window.Destroy()
			WinListMutex.RLock()
			defer WinListMutex.RUnlock()
			WindowList = append(WindowList[:wno], WindowList[wno+1:]...)
			WindowCount.Add(-1)
			return len(WindowList) > 0
		}
	}
	return true
}

func CaptureToFile(win *Window, filename string, x, y, w, h int) error {
	img := Capture(win, x, y, w, h) // 1057-300
	return gpu.SaveImage(filename, img)
}

func Capture(win *Window, x, y, w, h int) *image.RGBA {
	x = int(float32(x) * win.Gd.ScaleX)
	y = int(float32(y) * win.Gd.ScaleY)
	w = int(float32(w) * win.Gd.ScaleX)
	h = int(float32(h) * win.Gd.ScaleY)
	y = win.HeightPx - h - y
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	gl.PixelStorei(gl.PACK_ALIGNMENT, 1)
	gl.ReadPixels(int32(x), int32(y), int32(w), int32(h),
		gl.RGBA, gl.UNSIGNED_BYTE, unsafe.Pointer(&img.Pix[0]))
	gpu.GetErrors("Capture")
	//  Upside down
	for y := 0; y < h/2; y++ {
		for x := 0; x < (4*w - 1); x++ {
			tmp := img.Pix[x+y*img.Stride]
			img.Pix[x+y*img.Stride] = img.Pix[x+(h-y-1)*img.Stride]
			img.Pix[x+(h-y-1)*img.Stride] = tmp
		}
	}
	// Scale by alfa
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			ofs := x*4 + y*img.Stride
			// Set alpha to 1.0 (dvs 255). It is not used in files
			img.Pix[ofs+3] = 255
		}
	}
	return img
}
