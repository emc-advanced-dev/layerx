package framework_api_handlers

import (
	"github.com/emc-advanced-dev/layerx/layerx-core/layerx_rpi_client"
	"github.com/emc-advanced-dev/pkg/errors"
	"github.com/mesos/mesos-go/mesosproto"
	"github.com/mesos/mesos-go/scheduler"
)

func HandleOfferRescinded(core *layerx_rpi_client.LayerXRpi, driver scheduler.SchedulerDriver, id *mesosproto.OfferID) error {
	if err := core.RescindResource(id.GetValue()); err != nil {
		return errors.New("submitting rescind resource request to core", err)
	}
	//make sure we flush whatever isn't valid anymore
	go driver.DeclineOffer(id, &mesosproto.Filters{})
	return nil
}
