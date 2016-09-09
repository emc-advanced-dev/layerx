package framework_manager
import (
	"github.com/gogo/protobuf/proto"
	"github.com/mesos/mesos-go/mesosproto"
)

const very_high_float = 100000000.0

var scalarType = mesosproto.Value_SCALAR
var rangesType = mesosproto.Value_RANGES

func newPhonyOffer(frameworkId string, offerId string, slaveId string) *mesosproto.Offer {
	resources := []*mesosproto.Resource{
		&mesosproto.Resource{
			Name: proto.String("mem"),
			Type: &scalarType,
			Scalar: &mesosproto.Value_Scalar{
				Value: proto.Float64(very_high_float),
			},
		},
		&mesosproto.Resource{
			Name: proto.String("disk"),
			Type: &scalarType,
			Scalar: &mesosproto.Value_Scalar{
				Value: proto.Float64(very_high_float),
			},
		},
		&mesosproto.Resource{
			Name: proto.String("ports"),
			Type: &rangesType,
			Ranges: &mesosproto.Value_Ranges{
				Range: []*mesosproto.Value_Range{
					&mesosproto.Value_Range{
						Begin: proto.Uint64(31000),
						End:   proto.Uint64(32000),
					},
				},
			},
		},
		&mesosproto.Resource{
			Name: proto.String("cpus"),
			Type: &scalarType,
			Scalar: &mesosproto.Value_Scalar{
				Value: proto.Float64(very_high_float),
			},
		},
	}

	return &mesosproto.Offer{
		Id: &mesosproto.OfferID{
			Value: proto.String("PHONY-OFFER-ID"),
		},
		FrameworkId: &mesosproto.FrameworkID{
			Value: proto.String(frameworkId),
		},
		SlaveId: &mesosproto.SlaveID{
			Value: proto.String(slaveId),
		},
		Hostname:  proto.String("PHONY-SLAVE-HOSTNAME"),
		Resources: resources,
	}
}
