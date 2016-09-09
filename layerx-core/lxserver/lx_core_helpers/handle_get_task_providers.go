package lx_core_helpers

import (
	"github.com/emc-advanced-dev/layerx/layerx-core/lxstate"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"
	"github.com/layer-x/layerx-commons/lxerrors"
)

func GetTaskProviders(state *lxstate.State) ([]*lxtypes.TaskProvider, error) {
	taskProviderMap, err := state.TaskProviderPool.GetTaskProviders()
	if err != nil {
		return nil, lxerrors.New("getting task provider list from pool", err)
	}
	taskProviders := []*lxtypes.TaskProvider{}
	for _, taskProvider := range taskProviderMap {
		taskProviders = append(taskProviders, taskProvider)
	}
	return taskProviders, nil
}
