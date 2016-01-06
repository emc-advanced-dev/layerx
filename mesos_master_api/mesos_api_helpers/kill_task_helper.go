package mesos_api_helpers
import (
"github.com/layer-x/layerx-core_v2/layerx_tpi_client"
"github.com/layer-x/layerx-commons/lxlog"
	"github.com/Sirupsen/logrus"
	"github.com/layer-x/layerx-commons/lxerrors"
)

func HandleKillTaskRequest(tpi *layerx_tpi_client.LayerXTpi, frameworkId, taskId string) error {
	err := tpi.KillTask(frameworkId, taskId)
	if err != nil {
		lxlog.Errorf(logrus.Fields{
			"error": err.Error(),
			"tpi": tpi,
			"task_id": taskId,
		}, "submitting kill task " + taskId + " message to layer-x core")
		return lxerrors.New("submitting kill task " + taskId + " message to layer-x core", err)
	}
	return nil
}