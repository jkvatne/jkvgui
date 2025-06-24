package sys

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/gpu/font"
	"log/slog"
)

var (
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
	// Check all monitors and print size data
	ms := GetMonitors()
	// Select monitor as given, or use primary monitor.
	monitorNo = max(0, min(monitorNo-1, len(ms)-1))
	for i, m := range ms {
		SizeMmX, SizeMmY := m.GetPhysicalSize()
		ScaleX, ScaleY := m.GetContentScale()
		PosX, PosY, SizePxX, SizePxY := m.GetWorkarea()
		slog.Info("InitWindow()", "Monitor", i+1,
			"WidthMm", SizeMmX, "HeightMm", SizeMmY,
			"WidthPx", SizePxX, "HeightPx", SizePxY, "PosX", PosX, "PosY", PosY,
			"ScaleX", f32.F2S(ScaleX, 3), "ScaleY", f32.F2S(ScaleY, 3))
		if i == monitorNo {
			if wRequest == 0 {
				wRequest = float32(SizePxX)
			} else {
				wRequest = min(wRequest*ScaleX, float32(SizePxX))
			}
			if hRequest == 0 {
				hRequest = float32(SizePxY)
			} else {
				hRequest = min(hRequest*ScaleY, float32(SizePxY))
			}
			CurrentWindow = createWindow(int(wRequest), int(hRequest), name, nil)

			// Move the window to the selected monitor
			CurrentWindow.SetPos(PosX, PosY)
			_, top, _, _ := CurrentWindow.GetFrameSize()
			CurrentWindow.SetSize(int(wRequest), int(hRequest)-top)
			CurrentWindow.SetPos(PosX+(SizePxX-int(wRequest))/2, top+PosY+(SizePxY-int(hRequest))/2)
		}
	}

	// Now we can update size and scaling
	gpu.UserScale = userScale
	UpdateSize(CurrentWindow)
	SetupCursors()
	CurrentWindow.MakeContextCurrent()
	CurrentWindow.Show()
	CurrentWindow.Focus()

	slog.Info("New window", "ScaleX", f32.F2S(gpu.ScaleX, 2), "ScaleY", f32.F2S(gpu.ScaleY, 2), "Monitor", monitorNo, "UserScale", f32.F2S(userScale, 2),
		"W", wRequest, "H", hRequest, "WDp", int(WindowWidthDp), "HDp", int(WindowHeightDp))

	gpu.InitGpu()
	font.LoadDefaultFonts()
	gpu.LoadIcons()
	gpu.UpdateResolution()
	setCallbacks(CurrentWindow)
}

func SetVresizeCursor() {
	CurrentWindow.SetCursor(vResizeCursor)
}

func SetHresizeCursor() {
	CurrentWindow.SetCursor(hResizeCursor)
}

func resetCursor() {
	CurrentWindow.SetCursor(nil)
}

func Running() bool {
	return !CurrentWindow.ShouldClose()
}
