//go:build !noglfw

package sys

import (
	"flag"
	"github.com/jkvatne/jkvgui/buildinfo"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/gpu/font"
	"github.com/jkvatne/jkvgui/theme"
	"image"
	"log/slog"
	"runtime"
)

type Monitor struct {
	SizeMm image.Point
	SizePx image.Point
	ScaleX float32
	ScaleY float32
	Pos    image.Point
}

var (
	Monitors       []Monitor
	WindowWidthDp  float32
	WindowHeightDp float32
)

// InitWindow initializes glfw and returns a Window to use.
// MonitorNo is 1 or 0 for the primary monitor, 2 for secondary monitor etc.
// Size is given in dp (device independent pixels)
// Windows typically fills the screen in one of the following ways:
// - Constant aspect ratio, use as much of screen as possible (h=10000, w=10000)
// - Full screen. (Maximized window) (w=0, h=0)
// - Small window of a given size, shrunk if the screen is not big enough (h=200, w=200)
// - Use full screen height, but limit width (h=0, w=800)
// - Use full screen width, but limit height (h=800, w=0)
func InitWindow(wRequest, hRequest float32, name string, monitorNo int, userScale float32) {
	runtime.LockOSThread()
	flag.Parse()
	slog.SetLogLoggerLevel(slog.Level(*logLevel))
	InitializeProfiling()
	buildinfo.Get()
	if *maxFps {
		MaxDelay = 0
	}
	theme.SetDefaultPallete(true)
	InitGlfw()

	// Check all monitors and print size data
	ms := GetMonitors()
	for i, monitor := range ms {
		m := Monitor{}
		m.SizeMm.X, m.SizeMm.Y = monitor.GetPhysicalSize()
		_, _, m.SizePx.X, m.SizePx.Y = monitor.GetWorkarea()
		m.ScaleX, m.ScaleY = monitor.GetContentScale()
		m.Pos.X, m.Pos.Y, m.SizePx.X, m.SizePx.Y = monitor.GetWorkarea()
		slog.Info("InitWindow()", "Monitor", i+1,
			"WidthMm", m.SizeMm.X, "HeightMm", m.SizeMm.Y,
			"WidthPx", m.SizePx.X, "HeightPx", m.SizePx.Y, "PosX", m.Pos.X, "PosY", m.Pos.Y,
			"ScaleX", f32.F2S(m.ScaleX, 3), "ScaleY", f32.F2S(m.ScaleY, 3))
		if m.ScaleX == 0.0 {
			m.ScaleX = 1.0
		}
		if m.ScaleY == 0.0 {
			m.ScaleY = 1.0
		}
		Monitors = append(Monitors, m)
	}

	// Select monitor as given, or use primary monitor.
	monitorNo = max(0, min(monitorNo-1, len(Monitors)-1))
	m := Monitors[monitorNo]

	if wRequest == 0 {
		wRequest = float32(m.SizePx.X)
	} else {
		wRequest = min(wRequest*m.ScaleX, float32(m.SizePx.X))
	}
	if hRequest == 0 {
		hRequest = float32(m.SizePx.Y)
	} else {
		hRequest = min(hRequest*m.ScaleY, float32(m.SizePx.Y))
	}
	SetHints(int(wRequest), int(hRequest), name)

	// Move the window to the selected monitor
	Window.SetPos(m.Pos.X, m.Pos.Y)
	left, top, right, bottom := Window.GetFrameSize()
	slog.Info("Window.GetFrameSize()", "left", left, "top", top, "right", right, "bottom", bottom)

	Window.SetSize(int(wRequest), int(hRequest)-top)
	Window.SetPos(m.Pos.X+(m.SizePx.X-int(wRequest))/2, top+m.Pos.Y+(m.SizePx.Y-int(hRequest))/2)

	// Now we can update size and scaling
	gpu.UserScale = userScale
	UpdateSize(Window)
	Window.Show()
	slog.Info("New window", "ScaleX", f32.F2S(gpu.ScaleX, 2), "ScaleY", f32.F2S(gpu.ScaleY, 2), "Monitor", monitorNo, "UserScale", f32.F2S(userScale, 2),
		"W", wRequest, "H", hRequest, "WDp", int(WindowWidthDp), "HDp", int(WindowHeightDp))

	Window.MakeContextCurrent()
	WindowStart()

	gpu.InitGpu()
	font.LoadDefaultFonts()
	gpu.LoadIcons()
	gpu.UpdateResolution()
	setCallbacks(Window)
}

func SetVresizeCursor() {
	Window.SetCursor(vResizeCursor)
}

func SetHresizeCursor() {
	Window.SetCursor(hResizeCursor)
}

func resetCursor() {
	Window.SetCursor(nil)
}

func Running() bool {
	return !Window.ShouldClose()
}
