package fakes

import (
	"github.com/gogo/protobuf/proto"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"
	"github.com/mesos/mesos-go/mesosproto"
)

func FakeLXTask(taskId, taskName, slaveId, taskCommand string) *lxtypes.Task {
	mesosTask := FakeMesosTask(taskId, taskName, slaveId, taskCommand)
	return lxtypes.NewTaskFromMesos(mesosTask)
}

func FakeMesosTask(taskId, taskName, slaveId, taskCommand string) *mesosproto.TaskInfo {
	return &mesosproto.TaskInfo{
		Name: proto.String(taskName),
		TaskId: &mesosproto.TaskID{
			Value: proto.String(taskId),
		},
		SlaveId: &mesosproto.SlaveID{
			Value: proto.String(slaveId),
		},
		Resources: fakeMesosResources(),
		Command: &mesosproto.CommandInfo{
			Value: proto.String(taskCommand),
			Shell: proto.Bool(true),
		},
	}
}

func fakeMesosResources() []*mesosproto.Resource {
	var scalarType = mesosproto.Value_SCALAR
	var rangesType = mesosproto.Value_RANGES
	return []*mesosproto.Resource{
		&mesosproto.Resource{
			Name: proto.String("mem"),
			Type: &scalarType,
			Scalar: &mesosproto.Value_Scalar{
				Value: proto.Float64(1234),
			},
		},
		&mesosproto.Resource{
			Name: proto.String("disk"),
			Type: &scalarType,
			Scalar: &mesosproto.Value_Scalar{
				Value: proto.Float64(1234),
			},
		},
		&mesosproto.Resource{
			Name: proto.String("ports"),
			Type: &rangesType,
			Ranges: &mesosproto.Value_Ranges{
				Range: []*mesosproto.Value_Range{
					&mesosproto.Value_Range{
						Begin: proto.Uint64(1234),
						End:   proto.Uint64(12345),
					},
				},
			},
		},
		&mesosproto.Resource{
			Name: proto.String("cpus"),
			Type: &scalarType,
			Scalar: &mesosproto.Value_Scalar{
				Value: proto.Float64(0.1234),
			},
		},
	}
}
