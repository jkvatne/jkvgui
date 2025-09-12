package main

import (
	"log/slog"
	"os"

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
	// Configure slog to print date and time in standard format
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			a.Value = slog.StringValue(a.Value.Time().Format("2006-01-02 15:04:05.000"))
		}
		return a
	}})))
	slog.Info("Starting Resize demo")

	sys.Init()
	defer sys.Shutdown()
	w0 := sys.CreateWindow(100, 100, 500, 400, "Resizing1", 1, 1)
	sys.LoadOpenGl(w0)
	w1 := sys.CreateWindow(200, 200, 500, 400, "Resizing2", 2, 1)
	sys.LoadOpenGl(w1)
	image[0], _ = wid.NewImage("music.jpg")
	image[1], _ = wid.NewImage("ts.jpg")
	for w0.Running() && w1.Running() {
		w0.StartFrame(theme.Surface.Bg())
		wid.Show(Form(0))
		w0.EndFrame()
		w1.StartFrame(theme.Surface.Bg())
		wid.Show(Form(1))
		w1.EndFrame()
		if w0.Focused {
			w0.PollEvents()
		} else {
			w1.PollEvents()
		}
	}
}
