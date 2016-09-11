package lx_core_helpers

import (
	"github.com/emc-advanced-dev/layerx/layerx-core/lxstate"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"
	"github.com/emc-advanced-dev/pkg/errors"
)

func GetNodes(state *lxstate.State) ([]*lxtypes.Node, error) {
	nodeMap, err := state.NodePool.GetNodes()
	if err != nil {
		return nil, errors.New("getting node list from pool", err)
	}
	nodes := []*lxtypes.Node{}
	for _, node := range nodeMap {
		nodes = append(nodes, node)
	}
	return nodes, nil
}
