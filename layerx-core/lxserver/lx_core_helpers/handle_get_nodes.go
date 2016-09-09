package lx_core_helpers
import (
	"github.com/emc-advanced-dev/layerx/layerx-core/lxstate"
	"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"
)

func  GetNodes(state *lxstate.State) ([]*lxtypes.Node, error) {
	nodeMap, err := state.NodePool.GetNodes()
	if err != nil {
		return nil, lxerrors.New("getting node list from pool", err)
	}
	nodes := []*lxtypes.Node{}
	for _, node := range nodeMap {
		nodes = append(nodes, node)
	}
	return nodes, nil
}
