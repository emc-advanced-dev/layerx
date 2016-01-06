package lx_core_helpers
import (
	"github.com/layer-x/layerx-core_v2/lxstate"
	"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/layer-x/layerx-core_v2/layerx_brain_client"
	"github.com/layer-x/layerx-commons/lxlog"
	"github.com/Sirupsen/logrus"
)

func AssignTasks(state *lxstate.State, assignTasksMessage layerx_brain_client.BrainAssignTasksMessage) error {
	for _, taskId := range assignTasksMessage.TaskIds {
		err := assignTask(state, assignTasksMessage.NodeId, taskId)
		if err != nil {
			return lxerrors.New("assigning task "+taskId, err)
		}
	}
	return nil
}

func assignTask(state *lxstate.State, nodeId, taskId string) error {
	task, err := state.PendingTaskPool.GetTask(taskId)
	if err != nil {
		return lxerrors.New("retrieving task "+taskId+" from pending task pool", err)
	}
	err = state.PendingTaskPool.DeleteTask(taskId)
	if err != nil {
		return lxerrors.New("deleting task "+taskId+" from pending task pool", err)
	}
	task.SlaveId = nodeId
	err = state.StagingTaskPool.AddTask(task)
	if err != nil {
		return lxerrors.New("moving task "+taskId+" from pending task pool to staging task pool", err)
	}
	tasks, _ := state.StagingTaskPool.GetTasks()
	lxlog.Debugf(logrus.Fields{"staging_tasks": tasks, "task": task, "nodeId": nodeId}, "moved task into staging task pool")
	return nil
}