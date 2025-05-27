//go:build !noglfw

package sys

import (
	"flag"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"time"
)

var (
	MaxDelay     = time.Second
	redraws      int
	redrawStart  time.Time
	redrawsPrSec int
	minDelay     = time.Second / 25
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
	gpu.SetBackgroundColor(bg)
	gpu.Blinking.Store(false)
	resetCursor()
	resetFocus()
}

func SetFrameRate(maxFrameRate float32) {
	if maxFrameRate <= 0 {
		maxFrameRate = 25
	}
	if maxFrameRate > 1000 {
		maxFrameRate = 1000
	}
	minDelay = time.Second / time.Duration(maxFrameRate)
}

// EndFrame will do buffer swapping and focus updates
// Then it will loop and sleep until an event happens
// maxFrameRate is used to limit the use of CPU/GPU. A maxFrameRate of zero will run the GPU/CPU as fast as
// possible with very high power consumption. More than 1k frames pr second is possible.
// Minimum framerate is 1 fps, so we will allways redraw once pr second - just in case we missed an event.
func EndFrame() {
	gpu.RunDeferred()
	LastKey = 0
	FrameEnd()
	Window.SwapBuffers()
	t := time.Now()

	glfw.PollEvents()

	// Tight loop, waiting for events, checking for events every millisecond
	for len(gpu.InvalidateChan) == 0 && time.Since(t) < MaxDelay {
		time.Sleep(minDelay)
		glfw.PollEvents()
	}
	// Empty the invalidate channel.
	if len(gpu.InvalidateChan) > 0 {
		<-gpu.InvalidateChan
	}
}

var maxFps = flag.Bool("maxfps", false, "Set to force redrawing as fast as possible")
var logLevel = flag.Int("loglevel", 0, "Set log level (8=Error, 4=Warning, 0=Info(default), -4=Debug)")

func Shutdown() {
	glfw.Terminate()
	TerminateProfiling()
}
