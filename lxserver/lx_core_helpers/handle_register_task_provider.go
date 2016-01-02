package lx_core_helpers
import (
	"github.com/layer-x/layerx-core_v2/lxstate"
	"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/layer-x/layerx-core_v2/lxtypes"
)

func RegisterTaskProvider(state *lxstate.State, taskProvider *lxtypes.TaskProvider) error {
	err := state.TaskProviderPool.AddTaskProvider(taskProvider)
	if err != nil {
		return lxerrors.New("adding task provider to pool", err)
	}
	return nil
}
