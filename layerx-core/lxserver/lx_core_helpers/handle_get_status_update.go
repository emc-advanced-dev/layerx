package lx_core_helpers

import (
	"github.com/emc-advanced-dev/layerx/layerx-core/lxstate"
	"github.com/emc-advanced-dev/pkg/errors"
	"github.com/gogo/protobuf/proto"
	"github.com/mesos/mesos-go/mesosproto"
)

func GetStatusUpdate(state *lxstate.State, taskId string) (*mesosproto.TaskStatus, error) {
	status, err := state.StatusPool.GetStatus(taskId)
	if err == nil {
		return status, nil
	}
	taskPool, err := state.GetTaskPoolContainingTask(taskId)
	if err != nil {
		return generateTaskStatus(taskId, mesosproto.TaskState_TASK_LOST, "task not found"), nil
	}
	if taskPool == state.PendingTaskPool {
		return generateTaskStatus(taskId, mesosproto.TaskState_TASK_STAGING, "task is waiting to be scheduled"), nil
	}
	if taskPool == state.StagingTaskPool {
		return generateTaskStatus(taskId, mesosproto.TaskState_TASK_STARTING, "task has been assigned, waiting for status"), nil
	}
	return nil, errors.New("task exists on node but no status known yet?", nil)
}

func generateTaskStatus(taskId string, taskState mesosproto.TaskState, message string) *mesosproto.TaskStatus {
	return &mesosproto.TaskStatus{
		TaskId: &mesosproto.TaskID{
			Value: proto.String(taskId),
		},
		State:   &taskState,
		Message: proto.String(message),
	}
}
