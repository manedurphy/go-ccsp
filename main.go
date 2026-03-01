package main

/*
#cgo LDFLAGS: -lhal_mta
#include <ccsp/mta_hal.h>
*/
import "C"

import (
	"encoding/xml"
	"fmt"
	"os"
	"syscall"
	"time"

	"go-ccsp/common"
)

var (
	cfg       common.CcspComponentCfg
	subsystem string
)

func main() {
	// TODO: Initialize a logger
	componentCfgXML, err := os.ReadFile("/usr/ccsp/mta/CcspMta.cfg")
	if err != nil {
		panic(err)
	}

	err = xml.Unmarshal(componentCfgXML, &cfg)
	if err != nil {
		panic(err)
	}

	runAsDaemon := true
	for _, arg := range os.Args {
		if arg == "-subsys" {
			subsystem = os.Args[2]
		} else if arg == "-c" {
			runAsDaemon = false
		}
	}

	if runAsDaemon {
		// Daemonize the process
		// TODO: Implement in commmon package
	}

	if syscall.Getuid() == 0 {
		// Drop privileges to "non-root" user
		// runtime.LockOSThread()
		err = syscall.Setgid(950)
		if err != nil {
			panic(err)
		}

		err := syscall.Setuid(950)
		if err != nil {
			panic(err)
		}

		// runtime.UnlockOSThread()
	}

	fmt.Printf("Running MTA Agent: user_id=%d group_id=%d \n", syscall.Getuid(), syscall.Getgid())
	for {
		time.Sleep(10 * time.Second)
	}

	/*
		var mtaDHCPInfo C.MTAMGMT_MTA_DHCP_INFO
		C.mta_hal_GetDHCPInfo(&mtaDHCPInfo)
		fmt.Println("Called the HAL!")
		fmt.Printf("MTA DHCP Info: %+v\n", mtaDHCPInfo)

		goIP := net.IPv4(mtaDHCPInfo.IPAddress[0], mtaDHCPInfo.IPAddress[1], mtaDHCPInfo.IPAddress[2], mtaDHCPInfo.IPAddress[3])

		fmt.Printf("IP address: %s\n", goIP.String())
		fmt.Printf("MAC address: %s\n", C.GoString(&mtaDHCPInfo.MACAddress[0]))
	*/
}
