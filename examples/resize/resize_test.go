package main

import (
	"testing"
	"time"

	"github.com/jkvatne/jkvgui/sys"
)

func TestResize(t *testing.T) {
	go sys.AbortAfter(time.Second, 2)
	main()
}
