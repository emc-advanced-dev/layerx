package lx_core_helpers

import (
	"github.com/emc-advanced-dev/layerx/layerx-core/layerx_rpi_client"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxstate"
	"github.com/emc-advanced-dev/pkg/errors"
)

func RegisterRpi(state *lxstate.State, rpiRegister layerx_rpi_client.RpiInfo) error {
	if _, err := state.RpiPool.GetRpi(rpiRegister.Name); err == nil {
		return errors.New("rpi with name "+rpiRegister.Name+" already exists. try a different name", nil)
	}
	err := state.RpiPool.AddRpi(&rpiRegister)
	if err != nil {
		return errors.New("adding rpi info to state", err)
	}
	return nil
}
