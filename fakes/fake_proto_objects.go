package fakes
import (
	"github.com/gogo/protobuf/proto"
	"github.com/mesos/mesos-go/mesosproto"
)


func FakeFramework() *mesosproto.FrameworkInfo {
	capabilityType := mesosproto.FrameworkInfo_Capability_REVOCABLE_RESOURCES
	return &mesosproto.FrameworkInfo{
		User: proto.String("fake_framework_user"),
		Name: proto.String("fake_framework"),
		Id: &mesosproto.FrameworkID{
			Value: proto.String("fake_framework_id"),
		},
		FailoverTimeout: proto.Float64(0),
		Checkpoint:      proto.Bool(false),
		Role:            proto.String("*"),
		Hostname:        proto.String("fake_host"),
		Principal:       proto.String("fake_principal"),
		WebuiUrl:        proto.String("http://fakeip:fakeport"),
		Capabilities: []*mesosproto.FrameworkInfo_Capability{
			&mesosproto.FrameworkInfo_Capability{
				Type: &capabilityType,
			},
		},
		Labels: &mesosproto.Labels{
			Labels: []*mesosproto.Label{
				&mesosproto.Label{
					Key:   proto.String("FakeLabel"),
					Value: proto.String("FakeValue"),
				},
			},
		},
	}
}

func FakeRegisterFrameworkMessage() *mesosproto.RegisterFrameworkMessage {
	return &mesosproto.RegisterFrameworkMessage{
		Framework: FakeFramework(),
	}
}

func FakeReregisterFrameworkMessage() *mesosproto.ReregisterFrameworkMessage {
	return &mesosproto.ReregisterFrameworkMessage{
		Framework: FakeFramework(),
		Failover:  proto.Bool(true),
	}
}

func FakeUnregisterFrameworkMessage() *mesosproto.UnregisterFrameworkMessage {
	return &mesosproto.UnregisterFrameworkMessage{
		FrameworkId: FakeFramework().GetId(),
	}
}

func FakeLaunchTasksMessage(frameworkId string) *mesosproto.LaunchTasksMessage {
	fakeTask1 := FakeTask("fake_task_1")
	fakeTask2 := FakeTask("fake_task_2")
	fakeTask3 := FakeTask("fake_task_3")
	fakeTasks := []*mesosproto.TaskInfo{
		fakeTask1,
		fakeTask2,
		fakeTask3,
	}
	fakeOfferIds := []*mesosproto.OfferID{
		&mesosproto.OfferID{
			Value: proto.String("fake_offer_id"),
		},
	}
	return &mesosproto.LaunchTasksMessage{
		FrameworkId: &mesosproto.FrameworkID{
			Value: proto.String(frameworkId),
		},
		Tasks:    fakeTasks,
		Filters:  &mesosproto.Filters{},
		OfferIds: fakeOfferIds,
	}
}

func FakeTask(taskId string) *mesosproto.TaskInfo {
	return &mesosproto.TaskInfo{
		Name: proto.String("fake_task_name"),
		TaskId: &mesosproto.TaskID{
			Value: proto.String(taskId),
		},
		SlaveId: &mesosproto.SlaveID{
			Value: proto.String("fake_slave_id"),
		},
		Resources: FakeResources(),
		Command: &mesosproto.CommandInfo{
			Value:     proto.String("echo"),
			Arguments: []string{"fake_echo_message"},
			Shell:     proto.Bool(true),
		},
	}
}

func FakeResources() []*mesosproto.Resource {
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

func FakeTaskStatus(taskId string, fakeState mesosproto.TaskState) *mesosproto.TaskStatus {
	return &mesosproto.TaskStatus{
		TaskId: &mesosproto.TaskID{
			Value: proto.String(taskId),
		},
		State:   &fakeState,
		Message: proto.String("fake_message"),
	}
}


func FakeKillTaskMessage(frameworkId string, taskId string) *mesosproto.KillTaskMessage {
	return &mesosproto.KillTaskMessage{
		FrameworkId: &mesosproto.FrameworkID{
			Value: proto.String(frameworkId),
		},
		TaskId: &mesosproto.TaskID{
			Value: proto.String(taskId),
		},
	}
}