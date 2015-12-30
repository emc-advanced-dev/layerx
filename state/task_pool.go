package state
import (
	"github.com/layer-x/layerx-core_v2/lxtypes"
	"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/layer-x/layerx-commons/lxdatabase"
	"encoding/json"
)

type TaskPool struct {
	rootKey string
}

func (taskPool *TaskPool) GetKey() string {
	return taskPool.rootKey
}

func (taskPool *TaskPool) Initialize() error {
	err := lxdatabase.Mkdir(taskPool.rootKey)
	if err != nil {
		return lxerrors.New("initializing "+taskPool.rootKey +" directory", err)
	}
	return nil
}

func (taskPool *TaskPool) AddTask(task *lxtypes.Task) error {
	taskId := task.TaskId
	_, err := taskPool.GetTask(taskId)
	if err == nil {
		return lxerrors.New("task "+taskId+" already exists in database, try Modify()?", err)
	}
	taskData, err := json.Marshal(task)
	if err != nil {
		return lxerrors.New("could not marshal task to json", err)
	}
	err = lxdatabase.Set(taskPool.rootKey+"/"+taskId, string(taskData))
	if err != nil {
		return lxerrors.New("setting key/value pair for task", err)
	}
	return nil
}

func (taskPool *TaskPool) GetTask(taskId string) (*lxtypes.Task, error) {
	taskJson, err := lxdatabase.Get(taskPool.rootKey+"/"+taskId)
	if err != nil {
		return nil, lxerrors.New("retrieving task "+taskId+" from database", err)
	}
	var task lxtypes.Task
	err = json.Unmarshal([]byte(taskJson), &task)
	if err != nil {
		return nil, lxerrors.New("unmarshalling json into Task struct", err)
	}
	return &task, nil
}
func (taskPool *TaskPool) ModifyTask(taskId string, modifiedTask *lxtypes.Task) error {
	return lxerrors.New("not implemented", nil)
}

func (taskPool *TaskPool) GetTasks() (map[string]*lxtypes.Task, error) {
	return nil, lxerrors.New("not implemented", nil)
}