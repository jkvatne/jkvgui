package main

import (
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/input"
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
)

var pos1 wid.ResizerState
var pos2 wid.ResizerState

func main() {
	input.InitWindow(400, 200, "Resizing", 0, 2)
	defer sys.Shutdown()
	sys.InitializeWindow()
	sys.Initialize()
	input.SetCallbacks()
	for !input.ShouldClose() {

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
