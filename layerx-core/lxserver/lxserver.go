package lxserver

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/emc-advanced-dev/layerx/layerx-core/layerx_brain_client"
	"github.com/emc-advanced-dev/layerx/layerx-core/layerx_rpi_client"
	"github.com/emc-advanced-dev/layerx/layerx-core/layerx_tpi_client"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxserver/lx_core_helpers"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxstate"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"
	"github.com/emc-advanced-dev/pkg/errors"
	"github.com/go-martini/martini"
	"github.com/golang/protobuf/proto"
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
	GetPendingTasks = "/GetPendingTasks"
	GetStagingTasks = "/GetStagingTasks"
	GetNodes        = "/GetNodes"
	AssignTasks     = "/AssignTasks"
	MigrateTasks    = "/MigrateTasks"
)

var empty = []byte{}

type layerxCoreServerWrapper struct {
	state      *lxstate.State
	m          *martini.ClassicMartini
	driverErrc chan error
}

func NewLayerXCoreServerWrapper(state *lxstate.State, m *martini.ClassicMartini, driverErrc chan error) *layerxCoreServerWrapper {
	return &layerxCoreServerWrapper{
		state:      state,
		m:          m,
		driverErrc: driverErrc,
	}
}

func (wrapper *layerxCoreServerWrapper) WrapServer() *martini.ClassicMartini {
	registerTpiHandler := func(res http.ResponseWriter, req *http.Request) {
		fn := func() ([]byte, int, error) {
			data, err := ioutil.ReadAll(req.Body)
			if req.Body != nil {
				defer req.Body.Close()
			}
			if err != nil {
				return empty, 400, errors.New("parsing register tpi request", err)
			}
			var tpiRegistrationMessage layerx_tpi_client.TpiRegistrationMessage
			err = json.Unmarshal(data, &tpiRegistrationMessage)
			if err != nil {
				return empty, 500, errors.New("could not parse json to register tpi message", err)
			}
			err = lx_core_helpers.RegisterTpi(wrapper.state, tpiRegistrationMessage.TpiUrl)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Errorf("could not handle register tpi request")
				return empty, 500, errors.New("could not handle register tpi request", err)
			}
			logrus.WithFields(logrus.Fields{"tpi_url": tpiRegistrationMessage.TpiUrl}).Infof("Registered TPI to LayerX")
			return empty, 202, nil
		}
		_, statusCode, err := wrapper.doOperation(fn)
		if err != nil {
			res.WriteHeader(statusCode)
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Errorf("processing register tpi message")
			wrapper.driverErrc <- err
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
				return empty, 400, errors.New("parsing register rpi request", err)
			}
			var rpiRegistrationMessage layerx_rpi_client.RpiInfo
			err = json.Unmarshal(data, &rpiRegistrationMessage)
			if err != nil {
				return empty, 500, errors.New("could not parse json to register rpi message", err)
			}
			err = lx_core_helpers.RegisterRpi(wrapper.state, rpiRegistrationMessage)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Errorf("could not handle register rpi request")
				return empty, 500, errors.New("could not handle register rpi request", err)
			}
			logrus.WithFields(logrus.Fields{"rpi_url": rpiRegistrationMessage.Url}).Infof("Registered TPI to LayerX")
			return empty, 202, nil
		}
		_, statusCode, err := wrapper.doOperation(fn)
		if err != nil {
			res.WriteHeader(statusCode)
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Errorf("processing register rpi message")
			wrapper.driverErrc <- err
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
				return empty, 400, errors.New("parsing register TaskProvider request", err)
			}
			var taskProvider lxtypes.TaskProvider
			err = json.Unmarshal(data, &taskProvider)
			if err != nil {
				return empty, 500, errors.New("could not parse json to register TaskProvider message", err)
			}
			err = lx_core_helpers.RegisterTaskProvider(wrapper.state, &taskProvider)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Errorf("could not handle register TaskProvider request")
				return empty, 500, errors.New("could not handle register TaskProvider request", err)
			}
			logrus.WithFields(logrus.Fields{"task_provider": taskProvider}).Infof("Added new TaskProvider to LayerX")
			return empty, 202, nil
		}
		_, statusCode, err := wrapper.doOperation(fn)
		if err != nil {
			res.WriteHeader(statusCode)
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Errorf("processing register TaskProvider message")
			wrapper.driverErrc <- err
			return
		}
		res.WriteHeader(statusCode)
	}

	deregisterTaskProviderHandler := func(res http.ResponseWriter, req *http.Request, params martini.Params) {
		fn := func() ([]byte, int, error) {
			tpId := params["task_provider_id"]
			err := lx_core_helpers.DeregisterTaskProvider(wrapper.state, tpId)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Errorf("could not handle deregister TaskProvider request")
				return empty, 500, errors.New("could not handle deregister TaskProvider request", err)
			}
			logrus.WithFields(logrus.Fields{"task_provider_id": tpId}).Infof("removed task provider from LayerX")
			return empty, 202, nil
		}
		_, statusCode, err := wrapper.doOperation(fn)
		if err != nil {
			res.WriteHeader(statusCode)
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Errorf("processing deregister TaskProvider message")
			wrapper.driverErrc <- err
			return
		}
		res.WriteHeader(statusCode)
	}

	getTaskProvidersHandler := func(res http.ResponseWriter, req *http.Request) {
		fn := func() ([]byte, int, error) {
			taskProviders, err := lx_core_helpers.GetTaskProviders(wrapper.state)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Errorf("could not handle Get TaskProviders request")
				return empty, 500, errors.New("could not handle Get TaskProviders request", err)
			}
			data, err := json.Marshal(taskProviders)
			if err != nil {
				logrus.Errorf("could not marshal task providers to json")
				return empty, 500, errors.New("marshalling taskProviders to json", err)
			}
			logrus.WithFields(logrus.Fields{"task_providers": taskProviders}).Debugf("Added new TaskProvider to LayerX")
			return data, 200, nil
		}
		data, statusCode, err := wrapper.doOperation(fn)
		if err != nil {
			res.WriteHeader(statusCode)
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Errorf("processing Get TaskProviders message")
			wrapper.driverErrc <- err
			return
		}
		res.Write(data)
	}

	getStatusUpdatesHandler := func(res http.ResponseWriter, req *http.Request, params martini.Params) {
		fn := func() ([]byte, int, error) {
			tpId := params["task_provider_id"]
			statuses, err := lx_core_helpers.GetStatusUpdates(wrapper.state, tpId)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Errorf("could not handle Get Status Updates request")
				return empty, 500, errors.New("could not handle Get Status Updates request", err)
			}
			data, err := json.Marshal(statuses)
			if err != nil {
				logrus.Errorf("could not marshal Status Updates to json")
				return empty, 500, errors.New("marshalling Statuses to json", err)
			}
			logrus.WithFields(logrus.Fields{"statuses": statuses}).Debugf("Replying with Statuses")
			return data, 200, nil
		}
		data, statusCode, err := wrapper.doOperation(fn)
		if err != nil {
			res.WriteHeader(statusCode)
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Errorf("processing Get Status Updates message")
			wrapper.driverErrc <- err
			return
		}
		res.Write(data)
	}

	getStatusUpdateHandler := func(res http.ResponseWriter, req *http.Request, params martini.Params) {
		fn := func() ([]byte, int, error) {
			taskId := params["task_id"]
			status, err := lx_core_helpers.GetStatusUpdate(wrapper.state, taskId)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Errorf("could not handle Get Status Update request")
				return empty, 500, errors.New("could not handle Get Status Update request", err)
			}
			data, err := json.Marshal(status)
			if err != nil {
				logrus.Errorf("could not marshal Status Update to json")
				return empty, 500, errors.New("marshalling Statuses to json", err)
			}
			return data, 200, nil
		}
		data, statusCode, err := wrapper.doOperation(fn)
		if err != nil {
			res.WriteHeader(statusCode)
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Errorf("processing Get Status Update message")
			wrapper.driverErrc <- err
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
				return empty, 400, errors.New("parsing SubmitTask request", err)
			}
			var task lxtypes.Task
			err = json.Unmarshal(data, &task)
			if err != nil {
				return empty, 500, errors.New("could not parse json to task", err)
			}
			err = lx_core_helpers.SubmitTask(wrapper.state, tpId, &task)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Errorf("could not handle Get SubmitTask request")
				return empty, 500, errors.New("could not handle SubmitTask request", err)
			}
			logrus.WithFields(logrus.Fields{"task_provider_id": tpId, "task": task}).Infof("accepted task from task provider")
			return empty, 202, nil
		}
		_, statusCode, err := wrapper.doOperation(fn)
		if err != nil {
			res.WriteHeader(statusCode)
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Errorf("processing Get Status Updates message")
			wrapper.driverErrc <- err
			return
		}
		res.WriteHeader(statusCode)
	}

	killTaskHandler := func(res http.ResponseWriter, req *http.Request, params martini.Params) {
		fn := func() ([]byte, int, error) {
			taskId := params["task_id"]
			taskProviderId := params["task_provider_id"]
			err := lx_core_helpers.KillTask(wrapper.state, wrapper.state.GetTpiUrl(), taskProviderId, taskId)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Errorf("could not handle kill task request")
				return empty, 500, errors.New("could not handle kill task request", err)
			}
			logrus.WithFields(logrus.Fields{"task_id": taskId, "task_provider_id": taskProviderId}).Infof("killed task")
			return empty, 202, nil
		}
		_, statusCode, err := wrapper.doOperation(fn)
		if err != nil {
			res.WriteHeader(statusCode)
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Errorf("killing task")
			wrapper.driverErrc <- err
			return
		}
		res.WriteHeader(statusCode)
	}

	purgeTaskHandler := func(res http.ResponseWriter, req *http.Request, params martini.Params) {
		fn := func() ([]byte, int, error) {
			taskId := params["task_id"]
			err := lx_core_helpers.PurgeTask(wrapper.state, taskId)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Errorf("could not handle purge task request")
				return empty, 500, errors.New("could not handle purge task request", err)
			}
			logrus.WithFields(logrus.Fields{"task_id": taskId}).Infof("removed task from LayerX")
			return empty, 202, nil
		}
		_, statusCode, err := wrapper.doOperation(fn)
		if err != nil {
			res.WriteHeader(statusCode)
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Errorf("purging task")
			wrapper.driverErrc <- err
			return
		}
		res.WriteHeader(statusCode)
	}

	submitResourceHandler := func(res http.ResponseWriter, req *http.Request) {
		fn := func() ([]byte, int, error) {
			data, err := ioutil.ReadAll(req.Body)
			if req.Body != nil {
				defer req.Body.Close()
			}
			if err != nil {
				return empty, 400, errors.New("parsing submit resource request", err)
			}
			var resource lxtypes.Resource
			err = json.Unmarshal(data, &resource)
			if err != nil {
				return empty, 500, errors.New("could not parse json to resource", err)
			}
			err = lx_core_helpers.SubmitResource(wrapper.state, &resource)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Errorf("could not handle Submit Resource request")
				return empty, 500, errors.New("could not handle SubmitResource request", err)
			}
			logrus.WithFields(logrus.Fields{"resource": resource}).Infof("accepted resource from rpi")
			return empty, 202, nil
		}
		_, statusCode, err := wrapper.doOperation(fn)
		if err != nil {
			res.WriteHeader(statusCode)
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Errorf("processing Submit Resource request")
			wrapper.driverErrc <- err
			return
		}
		res.WriteHeader(statusCode)
	}

	submitStatusUpdateHandler := func(res http.ResponseWriter, req *http.Request) {
		fn := func() ([]byte, int, error) {
			data, err := ioutil.ReadAll(req.Body)
			if req.Body != nil {
				defer req.Body.Close()
			}
			if err != nil {
				return empty, 400, errors.New("parsing submit status request", err)
			}
			var status mesosproto.TaskStatus
			err = proto.Unmarshal(data, &status)
			if err != nil {
				return empty, 500, errors.New("could not parse protobuf to status", err)
			}
			err = lx_core_helpers.ProcessStatusUpdate(wrapper.state, wrapper.state.GetTpiUrl(), &status)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Errorf("could not handle Submit status request")
				return empty, 500, errors.New("could not handle submit status request", err)
			}
			logrus.WithFields(logrus.Fields{"status": status, "message": status.GetMessage()}).Infof("accepted status update from rpi")
			return empty, 202, nil
		}
		_, statusCode, err := wrapper.doOperation(fn)
		if err != nil {
			res.WriteHeader(statusCode)
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Errorf("processing submit status request")
			wrapper.driverErrc <- err
			return
		}
		res.WriteHeader(statusCode)
	}

	getNodesHandler := func(res http.ResponseWriter, req *http.Request) {
		fn := func() ([]byte, int, error) {
			nodes, err := lx_core_helpers.GetNodes(wrapper.state)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Errorf("could not handle Get Nodes request")
				return empty, 500, errors.New("could not handle Get Nodes request", err)
			}
			data, err := json.Marshal(nodes)
			if err != nil {
				logrus.Errorf("could not marshal nodes to json")
				return empty, 500, errors.New("marshalling nodes to json", err)
			}
			logrus.WithFields(logrus.Fields{"nodes": nodes}).Debugf("Replying with Nodes")
			return data, 200, nil
		}
		data, statusCode, err := wrapper.doOperation(fn)
		if err != nil {
			res.WriteHeader(statusCode)
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Errorf("processing Get Nodes message")
			wrapper.driverErrc <- err
			return
		}
		res.Write(data)
	}

	getPendingTasksHandler := func(res http.ResponseWriter, req *http.Request) {
		fn := func() ([]byte, int, error) {
			tasks, err := lx_core_helpers.GetPendingTasks(wrapper.state)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Errorf("could not handle Get Pending Tasks request")
				return empty, 500, errors.New("could not handle Get Pending Tasks request", err)
			}
			data, err := json.Marshal(tasks)
			if err != nil {
				logrus.Errorf("could not marshal tasks to json")
				return empty, 500, errors.New("marshalling tasks to json", err)
			}
			logrus.WithFields(logrus.Fields{"tasks": tasks}).Debugf("Replying with Pending Tasks")
			return data, 200, nil
		}
		data, statusCode, err := wrapper.doOperation(fn)
		if err != nil {
			res.WriteHeader(statusCode)
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Errorf("processing Get Pending Tasks message")
			wrapper.driverErrc <- err
			return
		}
		res.Write(data)
	}

	getStagingTasksHandler := func(res http.ResponseWriter, req *http.Request) {
		fn := func() ([]byte, int, error) {
			tasks, err := lx_core_helpers.GetStagingTasks(wrapper.state)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Errorf("could not handle Get Staging Tasks request")
				return empty, 500, errors.New("could not handle Get Staging Tasks request", err)
			}
			data, err := json.Marshal(tasks)
			if err != nil {
				logrus.Errorf("could not marshal tasks to json")
				return empty, 500, errors.New("marshalling tasks to json", err)
			}
			logrus.WithFields(logrus.Fields{"tasks": tasks}).Debugf("Replying with Staging Tasks")
			return data, 200, nil
		}
		data, statusCode, err := wrapper.doOperation(fn)
		if err != nil {
			res.WriteHeader(statusCode)
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Errorf("processing Get Staging Tasks message")
			wrapper.driverErrc <- err
			return
		}
		res.Write(data)
	}

	assignTasksHandler := func(res http.ResponseWriter, req *http.Request) {
		fn := func() ([]byte, int, error) {
			data, err := ioutil.ReadAll(req.Body)
			if req.Body != nil {
				defer req.Body.Close()
			}
			if err != nil {
				return empty, 400, errors.New("parsing assign tasks request", err)
			}
			var assignTasksMessage layerx_brain_client.BrainAssignTasksMessage
			err = json.Unmarshal(data, &assignTasksMessage)
			if err != nil {
				return empty, 500, errors.New("could not parse json to assign tasks message ("+string(data)+")", err)
			}
			err = lx_core_helpers.AssignTasks(wrapper.state, assignTasksMessage)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
					"data":  string(data),
				}).Errorf("could not handle assign tasks request")
				return empty, 500, errors.New("could not handle assign tasks message", err)
			}
			logrus.WithFields(logrus.Fields{"assignTasksMessage": assignTasksMessage}).Infof("accepted assign tasks message from brain")
			return empty, 202, nil
		}
		_, statusCode, err := wrapper.doOperation(fn)
		if err != nil {
			res.WriteHeader(statusCode)
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Errorf("processing assign tasks request")
			wrapper.driverErrc <- err
			return
		}
		res.WriteHeader(statusCode)
	}

	migrateTasksHandler := func(res http.ResponseWriter, req *http.Request) {
		fn := func() ([]byte, int, error) {
			data, err := ioutil.ReadAll(req.Body)
			if req.Body != nil {
				defer req.Body.Close()
			}
			if err != nil {
				return empty, 400, errors.New("parsing migrate tasks request", err)
			}
			var migrateTasksMessage layerx_brain_client.MigrateTaskMessage
			err = json.Unmarshal(data, &migrateTasksMessage)
			if err != nil {
				return empty, 500, errors.New("could not parse json to migrate tasks message", err)
			}
			err = lx_core_helpers.MigrateTasks(wrapper.state, migrateTasksMessage)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Errorf("could not handle migrate tasks request")
				return empty, 500, errors.New("could not handle migrate tasks message", err)
			}
			logrus.WithFields(logrus.Fields{"migrateTasksMessage": migrateTasksMessage}).Infof("accepted migrate tasks message from brain")
			return empty, 202, nil
		}
		_, statusCode, err := wrapper.doOperation(fn)
		if err != nil {
			res.WriteHeader(statusCode)
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Errorf("processing migrate tasks request")
			wrapper.driverErrc <- err
			return
		}
		res.WriteHeader(statusCode)
	}

	wrapper.m.Post(RegisterTpi, registerTpiHandler)
	wrapper.m.Post(RegisterRpi, registerRpiHandler)
	wrapper.m.Post(RegisterTaskProvider, registerTaskProviderHandler)
	wrapper.m.Post(DeregisterTaskProvider+"/:task_provider_id", deregisterTaskProviderHandler)
	wrapper.m.Get(GetTaskProviders, getTaskProvidersHandler)
	wrapper.m.Get(GetStatusUpdates+"/:task_provider_id", getStatusUpdatesHandler)
	wrapper.m.Get(GetStatusUpdate+"/:task_id", getStatusUpdateHandler)
	wrapper.m.Post(SubmitTask+"/:task_provider_id", submitTaskHandler)
	wrapper.m.Post(KillTask+"/:task_provider_id/:task_id", killTaskHandler)
	wrapper.m.Post(PurgeTask+"/:task_id", purgeTaskHandler)
	wrapper.m.Post(SubmitResource, submitResourceHandler)
	wrapper.m.Post(SubmitStatusUpdate, submitStatusUpdateHandler)
	wrapper.m.Get(GetNodes, getNodesHandler)
	wrapper.m.Get(GetStatusUpdates, getStatusUpdatesHandler)
	wrapper.m.Get(GetPendingTasks, getPendingTasksHandler)
	wrapper.m.Get(GetStagingTasks, getStagingTasksHandler)
	wrapper.m.Post(AssignTasks, assignTasksHandler)
	wrapper.m.Post(MigrateTasks, migrateTasksHandler)

	return wrapper.m
}

func (wrapper *layerxCoreServerWrapper) doOperation(f func() ([]byte, int, error)) ([]byte, int, error) {
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
