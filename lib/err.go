package lib

import (
	"fmt"
	"log"
	"os"
)

func ExitWithCode(code int, description string, args ...any) {
	// Check flag to avoid recursive calls from closers
	log.Printf(description, args...)
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
