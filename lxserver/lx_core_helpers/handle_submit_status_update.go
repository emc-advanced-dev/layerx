package lx_core_helpers
import (
	"github.com/layer-x/layerx-core_v2/lxstate"
	"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/mesos/mesos-go/mesosproto"
	"github.com/layer-x/layerx-core_v2/tpi_messenger"
	"github.com/layer-x/layerx-core_v2/lxtypes"
)

func ProcessStatusUpdate(state *lxstate.State, tpiUrl string, status *mesosproto.TaskStatus) error {
	taskId := status.GetTaskId().GetValue()
	taskPool, err := state.GetTaskPoolContainingTask(taskId)
	if err != nil {
		return lxerrors.New("getting task pool containing task "+taskId, err)
	}
	task, err := taskPool.GetTask(taskId)
	if err != nil {
		return lxerrors.New("fetching task for task_id "+taskId, err)
	}
	if isTerminal(status) {
		err = taskPool.DeleteTask(taskId)
		if err != nil {
			return lxerrors.New("deleting terminal task "+taskId+" with status "+status.GetState().String(), err)
		}
	}
	if taskPool == state.StagingTaskPool && status.GetState() == mesosproto.TaskState_TASK_RUNNING {
		nodeId := task.NodeId
		nodeTaskPool, err := state.NodePool.GetNodeTaskPool(nodeId)
		if err != nil {
			return lxerrors.New("retrieving node pool for node "+nodeId, err)
		}
		err = moveTaskBetweenPools(task, taskPool, nodeTaskPool)
		if err != nil {
			return lxerrors.New("migrating task from staging task pool to node "+nodeId+" task pook", err)
		}
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

func moveTaskBetweenPools(task *lxtypes.Task, sourceTaskPool, destinationTaskPool *lxstate.TaskPool) error {
	err := sourceTaskPool.DeleteTask(task.TaskId)
	if err != nil {
		return lxerrors.New("deleting task from source task pool ", err)
	}
	err = destinationTaskPool.AddTask(task)
	if err != nil {
		return lxerrors.New("moving task into destination task pool", err)
	}
	return nil
}

func isTerminal(status *mesosproto.TaskStatus) bool {
	state := status.GetState()
	return state == mesosproto.TaskState_TASK_ERROR ||
			state == mesosproto.TaskState_TASK_FAILED ||
			state == mesosproto.TaskState_TASK_FINISHED ||
			state == mesosproto.TaskState_TASK_KILLED ||
			state == mesosproto.TaskState_TASK_LOST
}