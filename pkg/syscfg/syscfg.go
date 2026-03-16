package syscfg

/*
#cgo LDFLAGS: -lsyscfg
#include <stdlib.h>
#include <syscfg/syscfg.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

func SyscfgGet(key string) (string, error) {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	buf := make([]byte, 256)
	cBufPtr := (*C.char)(unsafe.Pointer(&buf[0]))
	ret := C.syscfg_get(nil, cKey, cBufPtr, C.int(len(buf)))
	if ret != 0 {
		return "", fmt.Errorf("syscfg_get failed with code %d", ret)
	}

	return C.GoString(cBufPtr), nil
}

func SyscfgSet(key string, value string) error {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))

	ret := C.syscfg_set(nil, cKey, cValue)
	if ret != 0 {
		return fmt.Errorf("syscfg_set failed with code %d", ret)
	}
	return nil
}
