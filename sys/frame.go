package sys

import (
	"github.com/go-gl/glfw/v3.3/glfw"
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
	// Loop at max framerate, which should be >5
	for {
		dt := max(0, time.Second/time.Duration(maxFrameRate)-time.Since(startTime))
		time.Sleep(dt)
		startTime = time.Now()
		glfw.PollEvents()
		if time.Since(gpu.InvalidateAt) >= 0 {
			gpu.InvalidateAt = time.Now().Add(time.Second)
			break
		}
	}
}

var BlinkFrequency = 2.0
var BlinkState bool
var Blinking = true

func blinker() {
	for {
		time.Sleep(time.Microsecond * time.Duration(1e6/BlinkFrequency/2))
		BlinkState = !BlinkState
		if Blinking {
			gpu.InvalidateAt = time.Now()
		}
	}
}

func init() {
	go blinker()
}
