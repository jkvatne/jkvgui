package sys

import (
	"flag"
	"time"

	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gl"
	"github.com/jkvatne/jkvgui/gpu"
)

// UpdateResolution sets the resolution for all programs
func (w *Window) UpdateResolution() {
	ww := int32(w.WidthPx)
	hh := int32(w.HeightPx)
	gpu.SetResolution(w.Gd.FontProgram, ww, hh)
	gpu.SetResolution(w.Gd.RRprogram, ww, hh)
	gpu.SetResolution(w.Gd.ShaderProgram, ww, hh)
	gpu.SetResolution(w.Gd.ImgProgram, ww, hh)
}

func (w *Window) Clip(r f32.Rect) {
	ww := r.W * w.Gd.ScaleX
	hh := r.H * w.Gd.ScaleY
	xx := r.X * w.Gd.ScaleX
	yy := float32(w.HeightPx) - hh - r.Y*w.Gd.ScaleY
	gl.Scissor(int32(xx), int32(yy), int32(ww), int32(hh))
	gl.Enable(gl.SCISSOR_TEST)
}

func (w *Window) Fps() float64 {
	return w.fps
}

func (w *Window) StartFrame(bg f32.Color) {
	if w.Window.ShouldClose() {
		return
	}
	w.redraws++
	t := time.Since(w.redrawStart).Seconds()
	if t >= 1 {
		w.fps = float64(w.redraws) / t
		w.redrawStart = time.Now()
		w.redraws = 0
	}
	if WindowCount.Load() == 0 {
		panic("StartFrame() called, but no windows have been created")
	}
	w.MakeContextCurrent()
	w.UpdateSizeDp()
	w.UpdateResolution()
	SwapInterval(20)
	gpu.SetBackgroundColor(bg)
	w.Blinking.Store(false)
	w.Cursor = ArrowCursor
}

// EndFrame will do buffer swapping and focus updates
func (w *Window) EndFrame() {
	if w.Window.ShouldClose() {
		return
	}
	if !w.DialogVisible {
		w.SuppressEvents = false
	}
	w.RunDeferred()
	w.LastKey = 0
	w.LeftBtnClicked = false
	w.Window.SwapBuffers()
	switch w.Cursor {
	case VResizeCursor:
		w.Window.SetCursor(pVResizeCursor)
	case HResizeCursor:
		w.Window.SetCursor(pHResizeCursor)
	case CrosshairCursor:
		w.Window.SetCursor(pCrosshairCursor)
	case HandCursor:
		w.Window.SetCursor(pHandCursor)
	case IBeamCursor:
		w.Window.SetCursor(pIBeamCursor)
	default:
		w.Window.SetCursor(pArrowCursor)
	}
	DetachCurrentContext()
}

var logLevel = flag.Int("loglevel", 8, "Set log level (8=Error, 4=Warning, 0=Info(default), -4=Debug)")
