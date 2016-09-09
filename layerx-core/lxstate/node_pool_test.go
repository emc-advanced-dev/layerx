package lxstate_test

import (
	. "github.com/emc-advanced-dev/layerx/layerx-core/lxstate"

	"github.com/emc-advanced-dev/layerx/layerx-core/fakes"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NodePool", func() {
	Describe("AddNode(nodeId)", func() {
		Context("the node is new", func() {
			It(" tasks if it exists, else returns err", func() {
				state := NewState()
				state.InitializeState("http://127.0.0.1:4001")
				PurgeState()
				err := state.InitializeState("http://127.0.0.1:4001")
				Expect(err).To(BeNil())
				nodePool := state.NodePool
				fakeNode := fakes.FakeNode("fake_resource_id_1", "fake_node_id_1")
				node, err := nodePool.GetNode(fakeNode.Id)
				Expect(err).NotTo(BeNil())
				Expect(node).To(BeNil())
				err = nodePool.AddNode(fakeNode)
				Expect(err).To(BeNil())
				node, err = nodePool.GetNode(fakeNode.Id)
				Expect(err).To(BeNil())
				Expect(node).To(Equal(fakeNode))
			})
		})
	})
	Describe("GetNode(nodeId)", func() {
		Context("the node exists", func() {
			It("returns the node with all of its tasks", func() {
				state := NewState()
				state.InitializeState("http://127.0.0.1:4001")
				PurgeState()
				err := state.InitializeState("http://127.0.0.1:4001")
				Expect(err).To(BeNil())
				nodePool := state.NodePool
				fakeNode := fakes.FakeNode("fake_resource_id_1", "fake_node_id_1")
				err = nodePool.AddNode(fakeNode)
				Expect(err).To(BeNil())
				node, err := nodePool.GetNode(fakeNode.Id)
				Expect(err).To(BeNil())
				Expect(node).To(Equal(fakeNode))
			})
		})
		Context("the node does not exist", func() {
			It("returns err", func() {
				state := NewState()
				state.InitializeState("http://127.0.0.1:4001")
				PurgeState()
				err := state.InitializeState("http://127.0.0.1:4001")
				Expect(err).To(BeNil())
				nodePool := state.NodePool
				fakeNode := fakes.FakeNode("fake_resource_id_1", "fake_node_id_1")
				node, err := nodePool.GetNode(fakeNode.Id)
				Expect(err).NotTo(BeNil())
				Expect(node).To(BeNil())
				err = nodePool.AddNode(fakeNode)
				Expect(err).To(BeNil())
				node, err = nodePool.GetNode(fakeNode.Id)
				Expect(err).To(BeNil())
				Expect(node).To(Equal(fakeNode))
			})
		})
		Describe("GetNodes()", func() {
			It("returns all known nodes", func() {
				state := NewState()
				state.InitializeState("http://127.0.0.1:4001")
				PurgeState()
				err := state.InitializeState("http://127.0.0.1:4001")
				Expect(err).To(BeNil())
				nodePool := state.NodePool
				fakeNode1 := fakes.FakeNode("fake_resource_id_1", "fake_node_id_1")
				fakeNode2 := fakes.FakeNode("fake_resource_id_2", "fake_node_id_2")
				fakeNode3 := fakes.FakeNode("fake_resource_id_3", "fake_node_id_3")
				err = nodePool.AddNode(fakeNode1)
				Expect(err).To(BeNil())
				err = nodePool.AddNode(fakeNode2)
				Expect(err).To(BeNil())
				err = nodePool.AddNode(fakeNode3)
				Expect(err).To(BeNil())
				nodes, err := nodePool.GetNodes()
				Expect(err).To(BeNil())
				Expect(nodes).To(ContainElement(fakeNode1))
				Expect(nodes).To(ContainElement(fakeNode2))
				Expect(nodes).To(ContainElement(fakeNode3))
			})
		})
		Describe("DeleteNode(nodeId)", func() {
			Context("the node exists", func() {
				It("returns the node with all of its tasks", func() {
					state := NewState()
					state.InitializeState("http://127.0.0.1:4001")
					PurgeState()
					err := state.InitializeState("http://127.0.0.1:4001")
					Expect(err).To(BeNil())
					nodePool := state.NodePool
					fakeNode := fakes.FakeNode("fake_resource_id_1", "fake_node_id_1")
					err = nodePool.AddNode(fakeNode)
					Expect(err).To(BeNil())
					node, err := nodePool.GetNode(fakeNode.Id)
					Expect(err).To(BeNil())
					Expect(node).To(Equal(fakeNode))
					err = nodePool.DeleteNode(fakeNode.Id)
					Expect(err).To(BeNil())
					_, err = nodePool.GetNode(fakeNode.Id)
					Expect(err).NotTo(BeNil())
				})
			})
			Context("the node does not exist", func() {
				It("returns err", func() {
					state := NewState()
					state.InitializeState("http://127.0.0.1:4001")
					PurgeState()
					err := state.InitializeState("http://127.0.0.1:4001")
					Expect(err).To(BeNil())
					nodePool := state.NodePool
					fakeNode := fakes.FakeNode("fake_resource_id_1", "fake_node_id_1")
					err = nodePool.DeleteNode(fakeNode.Id)
					Expect(err).NotTo(BeNil())
				})
			}) //TODO: DELETENODE TEST, GETRESOURCEPOOL GETTASKPOOL
		})
		Describe("GetNodeResourcePool", func() {
			It("returns the resource pool for the specified nodeid", func() {
				state := NewState()
				state.InitializeState("http://127.0.0.1:4001")
				PurgeState()
				err := state.InitializeState("http://127.0.0.1:4001")
				Expect(err).To(BeNil())
				nodePool := state.NodePool
				fakeNode := fakes.FakeNode("fake_resource_id_1", "fake_node_id_1")
				err = nodePool.AddNode(fakeNode)
				Expect(err).To(BeNil())
				nodeResourcePool, err := nodePool.GetNodeResourcePool(fakeNode.Id)
				Expect(err).To(BeNil())
				Expect(nodeResourcePool).NotTo(BeNil())
				fakeResource := lxtypes.NewResourceFromMesos(fakes.FakeOffer("fake_resource_id_2", "fake_node_id_1"))
				err = nodeResourcePool.AddResource(fakeResource)
				Expect(err).To(BeNil())
				node, err := nodePool.GetNode(fakeNode.Id)
				Expect(err).To(BeNil())
				Expect(node.GetResources()).To(ContainElement(fakeResource))
			})
		})
		Describe("GetNodeTaskPool", func() {
			It("returns the task pool for the specified nodeid", func() {
				state := NewState()
				state.InitializeState("http://127.0.0.1:4001")
				PurgeState()
				err := state.InitializeState("http://127.0.0.1:4001")
				Expect(err).To(BeNil())
				nodePool := state.NodePool
				fakeNode := fakes.FakeNode("fake_resource_id_1", "fake_node_id_1")
				err = nodePool.AddNode(fakeNode)
				Expect(err).To(BeNil())
				nodeTaskPool, err := nodePool.GetNodeTaskPool(fakeNode.Id)
				Expect(err).To(BeNil())
				Expect(nodeTaskPool).NotTo(BeNil())
				fakeTask := fakes.FakeLXTask("fake_task_id_1", "fake_task", "fake_node_id_1", "echo FAKECOMMAND")
				fakeTask.TaskProvider = fakes.FakeTaskProvider("fake_task_provider_id", "tp@fakeip:fakeport")
				err = nodeTaskPool.AddTask(fakeTask)
				Expect(err).To(BeNil())
				node, err := nodePool.GetNode(fakeNode.Id)
				Expect(err).To(BeNil())
				Expect(node.GetTasks()).To(ContainElement(fakeTask))
			})
		})
	})
})
