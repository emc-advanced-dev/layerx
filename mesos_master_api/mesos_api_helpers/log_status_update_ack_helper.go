package mesos_api_helpers

import (
	"github.com/Sirupsen/logrus"
	"github.com/mesos/mesos-go/mesosproto"
)

func LogStatusUpdateAck(statusUpdateAck mesosproto.StatusUpdateAcknowledgementMessage) error {
	frameworkId := statusUpdateAck.GetFrameworkId().GetValue()
	slaveId := statusUpdateAck.GetSlaveId().GetValue()
	taskId := statusUpdateAck.GetTaskId().GetValue()
	msgUuid := statusUpdateAck.GetUuid()
	logrus.WithFields(logrus.Fields{
		"framework_id": frameworkId,
		"slave_id":     slaveId,
		"task_id":      taskId,
		"msg_uuid":     string(msgUuid),
	}).Debugf("received status update acknowledgement from framework %s", frameworkId)
	return nil
}
