package main

import (
	"testing"
	"time"

	"github.com/jkvatne/jkvgui/sys"
)

func TestDemo(t *testing.T) {
	go sys.AbortAfter(time.Second, 2)
	main()
}
