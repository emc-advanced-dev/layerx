package lx_core_helpers
import (
	"github.com/layer-x/layerx-core_v2/lxstate"
	"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/mesos/mesos-go/mesosproto"
	"github.com/layer-x/layerx-core_v2/tpi_messenger"
)

func SubmitStatusUpdate(state *lxstate.State, tpiUrl string, status *mesosproto.TaskStatus) error {
	err := state.StatusPool.AddStatus(status)
	if err != nil {
		return lxerrors.New("adding status for task "+status.GetTaskId().GetValue()+" to pool", err)
	}
	taskProvider, err := state.GetAllTasks()
	err := tpi_messenger.SendStatusUpdate()
	return nil
}