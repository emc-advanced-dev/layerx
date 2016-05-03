package lxstate

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/layer-x/layerx-commons/lxdatabase"
	"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/layer-x/layerx-core_v2/lxtypes"
	"github.com/mesos/mesos-go/mesosproto"
)

const (
	state_root            = "/state"
	nodes                 = state_root + "/nodes"
	pending_tasks         = state_root + "/pending_tasks"
	staging_tasks         = state_root + "/staging_tasks"
	task_providers        = state_root + "/task_providers"
	failed_task_providers = state_root + "/failed_task_providers"
	statuses              = state_root + "/statuses"
	tpi_url_key           = state_root + "/tpi_url"
	rpis                  = state_root + "/rpis"
)

type State struct {
	PendingTaskPool        *TaskPool
	StagingTaskPool        *TaskPool
	NodePool               *NodePool
	TaskProviderPool       *TaskProviderPool
	FailedTaskProviderPool *TaskProviderPool
	StatusPool             *StatusPool
	RpiPool                *RpiPool
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
		FailedTaskProviderPool: &TaskProviderPool{
			rootKey: failed_task_providers,
		},
		StatusPool: &StatusPool{
			rootKey: statuses,
		},
		RpiPool: &RpiPool{
			rootKey: rpis,
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
	state.FailedTaskProviderPool.Initialize()
	state.StatusPool.Initialize()
	state.RpiPool.Initialize()
	return nil
}

func (state *State) SetTpi(tpiUrl string) error {
	err := lxdatabase.Set(tpi_url_key, tpiUrl)
	if err != nil {
		return lxerrors.New("could not set tpi url", err)
	}
	return nil
}

func (state *State) GetTpi() (string, error) {
	tpiUrl, err := lxdatabase.Get(tpi_url_key)
	if err != nil {
		return "", lxerrors.New("could not get tpi url", err)
	}
	return tpiUrl, nil
}

func (state *State) GetAllTasks() (map[string]*lxtypes.Task, error) {
	allTasks := make(map[string]*lxtypes.Task)
	pendingTasks, err := state.PendingTaskPool.GetTasks()
	if err != nil {
		return nil, lxerrors.New("could not get tasks from pending task pool", err)
	}
	for _, task := range pendingTasks {
		allTasks[task.TaskId] = task
	}
	stagingTasks, err := state.StagingTaskPool.GetTasks()
	if err != nil {
		return nil, lxerrors.New("could not get tasks from staging task pool", err)
	}
	for _, task := range stagingTasks {
		allTasks[task.TaskId] = task
	}
	nodes, err := state.NodePool.GetNodes()
	if err != nil {
		return nil, lxerrors.New("getting list of nodes from node pool", err)
	}
	for _, node := range nodes {
		nodeId := node.Id
		nodeTaskPool, err := state.NodePool.GetNodeTaskPool(nodeId)
		if err != nil {
			return nil, lxerrors.New("getting task pool for node "+nodeId, err)
		}
		nodeTasks, err := nodeTaskPool.GetTasks()
		if err != nil {
			return nil, lxerrors.New("getting list of tasks from node "+nodeId+"task pool", err)
		}
		for _, task := range nodeTasks {
			allTasks[task.TaskId] = task
		}
	}
	return allTasks, nil
}

func (state *State) GetStatusUpdatesForTaskProvider(tpId string) (map[string]*mesosproto.TaskStatus, error) {
	taskProviders, err := state.TaskProviderPool.GetTaskProviders()
	if err != nil {
		return nil, lxerrors.New("could not get task provider list", err)
	}
	tpIds := []string{}
	for tpId := range taskProviders {
		tpIds = append(tpIds, tpId)
	}
	allTasks, err := state.GetAllTasks()
	if err != nil {
		return nil, lxerrors.New("getting all tasks from state", err)
	}
	targetTaskIds := []string{}
	for _, task := range allTasks {
		if containsString(tpIds, task.TaskProvider.Id) {
			targetTaskIds = append(targetTaskIds, task.TaskId)
		}
	}
	allStatuses, err := state.StatusPool.GetStatuses()
	if err != nil {
		return nil, lxerrors.New("getting all statuses from state", err)
	}
	statuses := make(map[string]*mesosproto.TaskStatus)
	for _, status := range allStatuses {
		if containsString(targetTaskIds, status.GetTaskId().GetValue()) {
			statuses[status.GetTaskId().GetValue()] = status
		}
	}
	return statuses, nil
}

func (state *State) GetStatusUpdates() (map[string]*mesosproto.TaskStatus, error) {
	allTasks, err := state.GetAllTasks()
	if err != nil {
		return nil, lxerrors.New("getting all tasks from state", err)
	}
	targetTaskIds := []string{}
	for _, task := range allTasks {
		targetTaskIds = append(targetTaskIds, task.TaskId)
	}
	allStatuses, err := state.StatusPool.GetStatuses()
	if err != nil {
		return nil, lxerrors.New("getting all statuses from state", err)
	}
	statuses := make(map[string]*mesosproto.TaskStatus)
	for _, status := range allStatuses {
		if containsString(targetTaskIds, status.GetTaskId().GetValue()) {
			statuses[status.GetTaskId().GetValue()] = status
		}
	}
	return statuses, nil
}

func (state *State) GetTaskFromAnywhere(taskId string) (*lxtypes.Task, error) {
	allTasks, err := state.GetAllTasks()
	if err != nil {
		return nil, lxerrors.New("could not get all tasks", err)
	}
	for _, task := range allTasks {
		if task.TaskId == taskId {
			return task, nil
		}
	}
	return nil, lxerrors.New("task was not found with id "+taskId, nil)
}

func (state *State) GetTaskPoolContainingTask(taskId string) (*TaskPool, error) {
	pendingTasks, err := state.PendingTaskPool.GetTasks()
	if err != nil {
		return nil, lxerrors.New("could not get tasks from pending task pool", err)
	}
	for _, task := range pendingTasks {
		if task.TaskId == taskId {
			return state.PendingTaskPool, nil
		}
	}
	stagingTasks, err := state.StagingTaskPool.GetTasks()
	if err != nil {
		return nil, lxerrors.New("could not get tasks from staging task pool", err)
	}
	for _, task := range stagingTasks {
		if task.TaskId == taskId {
			return state.StagingTaskPool, nil
		}
	}
	nodes, err := state.NodePool.GetNodes()
	if err != nil {
		return nil, lxerrors.New("getting list of nodes from node pool", err)
	}
	for _, node := range nodes {
		nodeId := node.Id
		nodeTaskPool, err := state.NodePool.GetNodeTaskPool(nodeId)
		if err != nil {
			return nil, lxerrors.New("getting task pool for node "+nodeId, err)
		}
		nodeTasks, err := nodeTaskPool.GetTasks()
		if err != nil {
			return nil, lxerrors.New("getting list of tasks from node "+nodeId+"task pool", err)
		}
		for _, task := range nodeTasks {
			if task.TaskId == taskId {
				return nodeTaskPool, nil
			}
		}
	}
	return nil, lxerrors.New("task pool not found that contains task "+taskId, nil)
}

func (state *State) GetTpiUrl() string {
	for {
		tpiUrl, err := state.GetTpi()
		if err != nil {
			logrus.WithFields(logrus.Fields{"err": err}).Errorf("Failed to retrieve rpis from state")
		} else {
			return tpiUrl
		}
		if tpiUrl != "" {
			logrus.WithFields(logrus.Fields{
				"tpiUrl": tpiUrl,
			}).Infof("TPI registered...")
			return tpiUrl
		} else {
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func (state *State) GetRpiUrls() []string {
	rpiUrls := []string{}
	for {
		rpis, err := state.RpiPool.GetRpis()
		if err != nil {
			logrus.WithError(err).Warnf("Failed to retrieve rpis from state")
		} else {
			for _, rpi := range rpis {
				rpiUrls = append(rpiUrls, rpi.Url)
			}
		}
		if len(rpiUrls) > 0 {
			return rpiUrls
		} else {
			logrus.Warnf("no RPIs registered...")
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func containsString(strArray []string, target string) bool {
	for _, str := range strArray {
		if str == target {
			return true
		}
	}
	return false
}
