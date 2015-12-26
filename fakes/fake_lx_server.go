package fakes

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/go-martini/martini"
	"github.com/layer-x/layerx-commons/lxlog"
	"github.com/layer-x/layerx-core_v2/lxtypes"
	"github.com/mesos/mesos-go/mesosproto"
	"io/ioutil"
	"net/http"
)

const (
	RegisterTaskProvider   = "/RegisterTaskProvider"
	DeregisterTaskProvider = "/DeregisterTaskProvider"
	GetTaskProviders       = "/GetTaskProviders"
	GetStatusUpdates       = "/GetStatusUpdates"
	SubmitTask             = "/SubmitTask"
	KillTask               = "/KillTask"
	PurgeTask              = "/PurgeTask"
)

func RunFakeLayerXServer(fakeStatuses []*mesosproto.TaskStatus, port int) {
	taskProviders := make(map[string]*lxtypes.TaskProvider)
	statusUpdates := make(map[string]*mesosproto.TaskStatus)
	tasks := make(map[string]*lxtypes.Task)

	for _, status := range fakeStatuses {
		statusUpdates[status.GetTaskId().GetValue()] = status
	}

	m := martini.Classic()
	m.Post(RegisterTaskProvider, func(res http.ResponseWriter, req *http.Request) {
		body, err := ioutil.ReadAll(req.Body)
		if req.Body != nil {
			defer req.Body.Close()
		}
		if err != nil {
			lxlog.Errorf(logrus.Fields{
				"error": err,
				"body":  string(body),
			}, "could not read  request body")
			res.WriteHeader(500)
			return
		}
		var tp lxtypes.TaskProvider
		err = json.Unmarshal(body, &tp)
		if err != nil {
			lxlog.Errorf(logrus.Fields{
				"error": err,
				"body":  string(body),
			}, "could parse json into tp")
			res.WriteHeader(500)
			return
		}
		taskProviders[tp.Id] = &tp
		res.WriteHeader(202)
	})
	m.Post(DeregisterTaskProvider+"/:task_provider_id", func(res http.ResponseWriter, req *http.Request, params martini.Params) {
		tpid := params["task_provider_id"]
		if _, ok := taskProviders[tpid]; !ok {
			lxlog.Errorf(logrus.Fields{
				"tpid": tpid,
			}, "task provider was not registered")
			res.WriteHeader(400)
			return
		}
		delete(taskProviders, tpid)
		res.WriteHeader(202)
	})
	m.Get(GetTaskProviders, func(res http.ResponseWriter, req *http.Request) {
		tps := []*lxtypes.TaskProvider{}
		for _, tp := range taskProviders {
			tps = append(tps, tp)
		}
		data, err := json.Marshal(tps)
		if err != nil {
			lxlog.Errorf(logrus.Fields{
				"error": err,
				"body":  string(data),
			}, "could parse tps into json")
			res.WriteHeader(500)
			return
		}
		res.Write(data)
	})
	m.Get(GetStatusUpdates, func(res http.ResponseWriter, req *http.Request) {
		statuses := []*mesosproto.TaskStatus{}
		for _, status := range statusUpdates {
			statuses = append(statuses, status)
		}
		data, err := json.Marshal(statuses)
		if err != nil {
			lxlog.Errorf(logrus.Fields{
				"error": err,
				"body":  string(data),
			}, "could parse statuses into json")
			res.WriteHeader(500)
			return
		}
		res.Write(data)
	})

	m.Post(SubmitTask+"/:task_provider_id", func(res http.ResponseWriter, req *http.Request, params martini.Params) {
		tpid := params["task_provider_id"]
		tp, ok := taskProviders[tpid]
		if !ok {
			lxlog.Errorf(logrus.Fields{
				"tp_id":  tpid,
			}, "task provider not found for tpid")
			res.WriteHeader(500)
		}
		body, err := ioutil.ReadAll(req.Body)
		if req.Body != nil {
			defer req.Body.Close()
		}
		if err != nil {
			lxlog.Errorf(logrus.Fields{
				"error": err,
				"body":  string(body),
			}, "could not read  request body")
			res.WriteHeader(500)
			return
		}
		var task lxtypes.Task
		err = json.Unmarshal(body, &task)
		if err != nil {
			lxlog.Errorf(logrus.Fields{
				"error": err,
				"body":  string(body),
			}, "could parse json into task")
			res.WriteHeader(500)
			return
		}
		task.TaskProvider = tp
		tasks[task.TaskId] = &task
		res.WriteHeader(202)
	})

	m.Post(KillTask+"/:task_id", func(res http.ResponseWriter, req *http.Request, params martini.Params) {
		taskid := params["task_id"]
		if _, ok := tasks[taskid]; !ok {
			lxlog.Errorf(logrus.Fields{
				"tpid": taskid,
			}, "task was not submitted")
			res.WriteHeader(400)
			return
		}
		tasks[taskid].KillRequested = true
		res.WriteHeader(202)
	})

	m.Post(PurgeTask+"/:task_id", func(res http.ResponseWriter, req *http.Request, params martini.Params) {
		taskid := params["task_id"]
		if _, ok := tasks[taskid]; !ok {
			lxlog.Errorf(logrus.Fields{
				"tpid": taskid,
			}, "task was not submitted")
			res.WriteHeader(400)
			return
		}
		delete(tasks, taskid)
		res.WriteHeader(202)
	})

	m.RunOnAddr(fmt.Sprintf(":%v", port))
}
