package main

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/theme"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {

}

func TestShadows(t *testing.T) {
	theme.SetDefaultPallete(true)
	_ = gpu.InitWindow(400, 100, "Test", 1)
	defer gpu.Shutdown()
	gpu.UpdateResolution()
	gpu.StartFrame(f32.White)
	r := f32.Rect{10, 10, 30, 20}
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
	err := gpu.CaptureToFile("./test-outputs/shadows.png", 0, 0, 400, 200)
	if err != nil {
		slog.Error("Capture to file failed, ", "file", "test-outputs/shadows.png", "error", err.Error())
	}
	img1, err := gpu.LoadImage("./test-assets/shadows.png")
	if err != nil {
		slog.Error("Load image failed, ", "file", "test-assets/shadows.png")
	}
	img2, err := gpu.LoadImage("./test-outputs/shadows.png")
	assert.Nil(t, err)
	assert.NotNil(t, img2)
	diff, err := gpu.Compare(img1, img2)
	slog.Info("shadows.png difference was", "diff", diff)
	assert.LessOrEqual(t, diff, int64(50))
}
