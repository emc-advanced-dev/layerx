package lx_core_helpers
import (
	"github.com/layer-x/layerx-core_v2/lxstate"
	"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/layer-x/layerx-core_v2/rpi_messenger"
)

func KillTask(state *lxstate.State, rpiUrl, taskId string) error {
	taskPool, err := state.GetTaskPoolContainingTask(taskId)
	if err != nil {
		return lxerrors.New("could not find task pool containing task "+taskId, err)
	}
	err = rpi_messenger.SendKillTaskRequest(rpiUrl, taskId)
	if err != nil {
		return lxerrors.New("sending kill task request to rpi", err)
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
