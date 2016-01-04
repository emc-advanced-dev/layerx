package lx_core_helpers
import (
	"github.com/layer-x/layerx-core_v2/lxstate"
	"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/mesos/mesos-go/mesosproto"
	"github.com/layer-x/layerx-core_v2/tpi_messenger"
)

func SubmitStatusUpdate(state *lxstate.State, tpiUrl string, status *mesosproto.TaskStatus) error {
	taskId := status.GetTaskId().GetValue()
	task, err := state.GetTaskFromAnywhere(taskId)
	if err != nil {
		return lxerrors.New("fetching task for task_id "+taskId, err)
	}
	err = tpi_messenger.SendStatusUpdate(tpiUrl, task.TaskProvider, status)
	if err != nil {
		return lxerrors.New("sending status update to tpi", err)
	}
	err = state.StatusPool.AddStatus(status)
	if err != nil {
		return lxerrors.New("adding status for task "+taskId+" to pool", err)
	}
	return nil
}