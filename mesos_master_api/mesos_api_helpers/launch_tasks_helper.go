package mesos_api_helpers

import (
	"github.com/Sirupsen/logrus"
	"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/layer-x/layerx-core_v2/layerx_tpi_client"
	"github.com/layer-x/layerx-core_v2/lxtypes"
	"github.com/mesos/mesos-go/mesosproto"
)

func HandleLaunchTasksRequest(tpi *layerx_tpi_client.LayerXTpi, frameworkId string, mesosTasks []*mesosproto.TaskInfo) error {
	for _, mesosTask := range mesosTasks {
		lxTask := lxtypes.NewTaskFromMesos(mesosTask)
		err := tpi.SubmitTask(frameworkId, lxTask)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error":       err.Error(),
				"frameworkId": frameworkId,
				"tpi":         tpi,
				"task":        lxTask,
			}).Errorf("submitting task to layer-x core")
			return lxerrors.New("submitting task "+lxTask.TaskId+" to layer-x core", err)
		}
	}
	return nil
}
