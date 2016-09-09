package lxstate_test

import (
	. "github.com/emc-advanced-dev/layerx/layerx-core/lxstate"

	"encoding/json"
	"github.com/emc-advanced-dev/layerx/layerx-core/fakes"
	"github.com/layer-x/layerx-commons/lxdatabase"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TaskPool", func() {
	Describe("GetTask(taskId)", func() {
		It("returns the task if it exists, else returns err", func() {
			state := NewState()
			state.InitializeState("http://127.0.0.1:4001")
			PurgeState()
			err := state.InitializeState("http://127.0.0.1:4001")
			Expect(err).To(BeNil())
			pendingTasks := state.PendingTaskPool
			fakeTask := fakes.FakeLXTask("fake_task_id_1", "fake_task_name", "fake_slave_id", "echo FAKE_COMMAND")
			fakeTask.TaskProvider = fakes.FakeTaskProvider("fake_task_provider_id", "tp@fakeip:fakeport")
			task, err := pendingTasks.GetTask(fakeTask.TaskId)
			Expect(err).NotTo(BeNil())
			Expect(task).To(BeNil())
			err = pendingTasks.AddTask(fakeTask)
			Expect(err).To(BeNil())
			task, err = pendingTasks.GetTask(fakeTask.TaskId)
			Expect(err).To(BeNil())
			Expect(task).To(Equal(fakeTask))
		})
	})
	Describe("AddTask", func() {
		Context("the task is new", func() {
			It("adds the task to etcd state", func() {
				state := NewState()
				state.InitializeState("http://127.0.0.1:4001")
				PurgeState()
				err := state.InitializeState("http://127.0.0.1:4001")
				Expect(err).To(BeNil())
				pendingTasks := state.PendingTaskPool
				fakeTask := fakes.FakeLXTask("fake_task_id_1", "fake_task_name", "fake_slave_id", "echo FAKE_COMMAND")
				fakeTask.TaskProvider = fakes.FakeTaskProvider("fake_task_provider_id", "tp@fakeip:fakeport")
				err = pendingTasks.AddTask(fakeTask)
				Expect(err).To(BeNil())
				expectedTaskJsonBytes, err := json.Marshal(fakeTask)
				Expect(err).To(BeNil())
				expectedTaskJson := string(expectedTaskJsonBytes)
				actualTaskJson, err := lxdatabase.Get(state.PendingTaskPool.GetKey() + "/" + fakeTask.TaskId)
				Expect(err).To(BeNil())
				Expect(actualTaskJson).To(Equal(expectedTaskJson))
			})
		})
		Context("the task is not new", func() {
			It("returns an error", func() {
				state := NewState()
				state.InitializeState("http://127.0.0.1:4001")
				PurgeState()
				err := state.InitializeState("http://127.0.0.1:4001")
				Expect(err).To(BeNil())
				pendingTasks := state.PendingTaskPool
				fakeTask := fakes.FakeLXTask("fake_task_id_1", "fake_task_name", "fake_slave_id", "echo FAKE_COMMAND")
				fakeTask.TaskProvider = fakes.FakeTaskProvider("fake_task_provider_id", "tp@fakeip:fakeport")
				err = pendingTasks.AddTask(fakeTask)
				Expect(err).To(BeNil())
				err = pendingTasks.AddTask(fakeTask)
				Expect(err).NotTo(BeNil())
			})
		})
	})
	Describe("ModifyTask", func() {
		Context("the exists", func() {
			It("modifies the task", func() {
				state := NewState()
				state.InitializeState("http://127.0.0.1:4001")
				PurgeState()
				err := state.InitializeState("http://127.0.0.1:4001")
				Expect(err).To(BeNil())
				pendingTasks := state.PendingTaskPool
				fakeTask := fakes.FakeLXTask("fake_task_id_1", "fake_task_name", "fake_slave_id", "echo FAKE_COMMAND")
				fakeTask.TaskProvider = fakes.FakeTaskProvider("fake_task_provider_id", "tp@fakeip:fakeport")
				err = pendingTasks.AddTask(fakeTask)
				Expect(err).To(BeNil())
				fakeTask.Mem = 666
				fakeTask.Cpus = 666
				fakeTask.Disk = 666
				err = pendingTasks.ModifyTask(fakeTask.TaskId, fakeTask)
				Expect(err).To(BeNil())
				expectedTaskJsonBytes, err := json.Marshal(fakeTask)
				Expect(err).To(BeNil())
				expectedTaskJson := string(expectedTaskJsonBytes)
				actualTaskJson, err := lxdatabase.Get(state.PendingTaskPool.GetKey() + "/" + fakeTask.TaskId)
				Expect(err).To(BeNil())
				Expect(actualTaskJson).To(Equal(expectedTaskJson))
			})
		})
		Context("the task doest exist", func() {
			It("returns an error", func() {
				state := NewState()
				state.InitializeState("http://127.0.0.1:4001")
				PurgeState()
				err := state.InitializeState("http://127.0.0.1:4001")
				Expect(err).To(BeNil())
				pendingTasks := state.PendingTaskPool
				fakeTask := fakes.FakeLXTask("fake_task_id_1", "fake_task_name", "fake_slave_id", "echo FAKE_COMMAND")
				fakeTask.TaskProvider = fakes.FakeTaskProvider("fake_task_provider_id", "tp@fakeip:fakeport")
				err = pendingTasks.ModifyTask(fakeTask.TaskId, fakeTask)
				Expect(err).NotTo(BeNil())
			})
		})
	})
	Describe("GetTasks()", func() {
		It("returns all known tasks in the pool", func() {
			state := NewState()
			state.InitializeState("http://127.0.0.1:4001")
			PurgeState()
			err := state.InitializeState("http://127.0.0.1:4001")
			Expect(err).To(BeNil())
			pendingTasks := state.PendingTaskPool
			fakeTask1 := fakes.FakeLXTask("fake_task_id_1", "fake_task_name", "fake_slave_id", "echo FAKE_COMMAND")
			fakeTask1.TaskProvider = fakes.FakeTaskProvider("fake_task_provider_id", "tp@fakeip:fakeport")
			fakeTask2 := fakes.FakeLXTask("fake_task_id_2", "fake_task_name", "fake_slave_id", "echo FAKE_COMMAND")
			fakeTask2.TaskProvider = fakes.FakeTaskProvider("fake_task_provider_id", "tp@fakeip:fakeport")
			fakeTask3 := fakes.FakeLXTask("fake_task_id_3", "fake_task_name", "fake_slave_id", "echo FAKE_COMMAND")
			fakeTask3.TaskProvider = fakes.FakeTaskProvider("fake_task_provider_id", "tp@fakeip:fakeport")
			err = pendingTasks.AddTask(fakeTask1)
			Expect(err).To(BeNil())
			err = pendingTasks.AddTask(fakeTask2)
			Expect(err).To(BeNil())
			err = pendingTasks.AddTask(fakeTask3)
			Expect(err).To(BeNil())
			tasks, err := pendingTasks.GetTasks()
			Expect(err).To(BeNil())
			Expect(tasks[fakeTask1.TaskId]).To(Equal(fakeTask1))
			Expect(tasks[fakeTask2.TaskId]).To(Equal(fakeTask2))
			Expect(tasks[fakeTask3.TaskId]).To(Equal(fakeTask3))
		})
	})
	Describe("DeleteTask(taskId)", func() {
		Context("task exists", func() {
			It("deletes the task", func() {
				state := NewState()
				state.InitializeState("http://127.0.0.1:4001")
				PurgeState()
				err := state.InitializeState("http://127.0.0.1:4001")
				Expect(err).To(BeNil())
				pendingTasks := state.PendingTaskPool
				fakeTask1 := fakes.FakeLXTask("fake_task_id_1", "fake_task_name", "fake_slave_id", "echo FAKE_COMMAND")
				fakeTask1.TaskProvider = fakes.FakeTaskProvider("fake_task_provider_id", "tp@fakeip:fakeport")
				fakeTask2 := fakes.FakeLXTask("fake_task_id_2", "fake_task_name", "fake_slave_id", "echo FAKE_COMMAND")
				fakeTask2.TaskProvider = fakes.FakeTaskProvider("fake_task_provider_id", "tp@fakeip:fakeport")
				fakeTask3 := fakes.FakeLXTask("fake_task_id_3", "fake_task_name", "fake_slave_id", "echo FAKE_COMMAND")
				fakeTask3.TaskProvider = fakes.FakeTaskProvider("fake_task_provider_id", "tp@fakeip:fakeport")
				err = pendingTasks.AddTask(fakeTask1)
				Expect(err).To(BeNil())
				err = pendingTasks.AddTask(fakeTask2)
				Expect(err).To(BeNil())
				err = pendingTasks.AddTask(fakeTask3)
				Expect(err).To(BeNil())
				err = pendingTasks.DeleteTask(fakeTask1.TaskId)
				Expect(err).To(BeNil())
				tasks, err := pendingTasks.GetTasks()
				Expect(err).To(BeNil())
				Expect(tasks[fakeTask1.TaskId]).To(BeNil())
				Expect(tasks[fakeTask2.TaskId]).To(Equal(fakeTask2))
				Expect(tasks[fakeTask3.TaskId]).To(Equal(fakeTask3))
			})
		})
		Context("task does not exist", func() {
			It("throws error", func() {
				state := NewState()
				state.InitializeState("http://127.0.0.1:4001")
				PurgeState()
				err := state.InitializeState("http://127.0.0.1:4001")
				Expect(err).To(BeNil())
				pendingTasks := state.PendingTaskPool
				err = pendingTasks.DeleteTask("nonexistent_task_id")
				Expect(err).NotTo(BeNil())
			})
		})
	})
})
