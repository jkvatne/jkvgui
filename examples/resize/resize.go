package main

import (
	"log"
	"log/slog"

	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
)

var posV [16]wid.ResizerState
var posH [16]wid.ResizerState
var image [16]*wid.Img

func Form(n int) wid.Wid {
	return wid.HorResizer(
		&posH[n], nil,
		wid.Image(image[n], nil, ""),
		wid.VertResizer(&posV[n], nil,
			wid.Btn("Left", nil, func() {}, nil, ""),
			wid.Btn("Right", nil, func() {}, nil, ""),
		),
	)
}

func main() {
	log.SetFlags(log.Lmicroseconds)
	slog.Info("Resize demo")
	sys.Init()
	defer sys.Shutdown()
	w0 := sys.CreateWindow(100, 100, 500, 400, "Resizing1", 1, 1)
	w1 := sys.CreateWindow(200, 200, 500, 400, "Resizing2", 2, 1)
	image[0], _ = wid.NewImage("music.jpg")
	image[1], _ = wid.NewImage("ts.jpg")
	for sys.Running() {
		w0.StartFrame(theme.Surface.Bg())
		wid.Show(Form(0))
		w0.EndFrame()
		w1.StartFrame(theme.Surface.Bg())
		wid.Show(Form(1))
		w1.EndFrame()
		sys.PollEvents()
	}
}
