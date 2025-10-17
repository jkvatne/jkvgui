package main

import (
	"log/slog"
	"testing"

	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
)

func TestColors(t *testing.T) {
	sys.Init()
	defer sys.Shutdown()
	slog.SetLogLoggerLevel(slog.LevelError)
	w := sys.CreateWindow(0, 0, 800, 600, "Rounded rectangle demo", 2, 2.0)
	w.StartFrame(theme.Surface.Bg())
	form1()(wid.NewCtx(w))
	w.EndFrame()
}

func BenchmarkColors(b *testing.B) {
	sys.Init()
	defer sys.Shutdown()
	b.ResetTimer()
	b.ReportAllocs()
	slog.SetLogLoggerLevel(slog.LevelError)
	w := sys.CreateWindow(0, 0, 800, 600, "Rounded rectangle demo", 2, 2.0)
	for i := 0; i < b.N; i++ {
		w.StartFrame(theme.Surface.Bg())
		form1()(wid.NewCtx(w))
		w.EndFrame()
	}
}
