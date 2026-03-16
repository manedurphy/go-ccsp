package rbus

/*
#cgo LDFLAGS: -lrbus
#include <rbus/rbus.h>
*/
import "C"

import (
	"fmt"
)

var rbusHandle C.rbusHandle_t = nil

func Open(componentName string) error {
	ret := C.rbus_open(&rbusHandle, C.CString(componentName))
	if ret != C.RBUS_ERROR_SUCCESS {
		return fmt.Errorf("rbus_open failed: code=%d err=%s", ret, C.GoString(C.rbusError_ToString(ret)))
	}

	return nil
}
