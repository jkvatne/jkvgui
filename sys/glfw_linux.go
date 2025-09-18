// sys is the only package that depends on glfw.
package sys

import (
	"flag"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jkvatne/jkvgui/f32"
)

var (
	Monitors      []*glfw.Monitor
	maxFps        = flag.Int("maxfps", 60, "Set to maximum alowed frames pr second. Default to 60")
	NoScaling     bool
	WindowList    []*Window
	WindowCount   atomic.Int32
	MinFrameDelay = time.Second / 50
	MaxFrameDelay = time.Second / 5
	OpenGlStarted bool
)

type HintDef struct {
	WidgetRect f32.Rect // Original widgets size
	Text       string
	T          time.Time
	Tag        any
}

// Pr window global variables.
type Window struct {
	Window               *glfw.Window
	Name                 string
	Wno                  int
	UserScale            float32
	Mutex                sync.Mutex
	Trigger              chan bool
	HintActive           bool
	Focused              bool
	Blinking             atomic.Bool
	Cursor               int
	CurrentTag           interface{}
	LastTag              interface{}
	MoveToNext           bool
	MoveToPrevious       bool
	ToNext               bool
	SuppressEvents       bool
	mousePos             f32.Pos
	LeftBtnIsDown        bool
	LeftBtnReleased      bool
	Dragging             bool
	LeftBtnDownTime      time.Time
	LeftBtnUpTime        time.Time
	LeftBtnDoubleClicked bool
	ScrolledDistY        float32
	DialogVisible        bool
	redraws              int
	fps                  float64
	redrawStart          time.Time
	LastRune             rune
	LastKey              glfw.Key
	LastMods             glfw.ModifierKey
	NoScaling            bool
	CurrentHint          HintDef
}

var (
	pVResizeCursor   *glfw.Cursor
	pHResizeCursor   *glfw.Cursor
	pArrowCursor     *glfw.Cursor
	pHandCursor      *glfw.Cursor
	pCrosshairCursor *glfw.Cursor
	pIBeamCursor     *glfw.Cursor
)

const (
	KeyRight     = glfw.KeyRight
	KeyLeft      = glfw.KeyLeft
	KeyUp        = glfw.KeyUp
	KeyDown      = glfw.KeyDown
	KeySpace     = glfw.KeySpace
	KeyEnter     = glfw.KeyEnter
	KeyKPEnter   = glfw.KeyKPEnter
	KeyEscape    = glfw.KeyEscape
	KeyBackspace = glfw.KeyBackspace
	KeyDelete    = glfw.KeyDelete
	KeyHome      = glfw.KeyHome
	KeyEnd       = glfw.KeyEnd
	KeyPageUp    = glfw.KeyPageUp
	KeyPageDown  = glfw.KeyPageDown
	KeyInsert    = glfw.KeyInsert
	KeyC         = glfw.KeyC
	KeyV         = glfw.KeyV
	KeyX         = glfw.KeyX
	ModShift     = glfw.ModShift
	ModControl   = glfw.ModControl
	ModAlt       = glfw.ModAlt
)

const (
	ArrowCursor     = int(glfw.ArrowCursor)
	IBeamCursor     = int(glfw.IBeamCursor)
	CrosshairCursor = int(glfw.CrosshairCursor)
	HandCursor      = int(glfw.HandCursor)
	HResizeCursor   = int(glfw.HResizeCursor)
	VResizeCursor   = int(glfw.VResizeCursor)
)

type Cursor glfw.Cursor

func (w *Window) Invalidate() {
	glfw.PostEmptyEvent()
	glfw.PostEmptyEvent()
}

func (w *Window) PollEvents() {
	w.ClearMouseBtns()
	// Tight loop, waiting for events, checking for events every minDelay
	// Break anyway if waiting more than MaxFrameDelay
	t := time.Now()
	for time.Since(t) < MaxFrameDelay {
		glfw.WaitEventsTimeout(float64(MaxFrameDelay) / 1e9)
	}
	if time.Since(t) < MinFrameDelay {
		time.Sleep(MinFrameDelay - time.Since(t))
	}
	glfw.PollEvents()
}

func PollEvents() {
	glfw.WaitEventsTimeout(float64(MaxFrameDelay) / 1e9)
	glfw.PollEvents()
}

func Shutdown() {
	for _, win := range WindowList {
		win.Window.Destroy()
	}
	WindowList = WindowList[0:0]
	WindowCount.Store(0)
	glfw.Terminate()
	TerminateProfiling()
}

func GetMonitors() []*glfw.Monitor {
	return glfw.GetMonitors()
}

func focusCallback(w *glfw.Window, focused bool) {
	win := GetWindow(w)
	if win != nil {
		win.Focused = focused
		if !focused {
			slog.Debug("Lost focus", "Wno", win.Wno+1)
			win.ClearMouseBtns()
		} else {
			slog.Debug("Got focus", "Wno", win.Wno+1)
		}
		win.Invalidate()
	} else {
		slog.Error("Focus callback without any window", "Wno", win.Wno+1)
	}
}

func setCallbacks(Window *glfw.Window) {
	Window.SetMouseButtonCallback(btnCallback)
	Window.SetCursorPosCallback(posCallback)
	Window.SetKeyCallback(keyCallback)
	Window.SetCharCallback(charCallback)
	Window.SetScrollCallback(scrollCallback)
	Window.SetContentScaleCallback(scaleCallback)
	Window.SetFocusCallback(focusCallback)
	Window.SetCloseCallback(closeCallback)
	Window.SetSizeCallback(sizeCallback)
}

