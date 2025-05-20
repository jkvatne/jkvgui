package sys

import (
	"fmt"
	"log/slog"
	"runtime/debug"
	"strings"
)

var (
	MainPath = "(developement build)"
	Tag      = "(developement build)"
	Hash     = "(developement build)"
)

// GetBuildInfo will read the build info from the go.mod file and set the variables
func GetBuildInfo() {
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
}

func PrintBuildInfo() {
	fmt.Printf("Buildinfo hash=%s, tag=%s, path=%s\n", Hash, Tag, MainPath)
}

func LogBuildInfo() {
	slog.Info("Buildinfo", "hash", Hash, "tag", Tag, "url", MainPath)
}
