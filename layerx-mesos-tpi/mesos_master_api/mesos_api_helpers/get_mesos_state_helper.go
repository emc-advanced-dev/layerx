package mesos_api_helpers

import (
	"encoding/json"
	"github.com/emc-advanced-dev/layerx/layerx-mesos-tpi/mesos_master_api/mesos_data"
	"github.com/emc-advanced-dev/pkg/errors"
)

var (
	empty = []byte{}
)

func GetMesosState(masterUpidString string) ([]byte, error) {
	state := &mesos_data.MesosState{
		Version: mesos_data.MESOS_VERSION,
		Leader:  masterUpidString,
	}
	data, err := json.Marshal(state)
	if err != nil {
		return empty, errors.New("marshalling master state to json", err)
	}
	return data, nil
}
