package lxtypes

import (
	"encoding/json"
	"github.com/gogo/protobuf/proto"
	"github.com/mesos/mesos-go/mesosproto"
	"github.com/mesos/mesos-go/mesosutil"
)

const (
	task_provider_key  = "layerx_mesos_tpi_task_provider_key"
	kill_requested_key = "layerx_mesos_tpi_kill_requested_key"
	checkpointed_key       = "layerx_mesos_tpi_checkpointed_key"
)

type Task struct {
	TaskProvider  *TaskProvider            `json:"task_provider"`
	Checkpointed bool                     `json:"checkpointed"`
	KillRequested bool                     `json:"kill_requested"`
	Name          string                   `json:"name,omitempty"`
	TaskId        string                   `json:"task_id,omitempty"`
	SlaveId       string                   `json:"slave_id,omitempty"`
	Cpus          float64                  `json:cpus`
	Mem           float64                  `json:mem`
	Disk          float64                  `json:disk`
	Ports         []PortRange              `json:ports`
	Executor      *mesosproto.ExecutorInfo `json:"executor,omitempty"`
	Command       *mesosproto.CommandInfo  `json:"command,omitempty"`
	// Task provided with a container will launch the container as part
	// of this task paired with the task's CommandInfo.
	Container *mesosproto.ContainerInfo `json:"container,omitempty"`
	Data      []byte                    `json:"data,omitempty"`
	// A health check for the task (currently in *alpha* and initial
	// support will only be for TaskInfo's that have a CommandInfo).
	HealthCheck *mesosproto.HealthCheck `json:"health_check,omitempty"`
	// Labels are free-form key value pairs which are exposed through
	// master and slave endpoints. Labels will not be interpreted or
	// acted upon by Mesos itself. As opposed to the data field, labels
	// will be kept in memory on master and slave processes. Therefore,
	// labels should be used to tag tasks with light-weight meta-data.
	Labels map[string]string `json:"labels,omitempty"`
	// Service discovery information for the task. It is not interpreted
	// or acted upon by Mesos. It is up to a service discovery system
	// to use this information as needed and to handle tasks without
	// service discovery information.
	Discovery *mesosproto.DiscoveryInfo `json:"discovery,omitempty"`
}

func (t *Task) ToMesos() *mesosproto.TaskInfo {
	var scalarType = mesosproto.Value_SCALAR
	var rangesType = mesosproto.Value_RANGES
	mesosPortRanges := []*mesosproto.Value_Range{}
	for _, portRange := range t.Ports {
		mesosPort := &mesosproto.Value_Range{
			Begin: proto.Uint64(portRange.Begin),
			End:   proto.Uint64(portRange.End),
		}
		mesosPortRanges = append(mesosPortRanges, mesosPort)
	}
	mesosPorts := &mesosproto.Value_Ranges{
		Range: mesosPortRanges,
	}
	var label_arr []*mesosproto.Label
	for key, val := range t.Labels {
		label := &mesosproto.Label{
			Key:   proto.String(key),
			Value: proto.String(val),
		}
		label_arr = append(label_arr, label)
	}
	taskProviderJson, err := json.Marshal(t.TaskProvider)
	//store current_status as a label
	if err != nil {
		tasks_label := &mesosproto.Label{
			Key:   proto.String(task_provider_key),
			Value: proto.String(string(taskProviderJson)),
		}
		label_arr = append(label_arr, tasks_label)
	}
	//store checkpointed as a label
	if t.Checkpointed {
		tasks_label := &mesosproto.Label{
			Key:   proto.String(checkpointed_key),
			Value: proto.String("true"),
		}
		label_arr = append(label_arr, tasks_label)
	} else {
		tasks_label := &mesosproto.Label{
			Key:   proto.String(checkpointed_key),
			Value: proto.String("false"),
		}
		label_arr = append(label_arr, tasks_label)
	}
	//store kill_requested as a label
	if t.KillRequested {
		tasks_label := &mesosproto.Label{
			Key:   proto.String(kill_requested_key),
			Value: proto.String("true"),
		}
		label_arr = append(label_arr, tasks_label)
	} else {
		tasks_label := &mesosproto.Label{
			Key:   proto.String(kill_requested_key),
			Value: proto.String("false"),
		}
		label_arr = append(label_arr, tasks_label)
	}

	labels := &mesosproto.Labels{
		Labels: label_arr,
	}
	return &mesosproto.TaskInfo{
		Name: proto.String(t.Name),
		TaskId: &mesosproto.TaskID{
			Value: proto.String(t.TaskId),
		},
		SlaveId: &mesosproto.SlaveID{
			Value: proto.String(t.SlaveId),
		},
		Resources: []*mesosproto.Resource{
			&mesosproto.Resource{
				Name: proto.String("cpus"),
				Type: &scalarType,
				Scalar: &mesosproto.Value_Scalar{
					Value: proto.Float64(t.Cpus),
				},
			},
			&mesosproto.Resource{
				Name: proto.String("mem"),
				Type: &scalarType,
				Scalar: &mesosproto.Value_Scalar{
					Value: proto.Float64(t.Mem),
				},
			},
			&mesosproto.Resource{
				Name: proto.String("disk"),
				Type: &scalarType,
				Scalar: &mesosproto.Value_Scalar{
					Value: proto.Float64(t.Disk),
				},
			},
			&mesosproto.Resource{
				Name:   proto.String("ports"),
				Type:   &rangesType,
				Ranges: mesosPorts,
			},
		},
		Executor:    t.Executor,
		Command:     t.Command,
		Container:   t.Container,
		Data:        t.Data,
		HealthCheck: t.HealthCheck,
		Labels:      labels,
		Discovery:   t.Discovery,
	}
}

