package lxtypes_test

import (
	. "github.com/emc-advanced-dev/layerx-core/lxtypes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/emc-advanced-dev/layerx-core/fakes"
	"github.com/mesos/mesos-go/mesosproto"
	"github.com/mesos/mesos-go/mesosutil"
)

var _ = Describe("Lxresource", func() {
	Describe("NewResourceFromMesos", func(){
		It("converts a mesos offer to a layerx resource", func(){
			fakeOffer := fakes.FakeOffer("fake_offer_id", "fake_slave_id")
			resource := NewResourceFromMesos(fakeOffer)
			Expect(resource.Cpus).To(Equal(getResourceScalar(fakeOffer.GetResources(), "cpus")))
			Expect(resource.Mem).To(Equal(getResourceScalar(fakeOffer.GetResources(), "mem")))
			Expect(resource.Disk).To(Equal(getResourceScalar(fakeOffer.GetResources(), "disk")))
			var ports []PortRange

			for _, resource := range fakeOffer.GetResources() {
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
			for _, portRange := range resource.Ports {
				Expect(ports).To(ContainElement(portRange))
			}
		})
	})
})

func getResourceScalar(resources []*mesosproto.Resource, name string) float64 {
	resources = mesosutil.FilterResources(resources, func(res *mesosproto.Resource) bool {
		return res.GetName() == name
	})

	value := 0.0
	for _, res := range resources {
		value += res.GetScalar().GetValue()
	}

	return value
}
