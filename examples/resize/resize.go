package main

import (
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
)

var pos1 wid.ResizerState
var pos2 wid.ResizerState

func main() {
	// Initialize the GUI system and parse arguments
	sys.Initialize()
	// Create a window with a title and size.
	gpu.InitWindow(400, 200, "Resizing", 0, 2)
	image, _ := wid.NewImage("music.jpg")
	// Initialize the window and the GUI system, including callbacks.
	sys.InitializeWindow()
	// Loop until the window is closed.
	for !gpu.ShouldClose() {
		sys.StartFrame(theme.Surface.Bg())
		ctx := wid.NewCtx()
		wid.HorResizer(
			&pos1, nil,
			wid.Image(image, nil, ""),
			wid.VertResizer(&pos2, nil,
				wid.Btn("Left", nil, func() {}, nil, ""),
				wid.Btn("Right", nil, func() {}, nil, ""),
			),
		)(ctx)
		// EndFrame will swap buffers and limit the maximum framerate.
		sys.EndFrame(50)
	}
}
