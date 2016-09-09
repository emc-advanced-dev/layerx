package lxstate_test

import (
	. "github.com/emc-advanced-dev/layerx/layerx-core/lxstate"

	"encoding/json"
	"github.com/emc-advanced-dev/layerx/layerx-core/fakes"
	"github.com/layer-x/layerx-commons/lxdatabase"
	"github.com/mesos/mesos-go/mesosproto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StatusPool", func() {
	Describe("GetStatus(statusId)", func() {
		It("returns the status if it exists, else returns err", func() {
			state := NewState()
			state.InitializeState("http://127.0.0.1:4001")
			PurgeState()
			err := state.InitializeState("http://127.0.0.1:4001")
			Expect(err).To(BeNil())
			statusPool := state.StatusPool
			fakeStatus := fakes.FakeTaskStatus("fake_task_id_1", mesosproto.TaskState_TASK_RUNNING)
			status, err := statusPool.GetStatus(fakeStatus.GetTaskId().GetValue())
			Expect(err).NotTo(BeNil())
			Expect(status).To(BeNil())
			err = statusPool.AddStatus(fakeStatus)
			Expect(err).To(BeNil())
			status, err = statusPool.GetStatus(fakeStatus.GetTaskId().GetValue())
			Expect(err).To(BeNil())
			Expect(status).To(Equal(fakeStatus))
		})
	})
	Describe("AddStatus", func() {
		Context("the status is new", func() {
			It("adds the status to etcd state", func() {
				state := NewState()
				state.InitializeState("http://127.0.0.1:4001")
				PurgeState()
				err := state.InitializeState("http://127.0.0.1:4001")
				Expect(err).To(BeNil())
				statusPool := state.StatusPool
				fakeStatus := fakes.FakeTaskStatus("fake_task_id_1", mesosproto.TaskState_TASK_RUNNING)
				err = statusPool.AddStatus(fakeStatus)
				Expect(err).To(BeNil())
				expectedStatusJsonBytes, err := json.Marshal(fakeStatus)
				Expect(err).To(BeNil())
				expectedStatusJson := string(expectedStatusJsonBytes)
				actualStatusJson, err := lxdatabase.Get(state.StatusPool.GetKey() + "/" + fakeStatus.GetTaskId().GetValue())
				Expect(err).To(BeNil())
				Expect(actualStatusJson).To(Equal(expectedStatusJson))
			})
		})
		Context("the status is not new", func() {
			It("modifies the original status", func() {
				state := NewState()
				state.InitializeState("http://127.0.0.1:4001")
				PurgeState()
				err := state.InitializeState("http://127.0.0.1:4001")
				Expect(err).To(BeNil())
				statusPool := state.StatusPool
				fakeStatus := fakes.FakeTaskStatus("fake_task_id_1", mesosproto.TaskState_TASK_RUNNING)
				err = statusPool.AddStatus(fakeStatus)
				Expect(err).To(BeNil())
				fakeStatus = fakes.FakeTaskStatus("fake_task_id_1", mesosproto.TaskState_TASK_FINISHED)
				err = statusPool.AddStatus(fakeStatus)
				Expect(err).To(BeNil())
				expectedStatusJsonBytes, err := json.Marshal(fakeStatus)
				Expect(err).To(BeNil())
				expectedStatusJson := string(expectedStatusJsonBytes)
				actualStatusJson, err := lxdatabase.Get(state.StatusPool.GetKey() + "/" + fakeStatus.GetTaskId().GetValue())
				Expect(err).To(BeNil())
				Expect(actualStatusJson).To(Equal(expectedStatusJson))
			})
		})
	})
	Describe("ModifyStatus", func() {
		Context("the exists", func() {
			It("modifies the status", func() {
				state := NewState()
				state.InitializeState("http://127.0.0.1:4001")
				PurgeState()
				err := state.InitializeState("http://127.0.0.1:4001")
				Expect(err).To(BeNil())
				statusPool := state.StatusPool
				fakeStatus := fakes.FakeTaskStatus("fake_task_id_1", mesosproto.TaskState_TASK_RUNNING)
				err = statusPool.AddStatus(fakeStatus)
				Expect(err).To(BeNil())
				newState := mesosproto.TaskState_TASK_FAILED
				fakeStatus.State = &newState
				err = statusPool.ModifyStatus(fakeStatus.GetTaskId().GetValue(), fakeStatus)
				Expect(err).To(BeNil())
				expectedStatusJsonBytes, err := json.Marshal(fakeStatus)
				Expect(err).To(BeNil())
				expectedStatusJson := string(expectedStatusJsonBytes)
				actualStatusJson, err := lxdatabase.Get(state.StatusPool.GetKey() + "/" + fakeStatus.GetTaskId().GetValue())
				Expect(err).To(BeNil())
				Expect(actualStatusJson).To(Equal(expectedStatusJson))
			})
		})
		Context("the status doest exist", func() {
			It("returns an error", func() {
				state := NewState()
				state.InitializeState("http://127.0.0.1:4001")
				PurgeState()
				err := state.InitializeState("http://127.0.0.1:4001")
				Expect(err).To(BeNil())
				statusPool := state.StatusPool
				fakeStatus := fakes.FakeTaskStatus("fake_task_id_1", mesosproto.TaskState_TASK_RUNNING)
				err = statusPool.ModifyStatus(fakeStatus.GetTaskId().GetValue(), fakeStatus)
				Expect(err).NotTo(BeNil())
			})
		})
	})
	Describe("GetStatuses()", func() {
		It("returns all known statuses in the pool", func() {
			state := NewState()
			state.InitializeState("http://127.0.0.1:4001")
			PurgeState()
			err := state.InitializeState("http://127.0.0.1:4001")
			Expect(err).To(BeNil())
			statusPool := state.StatusPool
			fakeStatus1 := fakes.FakeTaskStatus("fake_task_id_1", mesosproto.TaskState_TASK_RUNNING)
			fakeStatus2 := fakes.FakeTaskStatus("fake_task_id_2", mesosproto.TaskState_TASK_RUNNING)
			fakeStatus3 := fakes.FakeTaskStatus("fake_task_id_3", mesosproto.TaskState_TASK_RUNNING)
			err = statusPool.AddStatus(fakeStatus1)
			Expect(err).To(BeNil())
			err = statusPool.AddStatus(fakeStatus2)
			Expect(err).To(BeNil())
			err = statusPool.AddStatus(fakeStatus3)
			Expect(err).To(BeNil())
			statuses, err := statusPool.GetStatuses()
			Expect(err).To(BeNil())
			Expect(statuses[fakeStatus1.GetTaskId().GetValue()]).To(Equal(fakeStatus1))
			Expect(statuses[fakeStatus2.GetTaskId().GetValue()]).To(Equal(fakeStatus2))
			Expect(statuses[fakeStatus3.GetTaskId().GetValue()]).To(Equal(fakeStatus3))
		})
	})
	Describe("DeleteStatus(statusId)", func() {
		Context("status exists", func() {
			It("deletes the status", func() {
				state := NewState()
				state.InitializeState("http://127.0.0.1:4001")
				PurgeState()
				err := state.InitializeState("http://127.0.0.1:4001")
				Expect(err).To(BeNil())
				statusPool := state.StatusPool
				fakeStatus1 := fakes.FakeTaskStatus("fake_task_id_1", mesosproto.TaskState_TASK_RUNNING)
				fakeStatus2 := fakes.FakeTaskStatus("fake_task_id_2", mesosproto.TaskState_TASK_RUNNING)
				fakeStatus3 := fakes.FakeTaskStatus("fake_task_id_3", mesosproto.TaskState_TASK_RUNNING)
				err = statusPool.AddStatus(fakeStatus1)
				Expect(err).To(BeNil())
				err = statusPool.AddStatus(fakeStatus2)
				Expect(err).To(BeNil())
				err = statusPool.AddStatus(fakeStatus3)
				Expect(err).To(BeNil())
				err = statusPool.DeleteStatus(fakeStatus1.GetTaskId().GetValue())
				Expect(err).To(BeNil())
				statuses, err := statusPool.GetStatuses()
				Expect(err).To(BeNil())
				Expect(statuses[fakeStatus1.GetTaskId().GetValue()]).To(BeNil())
				Expect(statuses[fakeStatus2.GetTaskId().GetValue()]).To(Equal(fakeStatus2))
				Expect(statuses[fakeStatus3.GetTaskId().GetValue()]).To(Equal(fakeStatus3))
			})
		})
		Context("status does not exist", func() {
			It("throws error", func() {
				state := NewState()
				state.InitializeState("http://127.0.0.1:4001")
				PurgeState()
				err := state.InitializeState("http://127.0.0.1:4001")
				Expect(err).To(BeNil())
				statusPool := state.StatusPool
				err = statusPool.DeleteStatus("nonexistent_status_id")
				Expect(err).NotTo(BeNil())
			})
		})
	})
})
