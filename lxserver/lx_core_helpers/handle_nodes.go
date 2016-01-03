package lx_core_helpers
import (
	"github.com/layer-x/layerx-core_v2/lxstate"
	"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/layer-x/layerx-core_v2/lxtypes"
)

func  GetNodes(state *lxstate.State) ([]*lxtypes.Node, error) {
	nodeMap, err := state.NodePool.GetNodes()
	if err != nil {
		return nil, lxerrors.New("deleting task provider from pool", err)
	}
	nodes := []*lxtypes.Node{}
	for _, node := range nodeMap {
		nodes = append(nodes, node)
	}
	return nodes, nil
}
