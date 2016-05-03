package mesos_api_helpers
import (
	"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/mesos/mesos-go/mesosproto"
	"github.com/layer-x/layerx-core_v2/layerx_tpi_client"
	"github.com/layer-x/layerx-mesos-tpi_v2/framework_manager"
	"github.com/layer-x/layerx-mesos-tpi_v2/mesos_master_api/mesos_data"
	"github.com/Sirupsen/logrus"
)

func HandleReconcileTasksRequest(tpi *layerx_tpi_client.LayerXTpi, frameworkManager framework_manager.FrameworkManager, frameworkUpid *mesos_data.UPID, frameworkId string, taskIds []string) error {
	statusUpdates := []*mesosproto.TaskStatus{}
	for _, taskId := range taskIds {
		status, err := tpi.GetStatusUpdate(taskId)
		if err != nil {
			return lxerrors.New("getting status for task " + taskId + " from layerx core", err)
		}
		statusUpdates = append(statusUpdates, status)
		continue
	}

	for _, status := range statusUpdates {
		err := frameworkManager.SendStatusUpdate(frameworkId, frameworkUpid, status)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"framework_id": frameworkId,
				"framework_upid": frameworkUpid.String(),
				"status": status,
				"error": err,
			}).Errorf( "failed sending status update to framework")
			return lxerrors.New("sending status update to framework", err)
		}
	}
	return nil
}
