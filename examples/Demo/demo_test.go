package main

import (
	"testing"
	"time"

	"github.com/jkvatne/jkvgui/sys"
)

func TestDemoNonThreaded(t *testing.T) {
	go sys.AbortAfter(time.Second, 2)
	*threaded = false
	main()
}

func TestDemoThreaded(t *testing.T) {
	go sys.AbortAfter(time.Second, 2)
	*threaded = true
	main()
}
