package mesos_master_api
import (
"fmt"
"github.com/layer-x/layerx-commons/lxlog"
	"github.com/Sirupsen/logrus"
	"github.com/layer-x/layerx-commons/lxmartini"
"net/http"
	"github.com/layer-x/layerx-mesos-tpi_v2/mesos_master_api/handlers"
	"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/layer-x/layerx-commons/lxactionqueue"
"io/ioutil"
	"github.com/layer-x/layerx-mesos-tpi_v2/mesos_master_api/mesos_data"
	"github.com/layer-x/layerx-mesos-tpi_v2/framework_manager"
	"github.com/mesos/mesos-go/mesosproto"
	"github.com/gogo/protobuf/proto"
	"github.com/layer-x/layerx-core_v2/layerx_tpi"
)

const (
	GET_MASTER_STATE                      = "/master/state.json"
	GET_MASTER_STATE_DEPRECATED           = "/state.json"
	MESOS_SCHEDULER_CALL                  = "/master/mesos.scheduler.Call"
	REGISTER_FRAMEWORK_MESSAGE            = "/master/mesos.internal.RegisterFrameworkMessage"
	REREGISTER_FRAMEWORK_MESSAGE          = "/master/mesos.internal.ReregisterFrameworkMessage"
	UNREGISTER_FRAMEWORK_MESSAGE          = "/master/mesos.internal.UnregisterFrameworkMessage"
	LAUNCH_TASKS_MESSAGE                  = "/master/mesos.internal.LaunchTasksMessage"
	RECONCILE_TASKS_MESSAGE               = "/master/mesos.internal.ReconcileTasksMessage"
	KILL_TASK_MESSAGE                     = "/master/mesos.internal.KillTaskMessage"
	STATUS_UPDATE_ACKNOWLEDGEMENT_MESSAGE = "/master/mesos.internal.StatusUpdateAcknowledgementMessage"
	REVIVE_OFFERS_MESSAGE                 = "/master/mesos.internal.ReviveOffersMessage"
)

var empty = []byte{}

type mesosApiServer struct {
	actionQueue lxactionqueue.ActionQueue
	frameworkManager framework_manager.FrameworkManager
	tpi *layerx_tpi.LayerXTpi
}

func NewMesosApiServer(tpi *layerx_tpi.LayerXTpi, actionQueue lxactionqueue.ActionQueue, frameworkManager framework_manager.FrameworkManager) *mesosApiServer {
	return &mesosApiServer{
		actionQueue: actionQueue,
		frameworkManager: frameworkManager,
		tpi: tpi,
	}
}

func (server *mesosApiServer) queueOperation(f func() ([]byte, int, error)) ([]byte, int, error) {
	datac := make(chan []byte)
	statusCodec := make(chan int)
	errc := make(chan error)
	server.actionQueue.Push(
	func(){
		data, statusCode, err := f()
		datac <- data
		statusCodec <- statusCode
		errc <- err
	})
	return <-datac, <-statusCodec, <-errc
}

