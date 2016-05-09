package framework_api_handlers_test

import (
	. "github.com/emc-advanced-dev/layerx-mesos-rpi/mesos_framework_api/framework_api_handlers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/mesos/mesos-go/mesosproto"
core_fakes "github.com/emc-advanced-dev/layerx-core/fakes"
	"github.com/emc-advanced-dev/layerx-core/layerx_rpi_client"
)

var _ = Describe("HandleResourceOffers", func() {
	fakeStatus1 := core_fakes.FakeTaskStatus("fake_task_id_1", mesosproto.TaskState_TASK_RUNNING)

	fakeStatuses := []*mesosproto.TaskStatus{fakeStatus1}

	go core_fakes.RunFakeLayerXServer(fakeStatuses, 12346)
	lxRpi := &layerx_rpi_client.LayerXRpi{
		CoreURL: "127.0.0.1:12346",
	}

	Describe("SubmitResource", func() {
		It("submits a resource to the LX Server", func() {
			fakeOffer1 := core_fakes.FakeOffer("fake_offer_id1", "fake_slave_id1")
			fakeOffer2 := core_fakes.FakeOffer("fake_offer_id2", "fake_slave_id1")
			fakeOffer3 := core_fakes.FakeOffer("fake_offer_id3", "fake_slave_id2")
			fakeOffers := []*mesosproto.Offer{fakeOffer1, fakeOffer2, fakeOffer3}
			err := HandleResourceOffers(lxRpi, fakeOffers)
			Expect(err).To(BeNil())
		})
	})
})
