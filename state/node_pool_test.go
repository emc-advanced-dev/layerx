package state_test
//
//import (
//	. "github.com/layer-x/layerx-core_v2/state"
//
//	. "github.com/onsi/ginkgo"
//	. "github.com/onsi/gomega"
//	"github.com/layer-x/layerx-core_v2/fakes"
//)
//
//var _ = Describe("NodePool", func() {
//	Describe("GetNode(nodeId)", func(){
//		It("returns the node with all of its tasks if it exists, else returns err", func(){
//			state := NewState()
//			state.InitializeState("http://127.0.0.1:4001")
//			PurgeState()
//			err := state.InitializeState("http://127.0.0.1:4001")
//			Expect(err).To(BeNil())
//			nodePool := state.NodePool
//			fakeNode := fakes.FakeNode("fake_resource_id_1", "fake_node_id_1")
//			node, err := nodePool.GetNode(fakeNode.Id)
//			Expect(err).NotTo(BeNil())
//			Expect(node).To(BeNil())
//			err = nodePool.AddNode(fakeNode)
//			Expect(err).To(BeNil())
//			node, err = nodePool.GetNode(fakeNode.Id)
//			Expect(err).To(BeNil())
//			Expect(node).To(Equal(fakeNode))
//		})
//	})
//})
