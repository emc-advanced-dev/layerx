package tpi_api_helpers
import (
	"github.com/layer-x/layerx-core_v2/layerx_tpi_client"
	"github.com/layer-x/layerx-mesos-tpi_v2/framework_manager"
	"github.com/layer-x/layerx-mesos-tpi_v2/mesos_master_api/mesos_data"
	"github.com/layer-x/layerx-commons/lxerrors"
)

func UpdateTaskStatus(tpi *layerx_tpi.LayerXTpi, frameworkManager framework_manager.FrameworkManager, updateTaskStatusMessage layerx_tpi.UpdateTaskStatusMessage) error {
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