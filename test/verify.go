package main

import (
	"github.com/jkvatne/jkvgui/gpu"
	"log/slog"
	"testing"
)

func VerifyScreen(t *testing.T, testname string, w float32, h float32, setup bool) {
	err := gpu.CaptureToFile("./test-outputs/"+testname+".png", 0, 0, int(w), int(h))
	if err != nil {
		slog.Error("Capture to file failed, ", "file", "test-outputs/"+testname+".png", "error", err.Error())
	}
	if setup {
		err = gpu.CaptureToFile("./test-assets/"+testname+".png", 0, 0, int(w), int(h))
	}
	img1, err := gpu.LoadImage("./test-assets/" + testname + ".png")
	if err != nil {
		t.Errorf("Load image failed, file /test-assets/%s\n", testname+".png")
	}
	img2, err := gpu.LoadImage("./test-outputs/" + testname + ".png")
	if err != nil {
		t.Errorf("Load image failed, file ./test-outputs/%s\n", testname+".png")
	}
	if img1 == nil {
		t.Errorf("Load image failed, file ./test-assets/%s\n", testname+".png")
	}
	if img2 == nil {
		t.Errorf("Load image failed, file ./test-outputs/%s\n", testname+".png")
	}
	diff, err := gpu.Compare(img1, img2)
	if err != nil {
		t.Errorf("Compare failed, error %v\n", err.Error())
	}
	if diff > 50 {
		t.Errorf("shadows.png difference was %d\n", diff)
	}
	gpu.Window.SwapBuffers()
}
