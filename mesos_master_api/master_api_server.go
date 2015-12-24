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

func (server *mesosApiServer) RunMasterServer(port int, masterUpidString string, errc chan error) {
	portStr := fmt.Sprintf(":%v", port)
	lxlog.Infof(logrus.Fields{
		"port": port,
	}, "Master Server initialized")

	m := lxmartini.QuietMartini()

	m.Get(GET_MASTER_STATE, func(res http.ResponseWriter) {
		datac := make(chan []byte)
		server.actionQueue.Push(func(){
			data, err := handlers.GetMesosState(masterUpidString)
			if err != nil {
				errc <- lxerrors.New("retreiving master state", err)
				return
			}
			datac <- data
		})
		select {
		case data := <- datac:
			res.Write(data)
		case err := <- errc:
			res.WriteHeader(500)
			lxlog.Errorf(logrus.Fields{
				"port": port,
			}, "Retreiving master state: %s", err.Error())
		}
	})

	m.Get(GET_MASTER_STATE_DEPRECATED, func(res http.ResponseWriter) {
		datac := make(chan []byte)
		server.actionQueue.Push(func(){
			data, err := handlers.GetMesosState(masterUpidString)
			if err != nil {
				errc <- lxerrors.New("retreiving master state", err)
				return
			}
			datac <- data
		})
		select {
		case data := <- datac:
			res.Write(data)
		case err := <- errc:
			res.WriteHeader(500)
			lxlog.Errorf(logrus.Fields{
				"port": port,
			}, "Retreiving master state: %s", err.Error())
		}
	})

	m.Post(MESOS_SCHEDULER_CALL, func(res http.ResponseWriter, req *http.Request) {
		statusCodec := make(chan int)
		server.actionQueue.Push(func() {
			body, err := ioutil.ReadAll(req.Body)
			if req.Body != nil {
				defer req.Body.Close()
			}
			if err != nil {
				lxlog.Errorf(logrus.Fields{
					"error": err,
				}, "could not read  MESOS_SCHEDULER_CALL request body")
				statusCodec <- 500
				errc <- lxerrors.New("could not read  MESOS_SCHEDULER_CALL request body", err)
				return
			}
			requestingFramework := req.Header.Get("Libprocess-From")
			if requestingFramework == "" {
				lxlog.Errorf(logrus.Fields{}, "missing required header: %s", "Libprocess-From")
				statusCodec <- 400
				return
			}
			upid, err := mesos_data.UPIDFromString(requestingFramework)
			if err != nil {
				lxlog.Errorf(logrus.Fields{
					"error": err,
				}, "could not parse pid of requesting framework")
				statusCodec <- 500
				errc <- lxerrors.New("could not parse pid of requesting framework", err)
				return
			}
			err = server.processMesosCall(body, upid)
			if err != nil {
				lxlog.Errorf(logrus.Fields{
					"error": err,
				}, "could not read process scheduler call request")
				statusCodec <- 500
				errc <- lxerrors.New("could not read process scheduler call request", err)
				return
			}
			statusCodec <- 202
		})
		res.WriteHeader(<-statusCodec)
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
		err = handlers.HandleSubscribeRequest(server.tpi, server.frameworkManager, upid, subscribe)
		if err != nil {
			return lxerrors.New("processing subscribe request", err)
		}
		break
	default:
		return lxerrors.New("processing unknown call type: "+callType.String(), nil)
	}

	return nil
}
