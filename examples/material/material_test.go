package main

import (
	"testing"
	"time"

	"github.com/jkvatne/jkvgui/sys"
)

func TestMaterial(t *testing.T) {
	go sys.AbortAfter(time.Second, 1)
	main()
}
