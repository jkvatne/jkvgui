package test

import (
	"flag"
	"testing"

	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/sys"
)

var updateAssets bool

func init() {
	flag.BoolVar(&updateAssets, "test.update", false, "Save the captured screen to ./test-assets, making this the reference image.")
}

func VerifyScreen(t *testing.T, win *sys.Window, testName string, w float32, h float32, limit int64) {
	f32.AssertDir("test-outputs")
	err := sys.CaptureToFile(win, "./test-outputs/"+testName+".png", 0, 0, int(w), int(h))
	if err != nil {
		t.Error("Capture to file failed,", err)
	}
	if updateAssets {
		err = sys.CaptureToFile(win, "./test-assets/"+testName+".png", 0, 0, int(w), int(h))
		if err != nil {
			t.Error("Capture of asset failed,", err.Error())
			return
		}
	}
	img1, err := gpu.LoadImage("./test-assets/" + testName + ".png")
	if err != nil {
		t.Error("Load image from test-assets failed, file", testName+".png")
		return
	}
	img2, err := gpu.LoadImage("./test-outputs/" + testName + ".png")
	if err != nil {
		t.Error("Load image from test-outputs failed, file", testName+".png")
		return
	}
	diff, err := gpu.Compare(img1, img2)
	if err != nil {
		t.Error(testName+".png compare failed, error", err.Error())
	}
	if diff > limit {
		t.Error(testName+".png image difference was", diff)
	}
}
