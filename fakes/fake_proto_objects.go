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
