package fakes
import (
"github.com/mesos/mesos-go/mesosproto"
"github.com/gogo/protobuf/proto"
)

func FakeTaskStatus(taskId string, fakeState mesosproto.TaskState) *mesosproto.TaskStatus {
	return &mesosproto.TaskStatus{
		TaskId: &mesosproto.TaskID{
			Value: proto.String(taskId),
		},
		State:   &fakeState,
		Message: proto.String("fake_message"),
	}
}