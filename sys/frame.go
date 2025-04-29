package sys

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/focus"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/mouse"
	"time"
)

func StartFrame(color f32.Color) {
	startTime = time.Now()
	redraws++
	if time.Since(redrawStart).Seconds() >= 1 {
		RedrawsPrSec = redraws
		redrawStart = time.Now()
		redraws = 0
	}
	focus.StartFrame()
	gpu.BackgroundColor(color)
}

// EndFrame will do buffer swapping and focus updates
// Then it will loop and sleep until an event happens
// The event could be an invalidate call
func EndFrame(maxFrameRate int) {
	gpu.RunDefered()
	gpu.LastKey = 0
	mouse.FrameEnd()
	gpu.Window.SwapBuffers()
	gpu.WaitForEvent()
}
