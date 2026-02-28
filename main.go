package main

/*
#cgo CFLAGS: -I${SRCDIR}
#cgo LDFLAGS: -L${SRCDIR} -L${SRCDIR}/rootfs/usr/lib -lhal_mta -llattice -litc_rpc -lnetsnmp -lsysevent -lhal_util -lmocactl -lssl -lcrypto

#include "mta_hal.h"
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
