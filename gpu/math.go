package gpu

import (
	"unsafe"
)

type eface struct {
	typ, val unsafe.Pointer
}

func ptr(arg interface{}) unsafe.Pointer {
	return (*eface)(unsafe.Pointer(&arg)).val
}

func TagsEqual(a, b any) bool {
	return ptr(a) == ptr(b)
}
