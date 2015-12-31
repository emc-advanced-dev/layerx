package state
import (
	"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/layer-x/layerx-commons/lxdatabase"
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
		return lxerrors.New("initializing "+nodePool.GetKey() +" directory", err)
	}
	return nil
}

//
//func (nodePool *NodePool) AddNode(node *lxtypes.Node) error {
//	nodeId := node.Id
//	_, err := nodePool.GetNode(nodeId)
//	if err == nil {
//		return lxerrors.New("node "+nodeId+" already exists in database, try Modify()?", err)
//	}
//	nodeData, err := json.Marshal(node)
//	if err != nil {
//		return lxerrors.New("could not marshal node to json", err)
//	}
//	err = lxdatabase.Set(nodePool.GetKey()+"/"+nodeId, string(nodeData))
//	if err != nil {
//		return lxerrors.New("setting key/value pair for node", err)
//	}
//	return nil
//}
//
//func (nodePool *NodePool) GetNode(nodeId string) (*lxtypes.Node, error) {
//	nodeJson, err := lxdatabase.Get(nodePool.GetKey()+"/"+nodeId)
//	if err != nil {
//		return nil, lxerrors.New("retrieving node "+nodeId+" from database", err)
//	}
//	var node lxtypes.Node
//	err = json.Unmarshal([]byte(nodeJson), &node)
//	if err != nil {
//		return nil, lxerrors.New("unmarshalling json into Node struct", err)
//	}
//	return &node, nil
//}
//
//func (nodePool *NodePool) saveNode(node *lxtypes.Node) error {
//	nodeId := node.Id
//	node.Resources
//	err := lxdatabase.Set(nodePool.GetKey()+"/"+nodeId+"/", string(nodeData))
//	if err != nil {
//		return lxerrors.New("setting key/value pair for node", err)
//	}
//}