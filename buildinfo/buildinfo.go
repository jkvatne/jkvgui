package buildinfo

import (
	"log/slog"
	"runtime/debug"
	"strings"
)

var (
	MainPath = "(development build)"
	Tag      = "(development build)"
	Hash     = "(development build)"
)

// Get will read the build info from the go.mod file and set the variables
func Get() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		slog.Error("Could not read build info")
		return
	}
	s := info.Main.Version
	if s != "" {
		words := strings.Split(s, "-")
		Tag = words[0]
	}
	MainPath = info.Main.Path
	for _, setting := range info.Settings {
		key := setting.Key
		if key == "vcs.revision" {
			Hash = setting.Value[:8]
		}
	}
	slog.Info("BuildInfo", "hash", Hash, "tag", Tag, "url", MainPath)
}
