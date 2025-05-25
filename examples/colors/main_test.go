package main

import (
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/input"
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
	"log/slog"
	"testing"
)

func TestColors(t *testing.T) {
	input.Initialize()
	slog.SetLogLoggerLevel(slog.LevelError)
	gpu.InitWindow(0, 0, "Rounded rectangle demo", 2, 2.0)
	defer input.Shutdown()
	input.InitializeWindow()
	sys.StartFrame(theme.Surface.Bg())
	form1()(wid.NewCtx())
	sys.EndFrame(0)
}

func BenchmarkColors(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	input.Initialize()
	slog.SetLogLoggerLevel(slog.LevelError)
	gpu.InitWindow(0, 0, "Rounded rectangle demo", 2, 2.0)
	defer input.Shutdown()
	input.InitializeWindow()
	for i := 0; i < b.N; i++ {
		sys.StartFrame(theme.Surface.Bg())
		form1()(wid.NewCtx())
		sys.EndFrame(0)
	}
}
