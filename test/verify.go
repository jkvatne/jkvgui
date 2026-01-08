package test

import (
	"fmt"
	"testing"

	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/sys"
)

func VerifyScreen(t *testing.T, win *sys.Window, testName string, w float32, h float32, setup bool) error {
	f32.AssertDir("test-outputs")
	err := sys.CaptureToFile(win, "./test-outputs/"+testName+".png", 0, 0, int(w), int(h))
	if err != nil {
		return fmt.Errorf("Capture to file failed, %v", err)
	}
	if setup {
		err = sys.CaptureToFile(win, "./test-assets/"+testName+".png", 0, 0, int(w), int(h))
		if err != nil {
			return fmt.Errorf("Capture of asset failed, %s\n", err.Error())
		}
	}
	img1, err := gpu.LoadImage("./test-assets/" + testName + ".png")
	if err != nil {
		return fmt.Errorf("Load image failed, file /test-assets/%s\n", testName+".png")
	}
	img2, err := gpu.LoadImage("./test-outputs/" + testName + ".png")
	if err != nil {
		return fmt.Errorf("Load image failed, file ./test-outputs/%s\n", testName+".png")
	}
	if img1 == nil {
		return fmt.Errorf("Load image failed, file ./test-assets/%s\n", testName+".png")
	}
	if img2 == nil {
		return fmt.Errorf("Load image failed, file ./test-outputs/%s\n", testName+".png")
	}
	diff, err := gpu.Compare(img1, img2)
	if err != nil {
		return fmt.Errorf("Compare failed, error %v\n", err.Error())
	}
	if diff > 600 {
		return fmt.Errorf(testName+".png difference was %d\n", diff)
	}
	return nil
}
