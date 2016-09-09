package tpi_api_helpers
import (
	"github.com/emc-advanced-dev/layerx/layerx-core/layerx_tpi_client"
	"github.com/emc-advanced-dev/layerx/layerx-mesos-tpi/framework_manager"
	"github.com/emc-advanced-dev/layerx/layerx-mesos-tpi/mesos_master_api/mesos_data"
	"github.com/layer-x/layerx-commons/lxerrors"
)

func UpdateTaskStatus(tpi *layerx_tpi_client.LayerXTpi, frameworkManager framework_manager.FrameworkManager, updateTaskStatusMessage layerx_tpi_client.UpdateTaskStatusMessage) error {
	taskProvider := updateTaskStatusMessage.TaskProvider
	frameworkId := taskProvider.Id
	status := updateTaskStatusMessage.TaskStatus
	upid, err := mesos_data.UPIDFromString(taskProvider.Source)
	if err != nil {
		return lxerrors.New("parsing upid from string", err)
	}
	err = frameworkManager.SendStatusUpdate(frameworkId, upid, status)
	if err != nil {
		return lxerrors.New("sending status update offer to framework", err)
	}
	return nil
}