package layerx_brain_client_test

import (
	. "github.com/layer-x/layerx-core_v2/layerx_brain_client"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/layer-x/layerx-core_v2/fakes"
	"github.com/mesos/mesos-go/mesosproto"
	"github.com/layer-x/layerx-core_v2/lxtypes"
	"github.com/layer-x/layerx-core_v2/layerx_rpi_client"
	"github.com/layer-x/layerx-core_v2/layerx_tpi_client"
	"github.com/layer-x/layerx-commons/lxhttpclient"
	"github.com/layer-x/layerx-commons/lxerrors"
	"fmt"
)

func PurgeFakeServer(fakeLxUrl string) error {
	resp, _, err := lxhttpclient.Post(fakeLxUrl, "/Purge", nil, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return lxerrors.New(fmt.Sprintf("status code was %v",resp.StatusCode), nil)
	}
	return nil
}

var _ = Describe("LayerxBrainClient", func() {

	fakeStatus1 := fakes.FakeTaskStatus("fake_task_id_1", mesosproto.TaskState_TASK_RUNNING)

	fakeStatuses := []*mesosproto.TaskStatus{fakeStatus1}

	go fakes.RunFakeLayerXServer(fakeStatuses, 12349)
	brainClient := LayerXBrainClient{
		CoreURL: "127.0.0.1:12349",
	}
	lxTpi := layerx_tpi_client.LayerXTpi{
		CoreURL: "127.0.0.1:12349",
	}
	lxRpi := layerx_rpi_client.LayerXRpi{
		CoreURL: "127.0.0.1:12349",
	}

	Describe("GetNodes", func() {
		It("returns the list of known nodes", func() {
			purgeErr := PurgeFakeServer("127.0.0.1:12349")
			Expect(purgeErr).To(BeNil())
			fakeOffer1 := fakes.FakeOffer("fake_offer_id_1", "_1")
			fakeOffer2 := fakes.FakeOffer("fake_offer_id_2", "_1")
			fakeOffer3 := fakes.FakeOffer("fake_offer_id_3", "_2")
			fakeResource1 := lxtypes.NewResourceFromMesos(fakeOffer1)
			fakeResource2 := lxtypes.NewResourceFromMesos(fakeOffer2)
			fakeResource3 := lxtypes.NewResourceFromMesos(fakeOffer3)
			err := lxRpi.SubmitResource(fakeResource1)
			Expect(err).To(BeNil())
			err = lxRpi.SubmitResource(fakeResource2)
			Expect(err).To(BeNil())
			err = lxRpi.SubmitResource(fakeResource3)
			Expect(err).To(BeNil())
			fakeNode1 := lxtypes.NewNode("_1")
			err = fakeNode1.AddResource(fakeResource1)
			Expect(err).To(BeNil())
			err = fakeNode1.AddResource(fakeResource2)
			Expect(err).To(BeNil())
			fakeNode2 := lxtypes.NewNode("_2")
			err = fakeNode2.AddResource(fakeResource3)
			Expect(err).To(BeNil())
			//the actual test
			nodes, err := brainClient.GetNodes()
			Expect(err).To(BeNil())
			Expect(nodes).To(ContainElement(fakeNode1))
			Expect(nodes).To(ContainElement(fakeNode2))
		})
	})

	Describe("GetNodes", func() {
		It("returns the list of known nodes", func() {
			purgeErr := PurgeFakeServer("127.0.0.1:12349")
			Expect(purgeErr).To(BeNil())
			fakeOffer1 := fakes.FakeOffer("fake_offer_id_1", "_1")
			fakeOffer2 := fakes.FakeOffer("fake_offer_id_2", "_1")
			fakeOffer3 := fakes.FakeOffer("fake_offer_id_3", "_2")
			fakeResource1 := lxtypes.NewResourceFromMesos(fakeOffer1)
			fakeResource2 := lxtypes.NewResourceFromMesos(fakeOffer2)
			fakeResource3 := lxtypes.NewResourceFromMesos(fakeOffer3)
			err := lxRpi.SubmitResource(fakeResource1)
			Expect(err).To(BeNil())
			err = lxRpi.SubmitResource(fakeResource2)
			Expect(err).To(BeNil())
			err = lxRpi.SubmitResource(fakeResource3)
			Expect(err).To(BeNil())
			fakeNode1 := lxtypes.NewNode("_1")
			err = fakeNode1.AddResource(fakeResource1)
			Expect(err).To(BeNil())
			err = fakeNode1.AddResource(fakeResource2)
			Expect(err).To(BeNil())
			fakeNode2 := lxtypes.NewNode("_2")
			err = fakeNode2.AddResource(fakeResource3)
			Expect(err).To(BeNil())
			//the actual test
			nodes, err := brainClient.GetNodes()
			Expect(err).To(BeNil())
			Expect(nodes).To(ContainElement(fakeNode1))
			Expect(nodes).To(ContainElement(fakeNode2))
		})
	})

	Describe("GetPendingTasks", func() {
		It("returns the list of taks in the pending pool", func() {
			purgeErr := PurgeFakeServer("127.0.0.1:12349")
			Expect(purgeErr).To(BeNil())
			fakeLxTask := fakes.FakeLXTask("fake_task_id", "fake_task_name", "", "echo FAKE_COMMAND")
			taskProvider := &lxtypes.TaskProvider{
				Id:     "fake_task_provider_id",
				Source: "taskprovider@tphost:port",
			}
			err := lxTpi.RegisterTaskProvider(taskProvider)
			Expect(err).To(BeNil())
			err = lxTpi.SubmitTask("fake_task_provider_id", fakeLxTask)
			Expect(err).To(BeNil())

			fakeLxTask.TaskProvider = taskProvider

			tasks, err := brainClient.GetPendingTasks()
			Expect(err).To(BeNil())
			Expect(tasks).To(ContainElement(fakeLxTask))
		})
	})

	Describe("AssignTasks", func() {
		It("assigns the NodeId as the SlaveId on the specified tasks", func() {
			purgeErr := PurgeFakeServer("127.0.0.1:12349")
			Expect(purgeErr).To(BeNil())
			fakeTask1 := fakes.FakeLXTask("fake_task_id_1", "fake_task_name", "", "echo FAKE_COMMAND")
			fakeTask2 := fakes.FakeLXTask("fake_task_id_2", "fake_task_name", "", "echo FAKE_COMMAND")
			fakeTask3 := fakes.FakeLXTask("fake_task_id_3", "fake_task_name", "", "echo FAKE_COMMAND")
			taskProvider := &lxtypes.TaskProvider{
				Id:     "fake_task_provider_id",
				Source: "taskprovider@tphost:port",
			}
			err := lxTpi.RegisterTaskProvider(taskProvider)
			Expect(err).To(BeNil())
			err = lxTpi.SubmitTask("fake_task_provider_id", fakeTask1)
			Expect(err).To(BeNil())
			err = lxTpi.SubmitTask("fake_task_provider_id", fakeTask2)
			Expect(err).To(BeNil())
			err = lxTpi.SubmitTask("fake_task_provider_id", fakeTask3)
			Expect(err).To(BeNil())


			fakeTask1.TaskProvider = taskProvider
			fakeTask2.TaskProvider = taskProvider
			fakeTask3.TaskProvider = taskProvider

			fakeOffer1 := fakes.FakeOffer("fake_offer_id_1", "_1")
			fakeOffer2 := fakes.FakeOffer("fake_offer_id_2", "_1")
			fakeResource1 := lxtypes.NewResourceFromMesos(fakeOffer1)
			fakeResource2 := lxtypes.NewResourceFromMesos(fakeOffer2)
			err = lxRpi.SubmitResource(fakeResource1)
			Expect(err).To(BeNil())
			err = lxRpi.SubmitResource(fakeResource2)
			Expect(err).To(BeNil())
			nodes, err := lxRpi.GetNodes()
			Expect(err).To(BeNil())
			fakeNode := nodes[0]

			err = brainClient.AssignTasks(fakeNode.Id, fakeTask1.TaskId, fakeTask2.TaskId, fakeTask3.TaskId)
			Expect(err).To(BeNil())
			tasks, err := brainClient.GetPendingTasks()
			Expect(err).To(BeNil())
			Expect(tasks).To(BeEmpty())
		})
	})

	Describe("AssignTasks", func() {
		It("assigns the NodeId as the SlaveId on the specified tasks", func() {
			purgeErr := PurgeFakeServer("127.0.0.1:12349")
			Expect(purgeErr).To(BeNil())
			fakeTask1 := fakes.FakeLXTask("fake_task_id_1", "fake_task_name", "", "echo FAKE_COMMAND")
			fakeTask2 := fakes.FakeLXTask("fake_task_id_2", "fake_task_name", "", "echo FAKE_COMMAND")
			fakeTask3 := fakes.FakeLXTask("fake_task_id_3", "fake_task_name", "", "echo FAKE_COMMAND")
			taskProvider := &lxtypes.TaskProvider{
				Id:     "fake_task_provider_id",
				Source: "taskprovider@tphost:port",
			}
			err := lxTpi.RegisterTaskProvider(taskProvider)
			Expect(err).To(BeNil())
			err = lxTpi.SubmitTask("fake_task_provider_id", fakeTask1)
			Expect(err).To(BeNil())
			err = lxTpi.SubmitTask("fake_task_provider_id", fakeTask2)
			Expect(err).To(BeNil())
			err = lxTpi.SubmitTask("fake_task_provider_id", fakeTask3)
			Expect(err).To(BeNil())

			fakeTask1.TaskProvider = taskProvider
			fakeTask2.TaskProvider = taskProvider
			fakeTask3.TaskProvider = taskProvider

			fakeOffer1 := fakes.FakeOffer("fake_offer_id_1", "_1")
			fakeOffer2 := fakes.FakeOffer("fake_offer_id_2", "_1")
			fakeResource1 := lxtypes.NewResourceFromMesos(fakeOffer1)
			fakeResource2 := lxtypes.NewResourceFromMesos(fakeOffer2)
			err = lxRpi.SubmitResource(fakeResource1)
			Expect(err).To(BeNil())
			err = lxRpi.SubmitResource(fakeResource2)
			Expect(err).To(BeNil())
			Expect(err).To(BeNil())
			fakeOffer3 := fakes.FakeOffer("fake_offer_id_3", "_2")
			fakeOffer4 := fakes.FakeOffer("fake_offer_id_4", "_2")
			fakeResource3 := lxtypes.NewResourceFromMesos(fakeOffer3)
			fakeResource4 := lxtypes.NewResourceFromMesos(fakeOffer4)
			err = lxRpi.SubmitResource(fakeResource3)
			Expect(err).To(BeNil())
			err = lxRpi.SubmitResource(fakeResource4)
			Expect(err).To(BeNil())
			nodes, err := lxRpi.GetNodes()
			Expect(err).To(BeNil())
			Expect(len(nodes)).To(Equal(2))
			fakeNode1 := nodes[0]
			fakeNode2 := nodes[1]

			err = brainClient.AssignTasks(fakeNode1.Id, fakeTask1.TaskId, fakeTask2.TaskId, fakeTask3.TaskId)
			Expect(err).To(BeNil())
			tasks, err := brainClient.GetPendingTasks()
			Expect(err).To(BeNil())
			Expect(tasks).To(BeEmpty())
			err = brainClient.MigrateTasks(fakeNode2.Id, fakeTask1.TaskId, fakeTask2.TaskId, fakeTask3.TaskId)
			Expect(err).To(BeNil())
		})
	})
})
