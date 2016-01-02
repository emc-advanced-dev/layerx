package tpi_api_helpers
import (
	"github.com/layer-x/layerx-core_v2/layerx_tpi_client"
	"github.com/layer-x/layerx-mesos-tpi_v2/mesos_master_api/mesos_data"
	"github.com/layer-x/layerx-commons/lxerrors"
"github.com/layer-x/layerx-mesos-tpi_v2/framework_manager"
)

func CollectTasks(tpi *layerx_tpi_client.LayerXTpi, frameworkManager framework_manager.FrameworkManager, collectTasksMessage layerx_tpi_client.CollectTasksMessage) error {
	for _, taskProvider := range collectTasksMessage.TaskProivders {
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