package layerx_tpi

import (
	"encoding/json"
	"fmt"
	"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/layer-x/layerx-commons/lxhttpclient"
	"github.com/layer-x/layerx-core_v2/lxtypes"
	"github.com/mesos/mesos-go/mesosproto"
)

type LayerXTpi struct {
	CoreURL string
}

const (
	RegisterTaskProvider   = "/RegisterTaskProvider"
	DeregisterTaskProvider = "/DeregisterTaskProvider"
	GetTaskProviders       = "/GetTaskProviders"
	GetStatusUpdates       = "/GetStatusUpdates"
	SubmitTask             = "/SubmitTask"
	KillTask               = "/KillTask"
	PurgeTask              = "/PurgeTask"
)

//call this method when registering
// a new task provider to the tpi
func (tpi *LayerXTpi) RegisterTaskProvider(tp *lxtypes.TaskProvider) error {
	resp, _, err := lxhttpclient.Post(tpi.CoreURL, RegisterTaskProvider, nil, tp)
	if err != nil {
		return lxerrors.New("POSTing TaskProvider to LayerX core server", err)
	}
	if resp.StatusCode != 202 {
		msg := fmt.Sprintf("POSTing TaskProvider to LayerX core server; status code was %v, expected 202", resp.StatusCode)
		return lxerrors.New(msg, err)
	}
	return nil
}

//call this method when a non-failover task provider
// disconnects from the tpi
func (tpi *LayerXTpi) DeregisterTaskProvider(tpId string) error {
	resp, _, err := lxhttpclient.Post(tpi.CoreURL, DeregisterTaskProvider+"/"+tpId, nil, nil)
	if err != nil {
		return lxerrors.New("Requesting DeRegister of task provider "+tpId+" to LayerX core server", err)
	}
	if resp.StatusCode != 202 {
		msg := fmt.Sprintf("Requesting DeRegister of task provider "+tpId+" to LayerX core server; status code was %v, expected 202", resp.StatusCode)
		return lxerrors.New(msg, err)
	}
	return nil
}

//call this method to retrieve a
//specific task provider for the id
func (tpi *LayerXTpi) GetTaskProvider(tpid string) (*lxtypes.TaskProvider, error) {
	taskProviders, err := tpi.GetTaskProviders()
	if err != nil {
		return nil, lxerrors.New("retrieving list of task providers", err)
	}
	for _, taskProvider := range taskProviders {
		if taskProvider.Id == tpid {
			return taskProvider, nil
		}
	}
	return nil, lxerrors.New("task provider with id "+tpid+" not found", nil)
}

//call this method to retrieve a list of registered
//task providers. e.g. before polling task providers
//for pending tasks
func (tpi *LayerXTpi) GetTaskProviders() ([]*lxtypes.TaskProvider, error) {
	taskProviders := []*lxtypes.TaskProvider{}
	resp, body, err := lxhttpclient.Get(tpi.CoreURL, GetTaskProviders, nil)
	if err != nil {
		return nil, lxerrors.New("Requesting task provider list from LayerX core server", err)
	}
	if resp.StatusCode != 200 {
		msg := fmt.Sprintf("Requesting task provider list from LayerX core server; status code was %v, expected 202", resp.StatusCode)
		return nil, lxerrors.New(msg, err)
	}
	err = json.Unmarshal(body, &taskProviders)
	if err != nil {
		return nil, lxerrors.New("unmarshalling json to []*lxtypes.TaskProvider", err)
	}
	return taskProviders, nil
}

//call this method to retrieve a list of
// the most recent status updates for tasks
func (tpi *LayerXTpi) GetStatusUpdates() ([]*mesosproto.TaskStatus, error) {
	statusUpdates := []*mesosproto.TaskStatus{}
	resp, body, err := lxhttpclient.Get(tpi.CoreURL, GetStatusUpdates, nil)
	if err != nil {
		return nil, lxerrors.New("Requesting status update list from LayerX core server", err)
	}
	if resp.StatusCode != 200 {
		msg := fmt.Sprintf("Requesting status update list from LayerX core server; status code was %v, expected 202", resp.StatusCode)
		return nil, lxerrors.New(msg, err)
	}
	err = json.Unmarshal(body, &statusUpdates)
	if err != nil {
		return nil, lxerrors.New("unmarshalling json to []*mesosproto.TaskStatus", err)
	}
	return statusUpdates, nil
}

//call this method to submit
// a requested task to layer-x
func (tpi *LayerXTpi) SubmitTask(tpid string, task *lxtypes.Task) error {
	resp, _, err := lxhttpclient.Post(tpi.CoreURL, SubmitTask+"/"+tpid, nil, task)
	if err != nil {
		return lxerrors.New("POSTing Task to LayerX core server", err)
	}
	if resp.StatusCode != 202 {
		msg := fmt.Sprintf("POSTing Task to LayerX core server; status code was %v, expected 202", resp.StatusCode)
		return lxerrors.New(msg, err)
	}
	return nil
}

//call this method to submit
// a requested task to layer-x
func (tpi *LayerXTpi) KillTask(taskId string) error {
	resp, _, err := lxhttpclient.Post(tpi.CoreURL, KillTask+"/"+taskId, nil, nil)
	if err != nil {
		return lxerrors.New("Requesting KILL on task "+taskId+" to LayerX core server", err)
	}
	if resp.StatusCode != 202 {
		msg := fmt.Sprintf("Requesting KILL on task "+taskId+" to LayerX core server; status code was %v, expected 202", resp.StatusCode)
		return lxerrors.New(msg, err)
	}
	return nil
}

//call this method to purge
// a requested task from layer-x
func (tpi *LayerXTpi) PurgeTask(taskId string) error {
	resp, _, err := lxhttpclient.Post(tpi.CoreURL, PurgeTask+"/"+taskId, nil, nil)
	if err != nil {
		return lxerrors.New("Requesting Purge on task "+taskId+" to LayerX core server", err)
	}
	if resp.StatusCode != 202 {
		msg := fmt.Sprintf("Requesting Purge on task "+taskId+" to LayerX core server; status code was %v, expected 202", resp.StatusCode)
		return lxerrors.New(msg, err)
	}
	return nil
}
