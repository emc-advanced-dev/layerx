package lx_core_helpers

import (
	"github.com/emc-advanced-dev/layerx/layerx-core/lxstate"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"
	"github.com/layer-x/layerx-commons/lxerrors"
)

func RegisterTaskProvider(state *lxstate.State, taskProvider *lxtypes.TaskProvider) error {
	err := state.TaskProviderPool.AddTaskProvider(taskProvider)
	if err != nil {
		return lxerrors.New("adding task provider to pool", err)
	}
	return nil
}
