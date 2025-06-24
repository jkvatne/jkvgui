package main

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
	slog.SetLogLoggerLevel(slog.LevelError)
	theme.SetDefaultPallete(true)
	sys.InitWindow(400, 150, "Test", 1, 1.0)
	defer sys.Shutdown()
	gpu.SetBackgroundColor(f32.White)
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
	err := gpu.CaptureToFile("./test-outputs/shadows.png", 0, 0, 400, 100)
	if err != nil {
		slog.Error("Capture to file failed, ", "file", "test-outputs/shadows.png", "error", err.Error())
	}
	img1, err := gpu.LoadImage("./test-assets/shadows.png")
	if err != nil {
		slog.Error("Load image failed, ", "file", "test-assets/shadows.png")
	}
	img2, err := gpu.LoadImage("./test-outputs/shadows.png")
	if err != nil {
		t.Errorf("Load image failed, %s\n", "test-outputs/shadows.png")
	}
	if img2 == nil {
		t.Errorf("Load image failed, %s\n", "test-outputs/shadows.png")
	}
	diff, err := gpu.Compare(img1, img2)
	if diff > 50 {
		t.Errorf("shadows.png difference was %d", diff)
	}
	sys.WindowList.SwapBuffers()
	// Place breakpoint here in order to look at the screen output.
	time.Sleep(1 * time.Millisecond)

}
