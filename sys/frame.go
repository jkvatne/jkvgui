package sys

import (
	"time"

	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/theme"
)

func (win *Window) StartFrame() {
	if win.Window.ShouldClose() {
		return
	}
	win.redraws++
	t := time.Since(win.redrawStart).Seconds()
	if t >= 1 {
		win.fps = float64(win.redraws) / t
		win.redrawStart = time.Now()
		win.redraws = 0
	}
	if WindowCount.Load() == 0 {
		panic("StartFrame() called, but no windows have been created")
	}
	win.MakeContextCurrent()
	win.UpdateSizeDp()
	win.UpdateResolution()
	SwapInterval(20)
	gpu.SetBackgroundColor(theme.Canvas.Bg())
	win.Blinking.Store(false)
	win.Cursor = ArrowCursor
}

// EndFrame will do buffer swapping and focus updates
func (win *Window) EndFrame() {
	if win.Window.ShouldClose() {
		return
	}
	if !win.DialogVisible {
		win.SuppressEvents = false
	}
	win.RunDeferred()
	win.LastKey = 0
	win.LeftBtnClicked = false
	win.Window.SwapBuffers()
	switch win.Cursor {
	case VResizeCursor:
		win.Window.SetCursor(pVResizeCursor)
	case HResizeCursor:
		win.Window.SetCursor(pHResizeCursor)
	case CrosshairCursor:
		win.Window.SetCursor(pCrosshairCursor)
	case HandCursor:
		win.Window.SetCursor(pHandCursor)
	case IBeamCursor:
		win.Window.SetCursor(pIBeamCursor)
	default:
		win.Window.SetCursor(pArrowCursor)
	}
	DetachCurrentContext()
}
