package rpi_api_helpers

import (
	"fmt"
	"github.com/emc-advanced-dev/layerx/layerx-core/layerx_rpi_client"
	"github.com/golang/protobuf/proto"
	"github.com/emc-advanced-dev/pkg/errors"
	"github.com/mesos/mesos-go/mesosproto"
	"github.com/mesos/mesos-go/scheduler"
)

func LaunchTasks(driver scheduler.SchedulerDriver, launchTasksMessage layerx_rpi_client.LaunchTasksMessage) error {
	resources := launchTasksMessage.ResourcesToUse
	tasks := launchTasksMessage.TasksToLaunch
	mesosTasks := []*mesosproto.TaskInfo{}
	for _, task := range tasks {
		mesosTask := task.ToMesos()
		mesosTasks = append(mesosTasks, mesosTask)
	}
	offerIds := []*mesosproto.OfferID{}
	for _, resource := range resources {
		offerId := &mesosproto.OfferID{
			Value: proto.String(resource.Id),
		}
		offerIds = append(offerIds, offerId)
	}
	filters := &mesosproto.Filters{}
	status, err := driver.LaunchTasks(offerIds, mesosTasks, filters)
	if err != nil {
		errmsg := fmt.Sprintf("launching %v tasks on %v offers with mesos schedulerdriver", len(mesosTasks), len(offerIds))
		return errors.New(errmsg, err)
	}
	if status != mesosproto.Status_DRIVER_RUNNING {
		err = errors.New("expected status "+mesosproto.Status_DRIVER_RUNNING.String()+" but got "+status.String(), nil)
		errmsg := fmt.Sprintf("launching %v tasks on %v offers with mesos schedulerdriver", len(mesosTasks), len(offerIds))
		return errors.New(errmsg, err)
	}
	return nil
}
