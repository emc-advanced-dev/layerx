package layerx_tpi_api

import (
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"github.com/emc-advanced-dev/layerx/layerx-core/layerx_tpi_client"
	"github.com/emc-advanced-dev/layerx/layerx-mesos-tpi/framework_manager"
	"github.com/emc-advanced-dev/layerx/layerx-mesos-tpi/layerx_tpi_api/tpi_api_helpers"
	"github.com/emc-advanced-dev/pkg/errors"
	"github.com/go-martini/martini"
	"io/ioutil"
	"net/http"
)

const (
	COLLECT_TASKS              = "/collect_tasks"
	UPDATE_TASK_STATUS         = "/update_task_status"
	HEALTH_CHECK_TASK_PROVIDER = "/health_check_task_provider"
)

var empty = []byte{}

type tpiApiServerWrapper struct {
	frameworkManager framework_manager.FrameworkManager
	tpi              *layerx_tpi_client.LayerXTpi
}

func NewTpiApiServerWrapper(tpi *layerx_tpi_client.LayerXTpi, frameworkManager framework_manager.FrameworkManager) *tpiApiServerWrapper {
	return &tpiApiServerWrapper{
		frameworkManager: frameworkManager,
		tpi:              tpi,
	}
}

func (wrapper *tpiApiServerWrapper) WrapWithTpi(m *martini.ClassicMartini, masterUpidString string, driverErrc chan error) *martini.ClassicMartini {
	collectTasksHandler := func(req *http.Request, res http.ResponseWriter) {
		collectTasksFn := func() ([]byte, int, error) {
			data, err := ioutil.ReadAll(req.Body)
			if req.Body != nil {
				defer req.Body.Close()
			}
			if err != nil {
				return empty, 400, errors.New("parsing collect tasks request", err)
			}
			var collectTasksMessage layerx_tpi_client.CollectTasksMessage
			err = json.Unmarshal(data, &collectTasksMessage)
			if err != nil {
				return empty, 500, errors.New("could not parse json to collect tasks message", err)
			}
			err = tpi_api_helpers.CollectTasks(wrapper.tpi, wrapper.frameworkManager, collectTasksMessage)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Errorf("could not handle collect tasks request")
				return empty, 500, errors.New("could not handle collect tasks request", err)
			}
			return empty, 202, nil
		}
		_, statusCode, err := wrapper.queueOperation(collectTasksFn)
		if err != nil {
			res.WriteHeader(statusCode)
			logrus.WithFields(logrus.Fields{
				"error":           err.Error(),
				"request_sent_by": masterUpidString,
			}).Errorf("processing collect tasks message")
			//			driverErrc <- err
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
				return empty, 400, errors.New("parsing update task status request", err)
			}
			var updateTaskStatusMessage layerx_tpi_client.UpdateTaskStatusMessage
			err = json.Unmarshal(data, &updateTaskStatusMessage)
			if err != nil {
				return empty, 500, errors.New("could not parse json to update task status message", err)
			}
			err = tpi_api_helpers.UpdateTaskStatus(wrapper.tpi, wrapper.frameworkManager, updateTaskStatusMessage)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Errorf("could not handle collect tasks request")
				return empty, 500, errors.New("could not handle update task status request", err)
			}
			return empty, 202, nil
		}
		_, statusCode, err := wrapper.queueOperation(updateTaskStatusFn)
		if err != nil {
			res.WriteHeader(statusCode)
			logrus.WithFields(logrus.Fields{
				"error":           err.Error(),
				"request_sent_by": masterUpidString,
			}).Errorf("processing update task status message")
			driverErrc <- err
			return
		}
		res.WriteHeader(statusCode)
	}
	healthCheckFrameworkHandler := func(req *http.Request, res http.ResponseWriter) {
		fn := func() ([]byte, int, error) {
			data, err := ioutil.ReadAll(req.Body)
			if req.Body != nil {
				defer req.Body.Close()
			}
			if err != nil {
				return empty, 400, errors.New("parsing health check task provider request", err)
			}
			var healthCheckMessage layerx_tpi_client.HealthCheckTaskProviderMessage
			err = json.Unmarshal(data, &healthCheckMessage)
			if err != nil {
				return empty, 500, errors.New("could not parse json to health check task provider message", err)
			}
			healthy, err := tpi_api_helpers.HealthCheck(wrapper.tpi, wrapper.frameworkManager, healthCheckMessage)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Errorf("could not handle collect tasks request")
				return empty, 500, errors.New("could not handle health check task provider request", err)
			}
			statusCode := http.StatusGone
			if healthy {
				statusCode = http.StatusOK
			}
			return empty, statusCode, nil
		}
		_, statusCode, err := wrapper.queueOperation(fn)
		if err != nil {
			res.WriteHeader(statusCode)
			logrus.WithFields(logrus.Fields{
				"error":           err.Error(),
				"request_sent_by": masterUpidString,
			}).Errorf("processing health check task provider message")
			driverErrc <- err
			return
		}
		res.WriteHeader(statusCode)
	}

	m.Post(COLLECT_TASKS, collectTasksHandler)
	m.Post(UPDATE_TASK_STATUS, updateTaskStatusHandler)
	m.Post(HEALTH_CHECK_TASK_PROVIDER, healthCheckFrameworkHandler)
	return m
}

func (wrapper *tpiApiServerWrapper) queueOperation(f func() ([]byte, int, error)) ([]byte, int, error) {
	datac := make(chan []byte)
	statusCodec := make(chan int)
	errc := make(chan error)
	go func() {
		data, statusCode, err := f()
		datac <- data
		statusCodec <- statusCode
		errc <- err
	}()
	return <-datac, <-statusCodec, <-errc
}
