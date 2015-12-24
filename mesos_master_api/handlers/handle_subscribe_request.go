package handlers
import (
	"github.com/mesos/mesos-go/mesosproto"
	"github.com/layer-x/layerx-mesos-tpi_v2/mesos_master_api/mesos_data"
	"github.com/layer-x/layerx-commons/lxerrors"
"github.com/layer-x/layerx-commons/lxlog"
	"github.com/Sirupsen/logrus"
	"github.com/layer-x/layerx-mesos-tpi_v2/framework_manager"
)


func HandleSubscribeRequest(frameworkManager framework_manager.FrameworkManager, frameworkUpid *mesos_data.UPID, call *mesosproto.Call_Subscribe) error {
	frameworkInfo := call.GetFrameworkInfo()
	frameworkName := frameworkInfo.GetName()
	frameworkId := frameworkInfo.GetId().GetValue()
	err := frameworkManager.NotifyFrameworkRegistered(frameworkName, frameworkId, frameworkUpid)
	if err != nil {
		err = lxerrors.New("sending framework registered message to framework", err)
		lxlog.Errorf(logrus.Fields{
			"error": err.Error(),
			"frameworkName": frameworkName,
			"frameworkId": frameworkId,
			"frameworkUpid": frameworkUpid.String(),
		}, "handling subscribe call request", frameworkInfo.GetId().GetValue())
		return err
	}
	return nil
}
