package img

import (
	"image"
	"image/png"
	"log"
	"log/slog"
	"os"
)

func NewImage(filename string) err {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		slog.Error("Failed to open ", filename)
	}
	return img
}
