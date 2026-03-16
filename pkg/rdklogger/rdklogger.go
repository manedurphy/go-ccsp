package rdklogger

/*
#cgo LDFLAGS: -lrdkloggers
#include <stdlib.h>
#include <unistd.h>
#include <rdk_logger.h>

static int is_rdk_logger_enabled() {
#ifdef FEATURE_SUPPORT_RDKLOG
	return 1;
#else
    return 0;
#endif
}

static int init_rdk_logger(void) {
    return RDK_LOGGER_INIT();
}

static void log_message(rdk_LogLevel level, const char *module, const char *message) {
    rdk_logger_msg_printf(level, module, "%s", message);
}
*/
import "C"

import (
	"fmt"
	"unsafe"
)

type RDKLogLevel int

const (
	RDK_LOG_FATAL  RDKLogLevel = 0
	RDK_LOG_ERROR  RDKLogLevel = 1
	RDK_LOG_WARN   RDKLogLevel = 2
	RDK_LOG_NOTICE RDKLogLevel = 3
	RDK_LOG_INFO   RDKLogLevel = 4
	RDK_LOG_DEBUG  RDKLogLevel = 5
	RDK_LOG_TRACE  RDKLogLevel = 6
	RDK_LOG_NONE   RDKLogLevel = 7
)

func InitializeRDKLogger() error {
	if C.is_rdk_logger_enabled() == 0 {
		ret := C.init_rdk_logger()
		if ret != 0 {
			return fmt.Errorf("failed to initialize RDK Logger, error code: %d", ret)
		}
	}

	return nil
}

func RDKLog(level RDKLogLevel, module string, message string) {
	cModule := C.CString(module)
	defer C.free(unsafe.Pointer(cModule))

	cMessage := C.CString(message)
	defer C.free(unsafe.Pointer(cMessage))

	C.log_message(C.rdk_LogLevel(level), cModule, cMessage)
}
