package lx_core_helpers

import (
	"github.com/Sirupsen/logrus"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxstate"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"
	"github.com/emc-advanced-dev/layerx/layerx-core/tpi_messenger"
	"github.com/emc-advanced-dev/pkg/errors"
	"github.com/mesos/mesos-go/mesosproto"
)

func ProcessStatusUpdate(state *lxstate.State, tpiUrl string, status *mesosproto.TaskStatus) error {
	taskId := status.GetTaskId().GetValue()
	taskPool, err := state.GetTaskPoolContainingTask(taskId)
	if err != nil {
		return errors.New("getting task pool containing task "+taskId, err)
	}
	task, err := taskPool.GetTask(taskId)
	if err != nil {
		return errors.New("fetching task for task_id "+taskId, err)
	}
	if isTerminal(status) && !task.Checkpointed {
		err = taskPool.DeleteTask(taskId)
		if err != nil {
			return errors.New("deleting terminal task "+taskId+" with status "+status.GetState().String(), err)
		}
	}
	if taskPool == state.StagingTaskPool && status.GetState() == mesosproto.TaskState_TASK_RUNNING {
		nodeId := task.NodeId
		task.Checkpointed = false
		nodeTaskPool, err := state.NodePool.GetNodeTaskPool(nodeId)
		if err != nil {
			return errors.New("retrieving node pool for node "+nodeId, err)
		}
		err = moveTaskBetweenPools(task, taskPool, nodeTaskPool)
		if err != nil {
			return errors.New("migrating task from staging task pool to node "+nodeId+" task pook", err)
		}
	}
	if !task.Checkpointed {
		err = tpi_messenger.SendStatusUpdate(tpiUrl, task.TaskProvider, status)
		if err != nil {
			return errors.New("sending status update to tpi", err)
		}
	} else {
		logrus.WithFields(logrus.Fields{"task": task, "status": status}).Warnf("task is checkpointed, not bubbling status update")
	}
	err = state.StatusPool.AddStatus(status)
	if err != nil {
		return errors.New("adding status for task "+taskId+" to pool", err)
	}
	return nil
}

func moveTaskBetweenPools(task *lxtypes.Task, sourceTaskPool, destinationTaskPool *lxstate.TaskPool) error {
	err := sourceTaskPool.DeleteTask(task.TaskId)
	if err != nil {
		return errors.New("deleting task from source task pool ", err)
	}
	err = destinationTaskPool.AddTask(task)
	if err != nil {
		return errors.New("moving task into destination task pool", err)
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
