package mesos_api_helpers

import (
	"github.com/Sirupsen/logrus"
	"github.com/emc-advanced-dev/layerx/layerx-core/layerx_tpi_client"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"
	"github.com/emc-advanced-dev/layerx/layerx-mesos-tpi/framework_manager"
	"github.com/emc-advanced-dev/layerx/layerx-mesos-tpi/mesos_master_api/mesos_data"
	"github.com/emc-advanced-dev/pkg/errors"
	"github.com/mesos/mesos-go/mesosproto"
	"github.com/pborman/uuid"
)

func HandleRegisterRequest(tpi *layerx_tpi_client.LayerXTpi, frameworkManager framework_manager.FrameworkManager, frameworkUpid *mesos_data.UPID, frameworkInfo *mesosproto.FrameworkInfo) error {
	frameworkName := frameworkInfo.GetName()
	frameworkId := frameworkInfo.GetId().GetValue()
	failoverTimeout := frameworkInfo.GetFailoverTimeout()
	if frameworkId == "" {
		frameworkId = frameworkName + uuid.New()
	}

	taskProvider := &lxtypes.TaskProvider{
		Id:              frameworkId,
		Source:          frameworkUpid.String(),
		FailoverTimeout: failoverTimeout,
	}
	err := tpi.RegisterTaskProvider(taskProvider)
	if err != nil {
		err = errors.New("registering framework as new task provider with layer x", err)
		logrus.WithFields(logrus.Fields{
			"error":         err.Error(),
			"frameworkName": frameworkName,
			"frameworkId":   frameworkId,
			"tpi":           tpi,
		}).Errorf("handling subscribe call request")
		return err
	}

	err = frameworkManager.NotifyFrameworkRegistered(frameworkName, frameworkId, frameworkUpid)
	if err != nil {
		err = errors.New("sending framework registered message to framework", err)
		logrus.WithFields(logrus.Fields{
			"error":         err.Error(),
			"frameworkName": frameworkName,
			"frameworkId":   frameworkId,
			"frameworkUpid": frameworkUpid.String(),
		}).Errorf("handling subscribe call request")
		return err
	}
	return nil
}
