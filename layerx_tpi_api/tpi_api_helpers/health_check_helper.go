package tpi_api_helpers
import (
	"github.com/layer-x/layerx-core_v2/layerx_tpi_client"
	"github.com/layer-x/layerx-mesos-tpi_v2/mesos_master_api/mesos_data"
	"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/layer-x/layerx-mesos-tpi_v2/framework_manager"
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