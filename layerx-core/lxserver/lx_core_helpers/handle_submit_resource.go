package lx_core_helpers

import (
	"github.com/emc-advanced-dev/layerx/layerx-core/lxstate"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"
	"github.com/emc-advanced-dev/pkg/errors"
)

func SubmitResource(state *lxstate.State, resource *lxtypes.Resource) error {
	if resourcePool, _ := state.NodePool.GetNodeResourcePool(resource.NodeId); resourcePool != nil {
		return addResourceToNode(resourcePool, resource)
	}
	newNode := &lxtypes.Node{
		Id: resource.NodeId,
	}
	err := state.NodePool.AddNode(newNode)
	if err != nil {
		return errors.New("adding new node "+resource.NodeId+" to node pool", err)
	}
	resourcePool, err := state.NodePool.GetNodeResourcePool(resource.NodeId)
	if err != nil {
		return errors.New("retrieving resource pool for new node "+resource.NodeId, err)
	}
	return addResourceToNode(resourcePool, resource)
}

func addResourceToNode(resourcePool *lxstate.ResourcePool, resource *lxtypes.Resource) error {
	err := resourcePool.AddResource(resource)
	if err != nil {
		return errors.New("adding resource "+resource.Id+" to resource pool", err)
	}
	return nil
}
