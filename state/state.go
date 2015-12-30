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
}

func NewState() *State {
	return &State{
		PendingTaskPool: &TaskPool{
			rootKey: pending_tasks,
		},
		StagingTaskPool: &TaskPool{
			rootKey: staging_tasks,
		},
	}
}

func (state *State) InitializeState(etcdUrl string) error {
	err := lxdatabase.Init([]string{etcdUrl})
	if err != nil {
		return lxerrors.New("initializing etcd client failed", err)
	}
	lxdatabase.Mkdir(state_root)
	rootContents, err := lxdatabase.GetSubdirectories(state_root)
	if err != nil {
		return lxerrors.New("retrieving contents of state root dir", err)
	}
	state.PendingTaskPool.Initialize()
	state.StagingTaskPool.Initialize()
	err = initializeDirectoriesIfNotFound(rootContents, nodes, task_providers, statuses)
	if err != nil {
		return lxerrors.New("initializing state directories", err)
	}
	return nil
}

func initializeDirectoriesIfNotFound(rootContents []string, directoryNames ...string) error {
	for _, directoryName := range directoryNames {
		if !contains(rootContents, directoryName) {
			err := lxdatabase.Mkdir(directoryName)
			if err != nil {
				return lxerrors.New("initializing "+directoryName+" directory", err)
			}
		}
	}
	return nil
}

func contains(strArray []string, desired string) bool {
	for _, str := range strArray {
		if str == desired {
			return true
		}
	}
	return false
}