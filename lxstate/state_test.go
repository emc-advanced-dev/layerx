package lxstate_test

import (
	. "github.com/layer-x/layerx-core_v2/lxstate"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/layer-x/layerx-commons/lxdatabase"
"github.com/layer-x/layerx-core_v2/fakes"
)

func PurgeState() {
	lxdatabase.Rmdir("/state", true)
}

var _ = Describe("State", func() {
	Describe("InitializeState(etcdUrl)", func() {
		It("initializes client (lxdb), creates folders for nodes, tasks, statuses, tps", func() {
			state := NewState()
			state.InitializeState("http://127.0.0.1:4001")
			PurgeState()
			err := state.InitializeState("http://127.0.0.1:4001")
			Expect(err).To(BeNil())
			rootContents, err := lxdatabase.GetSubdirectories("/state")
			Expect(err).To(BeNil())
			Expect(rootContents).To(ContainElement("/state/nodes"))
			Expect(rootContents).To(ContainElement("/state/pending_tasks"))
			Expect(rootContents).To(ContainElement("/state/staging_tasks"))
			Expect(rootContents).To(ContainElement("/state/task_providers"))
			Expect(rootContents).To(ContainElement("/state/statuses"))
		})
	})
	Describe("Set/GetTpiUrl(tpiUrl)", func() {
		It("sets and gets the tpiurl", func() {
			state := NewState()
			state.InitializeState("http://127.0.0.1:4001")
			PurgeState()
			err := state.InitializeState("http://127.0.0.1:4001")
			Expect(err).To(BeNil())
			err = state.SetTpi("fake_url")
			Expect(err).To(BeNil())
			tpiUrl, err := state.GetTpi()
			Expect(err).To(BeNil())
			Expect(tpiUrl).To(Equal("fake_url"))
		})
	})
	Describe("Set/GetRpiUrl(tpiUrl)", func() {
		It("sets and gets the rpiurl", func() {
			state := NewState()
			state.InitializeState("http://127.0.0.1:4001")
			PurgeState()
			err := state.InitializeState("http://127.0.0.1:4001")
			Expect(err).To(BeNil())
			err = state.SetRpi("fake_url")
			Expect(err).To(BeNil())
			rpiUrl, err := state.GetRpi()
			Expect(err).To(BeNil())
			Expect(rpiUrl).To(Equal("fake_url"))
		})
	})
	Describe("GetAllTasks", func(){
		It("returns all known tasks from pending, staging, and node task pools", func(){
			state := NewState()
			state.InitializeState("http://127.0.0.1:4001")
			PurgeState()
			err := state.InitializeState("http://127.0.0.1:4001")
			Expect(err).To(BeNil())
			fakePendingTask := fakes.FakeLXTask("fake__pending_task_id", "fake_task_name", "fake_slave_id", "echo FAKE_COMMAND")
			err = state.PendingTaskPool.AddTask(fakePendingTask)
			Expect(err).To(BeNil())
			fakeStagingTask := fakes.FakeLXTask("fake_staging_task_id", "fake_task_name", "fake_slave_id", "echo FAKE_COMMAND")
			err = state.StagingTaskPool.AddTask(fakeStagingTask)
			Expect(err).To(BeNil())
			fakeNode := fakes.FakeNode("fake_resource_id_1", "fake_node_id_1")
			err = state.NodePool.AddNode(fakeNode)
			Expect(err).To(BeNil())
			nodeTaskPool, err := state.NodePool.GetNodeTaskPool(fakeNode.Id)
			Expect(err).To(BeNil())
			fakeNodeTask1 := fakes.FakeLXTask("fake__node_task_id_1", "fake_task", "fake_node_id_1", "echo FAKECOMMAND")
			err = nodeTaskPool.AddTask(fakeNodeTask1)
			Expect(err).To(BeNil())
			fakeNodeTask2 := fakes.FakeLXTask("fake__node_task_id_2", "fake_task", "fake_node_id_1", "echo FAKECOMMAND")
			err = nodeTaskPool.AddTask(fakeNodeTask2)
			Expect(err).To(BeNil())
			fakeNodeTask3 := fakes.FakeLXTask("fake__node_task_id_3", "fake_task", "fake_node_id_1", "echo FAKECOMMAND")
			err = nodeTaskPool.AddTask(fakeNodeTask3)
			Expect(err).To(BeNil())
			allTasks, err := state.GetAllTasks()
			Expect(err).To(BeNil())
			Expect(allTasks[fakePendingTask.TaskId]).To(Equal(fakePendingTask))
			Expect(allTasks[fakeStagingTask.TaskId]).To(Equal(fakeStagingTask))
			Expect(allTasks[fakeNodeTask1.TaskId]).To(Equal(fakeNodeTask1))
			Expect(allTasks[fakeNodeTask2.TaskId]).To(Equal(fakeNodeTask2))
			Expect(allTasks[fakeNodeTask3.TaskId]).To(Equal(fakeNodeTask3))
		})
	})
})
