package lxstate

import (
	"encoding/json"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"
	"github.com/emc-advanced-dev/pkg/errors"
	"github.com/layer-x/layerx-commons/lxdatabase"
)

type ResourcePool struct {
	nodeId  string
	rootKey string
}

func TempResourcePoolFunction(rootKey, nodeId string) *ResourcePool {
	return &ResourcePool{
		rootKey: rootKey,
		nodeId:  nodeId,
	}
}

func (resourcePool *ResourcePool) GetKey() string {
	return resourcePool.rootKey + "_" + resourcePool.nodeId
}

func (resourcePool *ResourcePool) Initialize() error {
	err := lxdatabase.Mkdir(resourcePool.GetKey())
	if err != nil {
		return errors.New("initializing "+resourcePool.GetKey()+" directory", err)
	}
	return nil
}

func (resourcePool *ResourcePool) AddResource(resource *lxtypes.Resource) error {
	if resourcePool.nodeId != resource.NodeId {
		return errors.New("resource given was for node "+resource.NodeId+" but this node is "+resourcePool.nodeId, nil)
	}
	resourceId := resource.Id
	_, err := resourcePool.GetResource(resourceId)
	if err == nil {
		return errors.New("resource "+resourceId+" already exists in database, try Modify()?", err)
	}
	resourceData, err := json.Marshal(resource)
	if err != nil {
		return errors.New("could not marshal resource to json", err)
	}
	err = lxdatabase.Set(resourcePool.GetKey()+"/"+resourceId, string(resourceData))
	if err != nil {
		return errors.New("setting key/value pair for resource", err)
	}
	return nil
}

func (resourcePool *ResourcePool) GetResource(resourceId string) (*lxtypes.Resource, error) {
	resourceJson, err := lxdatabase.Get(resourcePool.GetKey() + "/" + resourceId)
	if err != nil {
		return nil, errors.New("retrieving resource "+resourceId+" from database", err)
	}
	var resource lxtypes.Resource
	err = json.Unmarshal([]byte(resourceJson), &resource)
	if err != nil {
		return nil, errors.New("unmarshalling json into Resource struct", err)
	}
	return &resource, nil
}

func (resourcePool *ResourcePool) ModifyResource(resourceId string, modifiedResource *lxtypes.Resource) error {
	_, err := resourcePool.GetResource(resourceId)
	if err != nil {
		return errors.New("resource "+resourceId+" not found", err)
	}
	resourceData, err := json.Marshal(modifiedResource)
	if err != nil {
		return errors.New("could not marshal modified resource to json", err)
	}
	err = lxdatabase.Set(resourcePool.GetKey()+"/"+resourceId, string(resourceData))
	if err != nil {
		return errors.New("setting key/value pair for modified resource", err)
	}
	return nil

}

func (resourcePool *ResourcePool) GetResources() (map[string]*lxtypes.Resource, error) {
	resources := make(map[string]*lxtypes.Resource)
	knownResources, err := lxdatabase.GetKeys(resourcePool.GetKey())
	if err != nil {
		return nil, errors.New("retrieving list of known resources", err)
	}
	for _, resourceJson := range knownResources {
		var resource lxtypes.Resource
		err = json.Unmarshal([]byte(resourceJson), &resource)
		if err != nil {
			return nil, errors.New("unmarshalling json into Resource struct", err)
		}
		resources[resource.Id] = &resource
	}
	return resources, nil
}

func (resourcePool *ResourcePool) DeleteResource(resourceId string) error {
	_, err := resourcePool.GetResource(resourceId)
	if err != nil {
		return errors.New("resource "+resourceId+" not found", err)
	}
	err = lxdatabase.Rm(resourcePool.GetKey() + "/" + resourceId)
	if err != nil {
		return errors.New("removing resource "+resourceId+" from database", err)
	}
	return nil
}
