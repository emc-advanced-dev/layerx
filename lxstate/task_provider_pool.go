package lxstate
import (
	"github.com/emc-advanced-dev/layerx-core/lxtypes"
	"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/layer-x/layerx-commons/lxdatabase"
	"encoding/json"
)

type TaskProviderPool struct {
	rootKey string
}

func (taskProviderPool *TaskProviderPool) GetKey() string {
	return taskProviderPool.rootKey
}

func (taskProviderPool *TaskProviderPool) Initialize() error {
	err := lxdatabase.Mkdir(taskProviderPool.GetKey())
	if err != nil {
		return lxerrors.New("initializing "+taskProviderPool.GetKey() +" directory", err)
	}
	return nil
}

func (taskProviderPool *TaskProviderPool) AddTaskProvider(taskProvider *lxtypes.TaskProvider) error {
	taskProviderId := taskProvider.Id
	taskProviderData, err := json.Marshal(taskProvider)
	if err != nil {
		return lxerrors.New("could not marshal taskProvider to json", err)
	}
	err = lxdatabase.Set(taskProviderPool.GetKey()+"/"+taskProviderId, string(taskProviderData))
	if err != nil {
		return lxerrors.New("setting key/value pair for taskProvider", err)
	}
	return nil
}

func (taskProviderPool *TaskProviderPool) GetTaskProvider(taskProviderId string) (*lxtypes.TaskProvider, error) {
	taskProviderJson, err := lxdatabase.Get(taskProviderPool.GetKey()+"/"+taskProviderId)
	if err != nil {
		return nil, lxerrors.New("retrieving taskProvider "+taskProviderId+" from database", err)
	}
	var taskProvider lxtypes.TaskProvider
	err = json.Unmarshal([]byte(taskProviderJson), &taskProvider)
	if err != nil {
		return nil, lxerrors.New("unmarshalling json into TaskProvider struct", err)
	}
	return &taskProvider, nil
}

func (taskProviderPool *TaskProviderPool) GetTaskProviders() (map[string]*lxtypes.TaskProvider, error) {
	taskProviders := make(map[string]*lxtypes.TaskProvider)
	knownTaskProviders, err := lxdatabase.GetKeys(taskProviderPool.GetKey())
	if err != nil {
		return nil, lxerrors.New("retrieving list of known taskProviders", err)
	}
	for _, taskProviderJson := range knownTaskProviders {
		var taskProvider lxtypes.TaskProvider
		err = json.Unmarshal([]byte(taskProviderJson), &taskProvider)
		if err != nil {
			return nil, lxerrors.New("unmarshalling json into TaskProvider struct", err)
		}
		taskProviders[taskProvider.Id] = &taskProvider
	}
	return taskProviders, nil
}

func (taskProviderPool *TaskProviderPool) DeleteTaskProvider(taskProviderId string) error {
	_, err := taskProviderPool.GetTaskProvider(taskProviderId)
	if err != nil {
		return lxerrors.New("taskProvider "+taskProviderId+" not found", err)
	}
	err = lxdatabase.Rm(taskProviderPool.GetKey()+"/"+taskProviderId)
	if err != nil {
		return lxerrors.New("removing taskProvider "+taskProviderId+" from database", err)
	}
	return nil
}