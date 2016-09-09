package lx_core_helpers

import (
	"github.com/Sirupsen/logrus"
	"github.com/emc-advanced-dev/layerx/layerx-core/layerx_brain_client"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxstate"
	"github.com/layer-x/layerx-commons/lxerrors"
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
	if _, err := state.NodePool.GetNode(nodeId); err != nil {
		return lxerrors.New("retrieving node "+nodeId, err)
	}
	task, err := state.PendingTaskPool.GetTask(taskId)
	if err != nil {
		return lxerrors.New("retrieving task "+taskId+" from pending task pool", err)
	}
	err = state.PendingTaskPool.DeleteTask(taskId)
	if err != nil {
		return lxerrors.New("deleting task "+taskId+" from pending task pool", err)
	}
	task.NodeId = nodeId
	err = state.StagingTaskPool.AddTask(task)
	if err != nil {
		return lxerrors.New("moving task "+taskId+" from pending task pool to staging task pool", err)
	}
	tasks, _ := state.StagingTaskPool.GetTasks()
	logrus.WithFields(logrus.Fields{"staging_tasks": tasks, "task": task, "nodeId": nodeId}).Debugf("moved task into staging task pool")
	return nil
}
