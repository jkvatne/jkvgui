package sys

import (
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"runtime/pprof"
)

var (
	cpuprofile = flag.String("cpuprof", "", "write cpu profile to `file`, defaults no profile")
	memprofile = flag.String("memprof", "", "write memory profile to `file`, defaults to no profile")
)

// InitializeProfiling will initialize the profiling system.
// It will start the pprof server on http://localhost:6060/debug/pprof/heap
// The profiling files mem.prof and cpu.prof will be written to the current directory.
// To make pdf file, type the followin in the terminal: go tool pprof -pdf CPUprofile > prof.pdf
func InitializeProfiling() {
	if *cpuprofile != "" {
		f, err := os.Create("cpu.prof")
		if err != nil {
			log.Fatal(err)
		}
		_ = pprof.StartCPUProfile(f)
		go func() {
			err := http.ListenAndServe("localhost:6060", nil)
			if err != nil {
				log.Printf("pprof server failed: %v", err)
			}
		}()
	}
}

func TerminateProfiling() {
	if *cpuprofile != "" {
		pprof.StopCPUProfile()
		slog.Info("CPU profile written to", "file", *cpuprofile, "CMD", "go tool pprof -pdf cpu.prof > cpuprof.pdf")
	}
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				slog.Error("could not close memory profile", "error", err)
			}
		}(f)
		// invoke gc, in order to get up-to-date statistics
		runtime.GC()
		// Lookup("allocs") creates a profile similar to go test -memprofile.
		if err := pprof.Lookup("allocs").WriteTo(f, 0); err != nil {
			slog.Error("could not write memory profile", "error", err)
		} else {
			slog.Info("Memory profile written to", "file", *memprofile, "CMD", "go tool pprof -pdf mem.prof > cpuprof.pdf")
		}
	}
}
