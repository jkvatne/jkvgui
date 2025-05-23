package sys

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jkvatne/jkvgui/focus"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/mouse"
	"github.com/jkvatne/jkvgui/theme"
	"time"
)

var MaxDelay = time.Second

func StartFrame(role theme.UIRole) {
	redraws++
	if time.Since(redrawStart).Seconds() >= 1 {
		RedrawsPrSec = redraws
		redrawStart = time.Now()
		redraws = 0
	}
	focus.StartFrame()
	gpu.SetBackgroundColor(role.Bg())
	gpu.Blinking.Store(false)
	gpu.ResetCursor()
}

// EndFrame will do buffer swapping and focus updates
// Then it will loop and sleep until an event happens
// maxFrameRate is used to limit the use of CPU/GPU. A maxFrameRate of zero will run the GPU/CPU as fast as
// possible with very high power consumption. More than 1k frames pr second is possible.
// Minimum framerate is 1 fps, so we will allways redraw once pr second - just in case we missed an event.
func EndFrame(maxFrameRate int) {
	gpu.LastKey = 0
	mouse.FrameEnd()
	gpu.Window.SwapBuffers()
	t := time.Now()
	minDelay := time.Duration(0)
	if maxFrameRate != 0 {
		minDelay = time.Second / time.Duration(maxFrameRate)
	}
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
