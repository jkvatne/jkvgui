package buildinfo

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime/debug"
)

var (
	Hash        = "Unknown"
	CompileTime = "Unknown"
	Info        *debug.BuildInfo
	Dirty       bool
	Version     = "Unknown"
	GoVersion   = "Unknown"
	ExeName     = "Unknown"
)

// Get will read the build info from the go.mod file and set the variables
func Get() {
	var ok bool
	Info, ok = debug.ReadBuildInfo()
	if !ok {
		slog.Error("Could not read build info")
		return
	}
	exePath, err := os.Executable()
	if err == nil {
		ExeName = filepath.Base(exePath)
	}
	GoVersion = Info.GoVersion
	Version = Info.Main.Version
	for _, setting := range Info.Settings {
		key := setting.Key
		if key == "vcs.revision" {
			Hash = setting.Value[:8]
		}
		if setting.Key == "vcs.modified" {
			Dirty = setting.Value == "true"
		}
		if setting.Key == "vcs.time" {
			CompileTime = setting.Value
		}
	}
	if Dirty {
		Hash += "-dirty"
	}
}

func init() {
	Get()
	fmt.Printf("Running \"%s\" with hash=\"%s\", tag=%s, compiled %v\n", ExeName, Hash, Version, CompileTime)
}
