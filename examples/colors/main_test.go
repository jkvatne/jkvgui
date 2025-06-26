package main

import (
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
	"log/slog"
	"testing"
)

func TestColors(t *testing.T) {
	sys.Init()
	defer sys.Shutdown()
	slog.SetLogLoggerLevel(slog.LevelError)
	sys.CreateWindow(0, 0, 800, 600, "Rounded rectangle demo", 2, 2.0)
	sys.StartFrame(theme.Surface.Bg())
	form1()(wid.NewCtx())
	sys.EndFrame()
}

func BenchmarkColors(b *testing.B) {
	sys.Init()
	defer sys.Shutdown()
	b.ResetTimer()
	b.ReportAllocs()
	slog.SetLogLoggerLevel(slog.LevelError)
	sys.CreateWindow(0, 0, 800, 600, "Rounded rectangle demo", 2, 2.0)
	for i := 0; i < b.N; i++ {
		sys.StartFrame(theme.Surface.Bg())
		form1()(wid.NewCtx())
		sys.EndFrame()
	}
}
