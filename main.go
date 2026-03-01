package main

/*
#cgo LDFLAGS: -lhal_mta
#include <ccsp/mta_hal.h>
*/
import "C"

import (
	"fmt"
)

func main() {
	var mtaDHCPInfo C.MTAMGMT_MTA_DHCP_INFO
	C.mta_hal_GetDHCPInfo(&mtaDHCPInfo)
	fmt.Println("Called the HAL!")
}
