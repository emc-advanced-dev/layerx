package lx_core_helpers
import (
	"github.com/layer-x/layerx-core_v2/lxstate"
	"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/layer-x/layerx-core_v2/layerx_rpi_client"
)

func RegisterRpi(state *lxstate.State, rpiRegister layerx_rpi_client.RpiInfo) error {
	err := state.RpiPool.AddRpi(&rpiRegister)
	if err != nil {
		return lxerrors.New("adding rpi info to state", err)
	}
	return nil
}
