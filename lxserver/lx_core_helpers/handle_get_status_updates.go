package lx_core_helpers
import (
	"github.com/layer-x/layerx-core_v2/lxstate"
	"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/mesos/mesos-go/mesosproto"
)

func GetStatusUpdates(state *lxstate.State, tpId string) ([]*mesosproto.TaskStatus, error) {
	statusMap, err := state.GetStatusUpdatesForTaskProvider(tpId)
	if err != nil {
		return nil, lxerrors.New("getting statuses for task provider", err)
	}
	statuses := []*mesosproto.TaskStatus{}
	for _, status := range statusMap {
		statuses = append(statuses, status)
	}
	return statuses, nil
}
