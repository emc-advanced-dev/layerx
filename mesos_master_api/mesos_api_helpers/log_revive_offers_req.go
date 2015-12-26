package mesos_api_helpers
import (
"github.com/mesos/mesos-go/mesosproto"
"github.com/layer-x/layerx-commons/lxlog"
	"github.com/Sirupsen/logrus"
)

func LogReviveOffersMessage(reviveOffers mesosproto.ReviveOffersMessage) error {
	frameworkId := reviveOffers.GetFrameworkId().GetValue()
	lxlog.Debugf(logrus.Fields{
		"framework_id": frameworkId,
	}, "framework %s requested to revive offers", frameworkId)
	return nil
}