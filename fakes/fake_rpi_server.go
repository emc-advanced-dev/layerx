package fakes

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/go-martini/martini"
	"github.com/layer-x/layerx-commons/lxlog"
	"io/ioutil"
	"net/http"
	"github.com/layer-x/layerx-core_v2/layerx_rpi_client"
	"github.com/layer-x/layerx-commons/lxerrors"
"github.com/layer-x/layerx-core_v2/lxtypes"
)

const (
	COLLECT_RESOURCES = "/collect_resources"
	LAUNCH_TASKS = "/launch_resources"
	KILL_TASK = "/kill_task"
)

func RunFakeRpiServer(layerxUrl string, port int, driverErrc chan error) {

	m := martini.Classic()

	collectResourcesHandler := func(req *http.Request, res http.ResponseWriter) {
		collectResourcesFn := func() ([]byte, int, error) {
			err := fakeCollectResources(layerxUrl)
			if err != nil {
				lxlog.Errorf(logrus.Fields{
					"error": err,
				}, "could not handle collect resources request")
				return empty, 500, lxerrors.New("could not handle collect resources request", err)
			}
			return empty, 202, nil
		}
		_, statusCode, err := collectResourcesFn()
		if err != nil {
			res.WriteHeader(statusCode)
			lxlog.Errorf(logrus.Fields{
				"error": err.Error(),
			}, "processing collect resources message")
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
				return empty, 400, lxerrors.New("parsing launch task request", err)
			}
			var launchTaskMessage layerx_rpi_client.LaunchTasksMessage
			err = json.Unmarshal(data, &launchTaskMessage)
			if err != nil {
				return empty, 500, lxerrors.New("could not parse json to update launch task message", err)
			}
			err = fakeLaunchTasks(layerxUrl, launchTaskMessage)
			if err != nil {
				lxlog.Errorf(logrus.Fields{
					"error": err,
				}, "could not handle collect resources request")
				return empty, 500, lxerrors.New("could not handle update launch task request", err)
			}
			return empty, 202, nil
		}
		_, statusCode, err := launchTaskFn()
		if err != nil {
			res.WriteHeader(statusCode)
			lxlog.Errorf(logrus.Fields{
				"error": err.Error(),
			}, "processing update launch task message")
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
				lxlog.Errorf(logrus.Fields{
					"error": err,
				}, "could not kill task")
				return empty, 500, lxerrors.New("could not kill task", err)
			}
			return empty, 202, nil
		}
		_, statusCode, err := killTaskFn()
		if err != nil {
			res.WriteHeader(statusCode)
			lxlog.Errorf(logrus.Fields{
				"error": err.Error(),
			}, "processing update kill task message")
			driverErrc <- err
			return
		}
		res.WriteHeader(statusCode)
	}

	m.Post(COLLECT_RESOURCES, collectResourcesHandler)
	m.Post(LAUNCH_TASKS, launchTaskHandler)
	m.Post(KILL_TASK, killTaskHandler)

	m.RunOnAddr(fmt.Sprintf(":%v", port))
}

func fakeCollectResources(layerXUrl string) error {
	msg := fmt.Sprintf("accepted fake collect resources message")
	lxlog.Debugf(logrus.Fields{}, msg)
	rpiClient := layerx_rpi_client.LayerXRpi{
		CoreURL: layerXUrl,
	}
	fakeResource1 := lxtypes.NewResourceFromMesos(FakeOffer("fake_offer_id_1", "fake_slave_id_1"))
	fakeResource2 := lxtypes.NewResourceFromMesos(FakeOffer("fake_offer_id_2", "fake_slave_id_2"))
	fakeResource3 := lxtypes.NewResourceFromMesos(FakeOffer("fake_offer_id_3", "fake_slave_id_3"))
	err := rpiClient.SubmitResource(fakeResource1)
	if err != nil {
		return lxerrors.New("submitting resource", err)
	}
	err = rpiClient.SubmitResource(fakeResource2)
	if err != nil {
		return lxerrors.New("submitting resource", err)
	}
	err = rpiClient.SubmitResource(fakeResource3)
	if err != nil {
		return lxerrors.New("submitting resource", err)
	}
	return nil
}

func fakeLaunchTasks(layerXUrl string, launchTaskMessage layerx_rpi_client.LaunchTasksMessage) error {
//	msg := fmt.Sprintf("accepted fake collect resources message")
//	lxlog.Debugf(logrus.Fields{}, msg)
//	rpiClient := layerx_rpi_client.LayerXRpi{
//		CoreURL: layerXUrl,
//	}
//	for _
//	taskId := launcht
//	fakeStatus := FakeTaskStatus()
//	err := rpiClient.SubmitStatusUpdate()
//	if err != nil {
//		return lxerrors.New("submitting resource", err)
//	}
//	err = rpiClient.SubmitResource(fakeResource2)
//	if err != nil {
//		return lxerrors.New("submitting resource", err)
//	}
//	err = rpiClient.SubmitResource(fakeResource3)
//	if err != nil {
//		return lxerrors.New("submitting resource", err)
//	}
	for _, task := range launchTaskMessage.TasksToLaunch {
		lxlog.Debugf(logrus.Fields{"task": task}, "fake rpi launching fake task")
	}
	for _, resource := range launchTaskMessage.ResourcesToUse {
		lxlog.Debugf(logrus.Fields{"resource": resource}, "fake rpi launching fake tasks on resource")
	}
	return nil
}

func fakeKillTask(layerXUrl string, taskId string) error {
	lxlog.Debugf(logrus.Fields{"task_id": taskId}, "fake rpi killing fake task")
	return nil
}
