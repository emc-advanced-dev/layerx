package fakes

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/go-martini/martini"
	"github.com/layer-x/layerx-commons/lxlog"
	"io/ioutil"
	"net/http"
	"github.com/layer-x/layerx-core_v2/layerx_tpi_client"
	"github.com/layer-x/layerx-commons/lxerrors"
	"strings"
)

const (
	COLLECT_TASKS = "/collect_tasks"
	UPDATE_TASK_STATUS = "/update_task_status"
	HEALTH_CHECK_TASK_PROVIDER = "/health_check_task_provider"

	FAIL_ON_PURPOSE="failonpurpose"
)

var empty = []byte{}

func RunFakeTpiServer(layerxUrl string, port int, driverErrc chan error) {

	m := martini.Classic()

	collectTasksHandler := func(req *http.Request, res http.ResponseWriter) {
		collectTasksFn := func() ([]byte, int, error) {
			data, err := ioutil.ReadAll(req.Body)
			if req.Body != nil {
				defer req.Body.Close()
			}
			if err != nil {
				return empty, 400, lxerrors.New("parsing collect tasks request", err)
			}
			var collectTasksMessage layerx_tpi_client.CollectTasksMessage
			err = json.Unmarshal(data, &collectTasksMessage)
			if err != nil {
				return empty, 500, lxerrors.New("could not parse json to collect tasks message", err)
			}
			err = fakeCollectTasks(layerxUrl, collectTasksMessage)
			if err != nil {
				lxlog.Errorf(logrus.Fields{
					"error": err,
				}, "could not handle collect tasks request")
				return empty, 500, lxerrors.New("could not handle collect tasks request", err)
			}
			return empty, 202, nil
		}
		_, statusCode, err := collectTasksFn()
		if err != nil {
			res.WriteHeader(statusCode)
			lxlog.Errorf(logrus.Fields{
				"error": err.Error(),
			}, "processing collect tasks message")
			driverErrc <- err
			return
		}
		res.WriteHeader(statusCode)
	}
	updateTaskStatusHandler := func(req *http.Request, res http.ResponseWriter) {
		updateTaskStatusFn := func() ([]byte, int, error) {
			data, err := ioutil.ReadAll(req.Body)
			if req.Body != nil {
				defer req.Body.Close()
			}
			if err != nil {
				return empty, 400, lxerrors.New("parsing update task status request", err)
			}
			var updateTaskStatusMessage layerx_tpi_client.UpdateTaskStatusMessage
			err = json.Unmarshal(data, &updateTaskStatusMessage)
			if err != nil {
				return empty, 500, lxerrors.New("could not parse json to update task status message", err)
			}
			err = fakeUpdateTaskStatus(updateTaskStatusMessage)
			if err != nil {
				lxlog.Errorf(logrus.Fields{
					"error": err,
				}, "could not handle collect tasks request")
				return empty, 500, lxerrors.New("could not handle update task status request", err)
			}
			return empty, 202, nil
		}
		_, statusCode, err := updateTaskStatusFn()
		if err != nil {
			res.WriteHeader(statusCode)
			lxlog.Errorf(logrus.Fields{
				"error": err.Error(),
			}, "processing update task status message")
			driverErrc <- err
			return
		}
		res.WriteHeader(statusCode)
	}
	healthCheckTaskProviderHandler := func(req *http.Request, res http.ResponseWriter) {
		fn := func() ([]byte, int, error) {
			data, err := ioutil.ReadAll(req.Body)
			if req.Body != nil {
				defer req.Body.Close()
			}
			if err != nil {
				return empty, 400, lxerrors.New("parsing update health check task provider request", err)
			}
			var healthCheckTaskProviderMessage layerx_tpi_client.HealthCheckTaskProviderMessage
			err = json.Unmarshal(data, &healthCheckTaskProviderMessage)
			if err != nil {
				return empty, 500, lxerrors.New("could not parse json to health check task provider message", err)
			}
			//this is to make the request fail for testing
			failOnPurpose := strings.Contains(healthCheckTaskProviderMessage.TaskProvider.Id, FAIL_ON_PURPOSE)
			if failOnPurpose {
				return empty, http.StatusGone, nil
			}
			return empty, http.StatusOK, nil
		}
		_, statusCode, err := fn()
		if err != nil {
			res.WriteHeader(statusCode)
			lxlog.Errorf(logrus.Fields{
				"error": err.Error(),
			}, "processing update task status message")
			driverErrc <- err
			return
		}
		res.WriteHeader(statusCode)
	}

	m.Post(COLLECT_TASKS, collectTasksHandler)
	m.Post(UPDATE_TASK_STATUS, updateTaskStatusHandler)
	m.Post(HEALTH_CHECK_TASK_PROVIDER, healthCheckTaskProviderHandler)

	m.RunOnAddr(fmt.Sprintf(":%v", port))
}

func fakeCollectTasks(layerXUrl string, collectTasksMessage layerx_tpi_client.CollectTasksMessage) error {
	msg := fmt.Sprintf("accepted fake collect tasks message: %v", collectTasksMessage)
	lxlog.Debugf(logrus.Fields{}, msg)
	tpiClient := layerx_tpi_client.LayerXTpi{
		CoreURL: layerXUrl,
	}
	for _, taskProvider := range collectTasksMessage.TaskProviders {
		fakeTaskName := "fake_task_for_"+taskProvider.Id
		fakeTaskId := fakeTaskName+"_id"
		fakeSlaveId := "fake_slave_id"
		fakeCommand := `i=0; while true; do echo $i; i=$(expr $i + 1); sleep 1; done`
		fakeTaskForProvider := FakeLXTask(fakeTaskId, fakeTaskName, fakeSlaveId, fakeCommand)
		err := tpiClient.SubmitTask(taskProvider.Id, fakeTaskForProvider)
		if err != nil {
			return lxerrors.New("submitting fake task to lx core", err)
		}
	}
	return nil
}
func fakeUpdateTaskStatus(updateTaskStatusMessage layerx_tpi_client.UpdateTaskStatusMessage) error {
	msg := fmt.Sprintf("accepted fake task status update: %v", updateTaskStatusMessage)
	lxlog.Debugf(logrus.Fields{}, msg)
	return nil
}