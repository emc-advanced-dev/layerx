package mesos_api_helpers
import (
	"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/mesos/mesos-go/mesosproto"
	"github.com/layer-x/layerx-core_v2/layerx_tpi"
"github.com/layer-x/layerx-mesos-tpi_v2/framework_manager"
	"github.com/layer-x/layerx-mesos-tpi_v2/mesos_master_api/mesos_data"
	"github.com/layer-x/layerx-commons/lxlog"
	"github.com/Sirupsen/logrus"
)

func HandleReconcileTasksRequest(tpi *layerx_tpi.LayerXTpi, frameworkManager framework_manager.FrameworkManager, frameworkUpid *mesos_data.UPID, reconcileTasksMessage mesosproto.ReconcileTasksMessage) error {
	frameworkId := reconcileTasksMessage.GetFrameworkId().GetValue()
	statusUpdates := []*mesosproto.TaskStatus{}
	currentStatuses, err := tpi.GetStatusUpdates(frameworkId)
	if err != nil {
		return lxerrors.New("could not retrieve status updates from layer-x core", err)
	}
	for _, currentStatus := range currentStatuses {
		for _, desiredStatuses := range reconcileTasksMessage.GetStatuses() {
			if desiredStatuses.GetTaskId().GetValue() == currentStatus.GetTaskId().GetValue() {
				statusUpdates = append(statusUpdates, currentStatus)
				continue
			}
		}
	}
	for _, status := range statusUpdates {
		err = frameworkManager.SendStatusUpdate(frameworkId, frameworkUpid, status)
		if err != nil {
			lxlog.Errorf(logrus.Fields{
				"framework_id": frameworkId,
				"framework_upid": frameworkUpid.String(),
				"status": status,
				"error": err,
			}, "failed sending status update to framework")
			return lxerrors.New("sending status update to framework", err)
		}
	}
	return nil
}