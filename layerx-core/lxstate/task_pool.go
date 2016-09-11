package lxstate

import (
	"encoding/json"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"
	"github.com/emc-advanced-dev/pkg/errors"
	"github.com/layer-x/layerx-commons/lxdatabase"
)

type TaskPool struct {
	rootKey string
}

func (taskPool *TaskPool) GetKey() string {
	return taskPool.rootKey
}

func (taskPool *TaskPool) Initialize() error {
	err := lxdatabase.Mkdir(taskPool.GetKey())
	if err != nil {
		return errors.New("initializing "+taskPool.GetKey()+" directory", err)
	}
	return nil
}

func (taskPool *TaskPool) AddTask(task *lxtypes.Task) error {
	if task.TaskProvider == nil {
		return errors.New("cannot accept task "+task.TaskId+" with no task provider!", nil)
	}
	taskId := task.TaskId
	_, err := taskPool.GetTask(taskId)
	if err == nil {
		return errors.New("task "+taskId+" already exists in database, try Modify()?", err)
	}
	taskData, err := json.Marshal(task)
	if err != nil {
		return errors.New("could not marshal task to json", err)
	}
	err = lxdatabase.Set(taskPool.GetKey()+"/"+taskId, string(taskData))
	if err != nil {
		return errors.New("setting key/value pair for task", err)
	}
	return nil
}

func (taskPool *TaskPool) GetTask(taskId string) (*lxtypes.Task, error) {
	taskJson, err := lxdatabase.Get(taskPool.GetKey() + "/" + taskId)
	if err != nil {
		return nil, errors.New("retrieving task "+taskId+" from database", err)
	}
	var task lxtypes.Task
	err = json.Unmarshal([]byte(taskJson), &task)
	if err != nil {
		return nil, errors.New("unmarshalling json into Task struct", err)
	}
	return &task, nil
}

func (taskPool *TaskPool) ModifyTask(taskId string, modifiedTask *lxtypes.Task) error {
	_, err := taskPool.GetTask(taskId)
	if err != nil {
		return errors.New("task "+taskId+" not found", err)
	}
	taskData, err := json.Marshal(modifiedTask)
	if err != nil {
		return errors.New("could not marshal modified task to json", err)
	}
	err = lxdatabase.Set(taskPool.GetKey()+"/"+taskId, string(taskData))
	if err != nil {
		return errors.New("setting key/value pair for modified task", err)
	}
	return nil

}

func (taskPool *TaskPool) GetTasks() (map[string]*lxtypes.Task, error) {
	tasks := make(map[string]*lxtypes.Task)
	knownTasks, err := lxdatabase.GetKeys(taskPool.GetKey())
	if err != nil {
		return nil, errors.New("retrieving list of known tasks", err)
	}
	for _, taskJson := range knownTasks {
		var task lxtypes.Task
		err = json.Unmarshal([]byte(taskJson), &task)
		if err != nil {
			return nil, errors.New("unmarshalling json into Task struct", err)
		}
		tasks[task.TaskId] = &task
	}
	return tasks, nil
}

func (taskPool *TaskPool) DeleteTask(taskId string) error {
	_, err := taskPool.GetTask(taskId)
	if err != nil {
		return errors.New("task "+taskId+" not found", err)
	}
	err = lxdatabase.Rm(taskPool.GetKey() + "/" + taskId)
	if err != nil {
		return errors.New("removing task "+taskId+" from database", err)
	}
	return nil
}
