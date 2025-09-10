package sys

import (
	"flag"
	"time"

	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
)

var (
	MaxDelay     = time.Second
	redraws      int
	redrawStart  time.Time
	redrawsPrSec int
	minDelay     = time.Second / 50
)

func RedrawsPrSec() int {
	return redrawsPrSec
}

func StartFrame(bg f32.Color) {
	redraws++
	if time.Since(redrawStart).Seconds() >= 1 {
		redrawsPrSec = redraws
		redrawStart = time.Now()
		redraws = 0
	}
	MakeWindowCurrent(CurrentWno)
	gpu.SetBackgroundColor(bg)
	gpu.Info[CurrentWno].Blinking.Store(false)
	gpu.Info[CurrentWno].Cursor = ArrowCursor
}

// EndFrame will do buffer swapping and focus updates
// Then it will loop and sleep until an event happens
// maxFrameRate is used to limit the use of CPU/GPU. A maxFrameRate of zero will run the GPU/CPU as fast as
// possible with very high power consumption. More than 1k frames pr second is possible.
// Minimum framerate is 1 fps, so we will allways redraw once pr second - just in case we missed an event.
func EndFrame() {
	gpu.RunDeferred()
	LastKey = 0
	WindowList[CurrentWno].SwapBuffers()
	c := gpu.Info[CurrentWno].Cursor
	switch c {
	case VResizeCursor:
		WindowList[CurrentWno].SetCursor(pVResizeCursor)
	case HResizeCursor:
		WindowList[CurrentWno].SetCursor(pHResizeCursor)
	case CrosshairCursor:
		WindowList[CurrentWno].SetCursor(pCrosshairCursor)
	case HandCursor:
		WindowList[CurrentWno].SetCursor(pHandCursor)
	case IBeamCursor:
		WindowList[CurrentWno].SetCursor(pIBeamCursor)
	default:
		WindowList[CurrentWno].SetCursor(pArrowCursor)
	}
}

var logLevel = flag.Int("loglevel", 8, "Set log level (8=Error, 4=Warning, 0=Info(default), -4=Debug)")
