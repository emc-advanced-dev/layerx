package rpi_messenger_test

import (
	. "github.com/layer-x/layerx-core_v2/rpi_messenger"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/layer-x/layerx-commons/lxlog"
	"github.com/layer-x/layerx-core_v2/fakes"
	"fmt"
	"github.com/layer-x/layerx-commons/lxmartini"
	"github.com/layer-x/layerx-core_v2/lxserver"
	"github.com/layer-x/layerx-core_v2/lxstate"
	"github.com/layer-x/layerx-commons/lxdatabase"
"github.com/layer-x/layerx-core_v2/lxtypes"
	"github.com/layer-x/layerx-core_v2/layerx_rpi_client"
)

func PurgeState() {
	lxdatabase.Rmdir("/state", true)
}

var _ = Describe("RpiMessenger", func() {
	var serverErr error
	var state *lxstate.State

	Describe("setup", func() {
		It("sets up for the tests", func() {
			state = lxstate.NewState()
			err := state.InitializeState("http://127.0.0.1:4001")
			Expect(err).To(BeNil())
			driverErrc := make(chan error)
			coreServerWrapper := lxserver.NewLayerXCoreServerWrapper(state, lxmartini.QuietMartini(), driverErrc)

			err = state.SetTpi( "127.0.0.1:9955")
			Expect(err).To(BeNil())
			err = state.RpiPool.AddRpi(&layerx_rpi_client.RpiInfo{
				Name: "fake-rpi",
				Url: "127.0.0.1:9966",
			})

			m := coreServerWrapper.WrapServer()
			go m.RunOnAddr(fmt.Sprintf(":5566"))
			go fakes.RunFakeTpiServer("127.0.0.1:5566", 9955, driverErrc)
			go fakes.RunFakeRpiServer("127.0.0.1:5566", 9966, driverErrc)
			lxlog.ActiveDebugMode()

			go func() {
				for {
					serverErr = <-driverErrc
				}
			}()
		})
	})
	Describe("SendResurceCollectionRequest", func() {
		It("sends a POST to /collect_resources on the rpi server", func() {
			PurgeState()
			err := SendResourceCollectionRequest("127.0.0.1:9966")
			Expect(err).To(BeNil())
		})
	})
	Describe("SendLaunchTasksMessage", func() {
		It("sends a message containing a list of tasks and resources to the rpi server", func(){
			PurgeState()
			err2 := state.InitializeState("http://127.0.0.1:4001")
			Expect(err2).To(BeNil())
			err := state.SetTpi( "127.0.0.1:9955")
			Expect(err).To(BeNil())
			err = state.RpiPool.AddRpi(&layerx_rpi_client.RpiInfo{
				Name: "fake-rpi",
				Url: "127.0.0.1:9966",
			})
			fakeResource1 := lxtypes.NewResourceFromMesos(fakes.FakeOffer("fake_offer_id_1", "fake_slave_id_1"))
			fakeNode1 := lxtypes.NewNode(fakeResource1.NodeId)
			err = fakeNode1.AddResource(fakeResource1)
			Expect(err).To(BeNil())
			err = state.NodePool.AddNode(fakeNode1)
			fakeTask1 := fakes.FakeLXTask("fake_task_id_1", "fake_task_name", "fake_slave_id", "echo FAKE_COMMAND")
			fakeTaskProvider := &lxtypes.TaskProvider{
				Id:     "fake_task_provider_id_1",
				Source: "taskprovider1@tphost:port",
			}
			fakeTask1.NodeId = fakeNode1.Id
			fakeTask1.TaskProvider = fakeTaskProvider
			fakeTasks := []*lxtypes.Task{fakeTask1}
			fakeResources := []*lxtypes.Resource{fakeResource1}

			err = state.StagingTaskPool.AddTask(fakeTask1)
			Expect(err).To(BeNil())

			err = SendLaunchTasksMessage("127.0.0.1:9966", fakeTasks, fakeResources)
			Expect(err).To(BeNil())
		})
	})
	Describe("SendKillTaskRequest", func() {
		It("sends a POST to /kill_tasks/task_id on the rpi server", func() {
			err := SendKillTaskRequest("127.0.0.1:9966", "fake_task_id")
			Expect(err).To(BeNil())
		})
	})
})
