package rpi_api_helpers

import (
	"github.com/emc-advanced-dev/pkg/errors"
	"github.com/emc-advanced-dev/layerx/layerx-core/layerx_rpi_client"
	"github.com/layer-x/layerx-commons/lxhttpclient"
	"github.com/emc-advanced-dev/layerx/layerx-mesos-rpi/types"
	"encoding/json"
)

func CollectResources(core *layerx_rpi_client.LayerXRpi, masterAddr string) error {
	_, data, err := lxhttpclient.Get(masterAddr, "/slaves", nil)
	if err != nil {
		return errors.New("performing GET "+masterAddr+"/slaves.json", err)
	}
	var slaves types.Slaves
	if err := json.Unmarshal(data, &slaves); err != nil {
		return errors.New("unmarshalling slave data from mesos master", err)
	}
	for _, slave := range slaves.Slaves {
		resource, err := slave.ToResource()
		if err != nil {
			return errors.New("converting mesos slave to resource", err)
		}
		if err := core.SubmitResource(resource); err != nil {
			return errors.New("submitting resource to core", err)
		}
	}
	return nil
}