func (server *mesosApiServer) RunMasterServer(port int, masterUpidString string, driverErrc chan error) {
	portStr := fmt.Sprintf(":%v", port)
	lxlog.Infof(logrus.Fields{
		"port": port,
	}, "Master Server initialized")

	m := lxmartini.QuietMartini()

	m.Get(GET_MASTER_STATE, func(res http.ResponseWriter) {
		getStateFn := func() ([]byte, int, error) {
			data, err := handlers.GetMesosState(masterUpidString)
			if err != nil {
				return empty, 500, lxerrors.New("retreiving master state", err)
			}
			return data, 200, nil
		}
		data, statusCode, err := server.queueOperation(getStateFn)
		if err != nil {
			res.WriteHeader(statusCode)
			lxlog.Errorf(logrus.Fields{
				"request_sent_by": masterUpidString,
			}, "Retreiving master state")
			driverErrc <- err
			return
		}
		res.Write(data)
	})

	m.Get(GET_MASTER_STATE_DEPRECATED, func(res http.ResponseWriter) {
		getStateFn := func() ([]byte, int, error) {
			data, err := handlers.GetMesosState(masterUpidString)
			if err != nil {
				return empty, 500, lxerrors.New("retreiving master state", err)
			}
			return data, 200, nil
		}
		data, statusCode, err := server.queueOperation(getStateFn)
		if err != nil {
			res.WriteHeader(statusCode)
			lxlog.Errorf(logrus.Fields{
				"request_sent_by": masterUpidString,
			}, "Retreiving master state")
			driverErrc <- err
			return
		}
		res.Write(data)
	})

	m.Post(MESOS_SCHEDULER_CALL, func(res http.ResponseWriter, req *http.Request) {
		processMesosCallFn := func() ([]byte, int, error) {
			data, err := ioutil.ReadAll(req.Body)
			if req.Body != nil {
				defer req.Body.Close()
			}
			if err != nil {
				lxlog.Errorf(logrus.Fields{
					"error": err,
				}, "could not read  MESOS_SCHEDULER_CALL request body")
				return empty, 500, lxerrors.New("could not read  MESOS_SCHEDULER_CALL request body", err)
			}
			requestingFramework := req.Header.Get("Libprocess-From")
			if requestingFramework == "" {
				lxlog.Errorf(logrus.Fields{}, "missing required header: %s", "Libprocess-From")
				return	empty, 400, nil
			}
			upid, err := mesos_data.UPIDFromString(requestingFramework)
			if err != nil {
				lxlog.Errorf(logrus.Fields{
					"error": err,
				}, "could not parse pid of requesting framework")
				return empty, 500, lxerrors.New("could not parse pid of requesting framework", err)
			}
			err = server.processMesosCall(data, upid)
			if err != nil {
				lxlog.Errorf(logrus.Fields{
					"error": err,
				}, "could not read process scheduler call request")
				return empty, 500, lxerrors.New("could not read process scheduler call request", err)
			}
			return empty, 202, nil
		}
		_, statusCode, err := server.queueOperation(processMesosCallFn)
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
	})

	m.Post(REGISTER_FRAMEWORK_MESSAGE, func(res http.ResponseWriter, req *http.Request) {
		registerFrameworkFn := func() ([]byte, int, error) {
			data, err := ioutil.ReadAll(req.Body)
			if req.Body != nil {
				defer req.Body.Close()
			}
			if err != nil {
				lxlog.Errorf(logrus.Fields{
					"error": err,
				}, "could not read  REGISTER_FRAMEWORK_MESSAGE request body")
				return empty, 500, lxerrors.New("could not read  REGISTER_FRAMEWORK_MESSAGE request body", err)
			}
			requestingFramework := req.Header.Get("Libprocess-From")
			if requestingFramework == "" {
				lxlog.Errorf(logrus.Fields{}, "missing required header: %s", "Libprocess-From")
				return empty, 400, nil
			}
			upid, err := mesos_data.UPIDFromString(requestingFramework)
			if err != nil {
				lxlog.Errorf(logrus.Fields{
					"error": err,
				}, "could not parse pid of requesting framework")
				return empty, 500, lxerrors.New("could not parse pid of requesting framework", err)
			}
			var registerRequest mesosproto.RegisterFrameworkMessage
			err = proto.Unmarshal(data, &registerRequest)
			if err != nil {
				return empty, 500, lxerrors.New("could not parse data to protobuf msg Call", err)
			}
			err = handlers.HandleRegisterRequest(server.tpi, server.frameworkManager, upid, registerRequest.GetFramework())
			if err != nil {
				lxlog.Errorf(logrus.Fields{
					"error": err,
				}, "could not handle register framework request")
				return empty, 500, lxerrors.New("could not handle register framework request", err)
			}
			return empty, 202, nil
		}
		_, statusCode, err := server.queueOperation(registerFrameworkFn)
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
	})

		m.RunOnAddr(portStr)
}



func (server *mesosApiServer) processMesosCall(data []byte, upid *mesos_data.UPID) error {
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
		err = handlers.HandleRegisterRequest(server.tpi, server.frameworkManager, upid, subscribe.GetFrameworkInfo())
		if err != nil {
			return lxerrors.New("processing subscribe request", err)
		}
		break
	default:
		return lxerrors.New("processing unknown call type: "+callType.String(), nil)
	}

	return nil
}
