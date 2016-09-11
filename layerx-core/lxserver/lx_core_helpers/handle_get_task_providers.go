package lx_core_helpers

import (
	"github.com/emc-advanced-dev/layerx/layerx-core/lxstate"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"
	"github.com/emc-advanced-dev/pkg/errors"
)

func GetTaskProviders(state *lxstate.State) ([]*lxtypes.TaskProvider, error) {
	taskProviderMap, err := state.TaskProviderPool.GetTaskProviders()
	if err != nil {
		return nil, errors.New("getting task provider list from pool", err)
	}
	taskProviders := []*lxtypes.TaskProvider{}
	for _, taskProvider := range taskProviderMap {
		taskProviders = append(taskProviders, taskProvider)
	}
	return taskProviders, nil
}
