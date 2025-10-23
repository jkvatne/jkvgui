package main

import (
	"log/slog"
	"testing"
	"time"

	"github.com/jkvatne/jkvgui/sys"
)

func TestMaterial(t *testing.T) {
	go sys.AbortAfter(time.Second, 1)
	// NB We need to reset some global variables when doing repeated tests
	Cache = nil
	main()
	slog.Info("Exit TestMaterial")
}
