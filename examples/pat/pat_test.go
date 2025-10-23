package main

import (
	"testing"
	"time"

	"github.com/jkvatne/jkvgui/sys"
)

func TestPat(t *testing.T) {
	// NB: We need to reset the global variables when running repeated tests!
	Images = nil
	logText = nil
	go sys.AbortAfter(time.Second, 1)
	main()
}
