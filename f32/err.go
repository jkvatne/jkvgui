package f32

import (
	"errors"
	"log/slog"
	"os"
)

// Exit will abort the program with the given exit code
// Before exiting, it will log the description with slog.Error
func Exit(code int, description string, args ...any) {
	slog.Error(description, args...)
	os.Exit(code)
}

// ExitIf condition is true
func ExitIf(condition bool, description string, args ...any) {
	if condition {
		Exit(1, description, args...)
	}
}

// ExitOn an error
func ExitOn(err error, description string, args ...any) {
	if err != nil {
		description = description + ", " + err.Error()
		Exit(1, description, args...)
	}
}

// AssertDir will create the path if it does not already exist.
func AssertDir(path string) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err2 := os.Mkdir(path, os.ModePerm)
		if err2 != nil {
			slog.Error("AssertDir failed", "err", err.Error())
		}
	}
}
