package lxserver_test

import (
	. "github.com/layer-x/layerx-core_v2/lxserver"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/layer-x/layerx-core_v2/layerx_rpi_client"
	"github.com/layer-x/layerx-core_v2/layerx_tpi_client"
	"github.com/layer-x/layerx-commons/lxmartini"
	"github.com/layer-x/layerx-core_v2/lxstate"
	"github.com/layer-x/layerx-commons/lxlog"
	"fmt"
	"github.com/layer-x/layerx-commons/lxdatabase"
	"github.com/layer-x/layerx-core_v2/fakes"
	"github.com/mesos/mesos-go/mesosproto"
	"github.com/layer-x/layerx-core_v2/lxtypes"
	"github.com/layer-x/layerx-core_v2/layerx_brain_client"
)


func PurgeState() {
	lxdatabase.Rmdir("/state", true)
}

var _ = Describe("Lxserver", func() {
	var lxRpiClient *layerx_rpi_client.LayerXRpi
	var lxTpiClient *layerx_tpi_client.LayerXTpi
	var lxBrainClient *layerx_brain_client.LayerXBrainClient
	var state *lxstate.State
	var serverErr error

	Describe("setup", func() {
		It("sets up for the tests", func() {
			lxRpiClient = &layerx_rpi_client.LayerXRpi{
				CoreURL: "127.0.0.1:6677",
				RpiName: "fake-rpi",
			}
			lxTpiClient = &layerx_tpi_client.LayerXTpi{
				CoreURL: "127.0.0.1:6677",
			}
			lxBrainClient = &layerx_brain_client.LayerXBrainClient{
				CoreURL: "127.0.0.1:6677",
			}

			state = lxstate.NewState()
			err := state.InitializeState("http://127.0.0.1:4001")
			Expect(err).To(BeNil())

			driverErrc := make(chan error)
			go func() {
				for {
					serverErr = <-driverErrc
				}
			}()

			coreServerWrapper := NewLayerXCoreServerWrapper(state, lxmartini.QuietMartini(), driverErrc)

			err = state.SetTpi( "127.0.0.1:6688")
			Expect(err).To(BeNil())
			err = state.RpiPool.AddRpi(&layerx_rpi_client.RpiInfo{
				Name: "fake-rpi",
				Url: "127.0.0.1:6699",
			})
			Expect(err).To(BeNil())

			m := coreServerWrapper.WrapServer()
			go m.RunOnAddr(fmt.Sprintf(":6677"))
			go fakes.RunFakeTpiServer("127.0.0.1:6677", 6688, make(chan error))
			go fakes.RunFakeRpiServer("127.0.0.1:6677", 6699, make(chan error))
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
			err := lxRpiClient.RegisterRpi("fake-rpi", "127.0.0.1:6699")
			Expect(err).To(BeNil())
			rpis, err := state.RpiPool.GetRpis()
			Expect(err).To(BeNil())
			fakeRpi := rpis["fake-rpi"]
			Expect(fakeRpi).NotTo(BeNil())
			Expect(fakeRpi.Name).To(Equal("fake-rpi"))
			Expect(fakeRpi.Url).To(Equal("127.0.0.1:6699"))
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

	Describe("GetStatusUpdates(taskProvider)", func() {
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

	Describe("GetStatusUpdates", func() {
		It("gets all status updates", func() {
			PurgeState()
			err := state.InitializeState("http://127.0.0.1:4001")
			Expect(err).To(BeNil())
			fakeTaskProvider1 := fakes.FakeTaskProvider("fake_framework", "ff@fakeip:fakeport")
			err = lxTpiClient.RegisterTaskProvider(fakeTaskProvider1)
			Expect(err).To(BeNil())
			fakeTask1 := fakes.FakeLXTask("fake_task_id_1", "fake_task1", "fake_node_id_1", "echo FAKECOMMAND")
			fakeTask1.TaskProvider = fakes.FakeTaskProvider("fake_task_provider_id", "tp@fakeip:fakeport")
			fakeTask2 := fakes.FakeLXTask("fake_task_id_2", "fake_task2", "fake_node_id_1", "echo FAKECOMMAND")
			fakeTask2.TaskProvider = fakes.FakeTaskProvider("fake_task_provider_id", "tp@fakeip:fakeport")
			fakeTask3 := fakes.FakeLXTask("fake_task_id_3", "fake_task3", "fake_node_id_1", "echo FAKECOMMAND")
			fakeTask3.TaskProvider = fakes.FakeTaskProvider("fake_task_provider_id", "tp@fakeip:fakeport")
			fakeTask1.TaskProvider = fakeTaskProvider1
			fakeTask2.TaskProvider = fakeTaskProvider1
			fakeTask3.TaskProvider = fakeTaskProvider1
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

			fakeTaskProvider2 := fakes.FakeTaskProvider("fake_framework", "ff@fakeip:fakeport")
			err = lxTpiClient.RegisterTaskProvider(fakeTaskProvider2)
			Expect(err).To(BeNil())
			fakeTask4 := fakes.FakeLXTask("fake_task_id_4", "fake_task4", "fake_node_id_1", "echo FAKECOMMAND")
			fakeTask4.TaskProvider = fakes.FakeTaskProvider("fake_task_provider_id", "tp@fakeip:fakeport")
			fakeTask5 := fakes.FakeLXTask("fake_task_id_5", "fake_task5", "fake_node_id_1", "echo FAKECOMMAND")
			fakeTask5.TaskProvider = fakes.FakeTaskProvider("fake_task_provider_id", "tp@fakeip:fakeport")
			fakeTask6 := fakes.FakeLXTask("fake_task_id_6", "fake_task6", "fake_node_id_1", "echo FAKECOMMAND")
			fakeTask6.TaskProvider = fakes.FakeTaskProvider("fake_task_provider_id", "tp@fakeip:fakeport")
			fakeTask4.TaskProvider = fakeTaskProvider2
			fakeTask5.TaskProvider = fakeTaskProvider2
			fakeTask6.TaskProvider = fakeTaskProvider2
			err = state.StagingTaskPool.AddTask(fakeTask4)
			Expect(err).To(BeNil())
			err = state.StagingTaskPool.AddTask(fakeTask5)
			Expect(err).To(BeNil())
			err = state.StagingTaskPool.AddTask(fakeTask6)
			Expect(err).To(BeNil())
			fakeStatusUpdate4 := fakes.FakeTaskStatus("fake_task_id_4", mesosproto.TaskState_TASK_RUNNING)
			fakeStatusUpdate5 := fakes.FakeTaskStatus("fake_task_id_5", mesosproto.TaskState_TASK_KILLED)
			fakeStatusUpdate6 := fakes.FakeTaskStatus("fake_task_id_6", mesosproto.TaskState_TASK_ERROR)
			err = state.StatusPool.AddStatus(fakeStatusUpdate4)
			Expect(err).To(BeNil())
			err = state.StatusPool.AddStatus(fakeStatusUpdate5)
			Expect(err).To(BeNil())
			err = state.StatusPool.AddStatus(fakeStatusUpdate6)
			Expect(err).To(BeNil())

			statuses, err := lxBrainClient.GetStatusUpdates()
			Expect(err).To(BeNil())
			Expect(statuses).To(ContainElement(fakeStatusUpdate1))
			Expect(statuses).To(ContainElement(fakeStatusUpdate2))
			Expect(statuses).To(ContainElement(fakeStatusUpdate3))
			Expect(statuses).To(ContainElement(fakeStatusUpdate4))
			Expect(statuses).To(ContainElement(fakeStatusUpdate5))
			Expect(statuses).To(ContainElement(fakeStatusUpdate6))
		})
	})

	Describe("GetStatusUpdate", func() {
		Context("the task status exists in the database", func(){
			It("returns that status", func(){
				PurgeState()
				err := state.InitializeState("http://127.0.0.1:4001")
				Expect(err).To(BeNil())
				fakeTaskProvider := fakes.FakeTaskProvider("fake_framework", "ff@fakeip:fakeport")
				err = lxTpiClient.RegisterTaskProvider(fakeTaskProvider)
				Expect(err).To(BeNil())
				fakeTask1 := fakes.FakeLXTask("fake_task_id_1", "fake_task1", "fake_node_id_1", "echo FAKECOMMAND")
				fakeTask1.TaskProvider = fakes.FakeTaskProvider("fake_task_provider_id", "tp@fakeip:fakeport")
				fakeTask1.TaskProvider = fakeTaskProvider
				err = state.StagingTaskPool.AddTask(fakeTask1)
				Expect(err).To(BeNil())
				fakeStatusUpdate1 := fakes.FakeTaskStatus("fake_task_id_1", mesosproto.TaskState_TASK_RUNNING)
				err = state.StatusPool.AddStatus(fakeStatusUpdate1)
				Expect(err).To(BeNil())
				status, err := lxTpiClient.GetStatusUpdate("fake_task_id_1")
				Expect(err).To(BeNil())
				Expect(status).To(Equal(fakeStatusUpdate1))
			})
		})
		Context("the task is in the pending pool", func(){
			It("returns a status with TASK_STAGING", func(){
				PurgeState()
				err := state.InitializeState("http://127.0.0.1:4001")
				Expect(err).To(BeNil())
				fakeTaskProvider := fakes.FakeTaskProvider("fake_framework", "ff@fakeip:fakeport")
				err = lxTpiClient.RegisterTaskProvider(fakeTaskProvider)
				Expect(err).To(BeNil())
				fakeTask1 := fakes.FakeLXTask("fake_task_id_1", "fake_task1", "fake_node_id_1", "echo FAKECOMMAND")
				fakeTask1.TaskProvider = fakes.FakeTaskProvider("fake_task_provider_id", "tp@fakeip:fakeport")
				fakeTask1.TaskProvider = fakeTaskProvider
				err = state.PendingTaskPool.AddTask(fakeTask1)
				Expect(err).To(BeNil())
				status, err := lxTpiClient.GetStatusUpdate("fake_task_id_1")
				Expect(err).To(BeNil())
				Expect(status.GetState()).To(Equal(mesosproto.TaskState_TASK_STAGING))
			})
		})
		Context("the task is in the staging pool", func(){
			It("returns a status with TASK_STARTING", func(){
				PurgeState()
				err := state.InitializeState("http://127.0.0.1:4001")
				Expect(err).To(BeNil())
				fakeTaskProvider := fakes.FakeTaskProvider("fake_framework", "ff@fakeip:fakeport")
				err = lxTpiClient.RegisterTaskProvider(fakeTaskProvider)
				Expect(err).To(BeNil())
				fakeTask1 := fakes.FakeLXTask("fake_task_id_1", "fake_task1", "fake_node_id_1", "echo FAKECOMMAND")
				fakeTask1.TaskProvider = fakes.FakeTaskProvider("fake_task_provider_id", "tp@fakeip:fakeport")
				fakeTask1.TaskProvider = fakeTaskProvider
				err = state.StagingTaskPool.AddTask(fakeTask1)
				Expect(err).To(BeNil())
				status, err := lxTpiClient.GetStatusUpdate("fake_task_id_1")
				Expect(err).To(BeNil())
				Expect(status.GetState()).To(Equal(mesosproto.TaskState_TASK_STARTING))
			})
		})
		Context("the does not exist", func(){
			It("returns a status with TASK_LOST", func(){
				PurgeState()
				err := state.InitializeState("http://127.0.0.1:4001")
				Expect(err).To(BeNil())
				status, err := lxTpiClient.GetStatusUpdate("fake_task_id_1")
				Expect(err).To(BeNil())
				Expect(status.GetState()).To(Equal(mesosproto.TaskState_TASK_LOST))
			})
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
		Context("task does not exist", func(){
			It("sends TASK_LOST status to tpi", func() {
				PurgeState()
				err := state.InitializeState("http://127.0.0.1:4001")
				Expect(err).To(BeNil())
				err = state.SetTpi( "127.0.0.1:6688")
				Expect(err).To(BeNil())
				err = state.RpiPool.AddRpi(&layerx_rpi_client.RpiInfo{
					Name: "fake-rpi",
					Url: "127.0.0.1:6699",
				})
				fakeTaskProvider := fakes.FakeTaskProvider("fake_task_provider_id", "tp@fakeip:fakeport")
				err = lxTpiClient.RegisterTaskProvider(fakeTaskProvider)
				Expect(err).To(BeNil())
				err = lxTpiClient.KillTask(fakeTaskProvider.Id, "nonexistent_task_id")
				Expect(err).To(BeNil())
			})
		})
		Context("task staging is not complete", func(){
			It("deletes the task from staging or pending pool and sends TASK_KILLED status to tpi", func() {
				PurgeState()
				err := state.InitializeState("http://127.0.0.1:4001")
				err = state.SetTpi( "127.0.0.1:6688")
				Expect(err).To(BeNil())
				err = state.RpiPool.AddRpi(&layerx_rpi_client.RpiInfo{
					Name: "fake-rpi",
					Url: "127.0.0.1:6699",
				})
				Expect(err).To(BeNil())
				fakeResource1 := lxtypes.NewResourceFromMesos(fakes.FakeOffer("fake_offer_id_1", "fake_slave_id_1"))
				err = lxRpiClient.SubmitResource(fakeResource1)
				Expect(err).To(BeNil())
				fakeTaskProvider := fakes.FakeTaskProvider("fake_task_provider_id", "tp@fakeip:fakeport")
				err = lxTpiClient.RegisterTaskProvider(fakeTaskProvider)
				Expect(err).To(BeNil())
				fakeTask1 := fakes.FakeLXTask("fake_task_id_1", "fake_task1", "fake_node_id_1", "echo FAKECOMMAND")
				fakeTask1.TaskProvider = fakeTaskProvider
				err = state.PendingTaskPool.AddTask(fakeTask1)
				Expect(err).To(BeNil())
				err = lxBrainClient.AssignTasks(fakeResource1.NodeId, fakeTask1.TaskId)
				Expect(err).To(BeNil())
				err = lxTpiClient.KillTask(fakeTaskProvider.Id, fakeTask1.TaskId)
				Expect(err).To(BeNil())
				task1, err := state.GetTaskFromAnywhere(fakeTask1.TaskId)
				Expect(err).NotTo(BeNil())
				Expect(task1).To(BeNil())
			})
		})
		Context("task staging is complete", func(){
			It("sets the flag KillRequested to true on the task and sends KillTask request to RPI", func() {
				PurgeState()
				err := state.InitializeState("http://127.0.0.1:4001")
				Expect(err).To(BeNil())
				err = state.SetTpi( "127.0.0.1:6688")
				Expect(err).To(BeNil())
				err = state.RpiPool.AddRpi(&layerx_rpi_client.RpiInfo{
					Name: "fake-rpi",
					Url: "127.0.0.1:6699",
				})
				fakeResource1 := lxtypes.NewResourceFromMesos(fakes.FakeOffer("fake_offer_id_1", "fake_node_id_1"))
				err = lxRpiClient.SubmitResource(fakeResource1)
				Expect(err).To(BeNil())
				fakeTask1 := fakes.FakeLXTask("fake_task_id_1", "fake_task1", "fake_node_id_1", "echo FAKECOMMAND")
				fakeTask1.TaskProvider = fakes.FakeTaskProvider("fake_task_provider_id", "tp@fakeip:fakeport")
				nodeTaskPool, err := state.NodePool.GetNodeTaskPool(fakeResource1.NodeId)
				Expect(err).To(BeNil())
				err = nodeTaskPool.AddTask(fakeTask1)
				Expect(err).To(BeNil())
				err = lxTpiClient.KillTask(fakeTask1.TaskProvider.Id, fakeTask1.TaskId)
				Expect(err).To(BeNil())
				fakeTask1.KillRequested = true
				task1, err := nodeTaskPool.GetTask("fake_task_id_1")
				Expect(err).To(BeNil())
				Expect(task1).To(Equal(fakeTask1))
			})
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
			err := state.SetTpi( "127.0.0.1:6688")
			Expect(err).To(BeNil())
			err = state.RpiPool.AddRpi(&layerx_rpi_client.RpiInfo{
				Name: "fake-rpi",
				Url: "127.0.0.1:6699",
			})
			fakeTask1 := fakes.FakeLXTask("fake_task_id_1", "fake_task1", "fake_node_id_1", "echo FAKECOMMAND")
			fakeTask1.TaskProvider = fakes.FakeTaskProvider("fake_task_provider_id", "tp@fakeip:fakeport")
			err = state.StagingTaskPool.AddTask(fakeTask1)
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

	Describe("GetPendingTasks", func(){
		It("returns all tasks in the pending task pool", func(){
			PurgeState()
			fakeTaskProvider := fakes.FakeTaskProvider("fake_framework", "ff@fakeip:fakeport")
			err := lxTpiClient.RegisterTaskProvider(fakeTaskProvider)
			Expect(err).To(BeNil())
			fakeTask1 := fakes.FakeLXTask("fake_task_id_1", "fake_task1", "fake_node_id_1", "echo FAKECOMMAND")
			fakeTask1.TaskProvider = fakeTaskProvider
			err = state.PendingTaskPool.AddTask(fakeTask1)
			Expect(err).To(BeNil())
			fakeTask2 := fakes.FakeLXTask("fake_task_id_2", "fake_task2", "fake_node_id_1", "echo FAKECOMMAND")
			fakeTask2.TaskProvider = fakeTaskProvider
			err = state.PendingTaskPool.AddTask(fakeTask2)
			fakeTask3 := fakes.FakeLXTask("fake_task_id_3", "fake_task2", "fake_node_id_1", "echo FAKECOMMAND")
			fakeTask3.TaskProvider = fakeTaskProvider
			err = state.StagingTaskPool.AddTask(fakeTask3)
			Expect(err).To(BeNil())
			tasks, err := lxBrainClient.GetPendingTasks()
			Expect(err).To(BeNil())
			Expect(tasks).To(ContainElement(fakeTask1))
			Expect(tasks).To(ContainElement(fakeTask2))
			Expect(tasks).NotTo(ContainElement(fakeTask3))
		})
	})

	Describe("GetStagingTasks", func(){
		It("returns all tasks in the staging task pool", func(){
			PurgeState()
			fakeTaskProvider := fakes.FakeTaskProvider("fake_framework", "ff@fakeip:fakeport")
			err := lxTpiClient.RegisterTaskProvider(fakeTaskProvider)
			Expect(err).To(BeNil())
			fakeTask1 := fakes.FakeLXTask("fake_task_id_1", "fake_task1", "fake_node_id_1", "echo FAKECOMMAND")
			fakeTask1.TaskProvider = fakeTaskProvider
			err = state.PendingTaskPool.AddTask(fakeTask1)
			Expect(err).To(BeNil())
			fakeTask2 := fakes.FakeLXTask("fake_task_id_2", "fake_task2", "fake_node_id_1", "echo FAKECOMMAND")
			fakeTask2.TaskProvider = fakeTaskProvider
			err = state.PendingTaskPool.AddTask(fakeTask2)
			fakeTask3 := fakes.FakeLXTask("fake_task_id_3", "fake_task2", "fake_node_id_1", "echo FAKECOMMAND")
			fakeTask3.TaskProvider = fakeTaskProvider
			err = state.StagingTaskPool.AddTask(fakeTask3)
			Expect(err).To(BeNil())
			tasks, err := lxBrainClient.GetStagingTasks()
			Expect(err).To(BeNil())
			Expect(tasks).NotTo(ContainElement(fakeTask1))
			Expect(tasks).NotTo(ContainElement(fakeTask2))
			Expect(tasks).To(ContainElement(fakeTask3))
		})
	})

	Describe("AssignTasks", func(){
		It("moves a list of tasks to the staging pool, gives them the SlaveId of the target Node", func(){
			PurgeState()
			fakeTaskProvider := fakes.FakeTaskProvider("fake_framework", "ff@fakeip:fakeport")
			err := lxTpiClient.RegisterTaskProvider(fakeTaskProvider)
			Expect(err).To(BeNil())
			fakeTask1 := fakes.FakeLXTask("fake_task_id_1", "fake_task1", "fake_node_id_1", "echo FAKECOMMAND")
			fakeTask1.TaskProvider = fakeTaskProvider
			err = state.PendingTaskPool.AddTask(fakeTask1)
			Expect(err).To(BeNil())
			fakeTask2 := fakes.FakeLXTask("fake_task_id_2", "fake_task2", "fake_node_id_1", "echo FAKECOMMAND")
			fakeTask2.TaskProvider = fakeTaskProvider
			err = state.PendingTaskPool.AddTask(fakeTask2)
			fakeTask3 := fakes.FakeLXTask("fake_task_id_3", "fake_task2", "fake_node_id_1", "echo FAKECOMMAND")
			fakeTask3.TaskProvider = fakeTaskProvider
			err = state.PendingTaskPool.AddTask(fakeTask3)
			Expect(err).To(BeNil())
			fakeResource1 := lxtypes.NewResourceFromMesos(fakes.FakeOffer("fake_offer_id_1", "fake_slave_id_1"))
			err = lxRpiClient.SubmitResource(fakeResource1)
			Expect(err).To(BeNil())
			node, err := state.NodePool.GetNode(fakeResource1.NodeId)
			Expect(err).To(BeNil())
			Expect(node.Resources).To(ContainElement(fakeResource1))
			err = lxBrainClient.AssignTasks(fakeResource1.NodeId, fakeTask1.TaskId, fakeTask2.TaskId, fakeTask3.TaskId)
			Expect(err).To(BeNil())
			tasks, err := state.PendingTaskPool.GetTasks()
			Expect(err).To(BeNil())
			Expect(tasks).To(BeEmpty())
			tasks, err = state.StagingTaskPool.GetTasks()
			Expect(tasks).NotTo(ContainElement(fakeTask1))
			fakeTask1.NodeId = fakeResource1.NodeId
			fakeTask2.NodeId = fakeResource1.NodeId
			fakeTask3.NodeId = fakeResource1.NodeId
			Expect(tasks).To(ContainElement(fakeTask1))
			Expect(tasks).To(ContainElement(fakeTask2))
			Expect(tasks).To(ContainElement(fakeTask3))
		})
	})
	Describe("MigrateTasks", func(){
		It("moves a list of runnning tasks on various nodes back to the staging pool, gives them the SlaveId of the target Node, and sets Checkpointed=true", func(){
			PurgeState()
			purgeErr := state.InitializeState("http://127.0.0.1:4001")
			Expect(purgeErr).To(BeNil())
			err := state.SetTpi( "127.0.0.1:6688")
			Expect(err).To(BeNil())
			err = state.RpiPool.AddRpi(&layerx_rpi_client.RpiInfo{
				Name: "fake-rpi",
				Url: "127.0.0.1:6699",
			})
			Expect(err).To(BeNil())
			fakeResource1 := lxtypes.NewResourceFromMesos(fakes.FakeOffer("fake_offer_id_1", "fake_node_id_1"))
			err = lxRpiClient.SubmitResource(fakeResource1)
			Expect(err).To(BeNil())
			fakeResource2 := lxtypes.NewResourceFromMesos(fakes.FakeOffer("fake_offer_id_2", "fake_node_id_2"))
			err = lxRpiClient.SubmitResource(fakeResource2)
			Expect(err).To(BeNil())
			nodeTaskPool1, err := state.NodePool.GetNodeTaskPool(fakeResource1.NodeId)
			Expect(err).To(BeNil())
			nodeTaskPool2, err := state.NodePool.GetNodeTaskPool(fakeResource2.NodeId)
			Expect(err).To(BeNil())
			fakeTaskProvider := fakes.FakeTaskProvider("fake_framework", "ff@fakeip:fakeport")
			err = lxTpiClient.RegisterTaskProvider(fakeTaskProvider)
			Expect(err).To(BeNil())
			fakeTask1 := fakes.FakeLXTask("fake_task_id_1", "fake_task1", "fake_node_id_1", "echo FAKECOMMAND")
			fakeTask1.TaskProvider = fakeTaskProvider
			err = nodeTaskPool1.AddTask(fakeTask1)
			Expect(err).To(BeNil())
			fakeTask2 := fakes.FakeLXTask("fake_task_id_2", "fake_task2", "fake_node_id_1", "echo FAKECOMMAND")
			fakeTask2.TaskProvider = fakeTaskProvider
			err = nodeTaskPool1.AddTask(fakeTask2)
			fakeTask3 := fakes.FakeLXTask("fake_task_id_3", "fake_task2", "fake_node_id_1", "echo FAKECOMMAND")
			fakeTask3.TaskProvider = fakeTaskProvider
			err = nodeTaskPool2.AddTask(fakeTask3)
			Expect(err).To(BeNil())

			err = lxBrainClient.MigrateTasks(fakeResource2.NodeId, fakeTask1.TaskId, fakeTask2.TaskId, fakeTask3.TaskId)
			Expect(err).To(BeNil())
			tasks, err := nodeTaskPool1.GetTasks()
			Expect(err).To(BeNil())
			Expect(tasks).To(BeEmpty())
			tasks, err = state.StagingTaskPool.GetTasks()
			Expect(tasks).NotTo(ContainElement(fakeTask1))
			fakeTask1.NodeId = fakeResource2.NodeId
			fakeTask1.Checkpointed = true
			fakeTask2.NodeId = fakeResource2.NodeId
			fakeTask2.Checkpointed = true
			Expect(tasks).To(ContainElement(fakeTask1))
			Expect(tasks).To(ContainElement(fakeTask2))
			Expect(tasks).NotTo(ContainElement(fakeTask3))
		})
	})
})
