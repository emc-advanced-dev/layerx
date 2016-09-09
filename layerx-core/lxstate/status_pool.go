package lxstate
import (
	"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/layer-x/layerx-commons/lxdatabase"
	"encoding/json"
	"github.com/mesos/mesos-go/mesosproto"
)

type StatusPool struct {
	rootKey string
}

func (statusPool *StatusPool) GetKey() string {
	return statusPool.rootKey
}

func (statusPool *StatusPool) Initialize() error {
	err := lxdatabase.Mkdir(statusPool.GetKey())
	if err != nil {
		return lxerrors.New("initializing "+statusPool.GetKey() +" directory", err)
	}
	return nil
}

func (statusPool *StatusPool) AddStatus(status *mesosproto.TaskStatus) error {
	statusId := status.GetTaskId().GetValue()
	_, err := statusPool.GetStatus(statusId)
	if err == nil {
		return statusPool.ModifyStatus(statusId, status)
	}
	statusData, err := json.Marshal(status)
	if err != nil {
		return lxerrors.New("could not marshal status to json", err)
	}
	err = lxdatabase.Set(statusPool.GetKey()+"/"+statusId, string(statusData))
	if err != nil {
		return lxerrors.New("setting key/value pair for status", err)
	}
	return nil
}

func (statusPool *StatusPool) GetStatus(statusId string) (*mesosproto.TaskStatus, error) {
	statusJson, err := lxdatabase.Get(statusPool.GetKey()+"/"+statusId)
	if err != nil {
		return nil, lxerrors.New("retrieving status "+statusId+" from database", err)
	}
	var status mesosproto.TaskStatus
	err = json.Unmarshal([]byte(statusJson), &status)
	if err != nil {
		return nil, lxerrors.New("unmarshalling json into Status struct", err)
	}
	return &status, nil
}

func (statusPool *StatusPool) ModifyStatus(statusId string, modifiedStatus *mesosproto.TaskStatus) error {
	_, err := statusPool.GetStatus(statusId)
	if err != nil {
		return lxerrors.New("status "+statusId+" not found", err)
	}
	statusData, err := json.Marshal(modifiedStatus)
	if err != nil {
		return lxerrors.New("could not marshal modified status to json", err)
	}
	err = lxdatabase.Set(statusPool.GetKey()+"/"+statusId, string(statusData))
	if err != nil {
		return lxerrors.New("setting key/value pair for modified status", err)
	}
	return nil

}

func (statusPool *StatusPool) GetStatuses() (map[string]*mesosproto.TaskStatus, error) {
	statuses := make(map[string]*mesosproto.TaskStatus)
	knownStatuses, err := lxdatabase.GetKeys(statusPool.GetKey())
	if err != nil {
		return nil, lxerrors.New("retrieving list of known statuses", err)
	}
	for _, statusJson := range knownStatuses {
		var status mesosproto.TaskStatus
		err = json.Unmarshal([]byte(statusJson), &status)
		if err != nil {
			return nil, lxerrors.New("unmarshalling json into Status struct", err)
		}
		statuses[status.GetTaskId().GetValue()] = &status
	}
	return statuses, nil
}

func (statusPool *StatusPool) DeleteStatus(statusId string) error {
	_, err := statusPool.GetStatus(statusId)
	if err != nil {
		return lxerrors.New("status "+statusId+" not found", err)
	}
	err = lxdatabase.Rm(statusPool.GetKey()+"/"+statusId)
	if err != nil {
		return lxerrors.New("removing status "+statusId+" from database", err)
	}
	return nil
}