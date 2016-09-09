package lxstate_test

import (
	. "github.com/emc-advanced-dev/layerx-core/lxstate"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/layer-x/layerx-commons/lxdatabase"
	"encoding/json"
	"github.com/emc-advanced-dev/layerx-core/fakes"
)

var _ = Describe("TaskProviderPool", func() {
	Describe("GetTaskProvider(taskProviderId)", func(){
		It("returns the taskProvider if it exists, else returns err", func(){
			state := NewState()
			state.InitializeState("http://127.0.0.1:4001")
			PurgeState()
			err := state.InitializeState("http://127.0.0.1:4001")
			Expect(err).To(BeNil())
			taskProviderPool := state.TaskProviderPool
			fakeTaskProvider := fakes.FakeTaskProvider("fake_taskProvider_id_1", "faketp@fakeip")
			taskProvider, err := taskProviderPool.GetTaskProvider(fakeTaskProvider.Id)
			Expect(err).NotTo(BeNil())
			Expect(taskProvider).To(BeNil())
			err = taskProviderPool.AddTaskProvider(fakeTaskProvider)
			Expect(err).To(BeNil())
			taskProvider, err = taskProviderPool.GetTaskProvider(fakeTaskProvider.Id)
			Expect(err).To(BeNil())
			Expect(taskProvider).To(Equal(fakeTaskProvider))
		})
	})
	Describe("AddTaskProvider", func(){
		Context("the taskProvider is new", func(){
			It("adds the taskProvider to etcd state", func(){
				state := NewState()
				state.InitializeState("http://127.0.0.1:4001")
				PurgeState()
				err := state.InitializeState("http://127.0.0.1:4001")
				Expect(err).To(BeNil())
				taskProviderPool := state.TaskProviderPool
				fakeTaskProvider := fakes.FakeTaskProvider("fake_taskProvider_id_1", "faketp@fakeip")
				err = taskProviderPool.AddTaskProvider(fakeTaskProvider)
				Expect(err).To(BeNil())
				expectedTaskProviderJsonBytes, err := json.Marshal(fakeTaskProvider)
				Expect(err).To(BeNil())
				expectedTaskProviderJson := string(expectedTaskProviderJsonBytes)
				actualTaskProviderJson, err := lxdatabase.Get(state.TaskProviderPool.GetKey() + "/"+fakeTaskProvider.Id)
				Expect(err).To(BeNil())
				Expect(actualTaskProviderJson).To(Equal(expectedTaskProviderJson))
			})
		})
		Context("the taskProvider is not new", func(){
			It("accepts the task provider", func(){
				state := NewState()
				state.InitializeState("http://127.0.0.1:4001")
				PurgeState()
				err := state.InitializeState("http://127.0.0.1:4001")
				Expect(err).To(BeNil())
				taskProviderPool := state.TaskProviderPool
				fakeTaskProvider := fakes.FakeTaskProvider("fake_taskProvider_id_1", "faketp@fakeip")
				err = taskProviderPool.AddTaskProvider(fakeTaskProvider)
				Expect(err).To(BeNil())
				err = taskProviderPool.AddTaskProvider(fakeTaskProvider)
				Expect(err).To(BeNil())
			})
		})
	})
	Describe("GetTaskProviders()", func(){
		It("returns all known taskProviders in the pool", func(){
			state := NewState()
			state.InitializeState("http://127.0.0.1:4001")
			PurgeState()
			err := state.InitializeState("http://127.0.0.1:4001")
			Expect(err).To(BeNil())
			taskProviderPool := state.TaskProviderPool
			fakeTaskProvider1 := fakes.FakeTaskProvider("fake_taskProvider_id_1", "faketp@fakeip")
			fakeTaskProvider2 := fakes.FakeTaskProvider("fake_taskProvider_id_2", "faketp2@fakeip")
			fakeTaskProvider3 := fakes.FakeTaskProvider("fake_taskProvider_id_3", "faketp3@fakeip")
			err = taskProviderPool.AddTaskProvider(fakeTaskProvider1)
			Expect(err).To(BeNil())
			err = taskProviderPool.AddTaskProvider(fakeTaskProvider2)
			Expect(err).To(BeNil())
			err = taskProviderPool.AddTaskProvider(fakeTaskProvider3)
			Expect(err).To(BeNil())
			taskProviders, err := taskProviderPool.GetTaskProviders()
			Expect(err).To(BeNil())
			Expect(taskProviders[fakeTaskProvider1.Id]).To(Equal(fakeTaskProvider1))
			Expect(taskProviders[fakeTaskProvider2.Id]).To(Equal(fakeTaskProvider2))
			Expect(taskProviders[fakeTaskProvider3.Id]).To(Equal(fakeTaskProvider3))
		})
	})
	Describe("DeleteTaskProvider(taskProviderId)", func(){
		Context("taskProvider exists", func(){
			It("deletes the taskProvider", func(){
				state := NewState()
				state.InitializeState("http://127.0.0.1:4001")
				PurgeState()
				err := state.InitializeState("http://127.0.0.1:4001")
				Expect(err).To(BeNil())
				taskProviderPool := state.TaskProviderPool
				fakeTaskProvider1 := fakes.FakeTaskProvider("fake_taskProvider_id_1", "faketp@fakeip")
				fakeTaskProvider2 := fakes.FakeTaskProvider("fake_taskProvider_id_2", "faketp2@fakeip")
				fakeTaskProvider3 := fakes.FakeTaskProvider("fake_taskProvider_id_3", "faketp3@fakeip")
				err = taskProviderPool.AddTaskProvider(fakeTaskProvider1)
				Expect(err).To(BeNil())
				err = taskProviderPool.AddTaskProvider(fakeTaskProvider2)
				Expect(err).To(BeNil())
				err = taskProviderPool.AddTaskProvider(fakeTaskProvider3)
				Expect(err).To(BeNil())
				err = taskProviderPool.DeleteTaskProvider(fakeTaskProvider1.Id)
				Expect(err).To(BeNil())
				taskProviders, err := taskProviderPool.GetTaskProviders()
				Expect(err).To(BeNil())
				Expect(taskProviders[fakeTaskProvider1.Id]).To(BeNil())
				Expect(taskProviders[fakeTaskProvider2.Id]).To(Equal(fakeTaskProvider2))
				Expect(taskProviders[fakeTaskProvider3.Id]).To(Equal(fakeTaskProvider3))
			})
		})
		Context("taskProvider does not exist", func(){
			It("throws error", func(){
				state := NewState()
				state.InitializeState("http://127.0.0.1:4001")
				PurgeState()
				err := state.InitializeState("http://127.0.0.1:4001")
				Expect(err).To(BeNil())
				taskProviderPool := state.TaskProviderPool
				err = taskProviderPool.DeleteTaskProvider("nonexistent_taskProvider_id")
				Expect(err).NotTo(BeNil())
			})
		})
	})
})
