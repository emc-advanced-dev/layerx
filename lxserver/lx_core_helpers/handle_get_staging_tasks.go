package lx_core_helpers
import (
	"github.com/emc-advanced-dev/layerx-core/lxstate"
	"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/emc-advanced-dev/layerx-core/lxtypes"
)

func  GetStagingTasks(state *lxstate.State) ([]*lxtypes.Task, error) {
	taskMap, err := state.StagingTaskPool.GetTasks()
	if err != nil {
		return nil, lxerrors.New("getting list of tasks from staging task pool", err)
	}
	tasks := []*lxtypes.Task{}
	for _, task := range taskMap {
		tasks = append(tasks, task)
	}
	return tasks, nil
}
