package layerx_rpi_client
import (
	"github.com/layer-x/layerx-core_v2/lxtypes"
	"github.com/layer-x/layerx-commons/lxhttpclient"
	"github.com/layer-x/layerx-commons/lxerrors"
	"fmt"
	"github.com/mesos/mesos-go/mesosproto"
	"encoding/json"
)


type LayerXRpi struct {
	CoreURL string
}

const (
	RegisterRpi             = "/RegisterRpi"
	SubmitResource             = "/SubmitResource"
	SubmitStatusUpdate         = "/SubmitStatusUpdate"
	GetNodes         = "/GetNodes"
)

//call this method to register the RPI
//with layerx
func (rpi *LayerXRpi) RegisterRpi(rpiUrl string) error {
	reg := RpiRegistrationMessage{RpiUrl: rpiUrl}
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
	var jNodes []jsonNode
	err = json.Unmarshal(data, &jNodes)
	if err != nil {
		msg := fmt.Sprintf("unmarshalling data %s into node array", string(data))
		return nil, lxerrors.New(msg, err)
	}
	var nodes []*lxtypes.Node
	for _, jNode := range jNodes {
		node, err := jNode.toRealNode()
		if err != nil {
			return nil, lxerrors.New("could not convert json node to real node", err)
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

type jsonNode struct {
	Id string    `json:"id"`
	Resources    map[string]*lxtypes.Resource `json:"resources"`
	RunningTasks map[string]*lxtypes.Task `json:"tasks"`
}

func (jn *jsonNode) toRealNode() (*lxtypes.Node, error) {
	node := lxtypes.NewNode(jn.Id)
	for _, resource := range jn.Resources {
		err := node.AddResource(resource)
		if err != nil {
			return nil, lxerrors.New("unable to add resource to converted node", err)
		}
	}
	for _, task := range jn.RunningTasks {
		err := node.AddTask(task)
		if err != nil {
			return nil, lxerrors.New("unable to add task to converted node", err)
		}
	}
	return node, nil
}