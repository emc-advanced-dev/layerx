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
	"github.com/gogo/protobuf/proto"
	"github.com/layer-x/layerx-core_v2/layerx_rpi_client"
	"github.com/layer-x/layerx-core_v2/layerx_tpi_client"
	"github.com/layer-x/layerx-core_v2/layerx_brain_client"
)

const (
//tpi
	RegisterTpi = "/RegisterTpi"
	RegisterTaskProvider = "/RegisterTaskProvider"
	DeregisterTaskProvider = "/DeregisterTaskProvider"
	GetTaskProviders = "/GetTaskProviders"
	GetStatusUpdates = "/GetStatusUpdates"
	GetStatusUpdate = "/GetStatusUpdate"
	SubmitTask = "/SubmitTask"
	KillTask = "/KillTask"
	PurgeTask = "/PurgeTask"
//rpi
	RegisterRpi = "/RegisterRpi"
	SubmitResource = "/SubmitResource"
	SubmitStatusUpdate = "/SubmitStatusUpdate"
	//brain
	GetNodes = "/GetNodes"
	GetPendingTasks = "/GetPendingTasks"
	AssignTasks = "/AssignTasks"
	MigrateTask = "/MigrateTask"

	//for testing
	Purge = "/Purge"
)

func RunFakeLayerXServer(fakeStatuses []*mesosproto.TaskStatus, port int) {
	taskProviders := make(map[string]*lxtypes.TaskProvider)
	statusUpdates := make(map[string]*mesosproto.TaskStatus)
	tasks := make(map[string]*lxtypes.Task)
	nodes := make(map[string]*lxtypes.Node)

	for _, status := range fakeStatuses {
		statusUpdates[status.GetTaskId().GetValue()] = status
	}

	m := martini.Classic()

	//TPI
	m.Post(RegisterTpi, func(res http.ResponseWriter, req *http.Request) {
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
		var registrationMessage layerx_tpi_client.TpiRegistrationMessage
		err = json.Unmarshal(body, &registrationMessage)
		if err != nil {
			lxlog.Errorf(logrus.Fields{
				"error": err,
				"body":  string(body),
			}, "could parse json into resource")
			res.WriteHeader(500)
			return
		}
		res.WriteHeader(202)
	})

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
	m.Post(DeregisterTaskProvider + "/:task_provider_id", func(res http.ResponseWriter, req *http.Request, params martini.Params) {
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
	m.Get(GetStatusUpdates+"/:task_provider_id", func(res http.ResponseWriter, req *http.Request, params martini.Params) {
		tpid := params["task_provider_id"]
		statuses := []*mesosproto.TaskStatus{}
		for _, status := range statusUpdates {
			taskId := status.GetTaskId().GetValue()
			task, ok := tasks[taskId]
			if !ok {
				lxlog.Errorf(logrus.Fields{
					"task_id":  taskId,
				}, "could not find task for the id in the status")
				res.WriteHeader(500)
			}
			if task.TaskProvider.Id == tpid {
				statuses = append(statuses, status)
			}
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

	m.Get(GetStatusUpdate + "/:task_id", func(res http.ResponseWriter, req *http.Request, params martini.Params) {
		taskId := params["task_id"]
		status, ok := statusUpdates[taskId]
		if !ok {
			lxlog.Errorf(logrus.Fields{
				"task_id":  taskId,
			}, "could not find status for the id in the status")
			res.WriteHeader(500)
		}
		data, err := json.Marshal(status)
		if err != nil {
			lxlog.Errorf(logrus.Fields{
				"error": err,
				"body":  string(data),
			}, "could parse status into json")
			res.WriteHeader(500)
			return
		}
		res.Write(data)
	})

	m.Post(SubmitTask + "/:task_provider_id", func(res http.ResponseWriter, req *http.Request, params martini.Params) {
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

	m.Post(KillTask + "/:task_id", func(res http.ResponseWriter, req *http.Request, params martini.Params) {
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

	m.Post(PurgeTask + "/:task_id", func(res http.ResponseWriter, req *http.Request, params martini.Params) {
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

	//RPI
	m.Post(RegisterRpi, func(res http.ResponseWriter, req *http.Request) {
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
		var registrationMessage layerx_rpi_client.RpiRegistrationMessage
		err = json.Unmarshal(body, &registrationMessage)
		if err != nil {
			lxlog.Errorf(logrus.Fields{
				"error": err,
				"body":  string(body),
			}, "could parse json into resource")
			res.WriteHeader(500)
			return
		}
		res.WriteHeader(202)
	})

	m.Post(SubmitResource, func(res http.ResponseWriter, req *http.Request) {
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
		var resource lxtypes.Resource
		err = json.Unmarshal(body, &resource)
		if err != nil {
			lxlog.Errorf(logrus.Fields{
				"error": err,
				"body":  string(body),
			}, "could parse json into resource")
			res.WriteHeader(500)
			return
		}
		nodeId := resource.NodeId
		if knownNode, ok := nodes[nodeId]; ok {
			err = knownNode.AddResource(&resource)
			if err != nil {
				lxlog.Errorf(logrus.Fields{
					"error": err,
					"node":  knownNode,
					"resource":  resource,
				}, "could not add resource to node")
				res.WriteHeader(500)
				return
			}
			nodes[nodeId] = knownNode
		} else {
			newNode := lxtypes.NewNode(nodeId)
			err = newNode.AddResource(&resource)
			if err != nil {
				lxlog.Errorf(logrus.Fields{
					"error": err,
					"node":  newNode,
					"resource":  resource,
				}, "could not add resource to node")
				res.WriteHeader(500)
			}
			nodes[nodeId] = newNode
		}
		res.WriteHeader(202)
	})

	m.Post(SubmitStatusUpdate, func(res http.ResponseWriter, req *http.Request) {
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
		var status mesosproto.TaskStatus
		err = proto.Unmarshal(body, &status)
		if err != nil {
			lxlog.Errorf(logrus.Fields{
				"error": err,
				"body":  string(body),
			}, "could parse proto into resource")
			res.WriteHeader(500)
			return
		}
		taskId := status.GetTaskId().GetValue()
		statusUpdates[taskId] = &status
		res.WriteHeader(202)
	})

	m.Get(GetNodes, func(res http.ResponseWriter) {
		nodeArr := []*lxtypes.Node{}
		for _, node := range nodes {
			nodeArr = append(nodeArr, node)
		}
		data, err := json.Marshal(nodeArr)
		if err != nil {
			lxlog.Errorf(logrus.Fields{
				"error": err,
				"data":  string(data),
			}, "could marshal nodes to json")
			res.WriteHeader(500)
			return
		}
		res.Write(data)
	})

	m.Get(GetPendingTasks, func(res http.ResponseWriter) {
		taskArr := []*lxtypes.Task{}
		for _, task := range tasks {
			if task.SlaveId == "" {
				taskArr = append(taskArr, task)
			}
		}
		data, err := json.Marshal(taskArr)
		if err != nil {
			lxlog.Errorf(logrus.Fields{
				"error": err,
				"data":  string(data),
			}, "could marshal tasks to json")
			res.WriteHeader(500)
			return
		}
		res.Write(data)
	})

	m.Post(AssignTasks, func(res http.ResponseWriter, req *http.Request) {
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
		var brainAssignmentMessage layerx_brain_client.BrainAssignTasksMessage
		err = json.Unmarshal(body, &brainAssignmentMessage)
		if err != nil {
			lxlog.Errorf(logrus.Fields{
				"error": err,
				"body":  string(body),
			}, "could parse json into brainAssignmentMessage")
			res.WriteHeader(500)
			return
		}
		_, ok := nodes[brainAssignmentMessage.NodeId]
		if !ok {
			lxlog.Errorf(logrus.Fields{
				"node_id": brainAssignmentMessage.NodeId,
			}, "invalid node id")
			res.WriteHeader(400)
		}
		for _, taskId := range brainAssignmentMessage.TaskIds {
			task, ok := tasks[taskId]
			if !ok {
				lxlog.Errorf(logrus.Fields{
					"task_id": taskId,
				}, "invalid task id")
				res.WriteHeader(400)
			}
			task.SlaveId = brainAssignmentMessage.NodeId
		}
		res.WriteHeader(202)
	})

	m.Post(MigrateTask, func(res http.ResponseWriter, req *http.Request) {
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
		var migrateMessage layerx_brain_client.MigrateTaskMessage
		err = json.Unmarshal(body, &migrateMessage)
		if err != nil {
			lxlog.Errorf(logrus.Fields{
				"error": err,
				"body":  string(body),
			}, "could parse json into brainAssignmentMessage")
			res.WriteHeader(500)
			return
		}
		_, ok := nodes[migrateMessage.DestinationNodeId]
		if !ok {
			lxlog.Errorf(logrus.Fields{
				"node_id": migrateMessage.DestinationNodeId,
			}, "invalid destinationNodeId node id")
			res.WriteHeader(400)
		}
		for _, taskId := range migrateMessage.TaskIds {
			task, ok := tasks[taskId]
			if !ok {
				lxlog.Errorf(logrus.Fields{
					"task_id": taskId,
				}, "invalid task id")
				res.WriteHeader(400)
			}
			task.SlaveId = migrateMessage.DestinationNodeId
			res.WriteHeader(202)
		}
	})

	m.Post(Purge, func(){
		taskProviders = make(map[string]*lxtypes.TaskProvider)
		tasks = make(map[string]*lxtypes.Task)
		nodes = make(map[string]*lxtypes.Node)
	})

	m.RunOnAddr(fmt.Sprintf(":%v", port))
}
