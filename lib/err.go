package lib

import (
	"fmt"
	"log/slog"
	"os"
)

func ExitWithCode(code int, description string, args ...any) {
	slog.Error(description, args...)
	os.Exit(code)
}

func Exit(description string, args ...any) {
	ExitWithCode(1, description, args...)
}

func ExitIf(condition bool, description string, args ...any) {
	if condition {
		Exit(description, args...)
	}
}

func ExitOn(err error, description string, args ...any) {
	if err != nil {
		description = description + ", " + err.Error()
		Exit(description, args...)
	}
}

func PanicOn(err error, description string, args ...any) {
	if err != nil {
		s := fmt.Sprintf("%s, %s\n", description, args)
		panic(s + ", " + err.Error())
	}
}
