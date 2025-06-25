package sys

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/gpu/font"
	"log/slog"
	"time"
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
func CreateWindow(rect f32.Rect, name string, monitorNo int, userScale float32) {
	m := Monitors[max(0, min(monitorNo-1, len(Monitors)-1))]
	ScaleX, ScaleY := m.GetContentScale()
	PosX, PosY, SizePxX, SizePxY := m.GetWorkarea()
	if rect.W == 0 {
		rect.W = float32(SizePxX)
	} else {
		rect.W = min(rect.W*ScaleX, float32(SizePxX))
	}
	if rect.H == 0 {
		rect.H = float32(SizePxY)
	} else {
		rect.H = min(rect.H*ScaleY, float32(SizePxY))
	}
	w := createWindow(int(rect.W), int(rect.W), name, nil)
	WindowList = append(WindowList, w)
	info := gpu.WinInfo{}
	info.InvalidateChan = make(chan time.Duration, 10)
	wno := len(WindowList) - 1
	// Move the window to the selected monitor
	w.SetPos(PosX+int(rect.X), PosY+int(rect.Y))
	_, top, _, _ := WindowList[0].GetFrameSize()
	w.SetSize(int(rect.W), int(rect.H)-top)

	// Now we can update size and scaling
	info.UserScale = userScale
	gpu.Info = append(gpu.Info, info)
	gpu.CurrentInfo = &gpu.Info[wno]
	UpdateSize(len(WindowList) - 1)
	SetupCursors()
	w.MakeContextCurrent()
	w.Show()

	slog.Info("New window", "ScaleX", f32.F2S(gpu.Info[wno].ScaleX, 2), "ScaleY", f32.F2S(gpu.Info[wno].ScaleY, 2), "Monitor", monitorNo, "UserScale",
		f32.F2S(userScale, 2), "W", rect.W, "H", rect.H, "WDp", int(gpu.CurrentInfo.WindowRect.W),
		"HDp", int(gpu.CurrentInfo.WindowRect.H))

	setCallbacks(w)
	w.Focus()
	gpu.CurrentInfo.Focused = true
	gpu.InitGpu()
	font.LoadDefaultFonts()
	gpu.LoadIcons()
}

func resetCursor() {
	WindowList[0].SetCursor(nil)
}

func Running(wno int) bool {
	return !WindowList[wno].ShouldClose()
}
