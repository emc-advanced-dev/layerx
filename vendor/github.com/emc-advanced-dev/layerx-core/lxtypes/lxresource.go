package lxtypes
import "github.com/mesos/mesos-go/mesosproto"


type PortRange struct {
	Begin uint64 `json:"begin"`
	End   uint64 `json:"end"`
}

type Resource struct {
	Id    string            `json:"id"`
	NodeId string			`json:"node_id"`
	Cpus  float64           `json:"cpus"`
	Mem   float64           `json:"mem"`
	Disk  float64           `json:"disk"`
	Ports []PortRange       `json:"ports"`
	RpiName string			`json:"rpi_name"`
}

func NewResourceFromMesos(offer *mesosproto.Offer) *Resource {
	slaveId := offer.GetSlaveId().GetValue()
	offerId := offer.GetId().GetValue()
	var id = offerId
	var nodeId = slaveId
	var cpus float64
	var mem float64
	var disk float64
	var ports []PortRange

	cpus += getResourceScalar(offer.GetResources(), "cpus")
	mem += getResourceScalar(offer.GetResources(), "mem")
	disk += getResourceScalar(offer.GetResources(), "disk")
	for _, resource := range offer.GetResources() {
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

	return &Resource{
		Id:           id,
		NodeId:       nodeId,
		Cpus:         cpus,
		Mem:          mem,
		Disk:         disk,
		Ports:        ports,
	}
}
