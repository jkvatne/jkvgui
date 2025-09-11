package font_test

import (
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/gpu/font"
	"github.com/jkvatne/jkvgui/sys"
	"log/slog"
	"testing"
)

var testStrings = []string{
	"WlæøåÆØÅ$€ÆØÅ",
	"55 Some text with special characters æøåÆØÅ$€ and some more arbitary text to make a very long line that will be broken for wrap-around (or elipsis)",
	"Hello world, this is a test of the emergency broadcast system. This is only a test. The result will be evaluated. Do not panic.",
	"Averyveryveryverylongwordthatshouldbesplitted",
	"Hi",
	"one         two         three    four five",
}

var expected = []int{
	3,
	2,
	20,
	5,
	1,
	4,
}

var limit = []float32{
	55.0,
	829.3,
	55.0,
	55.0,
	55.0,
	55.0,
}

func TestSplit(t *testing.T) {
	sys.Init()
	defer sys.Shutdown()
	slog.SetLogLoggerLevel(slog.LevelError)
	sys.CreateWindow(0, 0, 800, 800, "Splittest", 2, 1.5)
	font.LoadFontBytes(gpu.Normal14, "RobotoNormal", font.Roboto400, 14, 400)
	for i, s := range testStrings {
		strings := font.Split(s, limit[i], font.Fonts[gpu.Normal14])
		if len(strings) != expected[i] {
			t.Errorf("Test %d expected %d strings, got %d\n", i, expected[i], len(strings))
		}
	}
}
