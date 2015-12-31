package state
import (
	"github.com/layer-x/layerx-core_v2/lxtypes"
	"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/layer-x/layerx-commons/lxdatabase"
	"encoding/json"
)

type ResourcePool struct {
	rootKey string
}

func TempResourcePoolFunction(rootKey string) *ResourcePool {
	return &ResourcePool{
		rootKey: rootKey,
	}
}

func (resourcePool *ResourcePool) GetKey() string {
	return resourcePool.rootKey
}

func (resourcePool *ResourcePool) Initialize() error {
	err := lxdatabase.Mkdir(resourcePool.GetKey())
	if err != nil {
		return lxerrors.New("initializing "+resourcePool.GetKey() +" directory", err)
	}
	return nil
}

func (resourcePool *ResourcePool) AddResource(resource *lxtypes.Resource) error {
	resourceId := resource.Id
	_, err := resourcePool.GetResource(resourceId)
	if err == nil {
		return lxerrors.New("resource "+resourceId+" already exists in database, try Modify()?", err)
	}
	resourceData, err := json.Marshal(resource)
	if err != nil {
		return lxerrors.New("could not marshal resource to json", err)
	}
	err = lxdatabase.Set(resourcePool.GetKey()+"/"+resourceId, string(resourceData))
	if err != nil {
		return lxerrors.New("setting key/value pair for resource", err)
	}
	return nil
}

func (resourcePool *ResourcePool) GetResource(resourceId string) (*lxtypes.Resource, error) {
	resourceJson, err := lxdatabase.Get(resourcePool.GetKey()+"/"+resourceId)
	if err != nil {
		return nil, lxerrors.New("retrieving resource "+resourceId+" from database", err)
	}
	var resource lxtypes.Resource
	err = json.Unmarshal([]byte(resourceJson), &resource)
	if err != nil {
		return nil, lxerrors.New("unmarshalling json into Resource struct", err)
	}
	return &resource, nil
}

func (resourcePool *ResourcePool) ModifyResource(resourceId string, modifiedResource *lxtypes.Resource) error {
	_, err := resourcePool.GetResource(resourceId)
	if err != nil {
		return lxerrors.New("resource "+resourceId+" not found", err)
	}
	resourceData, err := json.Marshal(modifiedResource)
	if err != nil {
		return lxerrors.New("could not marshal modified resource to json", err)
	}
	err = lxdatabase.Set(resourcePool.GetKey()+"/"+resourceId, string(resourceData))
	if err != nil {
		return lxerrors.New("setting key/value pair for modified resource", err)
	}
	return nil

}

func (resourcePool *ResourcePool) GetResources() (map[string]*lxtypes.Resource, error) {
	resources := make(map[string]*lxtypes.Resource)
	knownResources, err := lxdatabase.GetKeys(resourcePool.GetKey())
	if err != nil {
		return nil, lxerrors.New("retrieving list of known resources", err)
	}
	for _, resourceJson := range knownResources {
		var resource lxtypes.Resource
		err = json.Unmarshal([]byte(resourceJson), &resource)
		if err != nil {
			return nil, lxerrors.New("unmarshalling json into Resource struct", err)
		}
		resources[resource.Id] = &resource
	}
	return resources, nil
}

func (resourcePool *ResourcePool) DeleteResource(resourceId string) error {
	_, err := resourcePool.GetResource(resourceId)
	if err != nil {
		return lxerrors.New("resource "+resourceId+" not found", err)
	}
	err = lxdatabase.Rm(resourcePool.GetKey()+"/"+resourceId)
	if err != nil {
		return lxerrors.New("removing resource "+resourceId+" from database", err)
	}
	return nil
}