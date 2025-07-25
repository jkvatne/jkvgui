package sys

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/gpu/font"
	"log/slog"
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
	// sx, sy := Monitors[0].GetContentScale()
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
	WindowList = append(WindowList, win)
	info := gpu.WinInfo{}
	info.InvalidateCount.Store(0)
	wno := len(WindowList) - 1
	CurrentWindow = WindowList[wno]
	_, top, _, _ := CurrentWindow.GetFrameSize()
	// Move the window to the selected monitor
	win.SetPos(PosX+x, PosY+y+top)
	win.SetSize(w, h)
	// Now we can update size and scaling
	info.UserScale = userScale
	gpu.Info = append(gpu.Info, &info)
	gpu.CurrentInfo = gpu.Info[wno]
	UpdateSize(len(WindowList) - 1)
	SetupCursors()
	win.MakeContextCurrent()
	win.Show()
	slog.Info("CreateWindow()",
		"ScaleX", f32.F2S(gpu.Info[wno].ScaleX, 2), ""+
			"ScaleY", f32.F2S(gpu.Info[wno].ScaleY, 2),
		"Monitor", monitorNo, "UserScale",
		f32.F2S(userScale, 2), "W", w, "H", h,
		"WDp", int(gpu.CurrentInfo.WindowContentRectDp.W),
		"HDp", int(gpu.CurrentInfo.WindowContentRectDp.H))

	setCallbacks(win)
	win.Focus()
	gpu.CurrentInfo.Focused = true
	gpu.InitGpu()
	font.LoadDefaultFonts()
	gpu.LoadIcons()
}

func Running() bool {
	for wno, win := range WindowList {
		if win.ShouldClose() {
			WindowList = append(WindowList[:wno], WindowList[wno+1:]...)
			gpu.Info = append(gpu.Info[:wno], gpu.Info[wno+1:]...)
			// THis gives invalid handle: win.Destroy()
		}
	}
	return len(WindowList) > 0
}
