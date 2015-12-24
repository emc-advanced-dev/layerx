package fakes
import (
"github.com/mesos/mesos-go/mesosproto"
"github.com/gogo/protobuf/proto"
)


func FakeSubscribeCall() *mesosproto.Call {
	callType := mesosproto.Call_SUBSCRIBE
	return &mesosproto.Call{
		FrameworkId: &mesosproto.FrameworkID{
			Value: proto.String("fake_framework_id"),
		},
		Type: &callType,
		Subscribe: &mesosproto.Call_Subscribe{
			FrameworkInfo: FakeFramework(),
		},
	}
}