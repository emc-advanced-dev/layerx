package fakes
import (
	"github.com/gogo/protobuf/proto"
	"github.com/mesos/mesos-go/mesosproto"
	"github.com/emc-advanced-dev/layerx-core/lxtypes"
)

const fake_task_command = `i=0; while true; do echo $i; i=$(expr $i + 1); sleep 1; done`

func FakeTask(taskId string) *lxtypes.Task {
	return lxtypes.NewTaskFromMesos(fakeMesosTask(taskId))
}

func fakeMesosResources() []*mesosproto.Resource {
	var scalarType = mesosproto.Value_SCALAR
	return []*mesosproto.Resource{
		&mesosproto.Resource{
			Name: proto.String("mem"),
			Type: &scalarType,
			Scalar: &mesosproto.Value_Scalar{
				Value: proto.Float64(0.1),
			},
		},
		&mesosproto.Resource{
			Name: proto.String("disk"),
			Type: &scalarType,
			Scalar: &mesosproto.Value_Scalar{
				Value: proto.Float64(0.1),
			},
		},
		&mesosproto.Resource{
			Name: proto.String("cpus"),
			Type: &scalarType,
			Scalar: &mesosproto.Value_Scalar{
				Value: proto.Float64(0.1),
			},
		},
	}
}

func fakeMesosTask(taskId string) *mesosproto.TaskInfo {
	return &mesosproto.TaskInfo{
		Name: proto.String("fake_task_name"),
		TaskId: &mesosproto.TaskID{
			Value: proto.String(taskId),
		},
		SlaveId: &mesosproto.SlaveID{
			Value: proto.String("fake_slave_id"),
		},
		Resources: fakeMesosResources(),
		Command: &mesosproto.CommandInfo{
			Value: proto.String(fake_task_command),
			Shell: proto.Bool(true),
		},
	}
}
