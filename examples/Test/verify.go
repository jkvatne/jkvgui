package main

import (
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"testing"
)

func VerifyScreen(t *testing.T, testname string, w int, h int, setup bool) {
	err := gpu.CaptureToFile("./test-outputs/"+testname+".png", 0, 0, w, h)
	if err != nil {
		slog.Error("Capture to file failed, ", "file", "test-outputs/"+testname+".png", "error", err.Error())
	}
	if setup {
		err = gpu.CaptureToFile("./test-assets/"+testname+".png", 0, 0, w, h)
	}
	img1, err := gpu.LoadImage("./test-assets/" + testname + ".png")
	if err != nil {
		slog.Error("Load image failed, ", "file", "test-assets/"+testname+".png")
	}
	img2, err := gpu.LoadImage("./test-outputs/" + testname + ".png")
	assert.Nil(t, err)
	assert.NotNil(t, img1)
	assert.NotNil(t, img2)
	if img1 != nil && img2 != nil {
		diff, err := gpu.Compare(img1, img2)
		assert.Nil(t, err)
		slog.Info("shadows.png difference was", "diff", diff)
		assert.LessOrEqual(t, diff, int64(50))
	}
	gpu.Window.SwapBuffers()
}
