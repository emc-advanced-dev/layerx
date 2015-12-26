package mesos_api_helpers
import (
"github.com/mesos/mesos-go/mesosproto"
"github.com/layer-x/layerx-commons/lxlog"
	"github.com/Sirupsen/logrus"
)

func LogStatusUpdateAck(statusUpdateAck mesosproto.StatusUpdateAcknowledgementMessage) error {
	frameworkId := statusUpdateAck.GetFrameworkId().GetValue()
	slaveId := statusUpdateAck.GetSlaveId().GetValue()
	taskId := statusUpdateAck.GetTaskId().GetValue()
	msgUuid := statusUpdateAck.GetUuid()
	lxlog.Debugf(logrus.Fields{
		"framework_id": frameworkId,
		"slave_id":     slaveId,
		"task_id":      taskId,
		"msg_uuid":     string(msgUuid),
	}, "received status update acknowledgement from framework %s", frameworkId)
	return nil
}
