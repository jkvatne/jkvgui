package main

import (
	"testing"
	"time"

	"github.com/jkvatne/jkvgui/sys"
)

func TestPat(t *testing.T) {
	go sys.AbortAfter(time.Second, 1)
	main()
}
