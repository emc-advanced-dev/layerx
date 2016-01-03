package lxserver


import (
	"github.com/go-martini/martini"
	"net/http"
	"github.com/layer-x/layerx-commons/lxactionqueue"
	"github.com/layer-x/layerx-core_v2/lxstate"
	"github.com/Sirupsen/logrus"
"github.com/layer-x/layerx-commons/lxlog"
	"encoding/json"
	"github.com/layer-x/layerx-core_v2/layerx_tpi_client"
"io/ioutil"
	"github.com/layer-x/layerx-core_v2/lxserver/lx_core_helpers"
	"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/layer-x/layerx-core_v2/layerx_rpi_client"
"github.com/layer-x/layerx-core_v2/lxtypes"
)

const (
//tpi
	RegisterTpi             = "/RegisterTpi"
	RegisterTaskProvider   = "/RegisterTaskProvider"
	DeregisterTaskProvider = "/DeregisterTaskProvider"
	GetTaskProviders       = "/GetTaskProviders"
	GetStatusUpdates       = "/GetStatusUpdates"
	SubmitTask             = "/SubmitTask"
	KillTask               = "/KillTask"
	PurgeTask              = "/PurgeTask"
//rpi
	RegisterRpi             = "/RegisterRpi"
	SubmitResource             = "/SubmitResource"
	SubmitStatusUpdate         = "/SubmitStatusUpdate"
	GetNodes         = "/GetNodes"
)

var empty = []byte{}

type layerxCoreServerWrapper struct {
	actionQueue      lxactionqueue.ActionQueue
	state			*lxstate.State
}

func NewLayerXCoreServerWrapper(state *lxstate.State, actionQueue lxactionqueue.ActionQueue) *layerxCoreServerWrapper {
	return &layerxCoreServerWrapper{
		actionQueue: actionQueue,
		state: state,
	}
}

