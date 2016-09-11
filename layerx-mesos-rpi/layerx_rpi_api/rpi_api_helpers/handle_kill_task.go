package rpi_api_helpers

import (
	"fmt"
	"github.com/emc-advanced-dev/pkg/errors"
	"github.com/golang/protobuf/proto"
	"github.com/mesos/mesos-go/mesosproto"
	"github.com/mesos/mesos-go/scheduler"
)

func KillTask(driver scheduler.SchedulerDriver, taskId string) error {
	mesosTaskId := &mesosproto.TaskID{
		Value: proto.String(taskId),
	}
	status, err := driver.KillTask(mesosTaskId)
	if err != nil {
		errmsg := fmt.Sprintf("killing task %s with mesos schedulerdriver", taskId)
		return errors.New(errmsg, err)
	}
	if status != mesosproto.Status_DRIVER_RUNNING {
		err = errors.New("expected status "+mesosproto.Status_DRIVER_RUNNING.String()+" but got "+status.String(), nil)
		errmsg := fmt.Sprintf("killing task %s with mesos schedulerdriver", taskId)
		return errors.New(errmsg, err)
	}
	return nil
}
