package server

//import (
//	"github.com/emc-advanced-dev/pkg/errors"
//	"github.com/emc-advanced-dev/layerx/layerx-mesos-rpi/layerx_rpi_api/rpi_api_helpers"
//	"github.com/emc-advanced-dev/layerx/layerx-core/layerx_rpi_client"
//	"net/http"
//	"github.com/Sirupsen/logrus"
//	"io/ioutil"
//	"encoding/json"
//	"github.com/layer-x/layerx-commons/lxmartini"
//	"github.com/go-martini/martini"
//	"github.com/emc-advanced-dev/pkg/errors"
//)
//
//const (
//	COLLECT_RESOURCES = "/collect_resources"
//	LAUNCH_TASKS = "/launch_tasks"
//	KILL_TASK = "/kill_task"
//)
//
//var (
//	empty = []byte{}
//)
//
//func StartServer(port string) {
//	m := lxmartini.QuietMartini()
//
//	m.Post(COLLECT_RESOURCES, func(req *http.Request, res http.ResponseWriter) {
//		collectResourcesFn := func() ([]byte, int, error) {
//			if err := collectResources(); err != nil {
//				logrus.WithFields(logrus.Fields{
//					"error": err,
//				}).Errorf("could not handle collect resources request")
//				return empty, 500, errors.New("could not handle collect resources request", err)
//			}
//			return empty, 202, nil
//		}
//	})
//	m.Post(LAUNCH_TASKS, func(req *http.Request, res http.ResponseWriter) {
//		launchTasksFn := func() ([]byte, int, error) {
//			data, err := ioutil.ReadAll(req.Body)
//			if req.Body != nil {
//				defer req.Body.Close()
//			}
//			if err != nil {
//				return empty, 400, errors.New("parsing launch tasks request", err)
//			}
//			var launchTasksMessage layerx_rpi_client.LaunchTasksMessage
//			err = json.Unmarshal(data, &launchTasksMessage)
//			if err != nil {
//				return empty, 500, errors.New("could not parse json to update task status message", err)
//			}
//			err = rpi_api_helpers.LaunchTasks(wrapper.mesosSchedulerDriver, launchTasksMessage)
//			if err != nil {
//				logrus.WithFields(logrus.Fields{
//					"error": err,
//				}).Errorf("could not handle launch tasks request")
//				return empty, 500, errors.New("could not handle launch tasks request", err)
//			}
//			return empty, 202, nil
//		}
//		_, statusCode, err := wrapper.queueOperation(launchTasksFn)
//		if err != nil {
//			res.WriteHeader(statusCode)
//			logrus.WithFields(logrus.Fields{
//				"error": err.Error(),
//			}).Errorf("processing launch tasks message")
//			driverErrc <- err
//			return
//		}
//		res.WriteHeader(statusCode)
//	})
//	m.Post(KILL_TASK+"/:task_id", func(req *http.Request, res http.ResponseWriter, params martini.Params) {
//		killTaskFn := func() ([]byte, int, error) {
//			taskId := params["task_id"]
//			err := rpi_api_helpers.KillTask(wrapper.mesosSchedulerDriver, taskId)
//			if err != nil {
//				logrus.WithFields(logrus.Fields{
//					"error": err,
//				}).Errorf("could not handle kill task request")
//				return empty, 500, errors.New("could not handle kill task request", err)
//			}
//			return empty, 202, nil
//		}
//		_, statusCode, err := wrapper.queueOperation(killTaskFn)
//		if err != nil {
//			res.WriteHeader(statusCode)
//			logrus.WithFields(logrus.Fields{
//				"error": err.Error(),
//			}).Errorf("processing kill task message")
//			driverErrc <- err
//			return
//		}
//		res.WriteHeader(statusCode)
//	})
//
//	m.RunOnAddr(":" + *port)
//}
//
//func respond(res http.ResponseWriter, message interface{}) error {
//	switch message.(type) {
//	case string:
//		messageString := message.(string)
//		data := []byte(messageString)
//		_, err := res.Write(data)
//		if err != nil {
//			return errors.New("writing data", err)
//		}
//		return nil
//	case error:
//		responseError := message.(error)
//		_, err := res.Write([]byte(responseError.Error()))
//		if err != nil {
//			return errors.New("writing data", err)
//		}
//		return nil
//	}
//	data, err := json.Marshal(message)
//	if err != nil {
//		return errors.New("marshalling message to json", err)
//	}
//	_, err = res.Write(data)
//	if err != nil {
//		return errors.New("writing data", err)
//	}
//	return nil
//}
