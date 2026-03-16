package sysevent

/*
#cgo LDFLAGS: -lsysevent
#include <stdlib.h>
#include <sysevent/sysevent.h>
*/
import "C"

import "fmt"

const (
	SeServerWellKnownPort int    = 52367
	UDSPath               string = "/tmp/syseventd_connection"
	SeVersion             int    = 1

	TupleFlagNormal uint = 0x00000000
	TupleFlagSerial uint = 0x00000001
	TupleFlagEvent  uint = 0x00000002
)

type Sysevent struct {
	syseventFD    C.int
	syseventToken C.token_t
	asyncIDs      map[*AsyncID]C.async_id_t
}

type AsyncID struct{}

func New() *Sysevent {
	return &Sysevent{}
}

// Open estableshes a connection to the sysevent daemon
//
// ip: The IP address of the sysevent daemon. This may be a hostname.
// port: The port number on which the sysevent daemon is listening
// version: The version of the client
// id: The name of the client
//
// Returns an error if the underlying sysevent_open call fails
func (s *Sysevent) Open(ip string, port int, version int, id string) error {
	s.syseventFD = C.sysevent_open(C.CString(ip), C.ushort(port), C.int(version), C.CString(id), &s.syseventToken)
	if s.syseventFD < 0 {
		return fmt.Errorf("sysevent_open failed with code %d", s.syseventFD)
	}

	return nil
}

// SetOptions sends a request to the sysevent daemon to set options
//
// name: The tuple to set
// flags: The value to set for the tuple
//
// Returns an error if the underlying sysevent_set_options call fails
func (s *Sysevent) SetOptions(name string, flags uint) error {
	ret := C.sysevent_set_options(s.syseventFD, s.syseventToken, C.CString(name), C.uint(flags))
	if ret != 0 {
		return fmt.Errorf("sysevent_set_options failed with code %d", ret)
	}
	return nil
}

// int sysevent_setnotification(const int fd, const token_t token, char *subject, async_id_t *async_id);
func (s *Sysevent) SetNotification(subject string, asyncID *AsyncID) error {
	a, ok := s.asyncIDs[asyncID]
	if !ok {
		var asyncID_C C.async_id_t
		s.asyncIDs[asyncID] = asyncID_C
		a = asyncID_C
	}

	ret := C.sysevent_setnotification(s.syseventFD, s.syseventToken, C.CString(subject), &a)
	if ret != 0 {
		return fmt.Errorf("sysevent_setnotification failed with code %d", ret)
	}

	return nil
}
