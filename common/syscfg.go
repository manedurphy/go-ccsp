package common

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

func SyscfgGet(namespace string, key string) (string, error) {
	cNamespace := C.CString(namespace)
	defer C.free(unsafe.Pointer(cNamespace))

	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	buf := make([]byte, 256)
	cBufPtr := (*C.char)(unsafe.Pointer(&buf[0]))
	ret := C.syscfg_get(cNamespace, cKey, cBufPtr, C.int(len(buf)))
	if ret != 0 {
		return "", fmt.Errorf("syscfg_get failed with code %d", ret)
	}
	defer C.free(unsafe.Pointer(cBufPtr))

	return C.GoString(cBufPtr), nil
}
