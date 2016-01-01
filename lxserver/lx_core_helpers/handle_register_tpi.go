package lx_core_helpers
import (
	"github.com/layer-x/layerx-core_v2/lxstate"
	"github.com/layer-x/layerx-commons/lxerrors"
)

func RegisterTpi(state *lxstate.State, tpiUrl string) error {
	err := state.SetTpi(tpiUrl)
	if err != nil {
		return lxerrors.New("setting tpi url in state", err)
	}
	return nil
}
