package font

import (
	"github.com/jkvatne/jkvgui/gl"
	"github.com/jkvatne/jkvgui/gpu"
	"testing"
)

var testStrings = []string{
	"Hello world, this is a test of the emergency broadcast system. This is only a test. The result will be evaluated. Do not panic.",
	"Averyveryveryverylongwordthatshouldbesplitted",
	"Hi",
	"one         two         three    four five",
}

var expected = []int{
	20,
	5,
	1,
	4,
}

func TestSplit(t *testing.T) {
	err := gl.Init()
	if err != nil {
		t.Errorf("Failed to initialize OpenGL, %v", err)
	} else {
		LoadFontBytes(gpu.Normal14, "RobotoNormal", Roboto400, 14, 400)
		for i, s := range testStrings {
			strings := Split(s, 55.0, Fonts[gpu.Normal14])
			if len(strings) != expected[i] {
				t.Errorf("Test %d expected %d strings, got %d\n", i, expected[i], len(strings))
			}
		}
	}
}
