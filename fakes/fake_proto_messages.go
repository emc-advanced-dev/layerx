package fakes

import (
	"github.com/mesos/mesos-go/mesosproto"
	"github.com/gogo/protobuf/proto"
	"github.com/mesos/mesos-go/mesosproto/scheduler"
)


func FakeSubscribeCall() *scheduler.Call {
	callType := scheduler.Call_SUBSCRIBE
	return &scheduler.Call{
		FrameworkId: &mesosproto.FrameworkID{
			Value: proto.String("fake_framework_id"),
		},
		Type: &callType,
		Subscribe: &scheduler.Call_Subscribe{
			FrameworkInfo: FakeFramework(),
		},
	}
}