func closeCallback(w *glfw.Window) {
	slog.Debug("Close callback", "ShouldClose", w.ShouldClose())
}

// keyCallback see https://www.glfw.org/docs/latest/window_guide.html
func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	slog.Debug("keyCallback", "key", key, "scancode", scancode, "action", action, "mods", mods)
	win := GetWindow(w)
	win.Invalidate()
	if key == glfw.KeyTab && action == glfw.Release {
		win.MoveByKey(mods != glfw.ModShift)
	}
	if action == glfw.Release || action == glfw.Repeat {
		win.LastKey = key
	}
	win.LastMods = mods
}

func charCallback(w *glfw.Window, char rune) {
	slog.Debug("charCallback()", "Rune", int(char))
	win := GetWindow(w)
	win.Invalidate()
	win.LastRune = char
}

// btnCallback is called from the glfw window handler when mouse buttons change states.
func btnCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	win := GetWindow(w)
	win.Invalidate()
	win.LastMods = mods
	x, y := w.GetCursorPos()
	win.mousePos.X = float32(x) / win.Gd.ScaleX
	win.mousePos.Y = float32(y) / win.Gd.ScaleY
	slog.Debug("Mouse click:", "Button", button, "X", x, "Y", y, "Action", action, "FromWindow", wno)
	if button == glfw.MouseButtonLeft {
		if action == glfw.Release {
			win.LeftBtnIsDown = false
			win.LeftBtnReleased = true
			win.Dragging = false
			if time.Since(win.LeftBtnUpTime) < DoubleClickTime {
				win.LeftBtnDoubleClicked = true
			}
			win.LeftBtnUpTime = time.Now()
		} else if action == glfw.Press {
			win.LeftBtnIsDown = true
			win.LeftBtnDownTime = time.Now()
		}
	}
}

// posCallback is called from the glfw window handler when the mouse moves.
func posCallback(w *glfw.Window, xpos float64, ypos float64) {
	win := GetWindow(w)
	win.mousePos.X = float32(xpos) / win.Gd.ScaleX
	win.mousePos.Y = float32(ypos) / win.Gd.ScaleY
	win.Invalidate()
}

func scrollCallback(w *glfw.Window, xoff float64, yOff float64) {
	slog.Debug("Scroll", "dx", xoff, "dy", yOff)
	win := GetWindow(w)
	if win.LastMods == glfw.ModControl {
		// ctrl+scrollwheel will zoom the whole window by changing gpu.UserScale.
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

func GetWindow(w *glfw.Window) *Window {
	for i, _ := range WindowList {
		if WindowList[i].Window == w {
			return WindowList[i]
		}
	}
	return nil
}

func sizeCallback(w *glfw.Window, width int, height int) {
	slog.Debug("sizeCallback", "width", width, "height", height)
	win := GetWindow(w)
	win.UpdateSize(width, height)
	win.UpdateResolution()
	win.Invalidate()
}

func scaleCallback(w *glfw.Window, x float32, y float32) {
	slog.Debug("scaleCallback", "x", x, "y", y)
	win := GetWindow(w)
	win.UpdateSizeDp()
	win.UpdateResolution()
}

func SetDefaultHints() {
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Samples, 4)
	glfw.WindowHint(glfw.Floating, glfw.False) // True will keep the window on top
}

func SetMaximizedHint(maximized bool) {
	if maximized {
		glfw.WindowHint(glfw.Maximized, glfw.True)
	} else {
		glfw.WindowHint(glfw.Maximized, glfw.False)
	}

}

func createInvisibleWindow(w, h int, title string, monitor *glfw.Monitor) *glfw.Window {
	// Create invisible window so we can move it to correct monitor
	glfw.WindowHint(glfw.Visible, glfw.False)
	win, err := glfw.CreateWindow(w, h, title, monitor, nil)
	if err != nil || win == nil {
		panic(err)
	}
	return win
}

func SetupCursors() {
	pArrowCursor = glfw.CreateStandardCursor(glfw.ArrowCursor)
	pVResizeCursor = glfw.CreateStandardCursor(glfw.VResizeCursor)
	pHResizeCursor = glfw.CreateStandardCursor(glfw.HResizeCursor)
	pIBeamCursor = glfw.CreateStandardCursor(glfw.HResizeCursor)
	pCrosshairCursor = glfw.CreateStandardCursor(glfw.HResizeCursor)
	pHandCursor = glfw.CreateStandardCursor(glfw.HResizeCursor)
}

func SetClipboardString(s string) {
	glfw.SetClipboardString(s)
}

func GetClipboardString() (string, error) {
	return glfw.GetClipboardString(), nil
}

func MaximizeWindow(w *glfw.Window) {
	w.Maximize()
}

func MinimizeWindow(w *glfw.Window) {
	w.Iconify()
}

func PostEmptyEvent() {
	glfw.PostEmptyEvent()
}

func glfwInit() error {
	return glfw.Init()
}

func DetachCurrentContext() {
	glfw.DetachCurrentContext()
}

func SwapInterval(n int) {
	glfw.SwapInterval(n)
}

func GetCurrentContext() *glfw.Window {
	return glfw.GetCurrentContext()
}
