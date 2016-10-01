package fakes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/emc-advanced-dev/layerx/layerx-core/layerx_brain_client"
	"github.com/emc-advanced-dev/layerx/layerx-core/layerx_rpi_client"
	"github.com/emc-advanced-dev/layerx/layerx-core/layerx_tpi_client"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"
	"github.com/go-martini/martini"
	"github.com/gogo/protobuf/proto"
	"github.com/mesos/mesos-go/mesosproto"
)

const (
	//tpi
	RegisterTpi            = "/RegisterTpi"
	RegisterTaskProvider   = "/RegisterTaskProvider"
	DeregisterTaskProvider = "/DeregisterTaskProvider"
	GetTaskProviders       = "/GetTaskProviders"
	GetStatusUpdates       = "/GetStatusUpdates"
	GetStatusUpdate        = "/GetStatusUpdate"
	SubmitTask             = "/SubmitTask"
	KillTask               = "/KillTask"
	PurgeTask              = "/PurgeTask"
	//rpi
	RegisterRpi        = "/RegisterRpi"
	SubmitResource     = "/SubmitResource"
	SubmitStatusUpdate = "/SubmitStatusUpdate"
	//brain
	GetNodes        = "/GetNodes"
	GetPendingTasks = "/GetPendingTasks"
	GetStagingTasks = "/GetStagingTasks"
	AssignTasks     = "/AssignTasks"
	MigrateTasks    = "/MigrateTasks"

	//for testing
	Purge = "/Purge"
)

type FakeCore struct {
	TaskProviders map[string]*lxtypes.TaskProvider
	StatusUpdates map[string]*mesosproto.TaskStatus
	Tasks         map[string]*lxtypes.Task
	StagingTasks  map[string]*lxtypes.Task
	Nodes         map[string]*lxtypes.Node
}

func NewFakeCore() *FakeCore {
	return &FakeCore{
		TaskProviders: make(map[string]*lxtypes.TaskProvider),
		StatusUpdates: make(map[string]*mesosproto.TaskStatus),
		Tasks:         make(map[string]*lxtypes.Task),
		StagingTasks:  make(map[string]*lxtypes.Task),
		Nodes:         make(map[string]*lxtypes.Node),
	}
}

