package mesos_api_helpers
import (
	"github.com/layer-x/layerx-core_v2/layerx_tpi_client"
	"github.com/layer-x/layerx-commons/lxerrors"
"github.com/layer-x/layerx-commons/lxlog"
	"github.com/Sirupsen/logrus"
)

func HandleRemoveFramework(tpi *layerx_tpi.LayerXTpi, frameworkId string) error {
	err := tpi.DeregisterTaskProvider(frameworkId)
	if err != nil {
		err = lxerrors.New("registering framework as new task provider with layer x", err)
		lxlog.Errorf(logrus.Fields{
			"error": err.Error(),
			"frameworkId": frameworkId,
			"tpi": tpi,
		}, "handling removal of framework from layer-x")
		return err
	}
	return nil
}