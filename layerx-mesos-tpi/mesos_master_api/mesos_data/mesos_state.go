package mesos_data

import (
	"github.com/gogo/protobuf/proto"
	"github.com/emc-advanced-dev/pkg/errors"
	"github.com/layer-x/layerx-commons/lxutils"
	"github.com/mesos/mesos-go/mesosproto"
	"net"
	"strconv"
)

const (
	MESOS_VERSION = "0.22.1"
)

type MesosState struct {
	ActivatedSlaves     int    `json:"activated_slaves"`
	BuildDate           string `json:"build_date"`
	BuildTime           int    `json:"build_time"`
	BuildUser           string `json:"build_user"`
	CompletedFrameworks []struct {
		Active           bool          `json:"active"`
		Checkpoint       bool          `json:"checkpoint"`
		CompletedTasks   []interface{} `json:"completed_tasks"`
		FailoverTimeout  int64         `json:"failover_timeout"`
		Hostname         string        `json:"hostname"`
		ID               string        `json:"id"`
		Name             string        `json:"name"`
		OfferedResources struct {
			Cpus int `json:"cpus"`
			Disk int `json:"disk"`
			Mem  int `json:"mem"`
		} `json:"offered_resources"`
		Offers         []interface{} `json:"offers"`
		RegisteredTime float64       `json:"registered_time"`
		Resources      struct {
			Cpus int `json:"cpus"`
			Disk int `json:"disk"`
			Mem  int `json:"mem"`
		} `json:"resources"`
		Role             string        `json:"role"`
		Tasks            []interface{} `json:"tasks"`
		UnregisteredTime float64       `json:"unregistered_time"`
		UsedResources    struct {
			Cpus int `json:"cpus"`
			Disk int `json:"disk"`
			Mem  int `json:"mem"`
		} `json:"used_resources"`
		User     string `json:"user"`
		WebuiURL string `json:"webui_url"`
	} `json:"completed_frameworks"`
	DeactivatedSlaves int     `json:"deactivated_slaves"`
	ElectedTime       float64 `json:"elected_time"`
	FailedTasks       int     `json:"failed_tasks"`
	FinishedTasks     int     `json:"finished_tasks"`
	Flags             struct {
		AllocationInterval        string `json:"allocation_interval"`
		Authenticate              string `json:"authenticate"`
		AuthenticateSlaves        string `json:"authenticate_slaves"`
		Authenticators            string `json:"authenticators"`
		FrameworkSorter           string `json:"framework_sorter"`
		Help                      string `json:"help"`
		InitializeDriverLogging   string `json:"initialize_driver_logging"`
		LogAutoInitialize         string `json:"log_auto_initialize"`
		LogDir                    string `json:"log_dir"`
		Logbufsecs                string `json:"logbufsecs"`
		LoggingLevel              string `json:"logging_level"`
		Port                      string `json:"port"`
		Quiet                     string `json:"quiet"`
		Quorum                    string `json:"quorum"`
		RecoverySlaveRemovalLimit string `json:"recovery_slave_removal_limit"`
		Registry                  string `json:"registry"`
		RegistryFetchTimeout      string `json:"registry_fetch_timeout"`
		RegistryStoreTimeout      string `json:"registry_store_timeout"`
		RegistryStrict            string `json:"registry_strict"`
		RootSubmissions           string `json:"root_submissions"`
		SlaveReregisterTimeout    string `json:"slave_reregister_timeout"`
		UserSorter                string `json:"user_sorter"`
		Version                   string `json:"version"`
		WebuiDir                  string `json:"webui_dir"`
		WorkDir                   string `json:"work_dir"`
		Zk                        string `json:"zk"`
		ZkSessionTimeout          string `json:"zk_session_timeout"`
	} `json:"flags"`
	Frameworks []struct {
		Active           bool          `json:"active"`
		Checkpoint       bool          `json:"checkpoint"`
		CompletedTasks   []interface{} `json:"completed_tasks"`
		FailoverTimeout  int           `json:"failover_timeout"`
		Hostname         string        `json:"hostname"`
		ID               string        `json:"id"`
		Name             string        `json:"name"`
		OfferedResources struct {
			Cpus int `json:"cpus"`
			Disk int `json:"disk"`
			Mem  int `json:"mem"`
		} `json:"offered_resources"`
		Offers           []interface{} `json:"offers"`
		RegisteredTime   float64       `json:"registered_time"`
		ReregisteredTime float64       `json:"reregistered_time"`
		Resources        struct {
			Cpus int `json:"cpus"`
			Disk int `json:"disk"`
			Mem  int `json:"mem"`
		} `json:"resources"`
		Role             string        `json:"role"`
		Tasks            []interface{} `json:"tasks"`
		UnregisteredTime int           `json:"unregistered_time"`
		UsedResources    struct {
			Cpus int `json:"cpus"`
			Disk int `json:"disk"`
			Mem  int `json:"mem"`
		} `json:"used_resources"`
		User     string `json:"user"`
		WebuiURL string `json:"webui_url"`
	} `json:"frameworks"`
	GitSha      string        `json:"git_sha"`
	GitTag      string        `json:"git_tag"`
	Hostname    string        `json:"hostname"`
	ID          string        `json:"id"`
	KilledTasks int           `json:"killed_tasks"`
	Leader      string        `json:"leader"`
	LogDir      string        `json:"log_dir"`
	LostTasks   int           `json:"lost_tasks"`
	OrphanTasks []interface{} `json:"orphan_tasks"`
	Pid         string        `json:"pid"`
	Slaves      []struct {
		Active     bool `json:"active"`
		Attributes struct {
		} `json:"attributes"`
		Hostname       string  `json:"hostname"`
		ID             string  `json:"id"`
		Pid            string  `json:"pid"`
		RegisteredTime float64 `json:"registered_time"`
		Resources      struct {
			Cpus  int    `json:"cpus"`
			Disk  int    `json:"disk"`
			Mem   int    `json:"mem"`
			Ports string `json:"ports"`
		} `json:"resources"`
	} `json:"slaves"`
	StagedTasks            int           `json:"staged_tasks"`
	StartTime              float64       `json:"start_time"`
	StartedTasks           int           `json:"started_tasks"`
	UnregisteredFrameworks []interface{} `json:"unregistered_frameworks"`
	Version                string        `json:"version"`
}

func (state *MesosState) ToMasterInfo() (*mesosproto.MasterInfo, error) {
	upid, err := UPIDFromString(state.Leader)
	if err != nil {
		return nil, errors.New("parsing upid from state.Leader", err)
	}
	port, err := strconv.Atoi(upid.Port)
	if err != nil {
		return nil, errors.New("converting state.Port to int", err)
	}
	ipInt := lxutils.IpToI(net.ParseIP(upid.Host))
	return &mesosproto.MasterInfo{
		Id:   proto.String(state.ID),
		Ip:   proto.Uint32(ipInt),
		Port: proto.Uint32(uint32(port)),
		Pid:  proto.String(state.Leader),
		Address: &mesosproto.Address{
			Hostname: proto.String(upid.Host),
			Ip:       proto.String(upid.Host),
			Port:     proto.Int(port),
		},
	}, nil
}
