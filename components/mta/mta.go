package mta

/*
#cgo LDFLAGS: -lhal_mta
#include <ccsp/mta_hal.h>
*/
import "C"

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net"
	"os"
	"syscall"
	"time"

	"go-ccsp/internal/ccsp"
	"go-ccsp/pkg/rbus"
	"go-ccsp/pkg/rdklogger"
	"go-ccsp/pkg/syscfg"
	"go-ccsp/pkg/sysevent"
)

type LineTableEntry struct {
	InstanceNumber     uint32
	LineNumber         uint32
	Status             uint32
	HazardousPotential string
	ForeignEMF         string
	ResistiveFaults    string
	ReceiverOffHook    string
	RingerEquivalency  string
	CAName             string
	CAPort             uint32
	MWD                uint32
	CallsNumber        uint32
	Calls              struct {
		Codec                    string
		RemoteCodec              string
		CallStartTime            string
		CallEndTime              string
		CWErrorRate              string
		PktLossConcealment       string
		JitterBufferAdaptive     bool
		Originator               bool
		RemoteIPAddress          net.IP
		CallDuration             uint32
		CWErrors                 string
		SNR                      string
		MicroReflections         string
		DownstreamPower          string
		UpstreamPower            string
		EQIAverage               string
		EQIMinimum               string
		EQIMaximum               string
		EQIInstantaneous         string
		MOS_LQ                   string
		MOS_CQ                   string
		EchoReturnLoss           string
		SignalLevel              string
		NoiseLevel               string
		LossRate                 string
		DiscardRate              string
		BurstDensity             string
		GapDensity               string
		BurstDuration            string
		GapDuration              string
		RoundTripDelay           string
		Gmin                     string
		RFactor                  string
		ExternalRFactor          string
		JitterBufRate            string
		JBNominalDelay           string
		JBMaxDelay               string
		JBAbsMaxDelay            string
		TxPackets                string
		TxOctets                 string
		RxPackets                string
		RxOctets                 string
		PacketLoss               string
		IntervalJitter           string
		RemoteIntervalJitter     string
		RemoteMOS_LQ             string
		RemoteMOS_CQ             string
		RemoteEchoReturnLoss     string
		RemoteSignalLevel        string
		RemoteNoiseLevel         string
		RemoteLossRate           string
		RemotePktLossConcealment string
		RemoteDiscardRate        string
		RemoteBurstDensity       string
		RemoteGapDensity         string
		RemoteBurstDuration      string
		RemoteGapDuration        string
		RemoteRoundTripDelay     string
		RemoteGmin               string
		RemoteRFactor            string
		RemoteExternalRFactor    string
		RemoteJitterBufRate      string
		RemoteJBNominalDelay     string
		RemoteJBMaxDelay         string
		RemoteJBAbsMaxDelay      string
	}
	CallsUpdateTime  uint32
	OverCurrentFault uint32
}

type DataModel struct {
	Oid                uint32
	LineTable          []LineTableEntry
	LineTableCount     uint32
	ServiceClassNumber uint32
	ServiceClass       []struct {
		ServiceClassName string
	}
	ServiceClassUpdateTime uint32
	ServiceFlowNumber      uint32
	ServiceFlow            []struct {
		SFID               uint32
		ServiceClassName   string
		Direction          string
		ScheduleType       uint32
		DefaultFlow        bool
		NomGrantInterval   uint32
		UnsolicitGrantSize uint32
		TolGrantJitter     uint32
		NomPollInterval    uint32
		MinReservedPkt     uint32
		MaxTrafficRate     uint32
		MinReservedRate    uint32
		MaxTrafficBurst    uint32
		TrafficType        string
		NumberOfPackets    uint32
	}
	ServiceFlowUpdateTime uint32
	Handsets              []struct {
		InstanceNumber  uint32
		Status          bool
		LastActiveTime  string
		HandsetName     string
		HandsetFirmware string
		OperatingTN     string
		SupportedTN     string
	}
	HandsetsNumber     uint32
	HandsetsUpdateTime uint32
	Pktc               struct {
		MtaDevEnabled                bool
		SigDefCallSigTos             uint32
		SigDefMediaStreamTos         uint32
		MtaDevRealmOrgName           uint32
		MtaDevCmsKerbRealmName       uint32
		MtaDevCmsIpsecCtrl           uint32
		MtaDevCmsSolicitedKeyTimeout uint32
		MtaDevRealmPkinitGracePeriod uint32
	}
	Dect struct {
		RegisterDectHandset   uint32
		DeregisterDectHandset uint32
		HardwareVersion       string
		RFPI                  string
		SoftwareVersion       string
		PIN                   string
	}
	DSXLogNumber uint32
	DSXLog       []struct {
		Time        string
		Description string
		ID          uint32
		Level       uint32
	}
	DSXLogUpdateTime uint32
	MtaLogConfig     struct {
		EnableDECTLog bool
		EnableMTALog  bool
	}
	MtaLogNumber     uint32
	MtaLogUpdateTime uint32
	MtaLog           []struct {
		Index       uint32
		EventID     uint32
		EventLevel  string
		Time        string
		Description string
	}
	DectLogNumber     uint32
	DectLogUpdateTime uint32
	DectLog           []struct {
		Index       uint32
		EventID     uint32
		EventLevel  string
		Time        string
		Description string
	}
	BatteryInfo struct {
		ModelNumber             string
		SerialNumber            string
		PartNumber              string
		ChargerFirmwareRevision string
	}
	MTAProvInfo struct {
		StartupIPMode                  uint32
		IPv4PrimaryDhcpServerOptions   string
		IPv4SecondaryDhcpServerOptions string
		IPv6PrimaryDhcpServerOptions   string
		IPv6SecondaryDhcpServerOptions string
	}
}

