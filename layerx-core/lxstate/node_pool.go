package lxstate

import (
	"github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"
	"github.com/layer-x/layerx-commons/lxdatabase"
	"github.com/layer-x/layerx-commons/lxerrors"
	"strings"
)

type NodePool struct {
	rootKey string
}

func (nodePool *NodePool) GetKey() string {
	return nodePool.rootKey
}

func (nodePool *NodePool) Initialize() error {
	err := lxdatabase.Mkdir(nodePool.GetKey())
	if err != nil {
		return lxerrors.New("initializing "+nodePool.GetKey()+" directory", err)
	}
	return nil
}

func (nodePool *NodePool) AddNode(node *lxtypes.Node) error {
	nodeId := node.Id
	_, err := nodePool.GetNode(nodeId)
	if err == nil {
		return lxerrors.New("node "+nodeId+" already exists in database", err)
	}
	err = nodePool.saveNode(node)
	if err != nil {
		return lxerrors.New("saving node to database", err)
	}
	return nil
}

func (nodePool *NodePool) GetNode(nodeId string) (*lxtypes.Node, error) {
	_, err := lxdatabase.GetSubdirectories(nodePool.GetKey() + "/" + nodeId)
	if err != nil {
		return nil, lxerrors.New("retrieving node "+nodeId+" from database", err)
	}
	node, err := nodePool.loadNode(nodeId)
	if err != nil {
		return nil, lxerrors.New("loading node "+nodeId+" from stored information in database", err)
	}
	return node, nil
}

func (nodePool *NodePool) GetNodes() (map[string]*lxtypes.Node, error) {
	nodeKeys, err := lxdatabase.GetSubdirectories(nodePool.GetKey())
	if err != nil {
		return nil, lxerrors.New("retrieving node map from database", err)
	}
	nodes := make(map[string]*lxtypes.Node)
	for _, nodeKey := range nodeKeys {
		nodeId := strings.TrimPrefix(nodeKey, nodePool.GetKey()+"/")
		node, err := nodePool.loadNode(nodeId)
		if err != nil {
			return nil, lxerrors.New("loading node "+nodeId+" from stored information in database", err)
		}
		nodes[nodeId] = node
	}
	return nodes, nil
}

func (nodePool *NodePool) DeleteNode(nodeId string) error {
	_, err := lxdatabase.GetSubdirectories(nodePool.GetKey() + "/" + nodeId)
	if err != nil {
		return lxerrors.New("retrieving node "+nodeId+" from database", err)
	}
	err = lxdatabase.Rmdir(nodePool.GetKey()+"/"+nodeId, true)
	if err != nil {
		return lxerrors.New("recursivey removing directory "+nodePool.GetKey()+"/"+nodeId+" from database", err)
	}
	return nil
}

func (nodePool *NodePool) GetNodeResourcePool(nodeId string) (*ResourcePool, error) {
	if _, err := nodePool.GetNode(nodeId); err != nil {
		return nil, lxerrors.New("could not find node", err)
	}
	return &ResourcePool{
		nodeId:  nodeId,
		rootKey: nodePool.GetKey() + "/" + nodeId + "/resources",
	}, nil
}

func (nodePool *NodePool) GetNodeTaskPool(nodeId string) (*TaskPool, error) {
	if _, err := nodePool.GetNode(nodeId); err != nil {
		return nil, lxerrors.New("could not find node", err)
	}
	return &TaskPool{
		rootKey: nodePool.GetKey() + "/" + nodeId + "/running_tasks",
	}, nil
}

func (nodePool *NodePool) saveNode(node *lxtypes.Node) error {
	nodeId := node.Id
	err := lxdatabase.Mkdir(nodePool.GetKey() + "/" + nodeId)
	if err != nil {
		return lxerrors.New("initializing "+nodePool.GetKey()+"/"+nodeId+" directory", err)
	}
	nodeResourcePool := &ResourcePool{
		nodeId:  nodeId,
		rootKey: nodePool.GetKey() + "/" + nodeId + "/resources",
	}
	nodeResourcePool.Initialize()
	for _, resource := range node.Resources {
		err = nodeResourcePool.AddResource(resource)
		if err != nil {
			return lxerrors.New("adding resource "+resource.Id+" to node "+nodeId+" resource pool", err)
		}
	}
	nodeTaskPool := &TaskPool{
		rootKey: nodePool.GetKey() + "/" + nodeId + "/running_tasks",
	}
	nodeTaskPool.Initialize()
	for _, task := range node.RunningTasks {
		err = nodeTaskPool.AddTask(task)
		if err != nil {
			return lxerrors.New("adding task "+task.TaskId+" to node "+nodeId+" task pool", err)
		}
	}
	return nil
}

func (nodePool *NodePool) loadNode(nodeId string) (*lxtypes.Node, error) {
	node := lxtypes.NewNode(nodeId)
	nodeResourcePool := ResourcePool{
		nodeId:  nodeId,
		rootKey: nodePool.GetKey() + "/" + nodeId + "/resources",
	}
	resources, err := nodeResourcePool.GetResources()
	if err != nil {
		return nil, lxerrors.New("could not get list of resources for node "+nodeId, err)
	}
	for _, resource := range resources {
		err = node.AddResource(resource)
		if err != nil {
			return nil, lxerrors.New("could not add resource to node object", err)
		}
	}
	nodeTaskPool := TaskPool{
		rootKey: nodePool.GetKey() + "/" + nodeId + "/running_tasks",
	}
	tasks, err := nodeTaskPool.GetTasks()
	if err != nil {
		return nil, lxerrors.New("could not get list of tasks for node "+nodeId, err)
	}
	for _, task := range tasks {
		err = node.AddTask(task)
		if err != nil {
			return nil, lxerrors.New("could not add task to node object", err)
		}
	}
	return node, nil
}
