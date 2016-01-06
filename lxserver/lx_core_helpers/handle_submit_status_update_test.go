package lx_core_helpers_test
import (
	. "github.com/layer-x/layerx-core_v2/lxserver/lx_core_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/layer-x/layerx-commons/lxdatabase"
	"github.com/layer-x/layerx-core_v2/lxstate"
	"github.com/layer-x/layerx-commons/lxactionqueue"
"github.com/layer-x/layerx-core_v2/driver"
"github.com/layer-x/layerx-commons/lxlog"
"github.com/layer-x/layerx-core_v2/fakes"
	"fmt"
	"github.com/layer-x/layerx-commons/lxmartini"
	"github.com/layer-x/layerx-core_v2/lxserver"
	"github.com/mesos/mesos-go/mesosproto"
	"github.com/layer-x/layerx-core_v2/lxtypes"
	"github.com/layer-x/layerx-core_v2/layerx_rpi_client"
)

func PurgeState() {
	lxdatabase.Rmdir("/state", true)
}

var _ = Describe("HandleSubmitStatusUpdate", func() {
	var state *lxstate.State
	lxRpiClient := &layerx_rpi_client.LayerXRpi{
		CoreURL: "127.0.0.1:5675",
	}

	Describe("setup", func() {
		It("sets up for the tests", func() {
			actionQueue := lxactionqueue.NewActionQueue()
			state = lxstate.NewState()
			err := state.InitializeState("http://127.0.0.1:4001")
			Expect(err).To(BeNil())
			coreServerWrapper := lxserver.NewLayerXCoreServerWrapper(state, actionQueue, lxmartini.QuietMartini(), "127.0.0.1:5599", "127.0.0.1:4499", make(chan error))
			driver := driver.NewLayerXDriver(actionQueue)

			m := coreServerWrapper.WrapServer()
			go m.RunOnAddr(fmt.Sprintf(":5675"))
			go fakes.RunFakeTpiServer("127.0.0.1:5675", 5599, make(chan error))
			go fakes.RunFakeRpiServer("127.0.0.1:5675", 4499, make(chan error))
			go driver.Run()
			lxlog.ActiveDebugMode()
		})
	})



	Describe("ProcessStatusUpdate", func(){
		Context("status update is terminal", func(){
			It("deletes the task from the state", func(){
				PurgeState()
				purgeErr := state.InitializeState("http://127.0.0.1:4001")
				Expect(purgeErr).To(BeNil())
				fakeResource1 := lxtypes.NewResourceFromMesos(fakes.FakeOffer("fake_offer_id_1", "fake_node_id_1"))
				err := lxRpiClient.SubmitResource(fakeResource1)
				Expect(err).To(BeNil())
				nodeTaskPool, err := state.NodePool.GetNodeTaskPool(fakeResource1.NodeId)
				Expect(err).To(BeNil())
				fakeTask1 := fakes.FakeLXTask("fake_task_id_1", "fake_task1", "fake_node_id_1", "echo FAKECOMMAND")
				fakeTask1.TaskProvider = fakes.FakeTaskProvider("fake_task_provider_id", "tp@fakeip:fakeport")
				err = nodeTaskPool.AddTask(fakeTask1)
				Expect(err).To(BeNil())
				fakeStatus1 := fakes.FakeTaskStatus("fake_task_id_1", mesosproto.TaskState_TASK_KILLED)
				_, err = nodeTaskPool.GetTask(fakeTask1.TaskId)
				Expect(err).To(BeNil())
				err = ProcessStatusUpdate(state, "127.0.0.1:5599", fakeStatus1)
				Expect(err).To(BeNil())
				status, err := state.StatusPool.GetStatus(fakeStatus1.GetTaskId().GetValue())
				Expect(err).To(BeNil())
				Expect(status).To(Equal(fakeStatus1))
				_, err = nodeTaskPool.GetTask(fakeTask1.TaskId)
				Expect(err).NotTo(BeNil())
			})
		})
		Context("status update is running, and the task was staging", func(){
			It("moves the task from the staging pool to the node pool", func(){
				PurgeState()
				purgeErr := state.InitializeState("http://127.0.0.1:4001")
				Expect(purgeErr).To(BeNil())
				fakeResource1 := lxtypes.NewResourceFromMesos(fakes.FakeOffer("fake_offer_id_1", "fake_node_id_1"))
				err := lxRpiClient.SubmitResource(fakeResource1)
				Expect(err).To(BeNil())
				nodeTaskPool, err := state.NodePool.GetNodeTaskPool(fakeResource1.NodeId)
				Expect(err).To(BeNil())
				fakeTask1 := fakes.FakeLXTask("fake_task_id_1", "fake_task1", "fake_node_id_1", "echo FAKECOMMAND")
				fakeTask1.TaskProvider = fakes.FakeTaskProvider("fake_task_provider_id", "tp@fakeip:fakeport")
				err = state.StagingTaskPool.AddTask(fakeTask1)
				Expect(err).To(BeNil())
				fakeStatus1 := fakes.FakeTaskStatus("fake_task_id_1", mesosproto.TaskState_TASK_RUNNING)
				_, err = nodeTaskPool.GetTask(fakeTask1.TaskId)
				Expect(err).NotTo(BeNil())
				err = ProcessStatusUpdate(state, "127.0.0.1:5599", fakeStatus1)
				Expect(err).To(BeNil())
				status, err := state.StatusPool.GetStatus(fakeStatus1.GetTaskId().GetValue())
				Expect(err).To(BeNil())
				Expect(status).To(Equal(fakeStatus1))
				_, err = nodeTaskPool.GetTask(fakeTask1.TaskId)
				Expect(err).To(BeNil())
				_, err = state.StagingTaskPool.GetTask(fakeTask1.TaskId)
				Expect(err).NotTo(BeNil())
			})
		})
	})
})