type MTAAgentConfig struct {
	ccsp.CcspComponentCfg
	EthernetWANEnabled   bool // -DENABLE_ETH_WAN
	ErouterDHCPOptionMTA bool // -DEROUTER_DHCP_OPTION_MTA
	DMLConfigPath        string
}

type MTAAgent struct {
	DataModel DataModel
	cfg       *MTAAgentConfig
	sysevent  *sysevent.Sysevent
}

func New(cfg *MTAAgentConfig) *MTAAgent {
	return &MTAAgent{
		DataModel: DataModel{},
		cfg:       cfg,
		sysevent:  sysevent.New(),
	}
}

func (m *MTAAgent) InitializeDataModel() error {
	// Initialize the MTA HAL database
	ret := C.mta_hal_InitDB() // CosaDmlMTAInit
	if ret != C.RETURN_OK {
		return fmt.Errorf("mta_hal_InitDB failed with code %d", ret)
	}

	// TODO: Conditional logic for defined(VOICE_MTA_SUPPORT)

	// Initialize the line table
	lineTableNumberOfEntries := C.mta_hal_LineTableGetNumberOfEntries()
	for i := uint32(0); i < uint32(lineTableNumberOfEntries); i++ {
		m.DataModel.LineTable = append(m.DataModel.LineTable, LineTableEntry{
			InstanceNumber: uint32(i + 1),
		})
	}

	mtaDMLConfigBytes, err := os.ReadFile("/usr/ccsp/mta/mta_json_dml.json")
	if err != nil {
		return fmt.Errorf("failed to read MTA DML config file: err=%s", err)
	}

	var mtaDMLConfig map[string]any
	err = json.Unmarshal(mtaDMLConfigBytes, &mtaDMLConfig)
	if err != nil {
		return fmt.Errorf("failed to parse MTA DML config file: err=%s", err)
	}

	var processMap func(m map[string]any, key string)
	processMap = func(m map[string]any, key string) {
		val, ok := m[key]
		if !ok {
			return
		}

		switch v := val.(type) {
		case map[string]any:
			for k, subVal := range v {
				if subVal == "List_Of_Def" {
					// TODO: Register all parameters with RBUS
				} else {
					processMap(v, k)
				}
			}
		default:
			fmt.Printf("Value for key %s is not a map\n", key)
		}
	}

	processMap(mtaDMLConfig, "Device")

	return nil
}

func (d *DataModel) Create() error {
	return nil
}

func (d *DataModel) Remove() error {
	return nil
}

func (d *DataModel) Initialize() error {
	return nil
}

const (
	loggerModuleName string = "LOG.RDK.MTA"
	configFilePath   string = "/usr/ccsp/mta/CcspMta.cfg"
	// pidFilePath           string = "/var/run/CcspMtaAgentSsp.pid"
	pidFilePath string = "/var/tmp/go-ccsp.pid" // TODO: Revert to original path after testing
)

