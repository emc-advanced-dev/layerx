package mesos_master_api
import (
	"github.com/layer-x/layerx-commons/lxlog"
	"github.com/Sirupsen/logrus"
	"net/http"
	"github.com/layer-x/layerx-mesos-tpi_v2/mesos_master_api/mesos_api_helpers"
	"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/layer-x/layerx-commons/lxactionqueue"
	"github.com/layer-x/layerx-mesos-tpi_v2/mesos_master_api/mesos_data"
	"github.com/layer-x/layerx-mesos-tpi_v2/framework_manager"
	"github.com/mesos/mesos-go/mesosproto"
	"github.com/gogo/protobuf/proto"
	"github.com/layer-x/layerx-core_v2/layerx_tpi"
	"github.com/go-martini/martini"
)

const (
	GET_MASTER_STATE = "/master/state.json"
	GET_MASTER_STATE_DEPRECATED = "/state.json"
	MESOS_SCHEDULER_CALL = "/master/mesos.scheduler.Call"
	REGISTER_FRAMEWORK_MESSAGE = "/master/mesos.internal.RegisterFrameworkMessage"
	REREGISTER_FRAMEWORK_MESSAGE = "/master/mesos.internal.ReregisterFrameworkMessage"
	UNREGISTER_FRAMEWORK_MESSAGE = "/master/mesos.internal.UnregisterFrameworkMessage"
	LAUNCH_TASKS_MESSAGE = "/master/mesos.internal.LaunchTasksMessage"
	RECONCILE_TASKS_MESSAGE = "/master/mesos.internal.ReconcileTasksMessage"
	KILL_TASK_MESSAGE = "/master/mesos.internal.KillTaskMessage"
	STATUS_UPDATE_ACKNOWLEDGEMENT_MESSAGE = "/master/mesos.internal.StatusUpdateAcknowledgementMessage"
	REVIVE_OFFERS_MESSAGE = "/master/mesos.internal.ReviveOffersMessage"
)

var empty = []byte{}

type mesosApiServerWrapper struct {
	actionQueue      lxactionqueue.ActionQueue
	frameworkManager framework_manager.FrameworkManager
	tpi              *layerx_tpi.LayerXTpi
}

func NewMesosApiServerWrapper(tpi *layerx_tpi.LayerXTpi, actionQueue lxactionqueue.ActionQueue, frameworkManager framework_manager.FrameworkManager) *mesosApiServerWrapper {
	return &mesosApiServerWrapper{
		actionQueue: actionQueue,
		frameworkManager: frameworkManager,
		tpi: tpi,
	}
}

