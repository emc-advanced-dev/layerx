package fakes

import "github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"

func FakeNode(resourceId string, nodeId string) *lxtypes.Node {
	fakeResource := lxtypes.NewResourceFromMesos(FakeOffer(resourceId, nodeId))
	node := lxtypes.NewNode(nodeId)
	node.AddResource(fakeResource)
	return node
}
