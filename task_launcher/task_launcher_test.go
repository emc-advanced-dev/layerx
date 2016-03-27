package task_launcher_test

import (
	. "github.com/layer-x/layerx-core_v2/task_launcher"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/layer-x/layerx-core_v2/layerx_rpi_client"
	"github.com/layer-x/layerx-core_v2/layerx_tpi_client"
	"github.com/layer-x/layerx-core_v2/layerx_brain_client"
	"github.com/layer-x/layerx-core_v2/lxstate"
	"github.com/layer-x/layerx-commons/lxmartini"
	"fmt"
"github.com/layer-x/layerx-core_v2/fakes"
"github.com/layer-x/layerx-commons/lxlog"
	"github.com/layer-x/layerx-core_v2/lxserver"
	"github.com/layer-x/layerx-commons/lxdatabase"
"github.com/layer-x/layerx-core_v2/lxtypes"
)

func PurgeState() {
	lxdatabase.Rmdir("/state", true)
}

var _ = Describe("TaskLauncher", func() {

	var lxRpiClient *layerx_rpi_client.LayerXRpi
	var lxTpiClient *layerx_tpi_client.LayerXTpi
	var lxBrainClient *layerx_brain_client.LayerXBrainClient
	var state *lxstate.State
	var serverErr error

	Describe("setup", func() {
		It("sets up for the tests", func() {
			lxRpiClient = &layerx_rpi_client.LayerXRpi{
				CoreURL: "127.0.0.1:2277",
			}
			lxTpiClient = &layerx_tpi_client.LayerXTpi{
				CoreURL: "127.0.0.1:2277",
			}
			lxBrainClient = &layerx_brain_client.LayerXBrainClient{
				CoreURL: "127.0.0.1:2277",
			}

			state = lxstate.NewState()
			err := state.InitializeState("http://127.0.0.1:4001")
			Expect(err).To(BeNil())
			driverErrc := make(chan error)
			coreServerWrapper := lxserver.NewLayerXCoreServerWrapper(state, lxmartini.QuietMartini(), driverErrc)

			err = state.SetTpi( "127.0.0.1:2288")
			Expect(err).To(BeNil())
			err = state.RpiPool.AddRpi(&layerx_rpi_client.RpiInfo{
				Name: "fake-rpi",
				Url: "127.0.0.1:2299",
			})

			go func() {
				for {
					serverErr = <-driverErrc
				}
			}()

			m := coreServerWrapper.WrapServer()
			go m.RunOnAddr(fmt.Sprintf(":2277"))
			go fakes.RunFakeTpiServer("127.0.0.1:2277", 2288, make(chan error))
			go fakes.RunFakeRpiServer("127.0.0.1:2277", 2299, make(chan error))
			lxlog.ActiveDebugMode()
		})
	})
	Describe("LaunchStagedTasks", func(){
		It("sends LaunchTaskMessage to rpi for all tasks in the staging pool", func(){
			PurgeState()
			err2 := state.InitializeState("http://127.0.0.1:4001")
			Expect(err2).To(BeNil())
			err := state.SetTpi( "127.0.0.1:2288")
			Expect(err).To(BeNil())
			err = state.RpiPool.AddRpi(&layerx_rpi_client.RpiInfo{
				Name: "fake-rpi",
				Url: "127.0.0.1:2299",
			})
			fakeResource1 := lxtypes.NewResourceFromMesos(fakes.FakeOffer("fake_offer_id_1", "fake_slave_id_1"))
			fakeResource1.RpiName = "fake-rpi"
			fakeNode1 := lxtypes.NewNode(fakeResource1.NodeId)
			err = fakeNode1.AddResource(fakeResource1)
			Expect(err).To(BeNil())
			err = state.NodePool.AddNode(fakeNode1)
			fakeTask1 := fakes.FakeLXTask("fake_task_id_1", "fake_task_name", "fake_slave_id", "echo FAKE_COMMAND")
			fakeTaskProvider := &lxtypes.TaskProvider{
				Id:     "fake_task_provider_id_1",
				Source: "taskprovider1@tphost:port",
			}
			fakeTask1.SlaveId = fakeNode1.Id
			fakeTask1.TaskProvider = fakeTaskProvider
			fakeTask2 := fakes.FakeLXTask("fake_task_id_2", "fake_task_name", "fake_slave_id", "echo FAKE_COMMAND")
			fakeTask2.SlaveId = fakeNode1.Id
			fakeTask2.TaskProvider = fakeTaskProvider
			fakeTask3 := fakes.FakeLXTask("fake_task_id_3", "fake_task_name", "fake_slave_id", "echo FAKE_COMMAND")
			fakeTask3.SlaveId = fakeNode1.Id
			fakeTask3.TaskProvider = fakeTaskProvider

			err = state.StagingTaskPool.AddTask(fakeTask1)
			Expect(err).To(BeNil())

			err = state.StagingTaskPool.AddTask(fakeTask2)
			Expect(err).To(BeNil())

			err = state.StagingTaskPool.AddTask(fakeTask3)
			Expect(err).To(BeNil())

			taskLauncher := NewTaskLauncher(state)
			err = taskLauncher.LaunchStagedTasks()
			Expect(err).To(BeNil())

			stagingTasks, err := state.StagingTaskPool.GetTasks()
			Expect(err).To(BeNil())
			Expect(stagingTasks).To(BeEmpty())
			node1TaskPool, err := state.NodePool.GetNodeTaskPool(fakeNode1.Id)
			Expect(err).To(BeNil())
			node1Task1, err := node1TaskPool.GetTask(fakeTask1.TaskId)
			Expect(err).To(BeNil())
			Expect(node1Task1).To(Equal(fakeTask1))
			node1Task2, err := node1TaskPool.GetTask(fakeTask2.TaskId)
			Expect(node1Task2).To(Equal(fakeTask2))
			Expect(err).To(BeNil())
			node1Task3, err := node1TaskPool.GetTask(fakeTask3.TaskId)
			Expect(node1Task3).To(Equal(fakeTask3))
			Expect(err).To(BeNil())
			node1ResourcePool, err := state.NodePool.GetNodeResourcePool(fakeNode1.Id)
			Expect(err).To(BeNil())
			resourcesLeftAfterLaunch, err := node1ResourcePool.GetResources()
			Expect(err).To(BeNil())
			Expect(resourcesLeftAfterLaunch).To(BeEmpty())
		})
	})
})
