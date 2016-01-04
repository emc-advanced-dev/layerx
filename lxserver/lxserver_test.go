package lxserver_test

import (
	. "github.com/layer-x/layerx-core_v2/lxserver"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/layer-x/layerx-core_v2/layerx_rpi_client"
	"github.com/layer-x/layerx-core_v2/layerx_tpi_client"
	"github.com/layer-x/layerx-commons/lxactionqueue"
	"github.com/layer-x/layerx-commons/lxmartini"
	"github.com/layer-x/layerx-core_v2/driver"
	"github.com/layer-x/layerx-core_v2/lxstate"
	"github.com/layer-x/layerx-commons/lxlog"
	"fmt"
	"github.com/layer-x/layerx-commons/lxdatabase"
	"github.com/layer-x/layerx-core_v2/fakes"
	"github.com/mesos/mesos-go/mesosproto"
	"github.com/layer-x/layerx-core_v2/lxtypes"
)


func PurgeState() {
	lxdatabase.Rmdir("/state", true)
}

var _ = Describe("Lxserver", func() {
	var lxRpiClient *layerx_rpi_client.LayerXRpi
	var lxTpiClient *layerx_tpi_client.LayerXTpi
	var state *lxstate.State

	Describe("setup", func() {
		It("sets up for the tests", func() {
			lxRpiClient = &layerx_rpi_client.LayerXRpi{
				CoreURL: "127.0.0.1:6677",
			}
			lxTpiClient = &layerx_tpi_client.LayerXTpi{
				CoreURL: "127.0.0.1:6677",
			}

			actionQueue := lxactionqueue.NewActionQueue()
			state = lxstate.NewState()
			err := state.InitializeState("http://127.0.0.1:4001")
			Expect(err).To(BeNil())
			coreServerWrapper := NewLayerXCoreServerWrapper(state, actionQueue)
			driver := driver.NewLayerXDriver(actionQueue)

			m := coreServerWrapper.WrapServer(lxmartini.QuietMartini(), "127.0.0.1:6688", "127.0.0.1:6699", make(chan error))
			go m.RunOnAddr(fmt.Sprintf(":6677"))
			go fakes.RunFakeTpiServer("127.0.0.1:6677", 6688, make(chan error))
			go fakes.RunFakeRpiServer("127.0.0.1:6677", 6699, make(chan error))
			go driver.Run()
			lxlog.ActiveDebugMode()
		})
	})

	Describe("RegisterTpi", func() {
		It("adds the Tpi URL to the LX state", func() {
			PurgeState()
			err := lxTpiClient.RegisterTpi("127.0.0.1:6688")
			Expect(err).To(BeNil())
			tpiUrl, err := state.GetTpi()
			Expect(err).To(BeNil())
			Expect(tpiUrl).To(Equal("127.0.0.1:6688"))
		})
	})

	Describe("RegisterRpi", func() {
		It("adds the Rpi URL to the LX state", func() {
			PurgeState()
			err := lxRpiClient.RegisterRpi("127.0.0.1:6699")
			Expect(err).To(BeNil())
			rpiUrl, err := state.GetRpi()
			Expect(err).To(BeNil())
			Expect(rpiUrl).To(Equal("127.0.0.1:6699"))
		})
	})

	Describe("RegisterTaskProvider", func() {
		It("adds the task provider to the LX state", func() {
			PurgeState()
			fakeTaskProvider := fakes.FakeTaskProvider("fake_framework", "ff@fakeip:fakeport")
			err := lxTpiClient.RegisterTaskProvider(fakeTaskProvider)
			Expect(err).To(BeNil())
			taskProvider, err := state.TaskProviderPool.GetTaskProvider("fake_framework")
			Expect(err).To(BeNil())
			Expect(taskProvider).To(Equal(fakeTaskProvider))
		})
	})

	Describe("DeregisterTaskProvider", func() {
		It("removes the task provider from the LX state", func() {
			PurgeState()
			fakeTaskProvider := fakes.FakeTaskProvider("fake_framework", "ff@fakeip:fakeport")
			err := lxTpiClient.RegisterTaskProvider(fakeTaskProvider)
			Expect(err).To(BeNil())
			err = lxTpiClient.DeregisterTaskProvider(fakeTaskProvider.Id)
			Expect(err).To(BeNil())
			taskProvider, err := state.TaskProviderPool.GetTaskProvider("fake_framework")
			Expect(err).NotTo(BeNil())
			Expect(taskProvider).To(BeNil())
		})
	})

	Describe("GetTaskProviders", func() {
		It("gets the list of task providers that have been registered", func() {
			PurgeState()
			fakeTaskProvider1 := fakes.FakeTaskProvider("fake_framework1", "ff@fakeip:fakeport")
			fakeTaskProvider2 := fakes.FakeTaskProvider("fake_framework2", "ff@fakeip:fakeport")
			fakeTaskProvider3 := fakes.FakeTaskProvider("fake_framework3", "ff@fakeip:fakeport")
			err := lxTpiClient.RegisterTaskProvider(fakeTaskProvider1)
			Expect(err).To(BeNil())
			err = lxTpiClient.RegisterTaskProvider(fakeTaskProvider2)
			Expect(err).To(BeNil())
			err = lxTpiClient.RegisterTaskProvider(fakeTaskProvider3)
			Expect(err).To(BeNil())
			taskProviders, err := state.TaskProviderPool.GetTaskProviders()
			Expect(err).To(BeNil())
			Expect(taskProviders).To(ContainElement(fakeTaskProvider1))
			Expect(taskProviders).To(ContainElement(fakeTaskProvider2))
			Expect(taskProviders).To(ContainElement(fakeTaskProvider3))
			taskProviderArr, err := lxTpiClient.GetTaskProviders()
			Expect(err).To(BeNil())
			Expect(taskProviderArr).To(ContainElement(fakeTaskProvider1))
			Expect(taskProviderArr).To(ContainElement(fakeTaskProvider2))
			Expect(taskProviderArr).To(ContainElement(fakeTaskProvider3))
		})
	})

	Describe("GetStatusUpdates", func() {
		It("gets the list of status updates for the given task provider", func() {
			PurgeState()
			err := state.InitializeState("http://127.0.0.1:4001")
			Expect(err).To(BeNil())
			fakeTaskProvider := fakes.FakeTaskProvider("fake_framework", "ff@fakeip:fakeport")
			err = lxTpiClient.RegisterTaskProvider(fakeTaskProvider)
			Expect(err).To(BeNil())
			fakeTask1 := fakes.FakeLXTask("fake_task_id_1", "fake_task1", "fake_node_id_1", "echo FAKECOMMAND")
			fakeTask1.TaskProvider = fakes.FakeTaskProvider("fake_task_provider_id", "tp@fakeip:fakeport")
			fakeTask2 := fakes.FakeLXTask("fake_task_id_2", "fake_task2", "fake_node_id_1", "echo FAKECOMMAND")
			fakeTask2.TaskProvider = fakes.FakeTaskProvider("fake_task_provider_id", "tp@fakeip:fakeport")
			fakeTask3 := fakes.FakeLXTask("fake_task_id_3", "fake_task3", "fake_node_id_1", "echo FAKECOMMAND")
			fakeTask3.TaskProvider = fakes.FakeTaskProvider("fake_task_provider_id", "tp@fakeip:fakeport")
			fakeTask1.TaskProvider = fakeTaskProvider
			fakeTask2.TaskProvider = fakeTaskProvider
			fakeTask3.TaskProvider = fakeTaskProvider
			err = state.StagingTaskPool.AddTask(fakeTask1)
			Expect(err).To(BeNil())
			err = state.StagingTaskPool.AddTask(fakeTask2)
			Expect(err).To(BeNil())
			err = state.StagingTaskPool.AddTask(fakeTask3)
			Expect(err).To(BeNil())
			fakeStatusUpdate1 := fakes.FakeTaskStatus("fake_task_id_1", mesosproto.TaskState_TASK_RUNNING)
			fakeStatusUpdate2 := fakes.FakeTaskStatus("fake_task_id_2", mesosproto.TaskState_TASK_KILLED)
			fakeStatusUpdate3 := fakes.FakeTaskStatus("fake_task_id_3", mesosproto.TaskState_TASK_ERROR)
			err = state.StatusPool.AddStatus(fakeStatusUpdate1)
			Expect(err).To(BeNil())
			err = state.StatusPool.AddStatus(fakeStatusUpdate2)
			Expect(err).To(BeNil())
			err = state.StatusPool.AddStatus(fakeStatusUpdate3)
			Expect(err).To(BeNil())
			statuses, err := lxTpiClient.GetStatusUpdates("fake_framework")
			Expect(err).To(BeNil())
			Expect(statuses).To(ContainElement(fakeStatusUpdate1))
			Expect(statuses).To(ContainElement(fakeStatusUpdate2))
			Expect(statuses).To(ContainElement(fakeStatusUpdate3))
		})
	})

	Describe("SubmitTask", func() {
		It("adds the task to the pending task pool, sets the task provider info for the task", func() {
			PurgeState()
			fakeTaskProvider := fakes.FakeTaskProvider("fake_framework", "ff@fakeip:fakeport")
			err := lxTpiClient.RegisterTaskProvider(fakeTaskProvider)
			Expect(err).To(BeNil())
			fakeTask1 := fakes.FakeLXTask("fake_task_id_1", "fake_task1", "fake_node_id_1", "echo FAKECOMMAND")
			fakeTask1.TaskProvider = fakes.FakeTaskProvider("fake_task_provider_id", "tp@fakeip:fakeport")
			err = lxTpiClient.SubmitTask("fake_framework", fakeTask1)
			Expect(err).To(BeNil())
			task1, err := state.PendingTaskPool.GetTask("fake_task_id_1")
			Expect(err).To(BeNil())
			fakeTask1.TaskProvider = fakeTaskProvider
			Expect(task1).To(Equal(fakeTask1))
		})
	})

	Describe("KillTask", func() {
		It("sets the flag KillRequested to true on the task", func() {
			PurgeState()
			fakeTask1 := fakes.FakeLXTask("fake_task_id_1", "fake_task1", "fake_node_id_1", "echo FAKECOMMAND")
			fakeTask1.TaskProvider = fakes.FakeTaskProvider("fake_task_provider_id", "tp@fakeip:fakeport")
			err := state.PendingTaskPool.AddTask(fakeTask1)
			Expect(err).To(BeNil())
			err = lxTpiClient.KillTask(fakeTask1.TaskId)
			Expect(err).To(BeNil())
			fakeTask1.KillRequested = true
			task1, err := state.PendingTaskPool.GetTask("fake_task_id_1")
			Expect(err).To(BeNil())
			Expect(task1).To(Equal(fakeTask1))
		})
	})

	Describe("PurgeTask", func() {
		It("deletes the task from the task pool", func() {
			PurgeState()
			fakeTask1 := fakes.FakeLXTask("fake_task_id_1", "fake_task1", "fake_node_id_1", "echo FAKECOMMAND")
			fakeTask1.TaskProvider = fakes.FakeTaskProvider("fake_task_provider_id", "tp@fakeip:fakeport")
			err := state.PendingTaskPool.AddTask(fakeTask1)
			Expect(err).To(BeNil())
			err = lxTpiClient.PurgeTask(fakeTask1.TaskId)
			Expect(err).To(BeNil())
			task1, err := state.PendingTaskPool.GetTask("fake_task_id_1")
			Expect(err).NotTo(BeNil())
			Expect(task1).To(BeNil())
		})
	})

	Describe("SubmitResource", func() {
		Context("no node exists for the nodeId", func() {
			It("creates a new node, addds the resource to the node", func() {
				PurgeState()
				fakeResource1 := lxtypes.NewResourceFromMesos(fakes.FakeOffer("fake_offer_id_1", "fake_slave_id_1"))
				err := lxRpiClient.SubmitResource(fakeResource1)
				Expect(err).To(BeNil())
				node, err := state.NodePool.GetNode(fakeResource1.NodeId)
				Expect(err).To(BeNil())
				Expect(node.Resources).To(ContainElement(fakeResource1))
				Expect(node.GetFreeCpus()).To(Equal(fakeResource1.Cpus))
				Expect(node.GetFreeMem()).To(Equal(fakeResource1.Mem))
				Expect(node.GetFreeDisk()).To(Equal(fakeResource1.Disk))
			})
		})
		Context("a node already exists for the nodeId", func() {
			It("adds the resource to the node", func() {
				err := state.InitializeState("http://127.0.0.1:4001")
				PurgeState()
				fakeResource1 := lxtypes.NewResourceFromMesos(fakes.FakeOffer("fake_offer_id_1", "fake_slave_id_1"))
				err = lxRpiClient.SubmitResource(fakeResource1)
				Expect(err).To(BeNil())
				fakeResource2 := lxtypes.NewResourceFromMesos(fakes.FakeOffer("fake_offer_id_2", "fake_slave_id_1"))
				err = lxRpiClient.SubmitResource(fakeResource2)
				Expect(err).To(BeNil())
				node, err := state.NodePool.GetNode(fakeResource1.NodeId)
				Expect(err).To(BeNil())
				Expect(node.Resources).To(ContainElement(fakeResource1))
				Expect(node.Resources).To(ContainElement(fakeResource2))
				Expect(node.GetFreeCpus()).To(Equal(fakeResource1.Cpus + fakeResource2.Cpus))
				Expect(node.GetFreeMem()).To(Equal(fakeResource1.Mem + fakeResource2.Mem))
				Expect(node.GetFreeDisk()).To(Equal(fakeResource1.Disk + fakeResource2.Disk))
			})
		})
	})

	Describe("SubmitStatusUpdate", func() {
		It("adds the status to the lx state", func() {
			PurgeState()
			purgeErr := state.InitializeState("http://127.0.0.1:4001")
			Expect(purgeErr).To(BeNil())
			fakeTask1 := fakes.FakeLXTask("fake_task_id_1", "fake_task1", "fake_node_id_1", "echo FAKECOMMAND")
			fakeTask1.TaskProvider = fakes.FakeTaskProvider("fake_task_provider_id", "tp@fakeip:fakeport")
			err := state.StagingTaskPool.AddTask(fakeTask1)
			Expect(err).To(BeNil())
			fakeStatus1 := fakes.FakeTaskStatus("fake_task_id_1", mesosproto.TaskState_TASK_KILLED)
			err = lxRpiClient.SubmitStatusUpdate(fakeStatus1)
			Expect(err).To(BeNil())
			status, err := state.StatusPool.GetStatus(fakeStatus1.GetTaskId().GetValue())
			Expect(err).To(BeNil())
			Expect(status).To(Equal(fakeStatus1))
		})
	})
	Describe("GetNodes", func(){
		It("returns a list of nodes", func(){
			PurgeState()
			fakeResource1 := lxtypes.NewResourceFromMesos(fakes.FakeOffer("fake_offer_id_1", "fake_slave_id_1"))
			err := lxRpiClient.SubmitResource(fakeResource1)
			Expect(err).To(BeNil())
			fakeResource2 := lxtypes.NewResourceFromMesos(fakes.FakeOffer("fake_offer_id_2", "fake_slave_id_2"))
			err = lxRpiClient.SubmitResource(fakeResource2)
			Expect(err).To(BeNil())
			fakeResource3 := lxtypes.NewResourceFromMesos(fakes.FakeOffer("fake_offer_id_3", "fake_slave_id_3"))
			err = lxRpiClient.SubmitResource(fakeResource3)
			Expect(err).To(BeNil())
			nodes, err := lxRpiClient.GetNodes()
			Expect(err).To(BeNil())
			node1 := lxtypes.NewNode(fakeResource1.NodeId)
			err = node1.AddResource(fakeResource1)
			Expect(err).To(BeNil())
			node2 := lxtypes.NewNode(fakeResource2.NodeId)
			err = node2.AddResource(fakeResource2)
			Expect(err).To(BeNil())
			node3 := lxtypes.NewNode(fakeResource3.NodeId)
			err = node3.AddResource(fakeResource3)
			Expect(err).To(BeNil())

			Expect(nodes).To(ContainElement(node1))
			Expect(nodes).To(ContainElement(node2))
			Expect(nodes).To(ContainElement(node3))
		})
	})
})
