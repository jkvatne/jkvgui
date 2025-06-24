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

// CreateWindow initializes glfw and returns a Window to use.
// MonitorNo is 1 or 0 for the primary monitor, 2 for secondary monitor etc.
// Size is given in dp (device independent pixels)
// Windows typically fills the screen in one of the following ways:
// - Constant aspect ratio, use as much of screen as possible (h=10000, w=10000)
// - Full screen. (Maximized window) (w=0, h=0)
// - Small window of a given size, shrunk if the screen is not big enough (h=200, w=200)
// - Use full screen height, but limit width (h=0, w=800)
// - Use full screen width, but limit height (h=800, w=0)
func CreateWindow(wRequest, hRequest float32, name string, monitorNo int, userScale float32) {
	m := Monitors[max(0, min(monitorNo-1, len(Monitors)-1))]
	ScaleX, ScaleY := m.GetContentScale()
	PosX, PosY, SizePxX, SizePxY := m.GetWorkarea()
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
	w := createWindow(int(wRequest), int(hRequest), name, nil)
	WindowList = append(WindowList, w)
	info := gpu.WinInfo{}
	wno := len(WindowList) - 1
	// Move the window to the selected monitor
	w.SetPos(PosX, PosY)
	_, top, _, _ := WindowList[0].GetFrameSize()
	w.SetSize(int(wRequest), int(hRequest)-top)
	w.SetPos(PosX+(SizePxX-int(wRequest))/2, top+PosY+(SizePxY-int(hRequest))/2)

	// Now we can update size and scaling
	info.UserScale = userScale
	gpu.Info = append(gpu.Info, info)
	UpdateSize(len(WindowList) - 1)
	SetupCursors()
	w.MakeContextCurrent()
	w.Show()
	w.Focus()

	slog.Info("New window", "ScaleX", f32.F2S(gpu.Info[wno].ScaleX, 2), "ScaleY", f32.F2S(gpu.Info[wno].ScaleY, 2), "Monitor", monitorNo, "UserScale",
		f32.F2S(userScale, 2), "W", wRequest, "H", hRequest, "WDp", int(WindowWidthDp), "HDp", int(WindowHeightDp))

	gpu.InitGpu()
	font.LoadDefaultFonts()
	gpu.LoadIcons()
	gpu.UpdateResolution(wno)
	setCallbacks(w)
}

func resetCursor() {
	WindowList[0].SetCursor(nil)
}

func Running(wno int) bool {
	return !WindowList[wno].ShouldClose()
}
