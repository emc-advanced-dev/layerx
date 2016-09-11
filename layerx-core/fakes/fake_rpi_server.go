package fakes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/emc-advanced-dev/layerx/layerx-core/layerx_rpi_client"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"
	"github.com/emc-advanced-dev/pkg/errors"
	"github.com/go-martini/martini"
	"github.com/golang/protobuf/proto"
	"github.com/mesos/mesos-go/mesosproto"
	"github.com/pborman/uuid"
)

const (
	COLLECT_RESOURCES = "/collect_resources"
	LAUNCH_TASKS      = "/launch_tasks"
	KILL_TASK         = "/kill_task"
)

func RunFakeRpiServer(layerxUrl string, port int, driverErrc chan error) {

	m := martini.Classic()

	collectResourcesHandler := func(req *http.Request, res http.ResponseWriter) {
		collectResourcesFn := func() ([]byte, int, error) {
			err := fakeCollectResources(layerxUrl)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Errorf("could not handle collect resources request")
				return empty, 500, errors.New("could not handle collect resources request", err)
			}
			return empty, 202, nil
		}
		_, statusCode, err := collectResourcesFn()
		if err != nil {
			res.WriteHeader(statusCode)
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Errorf("processing collect resources message")
			driverErrc <- err
			return
		}
		res.WriteHeader(statusCode)
	}
	launchTaskHandler := func(req *http.Request, res http.ResponseWriter) {
		launchTaskFn := func() ([]byte, int, error) {
			data, err := ioutil.ReadAll(req.Body)
			if req.Body != nil {
				defer req.Body.Close()
			}
			if err != nil {
				return empty, 400, errors.New("parsing launch task request", err)
			}
			var launchTaskMessage layerx_rpi_client.LaunchTasksMessage
			err = json.Unmarshal(data, &launchTaskMessage)
			if err != nil {
				return empty, 500, errors.New("could not parse json to update launch task message", err)
			}
			err = fakeLaunchTasks(layerxUrl, launchTaskMessage)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Errorf("could not handle collect resources request")
				return empty, 500, errors.New("could not handle update launch task request", err)
			}
			return empty, 202, nil
		}
		_, statusCode, err := launchTaskFn()
		if err != nil {
			res.WriteHeader(statusCode)
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Errorf("processing update launch task message")
			driverErrc <- err
			return
		}
		res.WriteHeader(statusCode)
	}
	killTaskHandler := func(req *http.Request, res http.ResponseWriter, params martini.Params) {
		killTaskFn := func() ([]byte, int, error) {
			taskId := params["task_id"]
			err := fakeKillTask(layerxUrl, taskId)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Errorf("could not kill task")
				return empty, 500, errors.New("could not kill task", err)
			}
			return empty, 202, nil
		}
		_, statusCode, err := killTaskFn()
		if err != nil {
			res.WriteHeader(statusCode)
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Errorf("processing update kill task message")
			driverErrc <- err
			return
		}
		res.WriteHeader(statusCode)
	}

	m.Post(COLLECT_RESOURCES, collectResourcesHandler)
	m.Post(LAUNCH_TASKS, launchTaskHandler)
	m.Post(KILL_TASK+"/:task_id", killTaskHandler)

	m.RunOnAddr(fmt.Sprintf(":%v", port))
}

func fakeCollectResources(layerXUrl string) error {
	msg := fmt.Sprintf("accepted fake collect resources message")
	logrus.Debugf(msg)
	rpiClient := layerx_rpi_client.LayerXRpi{
		CoreURL: layerXUrl,
	}
	fakeResourceId1 := "fake_offer_id_" + uuid.New()
	fakeResourceId2 := "fake_offer_id_" + uuid.New()
	fakeResourceId3 := "fake_offer_id_" + uuid.New()
	fakeResource1 := lxtypes.NewResourceFromMesos(FakeOffer(fakeResourceId1, "fake_slave_id_1"))
	fakeResource2 := lxtypes.NewResourceFromMesos(FakeOffer(fakeResourceId2, "fake_slave_id_2"))
	fakeResource3 := lxtypes.NewResourceFromMesos(FakeOffer(fakeResourceId3, "fake_slave_id_3"))
	err := rpiClient.SubmitResource(fakeResource1)
	if err != nil {
		return errors.New("submitting resource", err)
	}
	err = rpiClient.SubmitResource(fakeResource2)
	if err != nil {
		return errors.New("submitting resource", err)
	}
	err = rpiClient.SubmitResource(fakeResource3)
	if err != nil {
		return errors.New("submitting resource", err)
	}
	return nil
}

func fakeLaunchTasks(layerXUrl string, launchTaskMessage layerx_rpi_client.LaunchTasksMessage) error {
	if len(launchTaskMessage.ResourcesToUse) < 1 {
		return errors.New("must specify at least one resource for fake launch!", nil)
	}
	nodeId := launchTaskMessage.ResourcesToUse[0].NodeId
	for _, task := range launchTaskMessage.TasksToLaunch {
		logrus.WithFields(logrus.Fields{"task": task}).Debugf("fake rpi launching fake task")
		rpiClient := layerx_rpi_client.LayerXRpi{
			CoreURL: layerXUrl,
		}
		fakeRunningStatus := FakeTaskStatus(task.TaskId, mesosproto.TaskState_TASK_RUNNING)
		fakeRunningStatus.SlaveId = &mesosproto.SlaveID{
			Value: proto.String(nodeId),
		}
		err := rpiClient.SubmitStatusUpdate(fakeRunningStatus)
		if err != nil {
			return errors.New("submitting fake TASK_RUNNING status update to layerx core", err)
		}
	}
	for _, task := range launchTaskMessage.TasksToLaunch {
		logrus.WithFields(logrus.Fields{"task": task}).Debugf("fake rpi launching fake tasks")
	}
	for _, resource := range launchTaskMessage.ResourcesToUse {
		logrus.WithFields(logrus.Fields{"resource": resource}).Debugf("fake rpi launching fake tasks on resource")
	}
	return nil
}

func fakeKillTask(layerXUrl string, taskId string) error {
	logrus.WithFields(logrus.Fields{"task_id": taskId}).Debugf("fake rpi killing fake task")
	return nil
}
