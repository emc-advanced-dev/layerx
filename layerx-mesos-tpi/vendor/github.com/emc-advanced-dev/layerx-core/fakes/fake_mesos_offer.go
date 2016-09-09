package fakes

import (
	"github.com/gogo/protobuf/proto"
	"github.com/mesos/mesos-go/mesosproto"
)

func FakeOffer(offerId string, slaveId string) *mesosproto.Offer {
	return &mesosproto.Offer{
		Id: &mesosproto.OfferID{
			Value: proto.String(offerId),
		},
		FrameworkId: &mesosproto.FrameworkID{
			Value: proto.String("fake_framework_id"),
		},
		SlaveId: &mesosproto.SlaveID{
			Value: proto.String(slaveId),
		},
		Hostname:  proto.String("fake_slave_hostname"),
		Resources: FakeResources(),
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