func (m *MTAAgent) Run(subSystem string) error {
	// Initialize RDK Logger
	rdklogger.InitializeRDKLogger()

	componentCfgXML, err := os.ReadFile(configFilePath)
	if err != nil {
		rdklogger.RDKLog(rdklogger.RDK_LOG_ERROR, loggerModuleName, fmt.Sprintf("Failed to read config file: err=%s\n", err))
		os.Exit(1)
	}

	err = xml.Unmarshal(componentCfgXML, &m.cfg)
	if err != nil {
		rdklogger.RDKLog(rdklogger.RDK_LOG_ERROR, loggerModuleName, fmt.Sprintf("Failed to parse config file: err=%s\n", err))
		os.Exit(1)
	}

	// Clean up
	componentCfgXML = nil

	rdklogger.RDKLog(rdklogger.RDK_LOG_INFO, loggerModuleName, fmt.Sprintf("Starting MTA Agent: component_name=%s subsystem=%s\n", m.cfg.ComponentName, subSystem))

	// if syscall.Getuid() == 0 {
	// 	// Drop privileges to "non-root" user
	// 	err = syscall.Setgid(950)
	// 	if err != nil {
	// 		rdklogger.RDKLog(rdklogger.RDK_LOG_ERROR, loggerModuleName, fmt.Sprintf("Failed to set GID: err=%s\n", err))
	// 		os.Exit(1)
	// 	}
	//
	// 	err = syscall.Setuid(950)
	// 	if err != nil {
	// 		rdklogger.RDKLog(rdklogger.RDK_LOG_ERROR, loggerModuleName, fmt.Sprintf("Failed to set UID: err=%s\n", err))
	// 		os.Exit(1)
	// 	}
	// }

	// Write PID to file
	err = os.WriteFile(pidFilePath, []byte(fmt.Sprintf("%d", os.Getpid())), 0o644)
	if err != nil {
		panic(err)
	}

	isEthWANEnabled, err := syscfg.SyscfgGet("eth_wan_enabled")
	if err != nil {
		rdklogger.RDKLog(rdklogger.RDK_LOG_ERROR, loggerModuleName, fmt.Sprintf("Failed to get key from syscfg: key=%s err=%s\n", "eth_wan_enabled", err))
		os.Exit(1)
	}

	if isEthWANEnabled == "true" {
		// TODO: Implement Ethernet WAN support
		rdklogger.RDKLog(rdklogger.RDK_LOG_INFO, loggerModuleName, "Box is in ETHWAN mode\n")
	} else {
		rdklogger.RDKLog(rdklogger.RDK_LOG_INFO, loggerModuleName, "Box is in DOCSIS mode\n")
	}

	rdklogger.RDKLog(rdklogger.RDK_LOG_INFO, loggerModuleName, "Initializing MTA Agent data model\n")

	// Initialize the MTA HAL database
	ret := C.mta_hal_InitDB() // CosaDmlMTAInit
	if ret != C.RETURN_OK {
		return fmt.Errorf("mta_hal_InitDB failed with code %d", ret)
	}

	if m.cfg.EthernetWANEnabled {
		err := m.sysevent.Open("127.0.0.1", sysevent.SeServerWellKnownPort, sysevent.SeVersion, "WAN State")
		if err != nil {
			rdklogger.RDKLog(rdklogger.RDK_LOG_ERROR, loggerModuleName, fmt.Sprintf("Failed to open sysevent connection: err=%s\n", err))
			return fmt.Errorf("failed to open sysevent connection: err=%s", err)
		}

		if m.cfg.ErouterDHCPOptionMTA {
			// Mta_Sysevent_thread_Dhcp_Option
			err := m.sysevent.SetOptions("current_wan_state", sysevent.TupleFlagEvent)
			if err != nil {
				rdklogger.RDKLog(rdklogger.RDK_LOG_ERROR, loggerModuleName, fmt.Sprintf("Failed to set sysevent options: err=%s\n", err))
				os.Exit(1)
			}

			// Goroutine for monitoring the WAN mode and WAN status
			go func() {
				var wanStateAsyncID sysevent.AsyncID
				err := m.sysevent.SetNotification("current_wan_state", &wanStateAsyncID)
				if err != nil {
					rdklogger.RDKLog(rdklogger.RDK_LOG_ERROR, loggerModuleName, fmt.Sprintf("Failed to set sysevent options: err=%s\n", err))
					os.Exit(1)
				}
			}()
		} else {
			// Mta_Sysevent_thread
		}
	}

	// TODO: Conditional logic for defined(VOICE_MTA_SUPPORT)

	// Initialize the line table
	lineTableNumberOfEntries := C.mta_hal_LineTableGetNumberOfEntries()
	for i := uint32(0); i < uint32(lineTableNumberOfEntries); i++ {
		m.DataModel.LineTable = append(m.DataModel.LineTable, LineTableEntry{
			InstanceNumber: uint32(i + 1),
		})
	}

	mtaDMLConfigBytes, err := os.ReadFile(m.cfg.DMLConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read MTA DML config file: err=%s", err)
	}

	var mtaDMLConfig map[string]any
	err = json.Unmarshal(mtaDMLConfigBytes, &mtaDMLConfig)
	if err != nil {
		return fmt.Errorf("failed to parse MTA DML config file: err=%s", err)
	}

	var processMap func(m map[string]any, key string)
	processMap = func(m map[string]any, key string) {
		val, ok := m[key]
		if !ok {
			return
		}

		switch v := val.(type) {
		case map[string]any:
			for k := range v {
				processMap(v, k)
			}
		case []any:
			// TODO: Register all parameters with RBUS
			for _, item := range v {
				if itemMap, ok := item.(map[string]any); ok {
					for k := range itemMap {
						fmt.Printf("Registering parameter with RBUS: %s\n", k)
					}
				}
			}
		default:
			fmt.Printf("Value for key %s is not a map\n", key)
		}
	}

	processMap(mtaDMLConfig, "Device")

	rdklogger.RDKLog(rdklogger.RDK_LOG_INFO, loggerModuleName, fmt.Sprintf("Initializing RBUS for component: component=%s\n", m.cfg.ComponentName))

	err = rbus.Open("CcspMtaAgentSsp")
	if err != nil {
		rdklogger.RDKLog(rdklogger.RDK_LOG_ERROR, loggerModuleName, fmt.Sprintf("Failed to open rbus connection: err=%s\n", err))
		os.Exit(1)
	}

	fmt.Printf("Running MTA Agent: user_id=%d group_id=%d \n", syscall.Getuid(), syscall.Getgid())
	for {
		time.Sleep(10 * time.Second)
	}
}
