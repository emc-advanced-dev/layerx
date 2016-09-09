package mesos_api_helpers

import (
	"github.com/Sirupsen/logrus"
	"github.com/mesos/mesos-go/mesosproto"
)

func LogReviveOffersMessage(reviveOffers mesosproto.ReviveOffersMessage) error {
	frameworkId := reviveOffers.GetFrameworkId().GetValue()
	logrus.WithFields(logrus.Fields{
		"framework_id": frameworkId,
	}).Debugf( "framework %s requested to revive offers", frameworkId)
	return nil
}
