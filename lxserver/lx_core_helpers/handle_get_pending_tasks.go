package lx_core_helpers
import (
	"github.com/layer-x/layerx-core_v2/lxstate"
	"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/layer-x/layerx-core_v2/lxtypes"
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
