package tpi_api_helpers
import (
	"github.com/emc-advanced-dev/layerx/layerx-core/layerx_tpi_client"
	"github.com/emc-advanced-dev/layerx/layerx-mesos-tpi/mesos_master_api/mesos_data"
	"github.com/layer-x/layerx-commons/lxerrors"
"github.com/emc-advanced-dev/layerx/layerx-mesos-tpi/framework_manager"
)

func CollectTasks(tpi *layerx_tpi_client.LayerXTpi, frameworkManager framework_manager.FrameworkManager, collectTasksMessage layerx_tpi_client.CollectTasksMessage) error {
	for _, taskProvider := range collectTasksMessage.TaskProviders {
		frameworkId := taskProvider.Id
		upid, err := mesos_data.UPIDFromString(taskProvider.Source)
		if err != nil {
			return lxerrors.New("parsing upid from string", err)
		}
		err = frameworkManager.SendTaskCollectionOffer(frameworkId,
			"phony_offer_id",
			"phony_slave",
			"phony_pid",
			upid)
		if err != nil {
			return lxerrors.New("sending task collection offer to framework", err)
		}
	}
	return nil
}