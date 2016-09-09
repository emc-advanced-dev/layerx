package layerx_rpi_client_test

import (
	. "github.com/emc-advanced-dev/layerx-core/layerx_rpi_client"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
"github.com/emc-advanced-dev/layerx-core/fakes"
	"github.com/mesos/mesos-go/mesosproto"
	"github.com/emc-advanced-dev/layerx-core/lxtypes"
)

var _ = Describe("LayerxRpi", func() {
	fakeStatus1 := fakes.FakeTaskStatus("fake_task_id_1", mesosproto.TaskState_TASK_RUNNING)

	fakeStatuses := []*mesosproto.TaskStatus{fakeStatus1}

	go fakes.RunFakeLayerXServer(fakeStatuses, 12346)
	lxRpi := LayerXRpi{
		CoreURL: "127.0.0.1:12346",
	}

	Describe("RegisterRpi", func() {
		It("registers the Rpi URL to the LX Server", func() {
			err := lxRpi.RegisterRpi("fake-mesos-rpi", "fake.rpi.ip:1234")
			Expect(err).To(BeNil())
		})
	})

	Describe("SubmitResource", func() {
		It("submits a resource to the LX Server", func() {
			fakeOffer := fakes.FakeOffer("fake_offer_id", "fake_slave_id")
			fakeResource := lxtypes.NewResourceFromMesos(fakeOffer)
			err := lxRpi.SubmitResource(fakeResource)
			Expect(err).To(BeNil())
		})
	})

	Describe("SubmitStatus", func() {
		It("submits a status to the LX Server", func() {
			fakeStatus2 := fakes.FakeTaskStatus("fake_task_id_2", mesosproto.TaskState_TASK_KILLED)
			fakeStatus3 := fakes.FakeTaskStatus("fake_task_id_3", mesosproto.TaskState_TASK_FINISHED)
			err := lxRpi.SubmitStatusUpdate(fakeStatus2)
			Expect(err).To(BeNil())
			err = lxRpi.SubmitStatusUpdate(fakeStatus3)
			Expect(err).To(BeNil())
		})
	})

	Describe("GetNodes", func(){
		It("returns the list of known nodes", func(){
			fakeOffer1 := fakes.FakeOffer("fake_offer_id_1", "fake_slave_id_1")
			fakeOffer2 := fakes.FakeOffer("fake_offer_id_2", "fake_slave_id_1")
			fakeOffer3 := fakes.FakeOffer("fake_offer_id_3", "fake_slave_id_2")
			fakeResource1 := lxtypes.NewResourceFromMesos(fakeOffer1)
			fakeResource2 := lxtypes.NewResourceFromMesos(fakeOffer2)
			fakeResource3 := lxtypes.NewResourceFromMesos(fakeOffer3)
			err := lxRpi.SubmitResource(fakeResource1)
			Expect(err).To(BeNil())
			err = lxRpi.SubmitResource(fakeResource2)
			Expect(err).To(BeNil())
			err = lxRpi.SubmitResource(fakeResource3)
			Expect(err).To(BeNil())
			fakeNode1 := lxtypes.NewNode("fake_slave_id_1")
			err = fakeNode1.AddResource(fakeResource1)
			Expect(err).To(BeNil())
			err = fakeNode1.AddResource(fakeResource2)
			Expect(err).To(BeNil())
			fakeNode2 := lxtypes.NewNode("fake_slave_id_2")
			err = fakeNode2.AddResource(fakeResource3)
			Expect(err).To(BeNil())
			//the actual test
			nodes, err := lxRpi.GetNodes()
			Expect(err).To(BeNil())
			Expect(nodes).To(ContainElement(fakeNode1))
			Expect(nodes).To(ContainElement(fakeNode2))
		})
	})
})
