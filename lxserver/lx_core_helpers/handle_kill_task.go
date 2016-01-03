package lx_core_helpers
import (
	"github.com/layer-x/layerx-core_v2/lxstate"
	"github.com/layer-x/layerx-commons/lxerrors"
)

func KillTask(state *lxstate.State, taskId string) error {
	taskPool, err := state.GetTaskPoolContainingTask(taskId)
	if err != nil {
		return lxerrors.New("could not find task pool containing task "+taskId, err)
	}
	taskToKill, err := taskPool.GetTask(taskId)
	if err != nil {
		return lxerrors.New("could not find task pool containing task "+taskId, err)
	}
	taskToKill.KillRequested = true
	err = taskPool.ModifyTask(taskId, taskToKill)
	if err != nil {
		return lxerrors.New("could not task with KillRequested set back into task pool", err)
	}
	return nil
}
