package mesos_api_helpers

import (
	"github.com/Sirupsen/logrus"
	"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/emc-advanced-dev/layerx/layerx-core/layerx_tpi_client"
)

func HandleRemoveFramework(tpi *layerx_tpi_client.LayerXTpi, frameworkId string) error {
	err := tpi.DeregisterTaskProvider(frameworkId)
	if err != nil {
		err = lxerrors.New("registering framework as new task provider with layer x", err)
		logrus.WithFields(logrus.Fields{
			"error":       err.Error(),
			"frameworkId": frameworkId,
			"tpi":         tpi,
		}).Errorf("handling removal of framework from layer-x")
		return err
	}
	return nil
}
