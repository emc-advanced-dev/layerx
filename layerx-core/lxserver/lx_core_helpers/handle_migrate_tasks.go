package lx_core_helpers

import (
	"github.com/emc-advanced-dev/layerx/layerx-core/layerx_brain_client"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxstate"
	"github.com/layer-x/layerx-commons/lxerrors"
)

func MigrateTasks(state *lxstate.State, migrateTasksMessage layerx_brain_client.MigrateTaskMessage) error {
	for _, taskId := range migrateTasksMessage.TaskIds {
		err := migrateTask(state, migrateTasksMessage.DestinationNodeId, taskId)
		if err != nil {
			return lxerrors.New("migrating task "+taskId, err)
		}
	}
	return nil
}

func migrateTask(state *lxstate.State, nodeId, taskId string) error {
	sourceTaskPool, err := state.GetTaskPoolContainingTask(taskId)
	if err != nil {
		return lxerrors.New("retrieving task pool containing task "+taskId, err)
	}
	targetTaskPool, err := state.NodePool.GetNodeTaskPool(nodeId)
	if err != nil {
		return lxerrors.New("getting task pool for destination node "+nodeId, err)
	}
	if sourceTaskPool == targetTaskPool {
		return nil
	}
	if sourceTaskPool == state.PendingTaskPool || sourceTaskPool == state.StagingTaskPool {
		return lxerrors.New("cannot migrate task "+taskId+" that is still pending or staging", nil)
	}
	task, err := sourceTaskPool.GetTask(taskId)
	if err != nil {
		return lxerrors.New("getting task for task "+taskId, err)
	}

	task.Checkpointed = true
	err = sourceTaskPool.ModifyTask(taskId, task)
	if err != nil {
		return lxerrors.New("setting task checkpointed to TRUE", err)
	}

	err = KillTask(state, state.GetTpiUrl(), task.TaskProvider.Id, taskId)
	if err != nil {
		return lxerrors.New("killing task "+taskId, err)
	}

	err = sourceTaskPool.DeleteTask(taskId)
	if err != nil {
		return lxerrors.New("deleting source task "+taskId, err)
	}
	task.NodeId = nodeId
	err = state.StagingTaskPool.AddTask(task)
	if err != nil {
		return lxerrors.New("moving task "+taskId+" from source to staging pool with new node "+nodeId, err)
	}
	return nil
}
