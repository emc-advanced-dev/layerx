package rpi_api_helpers
import (
	"github.com/mesos/mesos-go/scheduler"
	"github.com/layer-x/layerx-core_v2/layerx_rpi_client"
	"github.com/mesos/mesos-go/mesosproto"
	"github.com/golang/protobuf/proto"
	"github.com/layer-x/layerx-commons/lxerrors"
	"fmt"
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
	filters := &mesosproto.Filters{
		RefuseSeconds: proto.Float64(0),
	}
	status, err := driver.LaunchTasks(offerIds, mesosTasks, filters)
	if err != nil {
		errmsg := fmt.Sprintf("launching %v tasks on %v offers with mesos schedulerdriver", len(mesosTasks), len(offerIds))
		return lxerrors.New(errmsg, err)
	}
	if status != mesosproto.Status_DRIVER_RUNNING {
		err = lxerrors.New("expected status "+mesosproto.Status_DRIVER_RUNNING.String()+ " but got "+status.String(), nil)
		errmsg := fmt.Sprintf("launching %v tasks on %v offers with mesos schedulerdriver", len(mesosTasks), len(offerIds))
		return lxerrors.New(errmsg, err)
	}
	return nil
}