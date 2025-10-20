package main

import (
	"os"
	"testing"
	"time"

	"github.com/jkvatne/jkvgui/sys"
)

func stopper() {
	time.Sleep(100 * time.Millisecond)
	sys.WindowList = nil
}

func TestHello(t *testing.T) {
	go stopper()
	main()
}

func TestMain(m *testing.M) {
	exitcode := m.Run()
	os.Exit(exitcode)
}
