package layerx_tpi_api
import (
"github.com/layer-x/layerx-core_v2/layerx_tpi"
"github.com/layer-x/layerx-mesos-tpi_v2/framework_manager"
	"github.com/layer-x/layerx-commons/lxactionqueue"
	"github.com/go-martini/martini"
"net/http"
	"github.com/Sirupsen/logrus"
"github.com/layer-x/layerx-commons/lxlog"
	"github.com/layer-x/layerx-commons/lxerrors"
	"io/ioutil"
	"encoding/json"
	"github.com/layer-x/layerx-mesos-tpi_v2/layerx_tpi_api/tpi_api_helpers"
)

const (
	COLLECT_TASKS = "/collect_tasks"
	UPDATE_TASK_STATUS = "/update_task_status"
)

var empty = []byte{}

type tpiApiServerWrapper struct {
	actionQueue      lxactionqueue.ActionQueue
	frameworkManager framework_manager.FrameworkManager
	tpi              *layerx_tpi.LayerXTpi
}

func NewTpiApiServerWrapper(tpi *layerx_tpi.LayerXTpi, actionQueue lxactionqueue.ActionQueue, frameworkManager framework_manager.FrameworkManager) *tpiApiServerWrapper {
	return &tpiApiServerWrapper{
		actionQueue: actionQueue,
		frameworkManager: frameworkManager,
		tpi: tpi,
	}
}

func (wrapper *tpiApiServerWrapper) WrapWithTpi(m *martini.ClassicMartini, masterUpidString string, driverErrc chan error) *martini.ClassicMartini {
	collectTasksHandler := func(req *http.Request, res http.ResponseWriter) {
		registerFrameworkFn := func() ([]byte, int, error) {
			data, err := ioutil.ReadAll(req.Body)
			if req.Body != nil {
				defer req.Body.Close()
			}
			if err != nil {
				return empty, 400, lxerrors.New("parsing collect tasks request", err)
			}
			var collectTasksMessage layerx_tpi.CollectTasksMessage
			err = json.Unmarshal(data, &collectTasksMessage)
			if err != nil {
				return empty, 500, lxerrors.New("could not parse json to collect tasks message", err)
			}
			err = tpi_api_helpers.CollectTasks(wrapper.tpi, wrapper.frameworkManager, collectTasksMessage)
			if err != nil {
				lxlog.Errorf(logrus.Fields{
					"error": err,
				}, "could not handle collect tasks request")
				return empty, 500, lxerrors.New("could not handle collect tasks request", err)
			}
			return empty, 202, nil
		}
		_, statusCode, err := wrapper.queueOperation(registerFrameworkFn)
		if err != nil {
			res.WriteHeader(statusCode)
			lxlog.Errorf(logrus.Fields{
				"error": err.Error(),
				"request_sent_by": masterUpidString,
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
			var updateTaskStatusMessage layerx_tpi.UpdateTaskStatusMessage
			err = json.Unmarshal(data, &updateTaskStatusMessage)
			if err != nil {
				return empty, 500, lxerrors.New("could not parse json to update task status message", err)
			}
			err = tpi_api_helpers.UpdateTaskStatus(wrapper.tpi, wrapper.frameworkManager, updateTaskStatusMessage)
			if err != nil {
				lxlog.Errorf(logrus.Fields{
					"error": err,
				}, "could not handle collect tasks request")
				return empty, 500, lxerrors.New("could not handle update task status request", err)
			}
			return empty, 202, nil
		}
		_, statusCode, err := wrapper.queueOperation(updateTaskStatusFn)
		if err != nil {
			res.WriteHeader(statusCode)
			lxlog.Errorf(logrus.Fields{
				"error": err.Error(),
				"request_sent_by": masterUpidString,
			}, "processing update task status message")
			driverErrc <- err
			return
		}
		res.WriteHeader(statusCode)
	}

	m.Post(COLLECT_TASKS, collectTasksHandler)
	m.Post(UPDATE_TASK_STATUS, updateTaskStatusHandler)
	return m
}

func (wrapper *tpiApiServerWrapper) queueOperation(f func() ([]byte, int, error)) ([]byte, int, error) {
	datac := make(chan []byte)
	statusCodec := make(chan int)
	errc := make(chan error)
	wrapper.actionQueue.Push(
		func() {
			data, statusCode, err := f()
			datac <- data
			statusCodec <- statusCode
			errc <- err
		})
	return <-datac, <-statusCodec, <-errc
}
