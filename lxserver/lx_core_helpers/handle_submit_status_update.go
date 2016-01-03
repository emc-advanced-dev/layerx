package lx_core_helpers
import (
	"github.com/layer-x/layerx-core_v2/lxstate"
	"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/mesos/mesos-go/mesosproto"
)

func SubmitStatusUpdate(state *lxstate.State, status *mesosproto.TaskStatus) error {
	err := state.StatusPool.AddStatus(status)
	if err != nil {
		return lxerrors.New("adding status for task "+status.GetTaskId().GetValue()+" to pool", err)
	}
	return nil
}