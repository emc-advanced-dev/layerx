package fakes
import "github.com/layer-x/layerx-core_v2/lxtypes"

func FakeNode(resourceId string, nodeId string) *lxtypes.Node {
	fakeResource := lxtypes.NewResourceFromMesos(FakeOffer(resourceId, nodeId))
	node := lxtypes.NewNode(nodeId)
	node.AddResource(fakeResource)
	return node
}