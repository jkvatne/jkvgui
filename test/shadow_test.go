package test

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
	"log/slog"
	"testing"
	"time"
)

func TestShadows(t *testing.T) {
	sys.Init()
	defer sys.Shutdown()
	sys.NoScaling = true
	slog.SetLogLoggerLevel(slog.LevelError)
	theme.SetDefaultPallete(true)
	sys.CreateWindow(0, 0, 400, 150, "Test", 2, 1.0)
	sys.StartFrame(f32.White)
	r := f32.Rect{X: 10, Y: 10, W: 30, H: 20}
	gpu.RoundedRect(r, 0, 0.5, f32.Transparent, f32.Black)
	r.X += 50
	gpu.RoundedRect(r, 5, 0.5, f32.Transparent, f32.Black)
	r.X += 50
	gpu.RoundedRect(r, 5, 2, f32.Transparent, f32.Black)
	r.X += 50
	gpu.RoundedRect(r, 9999, 1, f32.Transparent, f32.Black)
	r.X += 50
	gpu.RoundedRect(r, 0, 0.5, f32.LightBlue, f32.Black)
	r.X += 50
	gpu.RoundedRect(r, 5, 0.51, f32.LightBlue, f32.Black)
	r.X += 50
	gpu.RoundedRect(r, 5, 2, f32.LightBlue, f32.Black)
	r.X += 50
	gpu.RoundedRect(r, 9999, 1, f32.LightBlue, f32.Black)
	r.X = 10
	r.Y += 50
	gpu.RoundedRect(r, 6, 0.5, f32.Transparent, f32.Black)
	gpu.Shade(r, 6, f32.Shade, 3)
	r.X += 50
	gpu.RoundedRect(r, 6, 0.5, f32.Transparent, f32.Black)
	gpu.Shade(r, 6, f32.Shade, 6)
	r.X += 50
	gpu.RoundedRect(r, 6, 0.5, f32.Transparent, f32.Black)
	gpu.Shade(r, 6, f32.Shade, 10)
	r.X += 50
	gpu.RoundedRect(r, 999, 0.5, f32.Transparent, f32.Black)
	gpu.Shade(r, 999, f32.Shade, 3)
	r.X += 50
	gpu.RoundedRect(r, 999, 0.5, f32.Transparent, f32.Black)
	gpu.Shade(r, 999, f32.Shade, 6)
	r.X += 50
	gpu.RoundedRect(r, 999, 0.5, f32.Transparent, f32.Black)
	gpu.Shade(r, 999, f32.Shade, 10)
	r.X += 50
	f32.AssertDir("test-outputs")
	VerifyScreen(t, "TestShadows", 400, 150, saveScreen)

	sys.WindowList[0].SwapBuffers()
	// Place breakpoint here in order to look at the screen output.
	time.Sleep(1 * time.Second)

}
