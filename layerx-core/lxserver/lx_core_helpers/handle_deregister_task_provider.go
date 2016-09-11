package lx_core_helpers

import (
	"github.com/emc-advanced-dev/layerx/layerx-core/lxstate"
	"github.com/emc-advanced-dev/pkg/errors"
)

func DeregisterTaskProvider(state *lxstate.State, tpId string) error {
	err := state.TaskProviderPool.DeleteTaskProvider(tpId)
	if err != nil {
		return errors.New("deleting task provider from pool", err)
	}
	return nil
}
