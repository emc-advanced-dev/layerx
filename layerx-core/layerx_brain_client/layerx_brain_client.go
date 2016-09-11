package layerx_brain_client

import (
	"encoding/json"
	"fmt"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"
	"github.com/emc-advanced-dev/pkg/errors"
	"github.com/layer-x/layerx-commons/lxhttpclient"
	"github.com/mesos/mesos-go/mesosproto"
)

type LayerXBrainClient struct {
	CoreURL string
}

const (
	GetPendingTasks  = "/GetPendingTasks"
	GetStagingTasks  = "/GetStagingTasks"
	GetNodes         = "/GetNodes"
	GetStatusUpdates = "/GetStatusUpdates"
	AssignTasks      = "/AssignTasks"
	MigrateTasks     = "/MigrateTasks"
)

//call this method to see unassigned tasks
func (brainClient *LayerXBrainClient) GetPendingTasks() ([]*lxtypes.Task, error) {
	resp, data, err := lxhttpclient.Get(brainClient.CoreURL, GetPendingTasks, nil)
	if err != nil {
		return nil, errors.New("GETing tasks from LayerX core server", err)
	}
	if resp.StatusCode != 200 {
		msg := fmt.Sprintf("GETing tasks from LayerX core server; status code was %v, expected 200", resp.StatusCode)
		return nil, errors.New(msg, err)
	}
	var tasks []*lxtypes.Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		msg := fmt.Sprintf("unmarshalling data %s into task array", string(data))
		return nil, errors.New(msg, err)
	}
	fmt.Printf("\n\n\n\nTASKS: %v\n\n\n\n", tasks)
	return tasks, nil
}

//call this method to see unassigned tasks
func (brainClient *LayerXBrainClient) GetStagingTasks() ([]*lxtypes.Task, error) {
	resp, data, err := lxhttpclient.Get(brainClient.CoreURL, GetStagingTasks, nil)
	if err != nil {
		return nil, errors.New("GETing tasks from LayerX core server", err)
	}
	if resp.StatusCode != 200 {
		msg := fmt.Sprintf("GETing tasks from LayerX core server; status code was %v, expected 200", resp.StatusCode)
		return nil, errors.New(msg, err)
	}
	var tasks []*lxtypes.Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		msg := fmt.Sprintf("unmarshalling data %s into task array", string(data))
		return nil, errors.New(msg, err)
	}
	fmt.Printf("\n\n\n\nTASKS: %v\n\n\n\n", tasks)
	return tasks, nil
}

//call this method to see submitted nodes
//and their resources
func (brainClient *LayerXBrainClient) GetNodes() ([]*lxtypes.Node, error) {
	resp, data, err := lxhttpclient.Get(brainClient.CoreURL, GetNodes, nil)
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

//call this method to see most recent status updates
//for all tasks
func (brainClient *LayerXBrainClient) GetStatusUpdates() ([]*mesosproto.TaskStatus, error) {
	resp, data, err := lxhttpclient.Get(brainClient.CoreURL, GetStatusUpdates, nil)
	if err != nil {
		return nil, errors.New("GETing task statuses from LayerX core server", err)
	}
	if resp.StatusCode != 200 {
		msg := fmt.Sprintf("GETing task statuses from LayerX core server; status code was %v, expected 200", resp.StatusCode)
		return nil, errors.New(msg, err)
	}
	var taskStatuses []*mesosproto.TaskStatus
	err = json.Unmarshal(data, &taskStatuses)
	if err != nil {
		msg := fmt.Sprintf("unmarshalling data %s into task status array", string(data))
		return nil, errors.New(msg, err)
	}
	return taskStatuses, nil
}

//call this method to assign tasks to a node
func (brainClient *LayerXBrainClient) AssignTasks(nodeId string, taskIds ...string) error {
	assignTasksMessage := BrainAssignTasksMessage{
		NodeId:  nodeId,
		TaskIds: taskIds,
	}
	resp, _, err := lxhttpclient.Post(brainClient.CoreURL, AssignTasks, nil, assignTasksMessage)
	if err != nil {
		return errors.New("POSTing assignTasksMessage to LayerX core server", err)
	}
	if resp.StatusCode != 202 {
		msg := fmt.Sprintf("POSTing assignTasksMessage to LayerX core server; status code was %v, expected 202", resp.StatusCode)
		return errors.New(msg, err)
	}
	return nil
}

//call this method to migrate running tasks from one node to another
func (brainClient *LayerXBrainClient) MigrateTasks(destinationNodeId string, taskIds ...string) error {
	assignTasksMessage := MigrateTaskMessage{
		DestinationNodeId: destinationNodeId,
		TaskIds:           taskIds,
	}
	resp, _, err := lxhttpclient.Post(brainClient.CoreURL, MigrateTasks, nil, assignTasksMessage)
	if err != nil {
		return errors.New("POSTing assignTasksMessage to LayerX core server", err)
	}
	if resp.StatusCode != 202 {
		msg := fmt.Sprintf("POSTing assignTasksMessage to LayerX core server; status code was %v, expected 202", resp.StatusCode)
		return errors.New(msg, err)
	}
	return nil
}
