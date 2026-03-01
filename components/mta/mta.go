package mta

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

func Run() error {
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

	if syscall.Getuid() == 0 {
		// Drop privileges to "non-root" user
		err = syscall.Setgid(950)
		if err != nil {
			panic(err)
		}

		err := syscall.Setuid(950)
		if err != nil {
			panic(err)
		}
	}

	if runAsDaemon {
		// Daemonize the process
		// TODO: Implement in commmon package, or maybe not at all since systemd is the default init system now
	}

	// Write PID to file
	err = os.WriteFile("/var/run/CcspMtaAgentSsp.pid", []byte(fmt.Sprintf("%d", os.Getpid())), 0o644)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Running MTA Agent: user_id=%d group_id=%d \n", syscall.Getuid(), syscall.Getgid())
	for {
		time.Sleep(10 * time.Second)
	}
}
