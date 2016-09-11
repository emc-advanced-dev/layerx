package lx_core_helpers

import (
	"github.com/emc-advanced-dev/layerx/layerx-core/lxstate"
	"github.com/emc-advanced-dev/pkg/errors"
)

func RegisterTpi(state *lxstate.State, tpiUrl string) error {
	err := state.SetTpi(tpiUrl)
	if err != nil {
		return errors.New("setting tpi url in state", err)
	}
	return nil
}
