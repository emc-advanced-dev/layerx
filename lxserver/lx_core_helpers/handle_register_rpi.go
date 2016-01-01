package lx_core_helpers
import (
	"github.com/layer-x/layerx-core_v2/lxstate"
	"github.com/layer-x/layerx-commons/lxerrors"
)

func RegisterRpi(state *lxstate.State, rpiUrl string) error {
	err := state.SetRpi(rpiUrl)
	if err != nil {
		return lxerrors.New("setting rpi url in state", err)
	}
	return nil
}
