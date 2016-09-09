package health_checker_test

import (
	. "github.com/emc-advanced-dev/layerx/layerx-core/health_checker"

	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/emc-advanced-dev/layerx/layerx-core/fakes"
	"github.com/emc-advanced-dev/layerx/layerx-core/layerx_brain_client"
	"github.com/emc-advanced-dev/layerx/layerx-core/layerx_rpi_client"
	"github.com/emc-advanced-dev/layerx/layerx-core/layerx_tpi_client"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxserver"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxstate"
	"github.com/layer-x/layerx-commons/lxdatabase"
	"github.com/layer-x/layerx-commons/lxmartini"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

func PurgeState() {
	lxdatabase.Rmdir("/state", true)
}

var _ = Describe("TaskProviderHealthChecker", func() {
	var lxRpiClient *layerx_rpi_client.LayerXRpi
	var lxTpiClient *layerx_tpi_client.LayerXTpi
	var lxBrainClient *layerx_brain_client.LayerXBrainClient
	var state *lxstate.State
	var serverErr error

	Describe("setup", func() {
		It("sets up for the tests", func() {
			lxRpiClient = &layerx_rpi_client.LayerXRpi{
				CoreURL: "127.0.0.1:2299",
			}
			lxTpiClient = &layerx_tpi_client.LayerXTpi{
				CoreURL: "127.0.0.1:2299",
			}
			lxBrainClient = &layerx_brain_client.LayerXBrainClient{
				CoreURL: "127.0.0.1:2299",
			}

			state = lxstate.NewState()
			err := state.InitializeState("http://127.0.0.1:4001")
			Expect(err).To(BeNil())
			driverErrc := make(chan error)
			coreServerWrapper := lxserver.NewLayerXCoreServerWrapper(state, lxmartini.QuietMartini(), driverErrc)

			err = state.SetTpi("127.0.0.1:3388")
			Expect(err).To(BeNil())
			err = state.RpiPool.AddRpi(&layerx_rpi_client.RpiInfo{
				Name: "fake-rpi",
				Url:  "127.0.0.1:3399",
			})

			go func() {
				for {
					serverErr = <-driverErrc
				}
			}()

			m := coreServerWrapper.WrapServer()
			go m.RunOnAddr(fmt.Sprintf(":2299"))
			go fakes.RunFakeTpiServer("127.0.0.1:2299", 3388, make(chan error))
			go fakes.RunFakeRpiServer("127.0.0.1:2299", 3399, make(chan error))
			logrus.SetLevel(logrus.DebugLevel)
		})
	})
	Describe("CheckTaskProviderHealth", func() {
		Context("a connected task provider is healthy", func() {
			It("does no-op", func() {
				PurgeState()
				err2 := state.InitializeState("http://127.0.0.1:4001")
				Expect(err2).To(BeNil())
				err := state.SetTpi("127.0.0.1:3388")
				Expect(err).To(BeNil())
				err = state.RpiPool.AddRpi(&layerx_rpi_client.RpiInfo{
					Name: "fake-rpi",
					Url:  "127.0.0.1:3399",
				})
				fakeTaskProvider1 := fakes.FakeTaskProvider("fake_framework_1", "ff@fakeip1:fakeport")
				err = state.TaskProviderPool.AddTaskProvider(fakeTaskProvider1)
				Expect(err).To(BeNil())
				fakePendingTask := fakes.FakeLXTask("fake__pending_task_id", "fake_task_name", "fake_slave_id", "echo FAKE_COMMAND")
				fakePendingTask.TaskProvider = fakeTaskProvider1
				err = state.PendingTaskPool.AddTask(fakePendingTask)
				Expect(err).To(BeNil())
				fakeStagingTask := fakes.FakeLXTask("fake_staging_task_id", "fake_task_name", "fake_slave_id", "echo FAKE_COMMAND")
				fakeStagingTask.TaskProvider = fakeTaskProvider1
				err = state.StagingTaskPool.AddTask(fakeStagingTask)
				Expect(err).To(BeNil())
				fakeNode := fakes.FakeNode("fake_resource_id_1", "fake_node_id_1")
				err = state.NodePool.AddNode(fakeNode)
				Expect(err).To(BeNil())
				nodeTaskPool, err := state.NodePool.GetNodeTaskPool(fakeNode.Id)
				Expect(err).To(BeNil())
				fakeNodeTask1 := fakes.FakeLXTask("fake__node_task_id_1", "fake_task", "fake_node_id_1", "echo FAKECOMMAND")
				fakeNodeTask1.TaskProvider = fakeTaskProvider1
				err = nodeTaskPool.AddTask(fakeNodeTask1)
				Expect(err).To(BeNil())
				healthChecker := NewHealthChecker(state)
				err = healthChecker.FailDisconnectedTaskProviders()
				Expect(err).To(BeNil())
				taskProviders, err := state.TaskProviderPool.GetTaskProviders()
				Expect(err).To(BeNil())
				Expect(taskProviders).To(ContainElement(fakeTaskProvider1))
				pendingTasks, err := state.PendingTaskPool.GetTasks()
				Expect(err).To(BeNil())
				Expect(pendingTasks).To(ContainElement(fakePendingTask))
				stagingTasks, err := state.StagingTaskPool.GetTasks()
				Expect(err).To(BeNil())
				Expect(stagingTasks).To(ContainElement(fakeStagingTask))
				nodeTasks, err := nodeTaskPool.GetTasks()
				Expect(err).To(BeNil())
				Expect(nodeTasks).To(ContainElement(fakeNodeTask1))
			})
		})
		Context("a connected task provider is not responding", func() {
			Context("the task provider does not support failover", func() {
				It("deletes the task provider from the task provider pool and kills all of its running tasks", func() {
					PurgeState()
					err2 := state.InitializeState("http://127.0.0.1:4001")
					Expect(err2).To(BeNil())
					err := state.SetTpi("127.0.0.1:3388")
					Expect(err).To(BeNil())
					err = state.RpiPool.AddRpi(&layerx_rpi_client.RpiInfo{
						Name: "fake-rpi",
						Url:  "127.0.0.1:3399",
					})
					//triggers failure
					fakeTaskProvider1 := fakes.FakeTaskProvider("fake_framework_1"+fakes.FAIL_ON_PURPOSE, "ff@fakeip1:fakeport")
					err = state.TaskProviderPool.AddTaskProvider(fakeTaskProvider1)
					Expect(err).To(BeNil())
					fakePendingTask := fakes.FakeLXTask("fake__pending_task_id", "fake_task_name", "fake_slave_id", "echo FAKE_COMMAND")
					fakePendingTask.TaskProvider = fakeTaskProvider1
					err = state.PendingTaskPool.AddTask(fakePendingTask)
					Expect(err).To(BeNil())
					fakeStagingTask := fakes.FakeLXTask("fake_staging_task_id", "fake_task_name", "fake_slave_id", "echo FAKE_COMMAND")
					fakeStagingTask.TaskProvider = fakeTaskProvider1
					err = state.StagingTaskPool.AddTask(fakeStagingTask)
					Expect(err).To(BeNil())
					fakeNode := fakes.FakeNode("fake_resource_id_1", "fake_node_id_1")
					err = state.NodePool.AddNode(fakeNode)
					Expect(err).To(BeNil())
					nodeTaskPool, err := state.NodePool.GetNodeTaskPool(fakeNode.Id)
					Expect(err).To(BeNil())
					fakeNodeTask1 := fakes.FakeLXTask("fake__node_task_id_1", "fake_task", "fake_node_id_1", "echo FAKECOMMAND")
					fakeNodeTask1.TaskProvider = fakeTaskProvider1
					err = nodeTaskPool.AddTask(fakeNodeTask1)
					Expect(err).To(BeNil())
					healthChecker := NewHealthChecker(state)
					err = healthChecker.FailDisconnectedTaskProviders()
					Expect(err).To(BeNil())
					taskProviders, err := state.TaskProviderPool.GetTaskProviders()
					Expect(err).To(BeNil())
					Expect(taskProviders).NotTo(ContainElement(fakeTaskProvider1))
					failedTaskProviders, err := state.FailedTaskProviderPool.GetTaskProviders()
					Expect(err).To(BeNil())
					Expect(failedTaskProviders).NotTo(ContainElement(fakeTaskProvider1))
					pendingTasks, err := state.PendingTaskPool.GetTasks()
					Expect(err).To(BeNil())
					Expect(pendingTasks).NotTo(ContainElement(fakePendingTask))
					stagingTasks, err := state.StagingTaskPool.GetTasks()
					Expect(err).To(BeNil())
					Expect(stagingTasks).NotTo(ContainElement(fakeStagingTask))
					nodeTasks, err := nodeTaskPool.GetTasks()
					Expect(err).To(BeNil())
					Expect(nodeTasks).NotTo(ContainElement(fakeNodeTask1))
				})
			})
			Context("the task provider supports failover", func() {
				It("moves the task provider to the failed pool and marks the time of failure but leaves the tasks running", func() {
					PurgeState()
					err2 := state.InitializeState("http://127.0.0.1:4001")
					Expect(err2).To(BeNil())
					err := state.SetTpi("127.0.0.1:3388")
					Expect(err).To(BeNil())
					err = state.RpiPool.AddRpi(&layerx_rpi_client.RpiInfo{
						Name: "fake-rpi",
						Url:  "127.0.0.1:3399",
					})
					//triggers failure
					fakeTaskProvider1 := fakes.FakeTaskProvider("fake_framework_1"+fakes.FAIL_ON_PURPOSE, "ff@fakeip1:fakeport")
					fakeTaskProvider1.FailoverTimeout = 1
					err = state.TaskProviderPool.AddTaskProvider(fakeTaskProvider1)
					Expect(err).To(BeNil())
					//for expectations, set failover timeout here
					fakeTaskProvider1.TimeFailed = float64(time.Now().Unix())
					fakePendingTask := fakes.FakeLXTask("fake__pending_task_id", "fake_task_name", "fake_slave_id", "echo FAKE_COMMAND")
					fakePendingTask.TaskProvider = fakeTaskProvider1
					err = state.PendingTaskPool.AddTask(fakePendingTask)
					Expect(err).To(BeNil())
					fakeStagingTask := fakes.FakeLXTask("fake_staging_task_id", "fake_task_name", "fake_slave_id", "echo FAKE_COMMAND")
					fakeStagingTask.TaskProvider = fakeTaskProvider1
					err = state.StagingTaskPool.AddTask(fakeStagingTask)
					Expect(err).To(BeNil())
					fakeNode := fakes.FakeNode("fake_resource_id_1", "fake_node_id_1")
					err = state.NodePool.AddNode(fakeNode)
					Expect(err).To(BeNil())
					nodeTaskPool, err := state.NodePool.GetNodeTaskPool(fakeNode.Id)
					Expect(err).To(BeNil())
					fakeNodeTask1 := fakes.FakeLXTask("fake__node_task_id_1", "fake_task", "fake_node_id_1", "echo FAKECOMMAND")
					fakeNodeTask1.TaskProvider = fakeTaskProvider1
					err = nodeTaskPool.AddTask(fakeNodeTask1)
					Expect(err).To(BeNil())
					healthChecker := NewHealthChecker(state)
					err = healthChecker.FailDisconnectedTaskProviders()
					Expect(err).To(BeNil())
					taskProviders, err := state.TaskProviderPool.GetTaskProviders()
					Expect(err).To(BeNil())
					Expect(taskProviders).NotTo(ContainElement(fakeTaskProvider1))
					failedTaskProviders, err := state.FailedTaskProviderPool.GetTaskProviders()
					Expect(err).To(BeNil())
					Expect(failedTaskProviders).To(ContainElement(fakeTaskProvider1))
					pendingTasks, err := state.PendingTaskPool.GetTasks()
					Expect(err).To(BeNil())
					Expect(pendingTasks).To(ContainElement(fakePendingTask))
					stagingTasks, err := state.StagingTaskPool.GetTasks()
					Expect(err).To(BeNil())
					Expect(stagingTasks).To(ContainElement(fakeStagingTask))
					nodeTasks, err := nodeTaskPool.GetTasks()
					Expect(err).To(BeNil())
					Expect(nodeTasks).To(ContainElement(fakeNodeTask1))
				})
			})
		})
	})
	Describe("ExpireTimedOutTaskProviders", func() {
		It("deletes all failed over task providers (and their tasks) whose failover timeout has expired", func() {
			PurgeState()
			err2 := state.InitializeState("http://127.0.0.1:4001")
			Expect(err2).To(BeNil())
			err := state.SetTpi("127.0.0.1:3388")
			Expect(err).To(BeNil())
			err = state.RpiPool.AddRpi(&layerx_rpi_client.RpiInfo{
				Name: "fake-rpi",
				Url:  "127.0.0.1:3399",
			})
			//triggers failure
			fakeTaskProvider1 := fakes.FakeTaskProvider("fake_framework_1"+fakes.FAIL_ON_PURPOSE, "ff@fakeip1:fakeport")
			fakeTaskProvider1.FailoverTimeout = 1
			err = state.TaskProviderPool.AddTaskProvider(fakeTaskProvider1)
			Expect(err).To(BeNil())
			//for expectations, set failover timeout here
			fakeTaskProvider1.TimeFailed = float64(time.Now().Unix())
			fakePendingTask := fakes.FakeLXTask("fake__pending_task_id", "fake_task_name", "fake_slave_id", "echo FAKE_COMMAND")
			fakePendingTask.TaskProvider = fakeTaskProvider1
			err = state.PendingTaskPool.AddTask(fakePendingTask)
			Expect(err).To(BeNil())
			fakeStagingTask := fakes.FakeLXTask("fake_staging_task_id", "fake_task_name", "fake_slave_id", "echo FAKE_COMMAND")
			fakeStagingTask.TaskProvider = fakeTaskProvider1
			err = state.StagingTaskPool.AddTask(fakeStagingTask)
			Expect(err).To(BeNil())
			fakeNode := fakes.FakeNode("fake_resource_id_1", "fake_node_id_1")
			err = state.NodePool.AddNode(fakeNode)
			Expect(err).To(BeNil())
			nodeTaskPool, err := state.NodePool.GetNodeTaskPool(fakeNode.Id)
			Expect(err).To(BeNil())
			fakeNodeTask1 := fakes.FakeLXTask("fake__node_task_id_1", "fake_task", "fake_node_id_1", "echo FAKECOMMAND")
			fakeNodeTask1.TaskProvider = fakeTaskProvider1
			err = nodeTaskPool.AddTask(fakeNodeTask1)
			Expect(err).To(BeNil())
			healthChecker := NewHealthChecker(state)
			err = healthChecker.FailDisconnectedTaskProviders()
			Expect(err).To(BeNil())
			taskProviders, err := state.TaskProviderPool.GetTaskProviders()
			Expect(err).To(BeNil())
			Expect(taskProviders).NotTo(ContainElement(fakeTaskProvider1))
			failedTaskProviders, err := state.FailedTaskProviderPool.GetTaskProviders()
			Expect(err).To(BeNil())
			Expect(failedTaskProviders).To(ContainElement(fakeTaskProvider1))
			pendingTasks, err := state.PendingTaskPool.GetTasks()
			Expect(err).To(BeNil())
			Expect(pendingTasks).To(ContainElement(fakePendingTask))
			stagingTasks, err := state.StagingTaskPool.GetTasks()
			Expect(err).To(BeNil())
			Expect(stagingTasks).To(ContainElement(fakeStagingTask))
			nodeTasks, err := nodeTaskPool.GetTasks()
			Expect(err).To(BeNil())
			Expect(nodeTasks).To(ContainElement(fakeNodeTask1))
			time.Sleep(2000 * time.Millisecond)
			err = healthChecker.ExpireTimedOutTaskProviders()
			Expect(err).To(BeNil())
			failedTaskProvidersAfterTimeout, err := state.FailedTaskProviderPool.GetTaskProviders()
			Expect(err).To(BeNil())
			Expect(failedTaskProvidersAfterTimeout).NotTo(ContainElement(fakeTaskProvider1))
			pendingTasksAfterTimeout, err := state.PendingTaskPool.GetTasks()
			Expect(err).To(BeNil())
			Expect(pendingTasksAfterTimeout).NotTo(ContainElement(fakePendingTask))
			stagingTasksAfterTimeout, err := state.StagingTaskPool.GetTasks()
			Expect(err).To(BeNil())
			Expect(stagingTasksAfterTimeout).NotTo(ContainElement(fakeStagingTask))
			nodeTasksAfterTimeout, err := nodeTaskPool.GetTasks()
			Expect(err).To(BeNil())
			Expect(nodeTasksAfterTimeout).NotTo(ContainElement(fakeNodeTask1))
		})
	})
})
