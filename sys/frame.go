package sys

import (
	"flag"
	"time"

	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
)

var (
	MinFrameDelay = time.Second / 50
	MaxFrameDelay = time.Second / 5
)

func (w *Window) Fps() float64 {
	return w.fps
}

func (w *Window) StartFrame(bg f32.Color) {
	if !OpenGlStarted {
		panic("OpenGl not started. Call sys.LoadOpenGl() before painting frames")
	}
	w.redraws++
	t := time.Since(w.redrawStart).Seconds()
	if t >= 1 {
		w.fps = float64(w.redraws) / t
		w.redrawStart = time.Now()
		w.redraws = 0
	}
	if len(WindowList) == 0 {
		panic("No windows have been created")
	}
	w.MakeContextCurrent()
	w.UpdateSize()
	gpu.UpdateResolution()
	SwapInterval(20)
	gpu.SetBackgroundColor(bg)
	w.Blinking.Store(false)
	w.Cursor = ArrowCursor
}

// EndFrame will do buffer swapping and focus updates
// Then it will loop and sleep until an event happens
// maxFrameRate is used to limit the use of CPU/GPU. A maxFrameRate of zero will run the GPU/CPU as fast as
// possible with very high power consumption. More than 1k frames pr second is possible.
// Minimum framerate is 1 fps, so we will allways redraw once pr second - just in case we missed an event.
func (w *Window) EndFrame() {
	w.SuppressEvents = false
	gpu.RunDeferred()
	LastKey = 0
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
	// w.ClearMouseBtns()
}

var logLevel = flag.Int("loglevel", 8, "Set log level (8=Error, 4=Warning, 0=Info(default), -4=Debug)")
