package lx_core_helpers

import (
	"github.com/emc-advanced-dev/layerx/layerx-core/lxstate"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"
	"github.com/emc-advanced-dev/pkg/errors"
	"github.com/Sirupsen/logrus"
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
	logrus.WithFields(logrus.Fields{"resource": resource, "node": newNode}).Infof("created new node for resource")
	resourcePool, err := state.NodePool.GetNodeResourcePool(resource.NodeId)
	if err != nil {
		return errors.New("retrieving resource pool for new node "+resource.NodeId, err)
	}
	return addResourceToNode(resourcePool, resource)
}

func addResourceToNode(resourcePool *lxstate.ResourcePool, resource *lxtypes.Resource) error {
	//if resource exists with the same id on the node, replace it
	if _, err := resourcePool.GetResource(resource.Id); err == nil {
		return resourcePool.ModifyResource(resource.Id, resource)
	}

	err := resourcePool.AddResource(resource)
	if err != nil {
		return errors.New("adding resource "+resource.Id+" to resource pool", err)
	}
	logrus.WithFields(logrus.Fields{"resource": resource}).Debugf("accepted resource from rpi")
	return nil
}
