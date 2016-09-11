package lx_core_helpers

import (
	"github.com/emc-advanced-dev/layerx/layerx-core/lxstate"
	"github.com/emc-advanced-dev/pkg/errors"
)

func PurgeTask(state *lxstate.State, taskId string) error {
	taskPool, err := state.GetTaskPoolContainingTask(taskId)
	if err != nil {
		return errors.New("could not find task pool containing task "+taskId, err)
	}
	err = taskPool.DeleteTask(taskId)
	if err != nil {
		return errors.New("could not delete task "+taskId, err)
	}
	return nil
}
