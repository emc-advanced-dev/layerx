package handlers
import (
	"github.com/layer-x/layerx-mesos-tpi_v2/mesos_master_api/mesos_data"
	"github.com/layer-x/layerx-commons/lxerrors"
	"encoding/json"
)

var (
	empty = []byte{}
)


func GetMesosState(masterUpidString string) ([]byte, error) {
	state := &mesos_data.MesosState{
		Version: mesos_data.MESOS_VERSION,
		Leader: masterUpidString,
	}
	data, err := json.Marshal(state)
	if err != nil {
		return empty, lxerrors.New("marshalling master state to json", err)
	}
	return data, nil
}