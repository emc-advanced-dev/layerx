package rpi_api_helpers

import (
	"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/mesos/mesos-go/mesosproto"
	"github.com/mesos/mesos-go/scheduler"
)

func CollectResources(driver scheduler.SchedulerDriver) error {
	status, err := driver.ReviveOffers()
	if err != nil {
		return lxerrors.New("reviving offers with mesos schedulerdriver", err)
	}
	if status != mesosproto.Status_DRIVER_RUNNING {
		err = lxerrors.New("expected status "+mesosproto.Status_DRIVER_RUNNING.String()+" but got "+status.String(), nil)
		return lxerrors.New("reviving offers with mesos schedulerdriver", err)
	}
	return nil
}