func (core *FakeCore) Start(fakeStatuses []*mesosproto.TaskStatus, port int) {
	for _, status := range fakeStatuses {
		core.StatusUpdates[status.GetTaskId().GetValue()] = status
	}

	m := martini.Classic()

	//TPI
	m.Post(RegisterTpi, func(res http.ResponseWriter, req *http.Request) {
		body, err := ioutil.ReadAll(req.Body)
		if req.Body != nil {
			defer req.Body.Close()
		}
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
				"body":  string(body),
			}).Errorf("could not read  request body")
			res.WriteHeader(500)
			return
		}
		var registrationMessage layerx_tpi_client.TpiRegistrationMessage
		err = json.Unmarshal(body, &registrationMessage)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
				"body":  string(body),
			}).Errorf("could parse json into resource")
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
			logrus.WithFields(logrus.Fields{
				"error": err,
				"body":  string(body),
			}).Errorf("could not read  request body")
			res.WriteHeader(500)
			return
		}
		var tp lxtypes.TaskProvider
		err = json.Unmarshal(body, &tp)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
				"body":  string(body),
			}).Errorf("could parse json into tp")
			res.WriteHeader(500)
			return
		}
		core.TaskProviders[tp.Id] = &tp
		res.WriteHeader(202)
	})
	m.Post(DeregisterTaskProvider+"/:task_provider_id", func(res http.ResponseWriter, req *http.Request, params martini.Params) {
		tpid := params["task_provider_id"]
		if _, ok := core.TaskProviders[tpid]; !ok {
			logrus.WithFields(logrus.Fields{
				"tpid": tpid,
			}).Errorf("task provider was not registered")
			res.WriteHeader(400)
			return
		}
		delete(core.TaskProviders, tpid)
		res.WriteHeader(202)
	})
	m.Get(GetTaskProviders, func(res http.ResponseWriter, req *http.Request) {
		tps := []*lxtypes.TaskProvider{}
		for _, tp := range core.TaskProviders {
			tps = append(tps, tp)
		}
		data, err := json.Marshal(tps)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
				"body":  string(data),
			}).Errorf("could parse tps into json")
			res.WriteHeader(500)
			return
		}
		res.Write(data)
	})
	m.Get(GetStatusUpdates+"/:task_provider_id", func(res http.ResponseWriter, req *http.Request, params martini.Params) {
		tpid := params["task_provider_id"]
		statuses := []*mesosproto.TaskStatus{}
		for _, status := range core.StatusUpdates {
			taskId := status.GetTaskId().GetValue()
			task, ok := core.Tasks[taskId]
			if !ok {
				logrus.WithFields(logrus.Fields{
					"task_id": taskId,
				}).Errorf("could not find task for the id in the status")
				res.WriteHeader(500)
			}
			if task.TaskProvider.Id == tpid {
				statuses = append(statuses, status)
			}
		}
		data, err := json.Marshal(statuses)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
				"body":  string(data),
			}).Errorf("could parse statuses into json")
			res.WriteHeader(500)
			return
		}
		res.Write(data)
	})
	m.Get(GetStatusUpdates, func(res http.ResponseWriter, req *http.Request, params martini.Params) {
		statuses := []*mesosproto.TaskStatus{}
		for _, status := range core.StatusUpdates {
			taskId := status.GetTaskId().GetValue()
			_, ok := core.Tasks[taskId]
			if !ok {
				logrus.WithFields(logrus.Fields{
					"task_id": taskId,
				}).Errorf("could not find task for the id in the status")
				res.WriteHeader(500)
			}
			statuses = append(statuses, status)
		}
		data, err := json.Marshal(statuses)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
				"body":  string(data),
			}).Errorf("could parse statuses into json")
			res.WriteHeader(500)
			return
		}
		res.Write(data)
	})

	m.Get(GetStatusUpdate+"/:task_id", func(res http.ResponseWriter, req *http.Request, params martini.Params) {
		taskId := params["task_id"]
		status, ok := core.StatusUpdates[taskId]
		if !ok {
			logrus.WithFields(logrus.Fields{
				"task_id": taskId,
			}).Errorf("could not find status for the id in the status")
			res.WriteHeader(500)
		}
		data, err := json.Marshal(status)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
				"body":  string(data),
			}).Errorf("could parse status into json")
			res.WriteHeader(500)
			return
		}
		res.Write(data)
	})

	m.Post(SubmitTask+"/:task_provider_id", func(res http.ResponseWriter, req *http.Request, params martini.Params) {
		tpid := params["task_provider_id"]
		tp, ok := core.TaskProviders[tpid]
		if !ok {
			logrus.WithFields(logrus.Fields{
				"tp_id": tpid,
			}).Errorf("task provider not found for tpid")
			res.WriteHeader(500)
		}
		body, err := ioutil.ReadAll(req.Body)
		if req.Body != nil {
			defer req.Body.Close()
		}
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
				"body":  string(body),
			}).Errorf("could not read  request body")
			res.WriteHeader(500)
			return
		}
		var task lxtypes.Task
		err = json.Unmarshal(body, &task)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
				"body":  string(body),
			}).Errorf("could parse json into task")
			res.WriteHeader(500)
			return
		}
		task.TaskProvider = tp
		core.Tasks[task.TaskId] = &task
		res.WriteHeader(202)
	})

	m.Post(KillTask+"/:tpid/:task_id", func(res http.ResponseWriter, req *http.Request, params martini.Params) {
		taskid := params["task_id"]
		tpid := params["framework_id"]
		if _, ok := core.Tasks[taskid]; !ok {
			logrus.WithFields(logrus.Fields{
				"taskid": taskid,
				"tpid":   tpid,
			}).Errorf("task was not submitted")
			res.WriteHeader(400)
			return
		}
		core.Tasks[taskid].KillRequested = true
		res.WriteHeader(202)
	})

	m.Post(PurgeTask+"/:task_id", func(res http.ResponseWriter, req *http.Request, params martini.Params) {
		taskid := params["task_id"]
		if _, ok := core.Tasks[taskid]; !ok {
			logrus.WithFields(logrus.Fields{
				"tpid": taskid,
			}).Errorf("task was not submitted")
			res.WriteHeader(400)
			return
		}
		delete(core.Tasks, taskid)
		res.WriteHeader(202)
	})

	//RPI
	m.Post(RegisterRpi, func(res http.ResponseWriter, req *http.Request) {
		body, err := ioutil.ReadAll(req.Body)
		if req.Body != nil {
			defer req.Body.Close()
		}
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
				"body":  string(body),
			}).Errorf("could not read  request body")
			res.WriteHeader(500)
			return
		}
		var registrationMessage layerx_rpi_client.RpiInfo
		err = json.Unmarshal(body, &registrationMessage)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
				"body":  string(body),
			}).Errorf("could parse json into resource")
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
			logrus.WithFields(logrus.Fields{
				"error": err,
				"body":  string(body),
			}).Errorf("could not read  request body")
			res.WriteHeader(500)
			return
		}
		var resource lxtypes.Resource
		err = json.Unmarshal(body, &resource)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
				"body":  string(body),
			}).Errorf("could parse json into resource")
			res.WriteHeader(500)
			return
		}
		nodeId := resource.NodeId
		if knownNode, ok := core.Nodes[nodeId]; ok {
			err = knownNode.AddResource(&resource)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error":    err,
					"node":     knownNode,
					"resource": resource,
				}).Errorf("could not add resource to node")
				res.WriteHeader(500)
				return
			}
			core.Nodes[nodeId] = knownNode
		} else {
			newNode := lxtypes.NewNode(nodeId)
			err = newNode.AddResource(&resource)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error":    err,
					"node":     newNode,
					"resource": resource,
				}).Errorf("could not add resource to node")
				res.WriteHeader(500)
			}
			core.Nodes[nodeId] = newNode
		}
		res.WriteHeader(202)
	})

	m.Post(SubmitStatusUpdate, func(res http.ResponseWriter, req *http.Request) {
		body, err := ioutil.ReadAll(req.Body)
		if req.Body != nil {
			defer req.Body.Close()
		}
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
				"body":  string(body),
			}).Errorf("could not read  request body")
			res.WriteHeader(500)
			return
		}
		var status mesosproto.TaskStatus
		err = proto.Unmarshal(body, &status)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
				"body":  string(body),
			}).Errorf("could parse proto into resource")
			res.WriteHeader(500)
			return
		}
		taskId := status.GetTaskId().GetValue()
		core.StatusUpdates[taskId] = &status
		res.WriteHeader(202)
	})

	m.Get(GetNodes, func(res http.ResponseWriter) {
		nodeArr := []*lxtypes.Node{}
		for _, node := range core.Nodes {
			nodeArr = append(nodeArr, node)
		}
		data, err := json.Marshal(nodeArr)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
				"data":  string(data),
			}).Errorf("could marshal core.Nodes to json")
			res.WriteHeader(500)
			return
		}
		res.Write(data)
	})

	m.Get(GetPendingTasks, func(res http.ResponseWriter) {
		taskArr := []*lxtypes.Task{}
		for _, task := range core.Tasks {
			taskArr = append(taskArr, task)
		}
		data, err := json.Marshal(taskArr)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
				"data":  string(data),
			}).Errorf("could marshal core.Tasks to json")
			res.WriteHeader(500)
			return
		}
		res.Write(data)
	})

	m.Get(GetStagingTasks, func(res http.ResponseWriter) {
		taskArr := []*lxtypes.Task{}
		logrus.WithFields(logrus.Fields{"core.StagingTasks": core.StagingTasks}).Infof("GETSTAGINGTASKS current staging core.Tasks pool")
		for _, task := range core.StagingTasks {
			taskArr = append(taskArr, task)
		}
		data, err := json.Marshal(taskArr)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
				"data":  string(data),
			}).Errorf("could marshal core.Tasks to json")
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
			logrus.WithFields(logrus.Fields{
				"error": err,
				"body":  string(body),
			}).Errorf("could not read  request body")
			res.WriteHeader(500)
			return
		}
		var brainAssignmentMessage layerx_brain_client.BrainAssignTasksMessage
		err = json.Unmarshal(body, &brainAssignmentMessage)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
				"body":  string(body),
			}).Errorf("could parse json into brainAssignmentMessage")
			res.WriteHeader(500)
			return
		}
		node, ok := core.Nodes[brainAssignmentMessage.NodeId]
		if !ok {
			logrus.WithFields(logrus.Fields{
				"node_id": brainAssignmentMessage.NodeId,
			}).Errorf("invalid node id")
			res.WriteHeader(400)
		}
		for _, taskId := range brainAssignmentMessage.TaskIds {
			task, ok := core.Tasks[taskId]
			if !ok {
				logrus.WithFields(logrus.Fields{
					"task_id": taskId,
				}).Errorf("invalid task id")
				res.WriteHeader(400)
			}
			task.NodeId = brainAssignmentMessage.NodeId
			core.StagingTasks[taskId] = task
			delete(core.Tasks, taskId)
			logrus.WithFields(logrus.Fields{"core.StagingTasks": core.StagingTasks}).Infof("current staging core.Tasks pool")
			err = node.AddTask(task)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"node_id": brainAssignmentMessage.NodeId,
				}).Errorf("could not add task to node")
				res.WriteHeader(400)
			} else {
				logrus.WithFields(logrus.Fields{"task": task, "node": node}).Infof("added task to node")
				go func() {
					//delay this for testing
					time.Sleep(3 * time.Second)
					delete(core.StagingTasks, taskId)
				}()
			}
		}
		res.WriteHeader(202)
	})

	m.Post(MigrateTasks, func(res http.ResponseWriter, req *http.Request) {
		body, err := ioutil.ReadAll(req.Body)
		if req.Body != nil {
			defer req.Body.Close()
		}
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
				"body":  string(body),
			}).Errorf("could not read  request body")
			res.WriteHeader(500)
			return
		}
		var migrateMessage layerx_brain_client.MigrateTaskMessage
		err = json.Unmarshal(body, &migrateMessage)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
				"body":  string(body),
			}).Errorf("could parse json into brainAssignmentMessage")
			res.WriteHeader(500)
			return
		}
		targetNode, ok := core.Nodes[migrateMessage.DestinationNodeId]
		if !ok {
			logrus.WithFields(logrus.Fields{
				"node_id": migrateMessage.DestinationNodeId,
			}).Errorf("invalid destinationNodeId node id")
			res.WriteHeader(400)
			return
		}
		for _, taskId := range migrateMessage.TaskIds {
			var task *lxtypes.Task
			var sourceNode *lxtypes.Node
			for _, node := range core.Nodes {
				logrus.WithFields(logrus.Fields{"task_id": taskId, "node": node}).Infof("searching node for task")
				task = node.GetTask(taskId)
				sourceNode = node
				if task != nil {
					break
				}
			}
			if task == nil {
				logrus.WithFields(logrus.Fields{"task_id": taskId, "core.Nodes": core.Nodes}).Errorf("invalid task id")
				res.WriteHeader(400)
				return
			}
			task.NodeId = migrateMessage.DestinationNodeId
			err = sourceNode.RemoveTask(taskId)
			if err != nil {
				logrus.WithFields(logrus.Fields{"task_id": taskId, "node": sourceNode}).Errorf("could not remove task from node")
				res.WriteHeader(400)
				return
			}
			core.StagingTasks[taskId] = task
			go func() {
				logrus.Debugf("in 3 seconds, moving from staging to running on the node")
				time.Sleep(1 * time.Second)
				err = targetNode.AddTask(task)
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"node_id": task.NodeId,
					}).Errorf("could not add task to node")
					res.WriteHeader(400)
					return
				} else {
					delete(core.StagingTasks, taskId)
				}
			}()

		}
		res.WriteHeader(202)
	})

	m.Post(Purge, func() {
		core.TaskProviders = make(map[string]*lxtypes.TaskProvider)
		core.Tasks = make(map[string]*lxtypes.Task)
		core.StagingTasks = make(map[string]*lxtypes.Task)
		core.Nodes = make(map[string]*lxtypes.Node)
	})

	m.RunOnAddr(fmt.Sprintf(":%v", port))
}
