package mesos_api_helpers

import (
	"github.com/Sirupsen/logrus"
	"github.com/emc-advanced-dev/layerx/layerx-core/layerx_tpi_client"
	"github.com/emc-advanced-dev/pkg/errors"
)

func HandleKillTaskRequest(tpi *layerx_tpi_client.LayerXTpi, frameworkId, taskId string) error {
	err := tpi.KillTask(frameworkId, taskId)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err.Error(),
			"tpi":     tpi,
			"task_id": taskId,
		}).Errorf("submitting kill task " + taskId + " message to layer-x core")
		return errors.New("submitting kill task "+taskId+" message to layer-x core", err)
	}
	return nil
}
