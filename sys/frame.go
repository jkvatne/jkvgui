package sys

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/focus"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/mouse"
	"sync/atomic"
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

func wait() {
	// Tight loop, waiting for events
	for {
		glfw.PollEvents()
		select {
		case <-gpu.InvalidateChan:
			return
		default:
			time.Sleep(time.Millisecond * time.Duration(50))
		}
	}
}

// EndFrame will do buffer swapping and focus updates
// Then it will loop and sleep until an event happens
// The event could be an invalidate call
func EndFrame(maxFrameRate int) {
	gpu.RunDefered()
	gpu.LastKey = 0
	mouse.FrameEnd()
	gpu.Window.SwapBuffers()
	wait()
}

var BlinkFrequency = 0.1
var BlinkState atomic.Bool

func blinker() {
	for {
		time.Sleep(time.Microsecond * time.Duration(1e6/BlinkFrequency/2))
		b := BlinkState.Load()
		BlinkState.Store(!b)
		gpu.InvalidateChan <- 0
	}
}

func init() {
	go blinker()
}
