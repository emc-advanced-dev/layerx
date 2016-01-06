package lx_core_helpers
import (
	"github.com/layer-x/layerx-core_v2/lxstate"
	"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/layer-x/layerx-core_v2/rpi_messenger"
	"github.com/layer-x/layerx-commons/lxlog"
	"github.com/Sirupsen/logrus"
	"github.com/layer-x/layerx-core_v2/tpi_messenger"
	"github.com/mesos/mesos-go/mesosproto"
)

func KillTask(state *lxstate.State, tpiUrl, rpiUrl, taskProviderId, taskId string) error {
	if _, err := state.GetTaskFromAnywhere(taskId); err != nil {
		lxlog.Warnf(logrus.Fields{"task_id": taskId, "task_provider": taskProviderId}, "requested to kill a task that Layer-X has no knowledge of, replying with TASK_LOST")
		err = sendTaskKilledStatus(state, mesosproto.TaskState_TASK_LOST, tpiUrl, taskProviderId, taskId)
		if err != nil {
			return lxerrors.New("sending TASK_KILLED status to task provider"+taskProviderId, err)
		}
		return nil
	}
	taskPool, err := state.GetTaskPoolContainingTask(taskId)
	if err != nil {
		return lxerrors.New("could not find task pool containing task "+taskId, err)
	}
	//if task is staging or pending, just delete it and say we did
	if taskPool == state.PendingTaskPool || taskPool == state.StagingTaskPool {
		lxlog.Warnf(logrus.Fields{"task_id": taskId, "task_provider": taskProviderId}, "requested to kill a task before task staging was complete, deleting from pool")
		err = taskPool.DeleteTask(taskId)
		if err != nil {
			return lxerrors.New("deleting task from staging or pending pool after kill was requested", err)
		}
		err = sendTaskKilledStatus(state, mesosproto.TaskState_TASK_KILLED, tpiUrl, taskProviderId, taskId)
		if err != nil {
			return lxerrors.New("sending TASK_KILLED status to task provider"+taskProviderId, err)
		}
		return nil
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

func sendTaskKilledStatus(state *lxstate.State, taskState mesosproto.TaskState, tpiUrl, taskProviderId, taskId string) error {
	taskProvider, err := state.TaskProviderPool.GetTaskProvider(taskProviderId)
	if err != nil {
		return lxerrors.New("finding task provider for kill request", err)
	}
	taskKilledStatus := generateTaskStatus(taskId, taskState, "Kill Task was requested before task staging was complete")
	err = tpi_messenger.SendStatusUpdate(tpiUrl, taskProvider, taskKilledStatus)
	if err != nil {
		return lxerrors.New("udpating tpi with TASK_KILLED status for task before task staging was complete", err)
	}
	return nil
}
