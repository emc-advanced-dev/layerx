package layerx_tpi_client_test

import (
	. "github.com/emc-advanced-dev/layerx/layerx-core/layerx_tpi_client"

	"github.com/emc-advanced-dev/layerx/layerx-core/fakes"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"
	"github.com/mesos/mesos-go/mesosproto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LayerxTpi", func() {

	fakeStatus1 := fakes.FakeTaskStatus("fake_task_id_1", mesosproto.TaskState_TASK_RUNNING)
	fakeStatus2 := fakes.FakeTaskStatus("fake_task_id_2", mesosproto.TaskState_TASK_KILLED)
	fakeStatus3 := fakes.FakeTaskStatus("fake_task_id_3", mesosproto.TaskState_TASK_FINISHED)

	fakeStatuses := []*mesosproto.TaskStatus{fakeStatus1, fakeStatus2, fakeStatus3}

	go fakes.RunFakeLayerXServer(fakeStatuses, 12345)
	lxTpi := LayerXTpi{
		CoreURL: "127.0.0.1:12345",
	}
	Describe("RegisterTpi", func() {
		It("registers the Tpi URL to the LX Server", func() {
			err := lxTpi.RegisterTpi("fake.tpi.ip:1234")
			Expect(err).To(BeNil())
		})
	})
	Describe("RegisterTaskProvider", func() {
		It("submits a new task provider to the LX Server", func() {
			taskProvider := &lxtypes.TaskProvider{
				Id:     "fake_task_provider_id",
				Source: "taskprovider@tphost:port",
			}
			err := lxTpi.RegisterTaskProvider(taskProvider)
			Expect(err).To(BeNil())
		})
	})
	Describe("DeregisterTaskProvider", func() {
		It("Requests the server to delete the task provider", func() {
			err := lxTpi.DeregisterTaskProvider("fake_task_provider_id")
			Expect(err).To(BeNil())
			err = lxTpi.DeregisterTaskProvider("fake_task_provider_id")
			Expect(err).ToNot(BeNil())
		})
	})
	Describe("GetTaskProvider(id)", func() {
		It("returns the task provider for the id, or error if it does not exist", func() {
			fakeTaskProvider := &lxtypes.TaskProvider{
				Id:     "fake_task_provider_id_1",
				Source: "taskprovider1@tphost:port",
			}
			err := lxTpi.RegisterTaskProvider(fakeTaskProvider)
			Expect(err).To(BeNil())

			taskProvider, err := lxTpi.GetTaskProvider("fake_task_provider_id_1")
			Expect(err).To(BeNil())
			Expect(taskProvider).To(Equal(fakeTaskProvider))
			taskProvider2, err := lxTpi.GetTaskProvider("fake_task_provider_id_2")
			Expect(err).NotTo(BeNil())
			Expect(taskProvider2).To(BeNil())
		})
	})
	Describe("GetTaskProviders", func() {
		It("returns a list of registered task providers", func() {
			taskProvider1 := &lxtypes.TaskProvider{
				Id:     "fake_task_provider_id_1",
				Source: "taskprovider1@tphost:port",
			}
			err := lxTpi.RegisterTaskProvider(taskProvider1)
			Expect(err).To(BeNil())

			taskProvider2 := &lxtypes.TaskProvider{
				Id:     "fake_task_provider_id_2",
				Source: "taskprovider2@tphost:port",
			}
			err = lxTpi.RegisterTaskProvider(taskProvider2)
			Expect(err).To(BeNil())

			taskProvider3 := &lxtypes.TaskProvider{
				Id:     "fake_task_provider_id_3",
				Source: "taskprovider2@tphost:port",
			}
			err = lxTpi.RegisterTaskProvider(taskProvider3)
			Expect(err).To(BeNil())

			taskProviders, err := lxTpi.GetTaskProviders()
			Expect(err).To(BeNil())
			Expect(taskProviders).To(ContainElement(taskProvider1))
			Expect(taskProviders).To(ContainElement(taskProvider2))
			Expect(taskProviders).To(ContainElement(taskProvider3))
		})
	})
	Describe("GetStatusUpdates(TPID)", func() {
		It("returns a list of status updates for the task provider", func() {
			taskProvider := &lxtypes.TaskProvider{
				Id:     "fake_task_provider_id",
				Source: "taskprovider@tphost:port",
			}
			err := lxTpi.RegisterTaskProvider(taskProvider)
			Expect(err).To(BeNil())
			fakeLxTask := fakes.FakeLXTask("fake_task_id_1", "fake_task_name", "fake_slave_id", "echo FAKE_COMMAND")
			err = lxTpi.SubmitTask("fake_task_provider_id", fakeLxTask)
			Expect(err).To(BeNil())
			fakeLxTask = fakes.FakeLXTask("fake_task_id_2", "fake_task_name", "fake_slave_id", "echo FAKE_COMMAND")
			err = lxTpi.SubmitTask("fake_task_provider_id", fakeLxTask)
			Expect(err).To(BeNil())
			fakeLxTask = fakes.FakeLXTask("fake_task_id_3", "fake_task_name", "fake_slave_id", "echo FAKE_COMMAND")
			err = lxTpi.SubmitTask("fake_task_provider_id", fakeLxTask)
			Expect(err).To(BeNil())
			Expect(err).To(BeNil())
			statuses, err := lxTpi.GetStatusUpdates("fake_task_provider_id")
			Expect(err).To(BeNil())
			Expect(statuses).To(ContainElement(fakeStatus1))
			Expect(statuses).To(ContainElement(fakeStatus2))
			Expect(statuses).To(ContainElement(fakeStatus3))
			err = lxTpi.DeregisterTaskProvider("fake_task_provider_id")
			Expect(err).To(BeNil())
		})
	})
	Describe("GetStatusUpdate(taskId)", func() {
		It("returns the current status of the given task", func() {
			taskProvider := &lxtypes.TaskProvider{
				Id:     "fake_task_provider_id",
				Source: "taskprovider@tphost:port",
			}
			err := lxTpi.RegisterTaskProvider(taskProvider)
			Expect(err).To(BeNil())
			fakeLxTask := fakes.FakeLXTask("fake_task_id_1", "fake_task_name", "fake_slave_id", "echo FAKE_COMMAND")
			err = lxTpi.SubmitTask("fake_task_provider_id", fakeLxTask)
			Expect(err).To(BeNil())
			status, err := lxTpi.GetStatusUpdate("fake_task_id_1")
			Expect(err).To(BeNil())
			Expect(status).To(Equal(fakeStatus1))
			err = lxTpi.DeregisterTaskProvider("fake_task_provider_id")
			Expect(err).To(BeNil())
		})
	})
	Describe("SubmitTask", func() {
		It("submits a task to the server", func() {
			fakeLxTask := fakes.FakeLXTask("fake_task_id", "fake_task_name", "fake_slave_id", "echo FAKE_COMMAND")
			err := lxTpi.SubmitTask("fake_task_provider_id", fakeLxTask)
			Expect(err).NotTo(BeNil())
			taskProvider := &lxtypes.TaskProvider{
				Id:     "fake_task_provider_id",
				Source: "taskprovider@tphost:port",
			}
			err = lxTpi.RegisterTaskProvider(taskProvider)
			Expect(err).To(BeNil())
			err = lxTpi.SubmitTask("fake_task_provider_id", fakeLxTask)
			Expect(err).To(BeNil())
		})
	})
	Describe("KillTask", func() {
		It("requests server to flag task with KillRequested", func() {
			err := lxTpi.KillTask("fake_task_provider_id", "fake_task_id")
			Expect(err).To(BeNil())
		})
	})
	Describe("PurgeTask", func() {
		It("requests server to flag remove the task from its database", func() {
			taskProvider := &lxtypes.TaskProvider{
				Id:     "fake_task_provider_id",
				Source: "taskprovider@tphost:port",
			}
			err := lxTpi.RegisterTaskProvider(taskProvider)
			Expect(err).To(BeNil())
			fakeLxTask := fakes.FakeLXTask("fake_task_id", "fake_task_name", "fake_slave_id", "echo FAKE_COMMAND")
			err = lxTpi.SubmitTask("fake_task_provider_id", fakeLxTask)
			Expect(err).To(BeNil())
			err = lxTpi.PurgeTask("fake_task_id")
			Expect(err).To(BeNil())
			err = lxTpi.PurgeTask("fake_task_id")
			Expect(err).ToNot(BeNil())
		})
	})

})
