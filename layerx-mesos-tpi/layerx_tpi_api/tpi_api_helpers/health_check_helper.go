package tpi_api_helpers

import (
	"github.com/emc-advanced-dev/layerx/layerx-core/layerx_tpi_client"
	"github.com/emc-advanced-dev/layerx/layerx-mesos-tpi/framework_manager"
	"github.com/emc-advanced-dev/layerx/layerx-mesos-tpi/mesos_master_api/mesos_data"
	"github.com/layer-x/layerx-commons/lxerrors"
)

func HealthCheck(tpi *layerx_tpi_client.LayerXTpi, frameworkManager framework_manager.FrameworkManager, healthCheckMessage layerx_tpi_client.HealthCheckTaskProviderMessage) (bool, error) {
	frameworkId := healthCheckMessage.TaskProvider.Id
	upid, err := mesos_data.UPIDFromString(healthCheckMessage.TaskProvider.Source)
	if err != nil {
		return false, lxerrors.New("parsing upid from string", err)
	}
	healthy, err := frameworkManager.HealthCheckFramework(frameworkId, upid)
	if err != nil {
		return false, lxerrors.New("sending health check GET to framework", err)
	}
	return healthy, nil
}
