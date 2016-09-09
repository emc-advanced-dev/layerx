package lx_core_helpers
import (
	"github.com/emc-advanced-dev/layerx/layerx-core/lxstate"
	"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"
)

func  GetPendingTasks(state *lxstate.State) ([]*lxtypes.Task, error) {
	taskMap, err := state.PendingTaskPool.GetTasks()
	if err != nil {
		return nil, lxerrors.New("getting list of tasks from pending task pool", err)
	}
	tasks := []*lxtypes.Task{}
	for _, task := range taskMap {
		tasks = append(tasks, task)
	}
	return tasks, nil
}
