package rpi_messenger_test

import (
	. "github.com/layer-x/layerx-core_v2/rpi_messenger"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/layer-x/layerx-commons/lxlog"
	"github.com/layer-x/layerx-core_v2/fakes"
	"fmt"
	"github.com/layer-x/layerx-commons/lxmartini"
	"github.com/layer-x/layerx-core_v2/driver"
	"github.com/layer-x/layerx-core_v2/lxserver"
	"github.com/layer-x/layerx-core_v2/lxstate"
	"github.com/layer-x/layerx-commons/lxactionqueue"
	"github.com/layer-x/layerx-commons/lxdatabase"
"github.com/layer-x/layerx-core_v2/lxtypes"
)

func PurgeState() {
	lxdatabase.Rmdir("/state", true)
}

var _ = Describe("RpiMessenger", func() {
	var serverErr error
	var state *lxstate.State

	Describe("setup", func() {
		It("sets up for the tests", func() {
			actionQueue := lxactionqueue.NewActionQueue()
			state = lxstate.NewState()
			err := state.InitializeState("http://127.0.0.1:4001")
			Expect(err).To(BeNil())
			coreServerWrapper := lxserver.NewLayerXCoreServerWrapper(state, actionQueue)
			driver := driver.NewLayerXDriver(actionQueue)

			driverErrc := make(chan error)

			m := coreServerWrapper.WrapServer(lxmartini.QuietMartini(), driverErrc)
			go m.RunOnAddr(fmt.Sprintf(":5566"))
			go fakes.RunFakeRpiServer("127.0.0.1:5566", 9966, driverErrc)
			go driver.Run()
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
			err := SendResourceCollectionRequest("127.0.0.1:9966")
			Expect(err).To(BeNil())
		})
	})
	Describe("SendLaunchTasksMessage", func() {
		It("sends a message containing a list of tasks and resources to the rpi server", func(){
			PurgeState()
			err2 := state.InitializeState("http://127.0.0.1:4001")
			Expect(err2).To(BeNil())
			fakeTask1 := fakes.FakeLXTask("fake_task_id_1", "fake_task_name", "fake_slave_id", "echo FAKE_COMMAND")
			fakeResoure1 := lxtypes.NewResourceFromMesos(fakes.FakeOffer("fake_offer_id_1", "fake_slave_id_1"))
			fakeTasks := []*lxtypes.Task{fakeTask1}
			fakeResources := []*lxtypes.Resource{fakeResoure1}
			err := SendLaunchTasksMessage("127.0.0.1:9966", fakeTasks, fakeResources)
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
