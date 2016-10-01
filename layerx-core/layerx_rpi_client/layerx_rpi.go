package layerx_rpi_client

import (
	"encoding/json"
	"fmt"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"
	"github.com/emc-advanced-dev/pkg/errors"
	"github.com/layer-x/layerx-commons/lxhttpclient"
	"github.com/mesos/mesos-go/mesosproto"
)

type LayerXRpi struct {
	CoreURL string
	RpiName string
}

const (
	RegisterRpi        = "/RegisterRpi"
	SubmitResource     = "/SubmitResource"
	RescindResource    = "/RescindResource"
	SubmitStatusUpdate = "/SubmitStatusUpdate"
	GetNodes           = "/GetNodes"
)

//call this method to register the RPI
//with layerx
func (rpi *LayerXRpi) RegisterRpi(name, rpiUrl string) error {
	reg := RpiInfo{
		Name: name,
		Url:  rpiUrl,
	}
	resp, _, err := lxhttpclient.Post(rpi.CoreURL, RegisterRpi, nil, reg)
	if err != nil {
		return errors.New("POSTing registration request to LayerX core server", err)
	}
	if resp.StatusCode != 202 {
		msg := fmt.Sprintf("POSTing registration request to LayerX core server; status code was %v, expected 202", resp.StatusCode)
		return errors.New(msg, err)
	}
	return nil
}

//call this method when submitting
// a new resource from the rpi
func (rpi *LayerXRpi) SubmitResource(resource *lxtypes.Resource) error {
	resource.RpiName = rpi.RpiName
	resp, _, err := lxhttpclient.Post(rpi.CoreURL, SubmitResource, nil, resource)
	if err != nil {
		return errors.New("POSTing resource to LayerX core server", err)
	}
	if resp.StatusCode != 202 {
		msg := fmt.Sprintf("POSTing resource to LayerX core server; status code was %v, expected 202", resp.StatusCode)
		return errors.New(msg, err)
	}
	return nil
}

//call this method if a resource has been rescinded
// or is no longer valid / available
func (rpi *LayerXRpi) RescindResource(resourceID string) error {
	resp, _, err := lxhttpclient.Post(rpi.CoreURL, RescindResource, nil, resourceID)
	if err != nil {
		return errors.New("Rescinding resource from LayerX core server", err)
	}
	if resp.StatusCode != 204 {
		msg := fmt.Sprintf("Rescinding to LayerX core server; status code was %v, expected 204", resp.StatusCode)
		return errors.New(msg, err)
	}
	return nil
}

//call this method when submitting
// a status update from the rpi
func (rpi *LayerXRpi) SubmitStatusUpdate(status *mesosproto.TaskStatus) error {
	resp, _, err := lxhttpclient.Post(rpi.CoreURL, SubmitStatusUpdate, nil, status)
	if err != nil {
		return errors.New("POSTing TaskStatus to LayerX core server", err)
	}
	if resp.StatusCode != 202 {
		msg := fmt.Sprintf("POSTing TaskStatus to LayerX core server; status code was %v, expected 202", resp.StatusCode)
		return errors.New(msg, err)
	}
	return nil
}

//call this method to see submitted nodes
//and their resources + tasks
//call this to see tasks that the specific RPI should care about.
//useful for knowing what tasks we should get status updates for
//TODO: flip all these things so that it's always one-way
// (i.e. the RPI should always be told, never have to ask)
func (rpi *LayerXRpi) GetNodes() ([]*lxtypes.Node, error) {
	resp, data, err := lxhttpclient.Get(rpi.CoreURL, GetNodes, nil)
	if err != nil {
		return nil, errors.New("GETing nodes from LayerX core server", err)
	}
	if resp.StatusCode != 200 {
		msg := fmt.Sprintf("GETing nodes from LayerX core server; status code was %v, expected 200", resp.StatusCode)
		return nil, errors.New(msg, err)
	}
	var nodes []*lxtypes.Node
	err = json.Unmarshal(data, &nodes)
	if err != nil {
		msg := fmt.Sprintf("unmarshalling data %s into node array", string(data))
		return nil, errors.New(msg, err)
	}
	return nodes, nil
}
