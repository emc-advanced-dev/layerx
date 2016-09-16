package server

import (
	"github.com/emc-advanced-dev/pkg/errors"
	"github.com/emc-advanced-dev/layerx/layerx-core/layerx_rpi_client"
	"net/http"
	"github.com/Sirupsen/logrus"
	"io/ioutil"
	"encoding/json"
	"github.com/layer-x/layerx-commons/lxmartini"
	"github.com/go-martini/martini"
	"github.com/emc-advanced-dev/layerx/layerx-k8s-rpi/kube"
)

const (
	COLLECT_RESOURCES = "/collect_resources"
	LAUNCH_TASKS = "/launch_tasks"
	KILL_TASK = "/kill_task"
)

func Start(port string, client *kube.Client, core *layerx_rpi_client.LayerXRpi) {
	m := lxmartini.QuietMartini()

	m.Post(COLLECT_RESOURCES, func(req *http.Request, res http.ResponseWriter) {
		handle(res, func() (interface{}, int, error) {
			resources, err := client.FetchResources()
			if err != nil {
				return nil, 500, errors.New("could not handle collect resources request", err)
			}
			for _, resource := range resources {
				go func(){
					if err := core.SubmitResource(resource); err != nil {
						logrus.WithError(err).Errorf("failed submitting resource %v to core %v", resource, core)
					}
				}()
			}
			return nil, 202, nil
		})
	})
	m.Post(LAUNCH_TASKS, func(req *http.Request, res http.ResponseWriter) {
		handle(res, func() (interface{}, int, error) {
			data, err := ioutil.ReadAll(req.Body)
			if req.Body != nil {
				defer req.Body.Close()
			}
			if err != nil {
				return nil, 400, errors.New("parsing launch tasks request", err)
			}
			var launchTasksMessage layerx_rpi_client.LaunchTasksMessage
			if err := json.Unmarshal(data, &launchTasksMessage); err != nil {
				return nil, 500, errors.New("could not parse json to update task status message", err)
			}
			if err := client.LaunchTasks(launchTasksMessage); err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Errorf("could not handle launch tasks request")
				return nil, 500, errors.New("could not handle launch tasks request", err)
			}
			return nil, 202, nil
		})
	})
	m.Post(KILL_TASK + "/:task_id", func(req *http.Request, res http.ResponseWriter, params martini.Params) {
		handle(res, func() (interface{}, int, error) {
			taskId := params["task_id"]
			if err := client.KillTask(taskId); err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Errorf("could not handle kill task request")
				return nil, 500, errors.New("could not handle kill task request", err)
			}
			return nil, 202, nil
		})
	})

	m.RunOnAddr(":" + port)
}

func handle(res http.ResponseWriter, action func() (interface{}, int, error)) {
	jsonObject, statusCode, err := action()
	res.WriteHeader(statusCode)
	if err != nil {
		if err := respond(res, err); err != nil {
			logrus.WithError(err).Errorf("failed to reply to http request")
		}
		logrus.WithError(err).Errorf("error handling request")
		return
	}
	if jsonObject != nil {
		if err := respond(res, jsonObject); err != nil {
			logrus.WithError(err).Errorf("failed to reply to http request")
		}
		logrus.WithField("result", jsonObject).Debugf("request finished")
	}
}

func respond(res http.ResponseWriter, message interface{}) error {
	switch message.(type) {
	case string:
		messageString := message.(string)
		data := []byte(messageString)
		_, err := res.Write(data)
		if err != nil {
			return errors.New("writing data", err)
		}
		return nil
	case error:
		responseError := message.(error)
		_, err := res.Write([]byte(responseError.Error()))
		if err != nil {
			return errors.New("writing data", err)
		}
		return nil
	}
	data, err := json.Marshal(message)
	if err != nil {
		return errors.New("marshalling message to json", err)
	}
	_, err = res.Write(data)
	if err != nil {
		return errors.New("writing data", err)
	}
	return nil
}
