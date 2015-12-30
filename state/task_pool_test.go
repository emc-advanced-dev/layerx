package state_test

import (
	. "github.com/layer-x/layerx-core_v2/state"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/layer-x/layerx-mesos-rpi_v2/fakes"
	"github.com/layer-x/layerx-commons/lxdatabase"
	"encoding/json"
)

var _ = Describe("TaskPool", func() {
	Describe("GetTask(taskId)", func(){
		It("returns the task if it exists, else returns err", func(){
			state := NewState()
			state.InitializeState("http://127.0.0.1:4001")
			PurgeState()
			err := state.InitializeState("http://127.0.0.1:4001")
			Expect(err).To(BeNil())
			pendingTasks := state.PendingTaskPool
			fakeTask := fakes.FakeTask("fake_task_id_1")
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
	Describe("AddTask", func(){
		It("adds the task to etcd state", func(){
			state := NewState()
			state.InitializeState("http://127.0.0.1:4001")
			PurgeState()
			err := state.InitializeState("http://127.0.0.1:4001")
			Expect(err).To(BeNil())
			pendingTasks := state.PendingTaskPool
			fakeTask := fakes.FakeTask("fake_task_id_1")
			err = pendingTasks.AddTask(fakeTask)
			Expect(err).To(BeNil())
			expectedTaskJsonBytes, err := json.Marshal(fakeTask)
			Expect(err).To(BeNil())
			expectedTaskJson := string(expectedTaskJsonBytes)
			actualTaskJson, err := lxdatabase.Get(state.PendingTaskPool.GetKey() + "/"+fakeTask.TaskId)
			Expect(err).To(BeNil())
			Expect(actualTaskJson).To(Equal(expectedTaskJson))
		})
	})
})
