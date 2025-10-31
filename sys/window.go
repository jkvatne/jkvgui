package sys

import (
	"flag"
	"image"
	"log/slog"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/jkvatne/jkvgui/buildinfo"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/gpu/font"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/purego-glfw/gl"
)

var logLevel = flag.Int("loglevel", 8, "Set log level (8=Error, 4=Warning, 0=Info(default), -4=Debug)")

// Window variables.
type Window struct {
	Window                *GlfwWindow
	Name                  string
	Wno                   int
	UserScale             float32
	Mutex                 sync.Mutex
	Trigger               chan bool
	HintActive            bool
	Focused               bool
	Blinking              atomic.Bool
	Cursor                int
	CurrentTag            interface{}
	LastTag               interface{}
	MoveToNext            bool
	MoveToPrevious        bool
	ToNext                bool
	SuppressEvents        bool
	mousePos              f32.Pos
	Dragging              bool
	DragStartPos          f32.Pos
	LeftBtnIsDown         bool
	LeftBtnDownTime       time.Time
	LeftBtnUpTime         time.Time
	LeftBtnDoubleClicked  bool
	LeftBtnClicked        bool
	RightBtnIsDown        bool
	RightBtnDownTime      time.Time
	RightBtnUpTime        time.Time
	RightBtnDoubleClicked bool
	RightBtnClicked       bool
	ScrolledDistY         float32
	DialogVisible         bool
	redraws               int
	fps                   float64
	redrawStart           time.Time
	LastRune              rune
	LastKey               Key
	LastMods              ModifierKey
	NoScaling             bool
	CurrentHint           HintDef
	DeferredFunctions     []func()
	HeightPx              int
	HeightDp              float32
	WidthPx               int
	WidthDp               float32
	Gd                    gpu.GlData
}

var (
	maxFps        = flag.Int("maxfps", 60, "Set to maximum allowed frames pr second. Default to 60")
	NoScaling     bool
	WindowList    []*Window
	WindowCount   atomic.Int32
	WinListMutex  sync.RWMutex
	MinFrameDelay = time.Second / 50
	MaxFrameDelay = time.Second / 5
	LastPollTime  time.Time
	OpenGlStarted bool
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
	win.Window.SetPos(PosX+x, PosY+y+tb)
	// Now we can update size and scaling
	win.UserScale = userScale
	w = min(w+lb+rb, SizePxX+lb+rb)
	h = min(h+tb+bb, SizePxY+bb)
	win.Window.SetSize(w, h)
	w, h = win.Window.GetSize()
	win.UpdateSize(w, h)
	WinListMutex.Lock()
	WindowList = append(WindowList, win)
	wno := len(WindowList) - 1
	WinListMutex.Unlock()
	win.Wno = wno
	WindowCount.Add(1)
	win.Name = name
	win.Trigger = make(chan bool, 1)
	setupCursors()
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
	win.ClearMouseBtns()
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
// It should be called once, at the very beginning of the application.
func Init() {
	runtime.LockOSThread()
	flag.Parse()
	slog.SetLogLoggerLevel(slog.Level(*logLevel))
	buildinfo.Get()
	if *maxFps == 0 {
		MinFrameDelay = 0
	} else {
		MinFrameDelay = time.Second / time.Duration(*maxFps)
	}
	// Initialize glfw
	if err := glfwInit(); err != nil {
		panic(err)
	}
	theme.SetDefaultPalette(true)
	SetDefaultHints()
	// Check all monitors and print size data
	Monitors = GetMonitors()
	// Select monitor as given, or use primary monitor.
	for i, m := range Monitors {
		SizeMmX, SizeMmY := m.GetPhysicalSize()
		mScaleX, mScaleY := m.GetContentScale()
		PosX, PosY, SizePxX, SizePxY := m.GetWorkarea()
		slog.Debug("Init()", "Monitor", i+1,
			"WidthMm", SizeMmX, "HeightMm", SizeMmY,
			"WidthPx", SizePxX, "HeightPx", SizePxY, "PosX", PosX, "PosY", PosY,
			"ScaleX", f32.F2S(mScaleX, 3, 4), "ScaleY", f32.F2S(mScaleY, 3, 4))
	}
	go Blinker()
}

func (win *Window) UpdateSizeDp() {
	if NoScaling {
		win.Gd.ScaleX = 1.0
		win.Gd.ScaleY = 1.0
		win.WidthDp = float32(win.WidthPx)
		win.HeightDp = float32(win.HeightPx)
	} else {
		win.Gd.ScaleX, win.Gd.ScaleY = win.Window.GetContentScale()
		win.Gd.ScaleX *= win.UserScale
		win.Gd.ScaleY *= win.UserScale
		win.WidthDp = float32(win.WidthPx) / win.Gd.ScaleX
		win.HeightDp = float32(win.HeightPx) / win.Gd.ScaleY
	}
}

func (win *Window) UpdateSize(width, height int) {
	win.WidthPx = width
	win.HeightPx = height
	win.UpdateSizeDp()
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

func (win *Window) ClientRectDp() f32.Rect {
	return f32.Rect{W: win.WidthDp, H: win.HeightDp}
}

func (win *Window) Defer(f func()) {
	for _, g := range win.DeferredFunctions {
		if &f == &g {
			return
		}
	}
	win.DeferredFunctions = append(win.DeferredFunctions, f)
}

func (win *Window) RunDeferred() {
	for _, f := range win.DeferredFunctions {
		f()
	}
	win.DeferredFunctions = win.DeferredFunctions[0:0]
}

func (win *Window) MakeContextCurrent() {
	win.Window.MakeContextCurrent()
}

func (win *Window) SetCursor(c int) {
	win.Cursor = c
}

// Invalidate will trigger all windows to paint their contents
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
			WinListMutex.Lock()
			WindowList = append(WindowList[:wno], WindowList[wno+1:]...)
			WinListMutex.Unlock()
			WindowCount.Add(-1)
			return len(WindowList) > 0
		}
	}
	WinListMutex.RLock()
	defer WinListMutex.RUnlock()
	return len(WindowList) > 0
}

