package layerx_rpi_client_test

import (
	. "github.com/layer-x/layerx-core_v2/layerx_rpi_client"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
"github.com/layer-x/layerx-core_v2/fakes"
	"github.com/mesos/mesos-go/mesosproto"
	"github.com/layer-x/layerx-core_v2/lxtypes"
)

var _ = Describe("LayerxRpi", func() {
	fakeStatus1 := fakes.FakeTaskStatus("fake_task_id_1", mesosproto.TaskState_TASK_RUNNING)

	fakeStatuses := []*mesosproto.TaskStatus{fakeStatus1}

	go fakes.RunFakeLayerXServer(fakeStatuses, 12346)
	lxRpi := LayerXRpi{
		CoreURL: "127.0.0.1:12346",
	}

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
})