func NewTaskFromMesos(taskInfo *mesosproto.TaskInfo) *Task {
	ports := []PortRange{}
	for _, resource := range taskInfo.GetResources() {
		if resource.GetName() == "ports" {
			for _, mesosRange := range resource.GetRanges().GetRange() {
				port := PortRange{
					Begin: mesosRange.GetBegin(),
					End:   mesosRange.GetEnd(),
				}
				ports = append(ports, port)
			}
		}
	}
	var taskProvider TaskProvider
	var killRequested bool
	var checkpointed bool
	labels := make(map[string]string)
	for _, label := range taskInfo.GetLabels().GetLabels() {
		//if label is task_provider_key, populate task provider instead
		if label.GetKey() == task_provider_key {
			json.Unmarshal([]byte(label.GetValue()), &taskProvider)
		} else if label.GetKey() == kill_requested_key {
			if label.GetValue() == "true" {
				killRequested = true
			}
		} else if label.GetKey() == checkpointed_key {
			if label.GetValue() == "true" {
				checkpointed = true
			}
		} else {
			labels[label.GetKey()] = label.GetValue()
		}
	}

	task := &Task{
		KillRequested: killRequested,
		Checkpointed:      checkpointed,
		Name:          taskInfo.GetName(),
		TaskId:        taskInfo.GetTaskId().GetValue(),
		SlaveId:       taskInfo.GetSlaveId().GetValue(),
		Cpus:          getResourceScalar(taskInfo.GetResources(), "cpus"),
		Mem:           getResourceScalar(taskInfo.GetResources(), "mem"),
		Disk:          getResourceScalar(taskInfo.GetResources(), "disk"),
		Ports:         ports,
		Executor:      taskInfo.Executor,
		Command:       taskInfo.Command,
		Data:          taskInfo.Data,
		HealthCheck:   taskInfo.HealthCheck,
	}

	if taskProvider.Id != "" {
		task.TaskProvider = &taskProvider
	}

	return task
}

func getResourceScalar(resources []*mesosproto.Resource, name string) float64 {
	resources = mesosutil.FilterResources(resources, func(res *mesosproto.Resource) bool {
		return res.GetName() == name
	})

	value := 0.0
	for _, res := range resources {
		value += res.GetScalar().GetValue()
	}

	return value
}
