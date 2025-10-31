package buildinfo

import (
	"log/slog"
	"runtime/debug"
)

var (
	Revision    = "Unknown"
	CompileTime = "Unknown"
	Info        *debug.BuildInfo
	Dirty       bool
)

// Get will read the build info from the go.mod file and set the variables
func Get() {
	var ok bool
	Info, ok = debug.ReadBuildInfo()
	if !ok {
		slog.Error("Could not read build info")
		return
	}
	for _, setting := range Info.Settings {
		key := setting.Key
		if key == "vcs.revision" {
			Revision = setting.Value[:8]
		}
		if setting.Key == "vcs.modified" {
			Dirty = setting.Value == "true"
		}
		if setting.Key == "vcs.time" {
			CompileTime = setting.Value
		}
		if Dirty {
			Revision += "-dirty"
		}
	}
	slog.Info("BuildInfo", "revision", Revision, "compiletime", CompileTime)
}
