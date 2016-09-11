package lx_core_helpers

import (
	"github.com/emc-advanced-dev/layerx/layerx-core/lxstate"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"
	"github.com/emc-advanced-dev/pkg/errors"
)

func RegisterTaskProvider(state *lxstate.State, taskProvider *lxtypes.TaskProvider) error {
	err := state.TaskProviderPool.AddTaskProvider(taskProvider)
	if err != nil {
		return errors.New("adding task provider to pool", err)
	}
	return nil
}
