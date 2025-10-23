package main

import (
	"log/slog"
	"testing"
	"time"

	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/wid"
)

func TestColors(t *testing.T) {
	go sys.AbortAfter(time.Second, 1)
	main()
}

func BenchmarkColors(b *testing.B) {
	sys.Init()
	defer sys.Shutdown()
	b.ResetTimer()
	b.ReportAllocs()
	slog.SetLogLoggerLevel(slog.LevelError)
	w := sys.CreateWindow(0, 0, 800, 600, "Rounded rectangle demo", 2, 2.0)
	for i := 0; i < b.N; i++ {
		w.StartFrame()
		form1(w)(wid.NewCtx(w))
		w.EndFrame()
	}
}
