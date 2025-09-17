package test

import (
	"log/slog"
	"testing"
	"time"

	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
)

func TestShadows(t *testing.T) {
	sys.Init()
	defer sys.Shutdown()
	sys.NoScaling = true
	slog.SetLogLoggerLevel(slog.LevelError)
	theme.SetDefaultPallete(true)
	w := sys.CreateWindow(0, 0, 400, 150, "Test", 2, 1.0)
	sys.LoadOpenGl(w)
	w.StartFrame(f32.White)
	r := f32.Rect{X: 10, Y: 10, W: 30, H: 20}
	w.Gd.RoundedRect(r, 0, 0.5, f32.Transparent, f32.Black)
	r.X += 50
	w.Gd.RoundedRect(r, 5, 0.5, f32.Transparent, f32.Black)
	r.X += 50
	w.Gd.RoundedRect(r, 5, 2, f32.Transparent, f32.Black)
	r.X += 50
	w.Gd.RoundedRect(r, 9999, 1, f32.Transparent, f32.Black)
	r.X += 50
	w.Gd.RoundedRect(r, 0, 0.5, f32.LightBlue, f32.Black)
	r.X += 50
	w.Gd.RoundedRect(r, 5, 0.51, f32.LightBlue, f32.Black)
	r.X += 50
	w.Gd.RoundedRect(r, 5, 2, f32.LightBlue, f32.Black)
	r.X += 50
	w.Gd.RoundedRect(r, 9999, 1, f32.LightBlue, f32.Black)
	r.X = 10
	r.Y += 50
	w.Gd.RoundedRect(r, 6, 0.5, f32.Transparent, f32.Black)
	w.Gd.Shade(r, 6, f32.Shade, 3)
	r.X += 50
	w.Gd.RoundedRect(r, 6, 0.5, f32.Transparent, f32.Black)
	w.Gd.Shade(r, 6, f32.Shade, 6)
	r.X += 50
	w.Gd.RoundedRect(r, 6, 0.5, f32.Transparent, f32.Black)
	w.Gd.Shade(r, 6, f32.Shade, 10)
	r.X += 50
	w.Gd.RoundedRect(r, 999, 0.5, f32.Transparent, f32.Black)
	w.Gd.Shade(r, 999, f32.Shade, 3)
	r.X += 50
	w.Gd.RoundedRect(r, 999, 0.5, f32.Transparent, f32.Black)
	w.Gd.Shade(r, 999, f32.Shade, 6)
	r.X += 50
	w.Gd.RoundedRect(r, 999, 0.5, f32.Transparent, f32.Black)
	w.Gd.Shade(r, 999, f32.Shade, 10)
	r.X += 50
	f32.AssertDir("test-outputs")
	VerifyScreen(t, w, "TestShadows", 400, 150, saveScreen)
	w.EndFrame()
	// Place breakpoint here in order to look at the screen output.
	time.Sleep(1 * time.Second)

}
