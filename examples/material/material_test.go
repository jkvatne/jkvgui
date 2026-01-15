package main

import (
	"log/slog"
	"testing"
	"time"

	"github.com/jkvatne/jkvgui/sys"
)

func TestMaterial(t *testing.T) {
	go sys.AbortAfter(time.Second, 1)
	main()
	slog.Info("Exit TestMaterial")
}
