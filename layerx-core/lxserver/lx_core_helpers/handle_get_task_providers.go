package lx_core_helpers
import (
	"github.com/emc-advanced-dev/layerx-core/lxstate"
	"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/emc-advanced-dev/layerx-core/lxtypes"
)

func  GetTaskProviders(state *lxstate.State) ([]*lxtypes.TaskProvider, error) {
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
