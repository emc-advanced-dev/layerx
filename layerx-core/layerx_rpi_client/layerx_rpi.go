package layerx_rpi_client

import (
	"encoding/json"
	"fmt"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"
	"github.com/layer-x/layerx-commons/lxerrors"
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
		return lxerrors.New("POSTing registration request to LayerX core server", err)
	}
	if resp.StatusCode != 202 {
		msg := fmt.Sprintf("POSTing registration request to LayerX core server; status code was %v, expected 202", resp.StatusCode)
		return lxerrors.New(msg, err)
	}
	return nil
}

//call this method when submitting
// a new resource from the rpi
func (rpi *LayerXRpi) SubmitResource(resource *lxtypes.Resource) error {
	resource.RpiName = rpi.RpiName
	resp, _, err := lxhttpclient.Post(rpi.CoreURL, SubmitResource, nil, resource)
	if err != nil {
		return lxerrors.New("POSTing resource to LayerX core server", err)
	}
	if resp.StatusCode != 202 {
		msg := fmt.Sprintf("POSTing resource to LayerX core server; status code was %v, expected 202", resp.StatusCode)
		return lxerrors.New(msg, err)
	}
	return nil
}

//call this method when submitting
// a status update from the rpi
func (rpi *LayerXRpi) SubmitStatusUpdate(status *mesosproto.TaskStatus) error {
	resp, _, err := lxhttpclient.Post(rpi.CoreURL, SubmitStatusUpdate, nil, status)
	if err != nil {
		return lxerrors.New("POSTing TaskStatus to LayerX core server", err)
	}
	if resp.StatusCode != 202 {
		msg := fmt.Sprintf("POSTing TaskStatus to LayerX core server; status code was %v, expected 202", resp.StatusCode)
		return lxerrors.New(msg, err)
	}
	return nil
}

//call this method to see submitted nodes
//and their resources
func (rpi *LayerXRpi) GetNodes() ([]*lxtypes.Node, error) {
	resp, data, err := lxhttpclient.Get(rpi.CoreURL, GetNodes, nil)
	if err != nil {
		return nil, lxerrors.New("GETing nodes from LayerX core server", err)
	}
	if resp.StatusCode != 200 {
		msg := fmt.Sprintf("GETing nodes from LayerX core server; status code was %v, expected 200", resp.StatusCode)
		return nil, lxerrors.New(msg, err)
	}
	var nodes []*lxtypes.Node
	err = json.Unmarshal(data, &nodes)
	if err != nil {
		msg := fmt.Sprintf("unmarshalling data %s into node array", string(data))
		return nil, lxerrors.New(msg, err)
	}
	return nodes, nil
}
