package state
import (
	"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/layer-x/layerx-commons/lxdatabase"
)

const (
	state_root = "/state"
	nodes = state_root + "/nodes"
	pending_tasks = state_root + "/pending_tasks"
	staging_tasks = state_root + "/staging_tasks"
	task_providers = state_root + "/task_providers"
	statuses = state_root + "/statuses"
)

type State struct {
	PendingTaskPool *TaskPool
	StagingTaskPool *TaskPool
	NodePool	*NodePool
	TaskProviderPool *TaskProviderPool
	StatusPool *StatusPool
}

func NewState() *State {
	return &State{
		PendingTaskPool: &TaskPool{
			rootKey: pending_tasks,
		},
		StagingTaskPool: &TaskPool{
			rootKey: staging_tasks,
		},
		NodePool: &NodePool{
			rootKey: nodes,
		},
		TaskProviderPool: &TaskProviderPool{
			rootKey: task_providers,
		},
		StatusPool: &StatusPool{
			rootKey: statuses,
		},
	}
}

func (state *State) InitializeState(etcdUrl string) error {
	err := lxdatabase.Init([]string{etcdUrl})
	if err != nil {
		return lxerrors.New("initializing etcd client failed", err)
	}
	lxdatabase.Mkdir(state_root)
	state.PendingTaskPool.Initialize()
	state.StagingTaskPool.Initialize()
	state.NodePool.Initialize()
	state.TaskProviderPool.Initialize()
	state.StatusPool.Initialize()
	return nil
}
