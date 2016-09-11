package lx_core_helpers

import (
	"github.com/emc-advanced-dev/layerx/layerx-core/lxstate"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"
	"github.com/emc-advanced-dev/pkg/errors"
)

func GetStagingTasks(state *lxstate.State) ([]*lxtypes.Task, error) {
	taskMap, err := state.StagingTaskPool.GetTasks()
	if err != nil {
		return nil, errors.New("getting list of tasks from staging task pool", err)
	}
	tasks := []*lxtypes.Task{}
	for _, task := range taskMap {
		tasks = append(tasks, task)
	}
	return tasks, nil
}
