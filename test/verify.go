package test

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/sys"
	"log/slog"
	"testing"
)

func VerifyScreen(t *testing.T, testName string, w float32, h float32, setup bool) {
	f32.AssertDir("test-outputs")
	err := gpu.CaptureToFile("./test-outputs/"+testName+".png", 0, 0, int(w), int(h))
	if err != nil {
		slog.Error("Capture to file failed, ", "file", "test-outputs/"+testName+".png", "error", err.Error())
	}
	if setup {
		err = gpu.CaptureToFile("./test-assets/"+testName+".png", 0, 0, int(w), int(h))
	}
	img1, err := gpu.LoadImage("./test-assets/" + testName + ".png")
	if err != nil {
		t.Errorf("Load image failed, file /test-assets/%s\n", testName+".png")
	}
	img2, err := gpu.LoadImage("./test-outputs/" + testName + ".png")
	if err != nil {
		t.Errorf("Load image failed, file ./test-outputs/%s\n", testName+".png")
	}
	if img1 == nil {
		t.Errorf("Load image failed, file ./test-assets/%s\n", testName+".png")
	}
	if img2 == nil {
		t.Errorf("Load image failed, file ./test-outputs/%s\n", testName+".png")
	}
	diff, err := gpu.Compare(img1, img2)
	if err != nil {
		t.Errorf("Compare failed, error %v\n", err.Error())
	}
	if diff > 50 {
		t.Errorf("shadows.png difference was %d\n", diff)
	}
	sys.WindowList[0].SwapBuffers()
}
