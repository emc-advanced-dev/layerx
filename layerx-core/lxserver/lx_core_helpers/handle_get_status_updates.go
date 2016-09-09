package lx_core_helpers
import (
	"github.com/emc-advanced-dev/layerx/layerx-core/lxstate"
	"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/mesos/mesos-go/mesosproto"
)

func GetStatusUpdates(state *lxstate.State, tpId string) ([]*mesosproto.TaskStatus, error) {
	var (
		statusMap map[string]*mesosproto.TaskStatus
		err error
	)
	if tpId == "" {
		statusMap, err = state.GetStatusUpdates()
		if err != nil {
			return nil, lxerrors.New("getting all statuses", err)
		}
	} else {
		statusMap, err = state.GetStatusUpdatesForTaskProvider(tpId)
		if err != nil {
			return nil, lxerrors.New("getting statuses for task provider", err)
		}
	}
	statuses := []*mesosproto.TaskStatus{}
	for _, status := range statusMap {
		statuses = append(statuses, status)
	}
	return statuses, nil
}