func (wrapper *layerxCoreServerWrapper) WrapServer(m *martini.ClassicMartini, driverErrc chan error) *martini.ClassicMartini {
	registerTpiHandler := func(res http.ResponseWriter, req *http.Request) {
		fn := func() ([]byte, int, error) {
			data, err := ioutil.ReadAll(req.Body)
			if req.Body != nil {
				defer req.Body.Close()
			}
			if err != nil {
				return empty, 400, lxerrors.New("parsing register tpi request", err)
			}
			var tpiRegistrationMessage layerx_tpi_client.TpiRegistrationMessage
			err = json.Unmarshal(data, &tpiRegistrationMessage)
			if err != nil {
				return empty, 500, lxerrors.New("could not parse json to register tpi message", err)
			}
			err = lx_core_helpers.RegisterTpi(wrapper.state, tpiRegistrationMessage.TpiUrl)
			if err != nil {
				lxlog.Errorf(logrus.Fields{
					"error": err,
				}, "could not handle register tpi request")
				return empty, 500, lxerrors.New("could not handle register tpi request", err)
			}
			lxlog.Infof(logrus.Fields{"tpi_url": tpiRegistrationMessage.TpiUrl}, "Registered TPI to LayerX")
			return empty, 202, nil
		}
		_, statusCode, err := wrapper.queueOperation(fn)
		if err != nil {
			res.WriteHeader(statusCode)
			lxlog.Errorf(logrus.Fields{
				"error": err.Error(),
			}, "processing register tpi message")
			driverErrc <- err
			return
		}
		res.WriteHeader(statusCode)
	}

	registerRpiHandler := func(res http.ResponseWriter, req *http.Request) {
		fn := func() ([]byte, int, error) {
			data, err := ioutil.ReadAll(req.Body)
			if req.Body != nil {
				defer req.Body.Close()
			}
			if err != nil {
				return empty, 400, lxerrors.New("parsing register rpi request", err)
			}
			var rpiRegistrationMessage layerx_rpi_client.RpiRegistrationMessage
			err = json.Unmarshal(data, &rpiRegistrationMessage)
			if err != nil {
				return empty, 500, lxerrors.New("could not parse json to register rpi message", err)
			}
			err = lx_core_helpers.RegisterRpi(wrapper.state, rpiRegistrationMessage.RpiUrl)
			if err != nil {
				lxlog.Errorf(logrus.Fields{
					"error": err,
				}, "could not handle register rpi request")
				return empty, 500, lxerrors.New("could not handle register rpi request", err)
			}
			lxlog.Infof(logrus.Fields{"rpi_url": rpiRegistrationMessage.RpiUrl}, "Registered TPI to LayerX")
			return empty, 202, nil
		}
		_, statusCode, err := wrapper.queueOperation(fn)
		if err != nil {
			res.WriteHeader(statusCode)
			lxlog.Errorf(logrus.Fields{
				"error": err.Error(),
			}, "processing register rpi message")
			driverErrc <- err
			return
		}
		res.WriteHeader(statusCode)
	}

	registerTaskProviderHandler := func(res http.ResponseWriter, req *http.Request) {
		fn := func() ([]byte, int, error) {
			data, err := ioutil.ReadAll(req.Body)
			if req.Body != nil {
				defer req.Body.Close()
			}
			if err != nil {
				return empty, 400, lxerrors.New("parsing register TaskProvider request", err)
			}
			var taskProvider lxtypes.TaskProvider
			err = json.Unmarshal(data, &taskProvider)
			if err != nil {
				return empty, 500, lxerrors.New("could not parse json to register TaskProvider message", err)
			}
			err = lx_core_helpers.RegisterTaskProvider(wrapper.state, &taskProvider)
			if err != nil {
				lxlog.Errorf(logrus.Fields{
					"error": err,
				}, "could not handle register TaskProvider request")
				return empty, 500, lxerrors.New("could not handle register TaskProvider request", err)
			}
			lxlog.Infof(logrus.Fields{"task_provider": taskProvider}, "Added new TaskProvider to LayerX")
			return empty, 202, nil
		}
		_, statusCode, err := wrapper.queueOperation(fn)
		if err != nil {
			res.WriteHeader(statusCode)
			lxlog.Errorf(logrus.Fields{
				"error": err.Error(),
			}, "processing register TaskProvider message")
			driverErrc <- err
			return
		}
		res.WriteHeader(statusCode)
	}

	deregisterTaskProviderHandler := func(res http.ResponseWriter, req *http.Request, params martini.Params) {
		fn := func() ([]byte, int, error) {
			tpId := params["task_provider_id"]
			err := lx_core_helpers.DeregisterTaskProvider(wrapper.state, tpId)
			if err != nil {
				lxlog.Errorf(logrus.Fields{
					"error": err,
				}, "could not handle deregister TaskProvider request")
				return empty, 500, lxerrors.New("could not handle deregister TaskProvider request", err)
			}
			lxlog.Infof(logrus.Fields{"task_provider_id": tpId}, "removed task provider from LayerX")
			return empty, 202, nil
		}
		_, statusCode, err := wrapper.queueOperation(fn)
		if err != nil {
			res.WriteHeader(statusCode)
			lxlog.Errorf(logrus.Fields{
				"error": err.Error(),
			}, "processing deregister TaskProvider message")
			driverErrc <- err
			return
		}
		res.WriteHeader(statusCode)
	}

	getTaskProvidersHandler := func(res http.ResponseWriter, req *http.Request) {
		fn := func() ([]byte, int, error) {
			taskProviders, err := lx_core_helpers.GetTaskProviders(wrapper.state)
			if err != nil {
				lxlog.Errorf(logrus.Fields{
					"error": err,
				}, "could not handle Get TaskProviders request")
				return empty, 500, lxerrors.New("could not handle Get TaskProviders request", err)
			}
			data, err := json.Marshal(taskProviders)
			if err != nil {
				lxlog.Errorf(logrus.Fields{}, "could not marshal task providers to json")
				return empty, 500, lxerrors.New("marshalling taskProviders to json", err)
			}
			lxlog.Debugf(logrus.Fields{"task_providers": taskProviders}, "Added new TaskProvider to LayerX")
			return data, 200, nil
		}
		data, statusCode, err := wrapper.queueOperation(fn)
		if err != nil {
			res.WriteHeader(statusCode)
			lxlog.Errorf(logrus.Fields{
				"error": err.Error(),
			}, "processing Get TaskProviders message")
			driverErrc <- err
			return
		}
		res.Write(data)
	}

	getStatusUpdatesHandler := func(res http.ResponseWriter, req *http.Request, params martini.Params) {
		fn := func() ([]byte, int, error) {
			tpId := params["task_provider_id"]
			statuses, err := lx_core_helpers.GetStatusUpdates(wrapper.state, tpId)
			if err != nil {
				lxlog.Errorf(logrus.Fields{
					"error": err,
				}, "could not handle Get Status Updates request")
				return empty, 500, lxerrors.New("could not handle Get Status Updates request", err)
			}
			data, err := json.Marshal(statuses)
			if err != nil {
				lxlog.Errorf(logrus.Fields{}, "could not marshal Status Updates to json")
				return empty, 500, lxerrors.New("marshalling Statuses to json", err)
			}
			return data, 200, nil
		}
		data, statusCode, err := wrapper.queueOperation(fn)
		if err != nil {
			res.WriteHeader(statusCode)
			lxlog.Errorf(logrus.Fields{
				"error": err.Error(),
			}, "processing Get Status Updates message")
			driverErrc <- err
			return
		}
		res.Write(data)
	}

	submitTaskHandler := func(res http.ResponseWriter, req *http.Request, params martini.Params) {
		fn := func() ([]byte, int, error) {
			tpId := params["task_provider_id"]
			data, err := ioutil.ReadAll(req.Body)
			if req.Body != nil {
				defer req.Body.Close()
			}
			if err != nil {
				return empty, 400, lxerrors.New("parsing SubmitTask request", err)
			}
			var task lxtypes.Task
			err = json.Unmarshal(data, &task)
			if err != nil {
				return empty, 500, lxerrors.New("could not parse json to task", err)
			}
			err = lx_core_helpers.SubmitTask(wrapper.state, tpId, &task)
			if err != nil {
				lxlog.Errorf(logrus.Fields{
					"error": err,
				}, "could not handle Get SubmitTask request")
				return empty, 500, lxerrors.New("could not handle SubmitTask request", err)
			}
			lxlog.Infof(logrus.Fields{"task_provider_id": tpId, "task": task}, "accepted task from task provider")
			return empty, 202, nil
		}
		_, statusCode, err := wrapper.queueOperation(fn)
		if err != nil {
			res.WriteHeader(statusCode)
			lxlog.Errorf(logrus.Fields{
				"error": err.Error(),
			}, "processing Get Status Updates message")
			driverErrc <- err
			return
		}
		res.WriteHeader(statusCode)
	}

	m.Post(RegisterTpi, registerTpiHandler)
	m.Post(RegisterRpi, registerRpiHandler)
	m.Post(RegisterTaskProvider, registerTaskProviderHandler)
	m.Post(DeregisterTaskProvider+"/:task_provider_id", deregisterTaskProviderHandler)
	m.Get(GetTaskProviders, getTaskProvidersHandler)
	m.Get(GetStatusUpdates+"/:task_provider_id", getStatusUpdatesHandler)
	m.Post(SubmitTask+"/:task_provider_id", submitTaskHandler)

	m.Post(KillTask+"/:task_id", func(res http.ResponseWriter, req *http.Request, params martini.Params) {
//		taskid := params["task_id"]
//		if _, ok := tasks[taskid]; !ok {
//			lxlog.Errorf(logrus.Fields{
//				"tpid": taskid,
//			}, "task was not submitted")
//			res.WriteHeader(400)
//			return
//		}
//		tasks[taskid].KillRequested = true
//		res.WriteHeader(202)
	})

	m.Post(PurgeTask+"/:task_id", func(res http.ResponseWriter, req *http.Request, params martini.Params) {
//		taskid := params["task_id"]
//		if _, ok := tasks[taskid]; !ok {
//			lxlog.Errorf(logrus.Fields{
//				"tpid": taskid,
//			}, "task was not submitted")
//			res.WriteHeader(400)
//			return
//		}
//		delete(tasks, taskid)
//		res.WriteHeader(202)
	})

	m.Post(SubmitResource, func(res http.ResponseWriter, req *http.Request) {
//		body, err := ioutil.ReadAll(req.Body)
//		if req.Body != nil {
//			defer req.Body.Close()
//		}
//		if err != nil {
//			lxlog.Errorf(logrus.Fields{
//				"error": err,
//				"body":  string(body),
//			}, "could not read  request body")
//			res.WriteHeader(500)
//			return
//		}
//		var resource lxtypes.Resource
//		err = json.Unmarshal(body, &resource)
//		if err != nil {
//			lxlog.Errorf(logrus.Fields{
//				"error": err,
//				"body":  string(body),
//			}, "could parse json into resource")
//			res.WriteHeader(500)
//			return
//		}
//		nodeId := resource.NodeId
//		if knownNode, ok := nodes[nodeId]; ok {
//			err = knownNode.AddResource(&resource)
//			if err != nil {
//				lxlog.Errorf(logrus.Fields{
//					"error": err,
//					"node":  knownNode,
//					"resource":  resource,
//				}, "could not add resource to node")
//				res.WriteHeader(500)
//				return
//			}
//			nodes[nodeId] = knownNode
//		} else {
//			newNode := lxtypes.NewNode(nodeId)
//			err = newNode.AddResource(&resource)
//			if err != nil {
//				lxlog.Errorf(logrus.Fields{
//					"error": err,
//					"node":  newNode,
//					"resource":  resource,
//				}, "could not add resource to node")
//				res.WriteHeader(500)
//			}
//			nodes[nodeId] = newNode
//		}
//		res.WriteHeader(202)
	})

	m.Post(SubmitStatusUpdate, func(res http.ResponseWriter, req *http.Request) {
//		body, err := ioutil.ReadAll(req.Body)
//		if req.Body != nil {
//			defer req.Body.Close()
//		}
//		if err != nil {
//			lxlog.Errorf(logrus.Fields{
//				"error": err,
//				"body":  string(body),
//			}, "could not read  request body")
//			res.WriteHeader(500)
//			return
//		}
//		var status mesosproto.TaskStatus
//		err = proto.Unmarshal(body, &status)
//		if err != nil {
//			lxlog.Errorf(logrus.Fields{
//				"error": err,
//				"body":  string(body),
//			}, "could parse proto into resource")
//			res.WriteHeader(500)
//			return
//		}
//		taskId := status.GetTaskId().GetValue()
//		statusUpdates[taskId] = &status
//		res.WriteHeader(202)
	})

	m.Get(GetNodes, func(res http.ResponseWriter){
//		nodeArr := []*lxtypes.Node{}
//		for _, node := range nodes {
//			nodeArr = append(nodeArr, node)
//		}
//		data, err := json.Marshal(nodeArr)
//		if err != nil {
//			lxlog.Errorf(logrus.Fields{
//				"error": err,
//				"data":  string(data),
//			}, "could marshal nodes to json")
//			res.WriteHeader(500)
//			return
//		}
//		res.Write(data)
	})

	return m
}



func (wrapper *layerxCoreServerWrapper) queueOperation(f func() ([]byte, int, error)) ([]byte, int, error) {
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