func (wrapper *mesosApiServerWrapper) WrapWithMesos(m *martini.ClassicMartini, masterUpidString string, driverErrc chan error) *martini.ClassicMartini {
	getMasterStateHandler := func(res http.ResponseWriter) {
		getStateFn := func() ([]byte, int, error) {
			data, err := mesos_api_helpers.GetMesosState(masterUpidString)
			if err != nil {
				return empty, 500, lxerrors.New("retreiving master state", err)
			}
			return data, 200, nil
		}
		data, statusCode, err := wrapper.queueOperation(getStateFn)
		if err != nil {
			res.WriteHeader(statusCode)
			lxlog.Errorf(logrus.Fields{
				"request_sent_by": masterUpidString,
			}, "Retreiving master state")
			driverErrc <- err
			return
		}
		res.Write(data)
	}
	mesosSchedulerCallHandler := func(req *http.Request, res http.ResponseWriter) {
		processMesosCallFn := func() ([]byte, int, error) {
			upid, data, statusCode, err := mesos_api_helpers.ProcessMesosHttpRequest(req)
			if err != nil {
				return empty, statusCode, lxerrors.New("parsing reregisterFramework request", err)
			}
			err = wrapper.processMesosCall(data, upid)
			if err != nil {
				lxlog.Errorf(logrus.Fields{
					"error": err,
				}, "could not read process scheduler call request")
				return empty, 500, lxerrors.New("could not read process scheduler call request", err)
			}
			return empty, 202, nil
		}
		_, statusCode, err := wrapper.queueOperation(processMesosCallFn)
		if err != nil {
			res.WriteHeader(statusCode)
			lxlog.Errorf(logrus.Fields{
				"error": err.Error(),
				"request_sent_by": masterUpidString,
			}, "processing mesos call message")
			driverErrc <- err
			return
		}
		res.WriteHeader(statusCode)
	}
	registerFrameworkMessageHandler := func(req *http.Request, res http.ResponseWriter) {
		registerFrameworkFn := func() ([]byte, int, error) {
			upid, data, statusCode, err := mesos_api_helpers.ProcessMesosHttpRequest(req)
			if err != nil {
				return empty, statusCode, lxerrors.New("parsing registerFramework request", err)
			}
			var registerRequest mesosproto.RegisterFrameworkMessage
			err = proto.Unmarshal(data, &registerRequest)
			if err != nil {
				return empty, 500, lxerrors.New("could not parse data to protobuf msg Call", err)
			}
			err = mesos_api_helpers.HandleRegisterRequest(wrapper.tpi, wrapper.frameworkManager, upid, registerRequest.GetFramework())
			if err != nil {
				lxlog.Errorf(logrus.Fields{
					"error": err,
				}, "could not handle register framework request")
				return empty, 500, lxerrors.New("could not handle register framework request", err)
			}
			return empty, 202, nil
		}
		_, statusCode, err := wrapper.queueOperation(registerFrameworkFn)
		if err != nil {
			res.WriteHeader(statusCode)
			lxlog.Errorf(logrus.Fields{
				"error": err.Error(),
				"request_sent_by": masterUpidString,
			}, "processing register framework message")
			driverErrc <- err
			return
		}
		res.WriteHeader(statusCode)
	}
	reregisterFrameworkMessageHandler := func(req *http.Request, res http.ResponseWriter) {
		reregisterFrameworkFn := func() ([]byte, int, error) {
			upid, data, statusCode, err := mesos_api_helpers.ProcessMesosHttpRequest(req)
			if err != nil {
				return empty, statusCode, lxerrors.New("parsing reregisterFramework request", err)
			}
			var reregisterRequest mesosproto.ReregisterFrameworkMessage
			err = proto.Unmarshal(data, &reregisterRequest)
			if err != nil {
				return empty, 500, lxerrors.New("could not unmarshal data to reregisterRequest", err)
			}
			err = mesos_api_helpers.HandleRegisterRequest(wrapper.tpi, wrapper.frameworkManager, upid, reregisterRequest.GetFramework())
			if err != nil {
				lxlog.Errorf(logrus.Fields{
					"error": err,
				}, "could not handle reregister framework request")
				return empty, 500, lxerrors.New("could not handle reregister framework request", err)
			}
			return empty, 202, nil
		}
		_, statusCode, err := wrapper.queueOperation(reregisterFrameworkFn)
		if err != nil {
			res.WriteHeader(statusCode)
			lxlog.Errorf(logrus.Fields{
				"error": err.Error(),
				"request_sent_by": masterUpidString,
			}, "processing reregister framework message")
			driverErrc <- err
			return
		}
		res.WriteHeader(statusCode)
	}
	unregisterFrameworkMessageHandler := func(req *http.Request, res http.ResponseWriter) {
		unregisterFrameworkFn := func() ([]byte, int, error) {
			_, data, statusCode, err := mesos_api_helpers.ProcessMesosHttpRequest(req)
			if err != nil {
				return empty, statusCode, lxerrors.New("parsing unregisterFramework request", err)
			}
			var unregisterRequest mesosproto.UnregisterFrameworkMessage
			err = proto.Unmarshal(data, &unregisterRequest)
			if err != nil {
				return empty, 500, lxerrors.New("could unmarshal data to unregister request", err)
			}
			err = mesos_api_helpers.HandleRemoveFramework(wrapper.tpi, unregisterRequest.GetFrameworkId().GetValue())
			if err != nil {
				lxlog.Errorf(logrus.Fields{
					"error": err,
				}, "could not handle register framework request")
				return empty, 500, lxerrors.New("could not handle unregister framework request", err)
			}
			return empty, 202, nil
		}
		_, statusCode, err := wrapper.queueOperation(unregisterFrameworkFn)
		if err != nil {
			res.WriteHeader(statusCode)
			lxlog.Errorf(logrus.Fields{
				"error": err.Error(),
				"request_sent_by": masterUpidString,
			}, "processing unregister framework message")
			driverErrc <- err
			return
		}
		res.WriteHeader(statusCode)
	}
	launchTasksMessageHandler := func(req *http.Request, res http.ResponseWriter) {
		reregisterFrameworkFn := func() ([]byte, int, error) {
			_, data, statusCode, err := mesos_api_helpers.ProcessMesosHttpRequest(req)
			if err != nil {
				return empty, statusCode, lxerrors.New("parsing launchTasks request", err)
			}
			var launchTasksMessage mesosproto.LaunchTasksMessage
			err = proto.Unmarshal(data, &launchTasksMessage)
			if err != nil {
				return empty, 500, lxerrors.New("could not unmarshal data to launchTasks", err)
			}
			frameworkId := launchTasksMessage.GetFrameworkId().GetValue()
			mesosTasks := launchTasksMessage.GetTasks()
			err = mesos_api_helpers.HandleLaunchTasksRequest(wrapper.tpi, frameworkId, mesosTasks)
			if err != nil {
				lxlog.Errorf(logrus.Fields{
					"error": err,
				}, "could not handle launch tasks request")
				return empty, 500, lxerrors.New("could not handle launchTasks request", err)
			}
			return empty, 202, nil
		}
		_, statusCode, err := wrapper.queueOperation(reregisterFrameworkFn)
		if err != nil {
			res.WriteHeader(statusCode)
			lxlog.Errorf(logrus.Fields{
				"error": err.Error(),
				"request_sent_by": masterUpidString,
			}, "processing launchTasks message")
			driverErrc <- err
			return
		}
		res.WriteHeader(statusCode)
	}
	reconcileTasksMessageHandler := func(req *http.Request, res http.ResponseWriter) {
		reconcileTasksFn := func() ([]byte, int, error) {
			upid, data, statusCode, err := mesos_api_helpers.ProcessMesosHttpRequest(req)
			if err != nil {
				return empty, statusCode, lxerrors.New("parsing reconcile tasks request", err)
			}
			var reconcileTasksMessage mesosproto.ReconcileTasksMessage
			err = proto.Unmarshal(data, &reconcileTasksMessage)
			if err != nil {
				return empty, 500, lxerrors.New("could not unmarshal data to reconcile tasks", err)
			}
			err = mesos_api_helpers.HandleReconcileTasksRequest(wrapper.tpi, wrapper.frameworkManager, upid, reconcileTasksMessage)
			if err != nil {
				lxlog.Errorf(logrus.Fields{
					"error": err,
				}, "could not handle reconcile tasks request")
				return empty, 500, lxerrors.New("could not handle reconcile tasks request", err)
			}
			return empty, 202, nil
		}
		_, statusCode, err := wrapper.queueOperation(reconcileTasksFn)
		if err != nil {
			res.WriteHeader(statusCode)
			lxlog.Errorf(logrus.Fields{
				"error": err.Error(),
				"request_sent_by": masterUpidString,
			}, "processing reconcileTasks message")
			driverErrc <- err
			return
		}
		res.WriteHeader(statusCode)
	}
	killTaskMessageHandler := func(req *http.Request, res http.ResponseWriter) {
		reconcileTasksFn := func() ([]byte, int, error) {
			_, data, statusCode, err := mesos_api_helpers.ProcessMesosHttpRequest(req)
			if err != nil {
				return empty, statusCode, lxerrors.New("parsing reconcile tasks request", err)
			}
			var killTaskMessage mesosproto.KillTaskMessage
			err = proto.Unmarshal(data, &killTaskMessage)
			if err != nil {
				return empty, 500, lxerrors.New("could not unmarshal data to killTaskMessage", err)
			}
			taskId := killTaskMessage.GetTaskId().GetValue()
			err = mesos_api_helpers.HandleKillTaskRequest(wrapper.tpi, taskId)
			if err != nil {
				lxlog.Errorf(logrus.Fields{
					"error": err,
				}, "could not handle kill task request")
				return empty, 500, lxerrors.New("could not handle kill task request", err)
			}
			return empty, 202, nil
		}
		_, statusCode, err := wrapper.queueOperation(reconcileTasksFn)
		if err != nil {
			res.WriteHeader(statusCode)
			lxlog.Errorf(logrus.Fields{
				"error": err.Error(),
				"request_sent_by": masterUpidString,
			}, "processing kill task message")
			driverErrc <- err
			return
		}
		res.WriteHeader(statusCode)
	}

	m.Get(GET_MASTER_STATE, getMasterStateHandler)
	m.Get(GET_MASTER_STATE_DEPRECATED, getMasterStateHandler)
	m.Post(MESOS_SCHEDULER_CALL, mesosSchedulerCallHandler)
	m.Post(REGISTER_FRAMEWORK_MESSAGE, registerFrameworkMessageHandler)
	m.Post(REREGISTER_FRAMEWORK_MESSAGE, reregisterFrameworkMessageHandler)
	m.Post(UNREGISTER_FRAMEWORK_MESSAGE, unregisterFrameworkMessageHandler)
	m.Post(LAUNCH_TASKS_MESSAGE, launchTasksMessageHandler)
	m.Post(RECONCILE_TASKS_MESSAGE, reconcileTasksMessageHandler)
	m.Post(KILL_TASK_MESSAGE, killTaskMessageHandler)
	return m
}

func (wrapper *mesosApiServerWrapper) queueOperation(f func() ([]byte, int, error)) ([]byte, int, error) {
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

func (wrapper *mesosApiServerWrapper) processMesosCall(data []byte, upid *mesos_data.UPID) error {
	var call mesosproto.Call
	err := proto.Unmarshal(data, &call)
	if err != nil {
		return lxerrors.New("could not parse data to protobuf msg Call", err)
	}
	callType := call.GetType()
	lxlog.Debugf(logrus.Fields{
		"call_type":     callType.String(),
		"framework_pid": upid.String(),
		"whole call":    call.String(),
	}, "Received mesosproto.Call")

	switch callType {
	case mesosproto.Call_SUBSCRIBE:
		subscribe := call.Subscribe
		err = mesos_api_helpers.HandleRegisterRequest(wrapper.tpi, wrapper.frameworkManager, upid, subscribe.GetFrameworkInfo())
		if err != nil {
			return lxerrors.New("processing subscribe request", err)
		}
		break
	default:
		return lxerrors.New("processing unknown call type: " + callType.String(), nil)
	}

	return nil
}