// AbortAfter is to be called as a go routine
// It closes all windows after the given delay
func AbortAfter(delay time.Duration, windowCount int) {
	// First wait until all windows are created (or timeout)
	t := time.Now()
	for len(WindowList) < windowCount && time.Since(t) < 10*time.Second {
		time.Sleep(10 * time.Millisecond)
	}
	// Show for given time
	time.Sleep(delay)
	// Close all windows
	for _, w := range WindowList {
		w.Window.SetShouldClose(true)
	}
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

// ClearMouseBtns is called when a window looses focus. It will reset the mouse button states.
func (win *Window) ClearMouseBtns() {
	win.LeftBtnIsDown = false
	win.Dragging = false
	win.ScrolledDistY = 0.0
	win.LeftBtnDoubleClicked = false
	win.LeftBtnClicked = false
	win.LeftBtnUpTime = time.Time{}
	win.RightBtnDoubleClicked = false
	win.RightBtnClicked = false
	win.RightBtnUpTime = time.Time{}
}

func (win *Window) leftBtnRelease() {
	win.LeftBtnIsDown = false
	win.Dragging = false
	if time.Since(win.LeftBtnUpTime) < DoubleClickTime {
		slog.Debug("MouseCb: - DoubleClick:")
		win.LeftBtnDoubleClicked = true
	} else {
		slog.Debug("MouseCb: - Click:")
		win.LeftBtnClicked = true
	}
	win.LeftBtnUpTime = time.Now()
}

func (win *Window) leftBtnPress() {
	win.LeftBtnIsDown = true
	win.LeftBtnClicked = false
	win.LeftBtnDownTime = time.Now()
}

func (win *Window) rightBtnRelease() {
	win.RightBtnIsDown = false
	if time.Since(win.RightBtnUpTime) < DoubleClickTime {
		slog.Debug("MouseCb: - Right DoubleClick:")
		win.RightBtnDoubleClicked = true
	} else {
		slog.Debug("MouseCb: - Right Click:")
		win.RightBtnClicked = true
	}
	win.RightBtnUpTime = time.Now()
}
func (win *Window) rightBtnPress() {
	win.RightBtnIsDown = true
	win.RightBtnClicked = false
	win.RightBtnDownTime = time.Now()
}

func (win *Window) SimPos(x, y float32) {
	win.mousePos.X = x
	win.mousePos.Y = y
}

func (win *Window) SimLeftBtnPress(x, y float32) {
	win.mousePos.X = x
	win.mousePos.Y = y
	win.leftBtnPress()
}

func (win *Window) SimLeftBtnRelease(x, y float32) {
	win.mousePos.X = x
	win.mousePos.Y = y
	win.leftBtnRelease()
}

// UpdateResolution sets the resolution for all programs
func (win *Window) UpdateResolution() {
	ww := int32(win.WidthPx)
	hh := int32(win.HeightPx)
	win.Gd.HeightPx = win.HeightPx
	win.Gd.WidthPx = win.WidthPx
	gpu.SetResolution(win.Gd.FontProgram, ww, hh)
	gpu.SetResolution(win.Gd.RRprogram, ww, hh)
	gpu.SetResolution(win.Gd.ShaderProgram, ww, hh)
	gpu.SetResolution(win.Gd.ImgProgram, ww, hh)
}

func (win *Window) Fps() float64 {
	return win.fps
}

func (win *Window) Destroy() {
	win.Window.Destroy()
}

func (win *Window) Invalidate() {
	PostEmptyEvent()
}

func PollEvents() {
	timeUsed := time.Now().Sub(LastPollTime)
	// If the drawing took less than the min frame delay...
	if timeUsed < MinFrameDelay {
		// Sleep the remaining time
		time.Sleep(MinFrameDelay - timeUsed)
	}
	// Then wait for an event
	WaitEventsTimeout(float32(MaxFrameDelay-MinFrameDelay) / 1e9)
	LastPollTime = time.Now()
}

func Shutdown() {
	WinListMutex.Lock()
	for _, win := range WindowList {
		win.Window.Destroy()
	}
	WindowList = WindowList[0:0]
	WindowCount.Store(0)
	WinListMutex.Unlock()
	Terminate()
	OpenGlStarted = false
}

func (win *Window) HandleFocus(focused bool) {
	win.Focused = focused
	if !focused {
		slog.Debug("Lost focus", "Wno ", win.Wno+1)
	} else {
		slog.Debug("Got focus", "Wno", win.Wno+1)
	}
	win.ClearMouseBtns()
	win.Invalidate()
}

func (win *Window) HandleKey(key Key, scancode int, action Action, mods ModifierKey) {
	slog.Debug("keyCallback", "key", key, "scancode", scancode, "action", action, "mods", mods)
	win.Invalidate()
	if key == KeyTab && action == Release {
		win.MoveByKey(mods != ModShift)
	}
	if action == Release || action == Repeat {
		win.LastKey = key
	}
	win.LastMods = mods
}

func (win *Window) HandleMouseButton(button MouseButton, action Action, mods ModifierKey) {
	win.LastMods = mods
	x, y := win.Window.GetCursorPos()
	win.mousePos.X = float32(x) / win.Gd.ScaleX
	win.mousePos.Y = float32(y) / win.Gd.ScaleY
	slog.Debug("MouseCb:", "Button", button, "X", x, "Y", y, "Action", action, "FromWindow", win.Wno, "Pos", win.mousePos)
	if button == MouseButtonLeft {
		if action == Release {
			win.leftBtnRelease()
		} else if action == Press {
			win.leftBtnPress()
		}
		win.Invalidate()
	}
}

func (win *Window) HandleMousePos(xPos float64, yPos float64) {
	win.mousePos.X = float32(xPos) / win.Gd.ScaleX
	win.mousePos.Y = float32(yPos) / win.Gd.ScaleY
	win.Invalidate()
}

func (win *Window) HandleMouseScroll(xOff float64, yOff float64) {
	slog.Debug("ScrollCb:", "dx", xOff, "dy", yOff)
	if win.LastMods == ModControl {
		// ctrl + scroll-wheel will zoom the whole window by changing gpu.UserScale.
		if yOff > 0 {
			win.UserScale *= ZoomFactor
		} else {
			win.UserScale /= ZoomFactor
		}
		win.UpdateSizeDp()
	} else {
		win.ScrolledDistY = float32(yOff)
	}
	win.Invalidate()
}

func (win *Window) HandleChar(char rune) {
	slog.Debug("charCallback()", "Rune", int(char))
	win.Invalidate()
	win.LastRune = char
}
